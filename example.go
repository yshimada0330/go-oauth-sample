package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"

	"github.com/yshimada0330/go-oauth-sample/handler"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})
	manager.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(manager)
	// srv.SetAllowGetAccessRequest(true)

	// NOTE: Authorization: Basic による認証の場合
	// $ curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -H "Authorization: Basic MDAwMDAwOjk5OTk5OQ==" -d "grant_type=client_credentials" http://localhost:8080/token
	// srv.SetClientInfoHandler(server.ClientBasicHandler)

	// NOTE: POSTパラメータ or GETパラメータ による認証の場合
	// $ curl --request POST \
	// --url 'http://localhost:8080/token' \
	// --header 'content-type: application/x-www-form-urlencoded' \
	// --data grant_type=client_credentials \
	// --data client_id=000000 \
	// --data client_secret=999999
	//
	// $ curl -X POST 'http://localhost:8080/token?grant_type=client_credentials&client_id=000000&client_secret=999999'
	// srv.SetClientInfoHandler(server.ClientFormHandler)

	srv.SetClientInfoHandler(clientHandler)

	srv.SetClientScopeHandler(clientScopeHandler)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	r.Use(func(srv *server.Server) gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("hoge", srv)
		}
	}(srv))

	r.POST("/token", func(c *gin.Context) {
		err := srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"message": err.Error()})
			return
		}
	})

	// curl -X GET "http://private/localhost:8080/test"  -H "Authorization: Bearer {TOKEN}"
	r.GET("/test", handler.Test)

	r.Run() // 0.0.0.0:8080 でサーバーを立てます。
}

// NOTE:
// 以下のように、ClientBasicHandler と ClientFormHandler を組み合わせるで、
// - Authorization: Basic による認証
// - POSTパラメータ or GETパラメータ による認証
// の両方に対応できるようになる。
func clientHandler(r *http.Request) (string, string, error) {
	clientID, clientSecret, err := server.ClientBasicHandler(r)
	if err == nil {
		return clientID, clientSecret, nil
	}

	clientID, clientSecret, err = server.ClientFormHandler(r)
	if err == nil {
		return clientID, clientSecret, nil
	}

	return "", "", err
}

// NOTE:
// - HandleAuthorizeRequest
// - パスワード、クライアントクレデンシャルのトークン発行時
// に呼ばれる
// ValidationBearerToken では呼ばれない
func clientScopeHandler(tgr *oauth2.TokenGenerateRequest) (bool, error) {
	log.Println("URL", tgr.Request.URL)
	log.Println("ClientID:", tgr.ClientID)
	log.Println("scope:", tgr.Scope)
	return true, nil
}
