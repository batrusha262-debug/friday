package values

import "github.com/google/uuid"

type GameID uuid.UUID

func NewGameID() GameID {
	return GameID(uuid.New())
}

type GameTeamID uuid.UUID

func NewGameTeamID() GameTeamID {
	return GameTeamID(uuid.New())
}
