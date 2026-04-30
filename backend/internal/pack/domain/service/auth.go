package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"git.appkode.ru/pub/go/failure"
	"github.com/google/uuid"

	"friday/internal/pack/domain/values"
)

var adminWhitelist = map[string]bool{
	"alina.zibarova.zb@gmail.com":  true,
	"korenda.mariia.kmr@gmail.com": true,
	"baturin.andrei.adt@gmail.com": true,
	"batrusha262@gmail.com":        true,
}

func (s *Service) RequestCode(ctx context.Context, email string) error {
	if !adminWhitelist[email] {
		return failure.NewForbiddenError("email not authorized")
	}

	code := fmt.Sprintf("%06d", rand.Intn(1_000_000))
	expiresAt := time.Now().Add(10 * time.Minute)

	_, err := s.repo.CreateAuthCode(ctx, email, code, expiresAt)
	if err != nil {
		return err
	}

	return s.mailer.Send(
		email,
		"Код входа в Своя игра",
		fmt.Sprintf("Ваш код подтверждения: %s\n\nКод действителен 10 минут.", code),
	)
}

func (s *Service) VerifyCode(ctx context.Context, email, code string) (values.Session, error) {
	if !adminWhitelist[email] {
		return values.Session{}, failure.NewForbiddenError("email not authorized")
	}

	_, err := s.repo.UseAuthCode(ctx, email, code)
	if err != nil {
		return values.Session{}, err
	}

	user, err := s.repo.UpsertAdminUser(ctx, email)
	if err != nil {
		return values.Session{}, err
	}

	token := uuid.New().String()

	_, err = s.repo.CreateSession(ctx, user.ID, token)
	if err != nil {
		return values.Session{}, err
	}

	return values.Session{
		Token: token,
		User:  user.ToDomain(),
	}, nil
}

func (s *Service) CreateGuestSession(ctx context.Context, name string) (values.Session, error) {
	if name == "" {
		return values.Session{}, failure.NewInvalidArgumentError("name is required")
	}

	user, err := s.repo.CreateGuestUser(ctx, name)
	if err != nil {
		return values.Session{}, err
	}

	token := uuid.New().String()

	_, err = s.repo.CreateSession(ctx, user.ID, token)
	if err != nil {
		return values.Session{}, err
	}

	return values.Session{
		Token: token,
		User:  user.ToDomain(),
	}, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.repo.DeleteSession(ctx, token)
}

func (s *Service) GetSessionUser(ctx context.Context, token string) (values.User, error) {
	e, err := s.repo.GetSessionUser(ctx, token)
	if err != nil {
		return values.User{}, err
	}

	return e.ToDomain(), nil
}
