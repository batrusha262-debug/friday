package values

import (
	"time"

	"github.com/google/uuid"

	"friday/internal/pack/domain/enum"
)

type User struct {
	ID        UserID    `json:"id"`
	Username  string    `json:"username"`
	Email     *string   `json:"email,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type Session struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type Pack struct {
	ID        PackID    `json:"id"`
	Title     string    `json:"title"`
	AuthorID  uuid.UUID `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Round struct {
	ID       RoundID            `json:"id"`
	PackID   PackID             `json:"pack_id"`
	Name     string             `json:"name"`
	Type     enum.RoundTypeEnum `json:"type"`
	OrderNum int16              `json:"order_num"`
}

type Category struct {
	ID       CategoryID `json:"id"`
	RoundID  RoundID    `json:"round_id"`
	Name     string     `json:"name"`
	OrderNum int16      `json:"order_num"`
}

type Question struct {
	ID            QuestionID            `json:"id"`
	CategoryID    CategoryID            `json:"category_id"`
	Price         int                   `json:"price"`
	Type          enum.QuestionTypeEnum `json:"type"`
	Question      string                `json:"question"`
	Answer        string                `json:"answer"`
	Comment       *string               `json:"comment,omitempty"`
	MediaURL      *string               `json:"media_url,omitempty"`
	OrderNum      int16                 `json:"order_num"`
	Options       []string              `json:"options"`
	CorrectOption int16                 `json:"correct_option"`
}

type Game struct {
	ID              GameID              `json:"id"`
	PackID          PackID              `json:"pack_id"`
	HostID          uuid.UUID           `json:"host_id"`
	Status          enum.GameStatusEnum `json:"status"`
	IsOpen          bool                `json:"is_open"`
	CreatedAt       time.Time           `json:"created_at"`
	StartedAt       *time.Time          `json:"started_at,omitempty"`
	FinishedAt      *time.Time          `json:"finished_at,omitempty"`
	CurrentPickerID *GameTeamID         `json:"current_picker_id,omitempty"`
}

type GameTeam struct {
	ID       GameTeamID `json:"id"`
	GameID   GameID     `json:"game_id"`
	Name     string     `json:"name"`
	Score    int        `json:"score"`
	OrderNum int16      `json:"order_num"`
}

type GameQuestionState struct {
	ID             uuid.UUID   `json:"id"`
	GameID         GameID      `json:"game_id"`
	QuestionID     QuestionID  `json:"question_id"`
	AnsweredBy     *GameTeamID `json:"answered_by,omitempty"`
	AnsweredAt     *time.Time  `json:"answered_at,omitempty"`
	RevealedCount  int16       `json:"revealed_count"`
	TimerStartedAt *time.Time  `json:"timer_started_at,omitempty"`
	WrongOptions   []int16     `json:"wrong_options"`
}

type AnswerClaim struct {
	ID         uuid.UUID   `json:"id"`
	GameID     GameID      `json:"game_id"`
	QuestionID QuestionID  `json:"question_id"`
	TeamID     GameTeamID  `json:"team_id"`
	ClaimedAt  time.Time   `json:"claimed_at"`
	Status     string      `json:"status"`
	ReviewedAt *time.Time  `json:"reviewed_at,omitempty"`
}

type GameBoard struct {
	Teams         []GameTeam          `json:"teams"`
	States        []GameQuestionState `json:"states"`
	PendingClaims []AnswerClaim       `json:"pending_claims"`
	MiniGame      *MiniGame           `json:"mini_game,omitempty"`
}

type MiniGame struct {
	ID             uuid.UUID  `json:"id"`
	GameID         GameID     `json:"game_id"`
	QuestionID     QuestionID `json:"question_id"`
	ExcludedTeamID *uuid.UUID `json:"excluded_team_id,omitempty"`
	PosX           int16      `json:"pos_x"`
	PosY           int16      `json:"pos_y"`
	StartedAt      time.Time  `json:"started_at"`
	AppearsAt      time.Time  `json:"appears_at"`
	WinnerTeamID   *uuid.UUID `json:"winner_team_id,omitempty"`
	FinishedAt     *time.Time `json:"finished_at,omitempty"`
}
