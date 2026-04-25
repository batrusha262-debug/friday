package values

import "github.com/google/uuid"

type PackID = uuid.UUID

func NewPackID() PackID {
	return uuid.New()
}
