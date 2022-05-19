package service

import "monitor-apiserver/internal/repository"

var Svc Service

type Service interface {
	Users() UserService
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) Users() UserService {
	return newUserService(s)
}
