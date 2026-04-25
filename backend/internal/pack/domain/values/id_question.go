package values

import "github.com/google/uuid"

type QuestionID = uuid.UUID

func NewQuestionID() QuestionID {
	return uuid.New()
}
