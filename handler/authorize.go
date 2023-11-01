package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"

	"github.com/yshimada0330/go-oauth-sample/repository"
)

type clientStore interface {
	FindByClientId(clientId string) (repository.Client, error)
}

type authorizeHandler struct {
	store clientStore
}

func NewAuthorizeHandler(s clientStore) *authorizeHandler {
	return &authorizeHandler{store: s}
}

func (h authorizeHandler) Test(c *gin.Context) {
	srv, ok := c.Get("hoge")
	if !ok {
		if s, ok := srv.(*server.Server); ok && s != nil {
			fmt.Println("s:", s)
			fmt.Println("error")
			return
		}
	}

	fmt.Println("httpMethod: ", c.Request.Method)
	fmt.Println("url: ", c.Request.URL)
	fmt.Println("srv:", srv)

	token, err := srv.(*server.Server).ValidationBearerToken(c.Request)
	// NOTE: scopeのチェックは、自前で実装して、errors.ErrInvalidScope のようなエラーを返すようにしないといけない

	client, err := h.store.FindByClientId(token.GetClientID())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	scope := strings.Split(client.Scope, " ")
	for _, s := range scope {
		fmt.Printf("%s\n", s)
	}

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

	c.JSON(http.StatusOK, gin.H{"message": "Authorized successfully "})
}
