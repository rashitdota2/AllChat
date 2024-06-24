package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"workwithimages/domain/models"
	"workwithimages/internalls/service"
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
	info := ctx.GetHeader("Authorization")
	if info == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	infos := strings.Split(info, " ")
	prefix, tokenstr := infos[0], infos[1]
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
	if prefix == "Refresh" {
		acs, rt, err := service.GetTokens(claims.UserId, claims.Name)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{"status": "error"})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusResetContent, gin.H{
			"access":  acs,
			"refresh": rt,
		})
		return
	}
	ctx.Set("claims", claims)
	ctx.Next()
}
