package entity

import (
	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
)

type Round struct {
	ID       values.RoundID     `db:"id"`
	PackID   values.PackID      `db:"pack_id"`
	Name     string             `db:"name"`
	Type     enum.RoundTypeEnum `db:"type"`
	OrderNum int16              `db:"order_num"`
}

func (e Round) ToDomain() values.Round {
	return values.Round{
		ID:       e.ID,
		PackID:   e.PackID,
		Name:     e.Name,
		Type:     e.Type,
		OrderNum: e.OrderNum,
	}
}
