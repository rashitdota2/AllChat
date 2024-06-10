package app

import (
	"workwithimages/internalls/handler"
	"workwithimages/internalls/models"
	"workwithimages/internalls/repository"
	"workwithimages/internalls/service"
	"workwithimages/pkg/pgx"
)

func Run() {

	db := pgx.NewDB()
	Repo := repository.NewRepo(db)
	Service := service.NewService(Repo)
	h := handler.NewHandler(Service, &models.Claims{})

	h.Init()
}
