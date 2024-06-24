package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"workwithimages/domain/models"
)

func Cors(ctx *gin.Context) {
	if ctx.Request.Method == "OPTIONS" {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Content-Type", "application/json, image/png, image/jpeg")
		ctx.Header("Content-Security-Policy", "default-src 'self';")
		ctx.AbortWithStatus(200)
	}
	ctx.Next()
}

func Auth(ctx *gin.Context) {
	tokenstr := ctx.GetHeader("Authorization")
	if tokenstr == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims := models.TokenClaims{}
	token, err := jwt.ParseWithClaims(tokenstr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || token == nil {
		ctx.AbortWithStatusJSON(401, gin.H{"status": "error"})
		return
	}
	if !token.Valid {
		ctx.AbortWithStatus(401)
		return
	}
	ctx.Set("claims", claims)
	ctx.Next()
}
