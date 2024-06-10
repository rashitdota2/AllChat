package handler

import (
	"workwithimages/internalls/models"
	"workwithimages/internalls/service"
)

type Handler struct {
	Serv   *service.Service
	Claims *models.Claims
}

func NewHandler(s *service.Service, c *models.Claims) *Handler {
	return &Handler{
		Serv:   s,
		Claims: c,
	}
}
