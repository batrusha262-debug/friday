package pack

//go:generate go run git.appkode.ru/pub/go/genum/cmd/genum@v0.1.3 --path=./domain/enum

import (
	"context"

	"github.com/google/uuid"

	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
	"friday/internal/pack/entity"
)

type Repository interface {
	CreateUser(context.Context, string) (entity.User, error)
	ListUsers(context.Context) ([]entity.User, error)

	CreatePack(context.Context, string, uuid.UUID) (entity.Pack, error)
	ListPacks(context.Context) ([]entity.Pack, error)
	GetPack(context.Context, uuid.UUID) (entity.Pack, error)
	DeletePack(context.Context, uuid.UUID) error

	CreateRound(context.Context, uuid.UUID, string, enum.RoundTypeEnum) (entity.Round, error)
	ListRounds(context.Context, uuid.UUID) ([]entity.Round, error)
	DeleteRound(context.Context, uuid.UUID) error

	CreateCategory(context.Context, uuid.UUID, string) (entity.Category, error)
	ListCategories(context.Context, uuid.UUID) ([]entity.Category, error)
	DeleteCategory(context.Context, uuid.UUID) error

	GetQuestion(context.Context, uuid.UUID) (entity.Question, error)
	CreateQuestion(context.Context, uuid.UUID, values.Question) (entity.Question, error)
	ListQuestions(context.Context, uuid.UUID) ([]entity.Question, error)
	UpdateQuestion(context.Context, uuid.UUID, values.Question) (entity.Question, error)
	DeleteQuestion(context.Context, uuid.UUID) error

	CreateGame(context.Context, uuid.UUID, uuid.UUID) (entity.Game, error)
	GetGame(context.Context, uuid.UUID) (entity.Game, error)
	DeleteGame(context.Context, uuid.UUID) error
	StartGame(context.Context, uuid.UUID) (entity.Game, error)
	FinishGame(context.Context, uuid.UUID) (entity.Game, error)

	AddGameTeam(context.Context, uuid.UUID, string) (entity.GameTeam, error)
	ListGameTeams(context.Context, uuid.UUID) ([]entity.GameTeam, error)
	RemoveGameTeam(context.Context, uuid.UUID) error
	AwardTeamPoints(context.Context, uuid.UUID, int) error

	SetCurrentPicker(context.Context, uuid.UUID, *uuid.UUID) error

	MarkQuestionAnswered(context.Context, uuid.UUID, uuid.UUID, *uuid.UUID) (entity.GameQuestionState, error)
	ListGameQuestionStates(context.Context, uuid.UUID) ([]entity.GameQuestionState, error)
}
