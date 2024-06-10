package service

import "workwithimages/internalls/repository"

type Service struct {
	Repo *repository.Repo
}

func NewService(repo *repository.Repo) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Serv() error {
	return nil
}
