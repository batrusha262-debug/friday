package values

import "github.com/google/uuid"

type CategoryID uuid.UUID

func NewCategoryID() CategoryID {
	return CategoryID(uuid.New())
}
