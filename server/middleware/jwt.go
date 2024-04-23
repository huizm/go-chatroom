package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/huizm/go-chatroom/logic"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if len(token) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "token not provided",
			})
			return
		}

		claims, err := logic.ParseToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"error":   err.Error(),
			})
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
