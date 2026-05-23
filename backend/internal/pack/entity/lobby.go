package entity

import (
	"friday/internal/pack/domain/values"
)

type Lobby struct {
	GameID    values.GameID `db:"game_id"`
	PackID    values.PackID `db:"pack_id"`
	PackTitle string        `db:"pack_title"`
	TeamCount int           `db:"team_count"`
	IsOpen    bool          `db:"is_open"`
}

func (e Lobby) ToDomain() values.Lobby {
	return values.Lobby{
		GameID:    e.GameID,
		PackID:    e.PackID,
		PackTitle: e.PackTitle,
		TeamCount: e.TeamCount,
		IsOpen:    e.IsOpen,
	}
}
