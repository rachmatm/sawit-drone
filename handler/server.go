package handler

import "github.com/SawitProRecruitment/UserService/repository"

type Server struct {
	repo repository.EstateRepository
}

func NewServer(repo repository.EstateRepository) *Server {
	return &Server{
		repo: repo,
	}
}
