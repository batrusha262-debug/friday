package service

import "friday/internal/pack"

type Mailer interface {
	Send(to, subject, body string) error
}

type Service struct {
	repo   pack.Repository
	mailer Mailer
}

func NewService(repo pack.Repository, mailer Mailer) *Service {
	return &Service{repo: repo, mailer: mailer}
}
