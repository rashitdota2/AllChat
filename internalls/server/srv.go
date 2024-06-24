package server

import (
	"net/http"
	"workwithimages/configs"
)

func RUN(h configs.HTTP, handler http.Handler) error {
	srv := &http.Server{
		Addr:           h.Host + h.Port,
		Handler:        handler,
		ReadTimeout:    h.ReadTimeout,
		WriteTimeout:   h.WriteTimeout,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
