package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/server"
)

func Test(c *gin.Context) {
	srv, ok := c.Get("hoge")
	if !ok {
		if s, ok := srv.(*server.Server); ok && s != nil {
			fmt.Println("s:", s)
			fmt.Println("error")
			return
		}
	}

	fmt.Println("srv:", srv)

	token, err := srv.(*server.Server).ValidationBearerToken(c.Request)
	// NOTE: scopeのチェックは、自前で実装して、errors.ErrInvalidScope のようなエラーを返すようにしないといけない

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

	c.JSON(http.StatusOK, gin.H{"message": "Authorized successfully "})
}
