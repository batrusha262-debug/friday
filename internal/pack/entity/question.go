package entity

import (
	"friday/internal/pack/domain/enum"
	"friday/internal/pack/domain/values"
)

type Question struct {
	ID         values.QuestionID     `db:"id"`
	CategoryID values.CategoryID     `db:"category_id"`
	Price      int                   `db:"price"`
	Type       enum.QuestionTypeEnum `db:"type"`
	Question   string                `db:"question"`
	Answer     string                `db:"answer"`
	Comment    *string               `db:"comment"`
	MediaURL   *string               `db:"media_url"`
	OrderNum   int16                 `db:"order_num"`
}

func (e Question) ToDomain() values.Question {
	return values.Question{
		ID:         e.ID,
		CategoryID: e.CategoryID,
		Price:      e.Price,
		Type:       e.Type,
		Question:   e.Question,
		Answer:     e.Answer,
		Comment:    e.Comment,
		MediaURL:   e.MediaURL,
		OrderNum:   e.OrderNum,
	}
}
