package service

import "friday/internal/pack"

type Service struct {
	repo pack.Repository
}

func NewService(repo pack.Repository) *Service {
	return &Service{repo: repo}
}
