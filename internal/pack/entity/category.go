package entity

import "friday/internal/pack/domain/values"

type Category struct {
	ID       values.CategoryID `db:"id"`
	RoundID  values.RoundID    `db:"round_id"`
	Name     string            `db:"name"`
	OrderNum int16             `db:"order_num"`
}

func (e Category) ToDomain() values.Category {
	return values.Category{
		ID:       e.ID,
		RoundID:  e.RoundID,
		Name:     e.Name,
		OrderNum: e.OrderNum,
	}
}
