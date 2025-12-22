package middleware

import (
	"Trade-app/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Auth := ctx.GetHeader("Authorization")
		if Auth == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header missing"})
			ctx.Abort()
			return
		}
		prefix := "Bearer "
		if !strings.HasPrefix(Auth, prefix) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid Authorization"})
			ctx.Abort()
			return
		}
		tokenstr := strings.TrimPrefix(Auth, prefix)
		claim := &token.Claim{}

		token, err := jwt.ParseWithClaims(tokenstr, claim, func(t *jwt.Token) (any, error) {
			return token.Secret_Key, nil
		})
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid or token expired"})
			ctx.Abort()
			return
		}
		ctx.Set("Email", claim.Email)
		ctx.Next()
	}
}
