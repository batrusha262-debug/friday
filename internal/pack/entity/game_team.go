package entity

import "friday/internal/pack/domain/values"

type GameTeam struct {
	ID       values.GameTeamID `db:"id"`
	GameID   values.GameID     `db:"game_id"`
	Name     string            `db:"name"`
	Score    int               `db:"score"`
	OrderNum int16             `db:"order_num"`
}

func (e GameTeam) ToDomain() values.GameTeam {
	return values.GameTeam{
		ID:       e.ID,
		GameID:   e.GameID,
		Name:     e.Name,
		Score:    e.Score,
		OrderNum: e.OrderNum,
	}
}
