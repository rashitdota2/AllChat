package handler

import (
	"github.com/gin-gonic/gin"
	"workwithimages/internalls/middleware"
)

func (h *Handler) Init() {
	r := gin.Default()
	r.Use(middleware.Token)

}
