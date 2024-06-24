package handler

import (
	"github.com/gin-gonic/gin"
	"workwithimages/internalls/middleware"
)

func (h *Handler) Init() *gin.Engine {
	r := gin.Default()
	h.Rout(r)
	return r
}
func (h *Handler) Rout(r *gin.Engine) {
	r.Use(middleware.Cors)

	r.POST("/sign-in", h.Sign)
	r.POST("/login", h.Login)
	//
	r.Use(middleware.Auth, h.GetClaims)
	r.POST("/avatar", h.UpdAvatar)
	r.GET("/avatar", h.GetAvatar)
	r.GET("/profile", h.GetProfile)
	r.POST("/profile/update", h.UpdateProfile)
	r.GET("/socket", h.GiveSocket)

}
