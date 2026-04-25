package values

import "github.com/google/uuid"

type RoundID uuid.UUID

func NewRoundID() RoundID {
	return RoundID(uuid.New())
}
