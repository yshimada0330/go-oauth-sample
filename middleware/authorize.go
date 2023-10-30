package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
)

func AuthorizeMiddleware(srv *server.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := srv.ValidationBearerToken(c.Request)
		// NOTE: scopeのチェックは、自前で実装して、errors.ErrInvalidScope のようなエラーを返すようにしないといけない

		// NOTE: scopeと該当リソースパス・メソッドのマッピングはmiddlewareでやると辛いかもしれない
		// 以下のように、リソースパス・メソッドを取得することはできるが、これを元に処理するぐらいなら各handlerの関数内で処理した方が良いかもしれない
		fmt.Println("httpMethod: ", c.Request.Method)
		fmt.Println("url: ", c.Request.URL)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"scope":      token.GetScope(),
		}
		e := json.NewEncoder(c.Writer)
		e.SetIndent("", "  ")
		e.Encode(data)

		fmt.Println("token:", data)

		return
	}
}
