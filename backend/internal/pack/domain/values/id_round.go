package values

import "github.com/google/uuid"

type RoundID = uuid.UUID

func NewRoundID() RoundID {
	return uuid.New()
}
