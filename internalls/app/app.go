package app

import (
	"github.com/gorilla/websocket"
	"log"
	"workwithimages/configs"
	"workwithimages/internalls/handler"
	"workwithimages/internalls/repository"
	"workwithimages/internalls/server"
	"workwithimages/internalls/service"
	"workwithimages/internalls/ws"
	"workwithimages/pkg/logger"
	postgres "workwithimages/pkg/pgx"

	//postgres "workwithimages/pkg/pgx"
	"workwithimages/pkg/redis"
)

func Run() {

	db, err := postgres.NewPostgresDB(configs.Psql)
	if err != nil {
		log.Fatal(err)
	}

	redix, err := redis.NewRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	logs, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewRepo(db, redix)

	service := service.NewService(repo, logs, redix)

	var upgrader = websocket.Upgrader{}

	handl := handler.NewHandler(service, &upgrader)

	go ws.WSRun()

	err = server.RUN(configs.HOST, handl.Init())
	if err != nil {
		log.Fatal(err)
	}

}
