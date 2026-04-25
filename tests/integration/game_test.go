//go:build integration

package integration_test

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (s *Suite) TestCreateGame() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)

	testCases := []struct {
		name    string
		request map[string]any
		check   func(*http.Response)
	}{
		{
			name:    "success",
			request: map[string]any{"pack_id": packID, "host_id": userID},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusCreated, resp.StatusCode)
			},
		},
		{
			name:    "missing pack_id",
			request: map[string]any{"pack_id": uuid.Nil, "host_id": userID},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "missing host_id",
			request: map[string]any{"pack_id": packID, "host_id": uuid.Nil},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "pack not found",
			request: map[string]any{"pack_id": uuid.New(), "host_id": userID},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "host not found",
			request: map[string]any{"pack_id": packID, "host_id": uuid.New()},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Post(ctx, "/admin/games", nil, tc.request)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestGetGame() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	gameID := s.seedGame(ctx, packID, userID)

	testCases := []struct {
		name   string
		gameID string
		check  func(*http.Response)
	}{
		{
			name:   "success",
			gameID: gameID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:   "not found",
			gameID: uuid.New().String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:   "invalid id",
			gameID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Get(ctx, "/admin/games/"+tc.gameID, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestDeleteGame() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	gameID := s.seedGame(ctx, packID, userID)

	testCases := []struct {
		name   string
		gameID string
		check  func(*http.Response)
	}{
		{
			name:   "success",
			gameID: gameID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNoContent, resp.StatusCode)
			},
		},
		{
			name:   "invalid id",
			gameID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Delete(ctx, "/admin/games/"+tc.gameID, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestStartGame() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	gameID := s.seedGame(ctx, packID, userID)
	alreadyStartedID := s.seedStartedGame(ctx, packID, userID)

	testCases := []struct {
		name   string
		gameID string
		check  func(*http.Response)
	}{
		{
			name:   "success",
			gameID: gameID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:   "not found",
			gameID: uuid.New().String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:   "already started",
			gameID: alreadyStartedID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:   "invalid id",
			gameID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Post(ctx, "/admin/games/"+tc.gameID+"/start", nil, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestFinishGame() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	waitingGameID := s.seedGame(ctx, packID, userID)
	activeGameID := s.seedStartedGame(ctx, packID, userID)

	testCases := []struct {
		name   string
		gameID string
		check  func(*http.Response)
	}{
		{
			name:   "success",
			gameID: activeGameID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:   "not found",
			gameID: uuid.New().String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:   "not started",
			gameID: waitingGameID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:   "invalid id",
			gameID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Post(ctx, "/admin/games/"+tc.gameID+"/finish", nil, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestAddTeam() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	gameID := s.seedGame(ctx, packID, userID)

	testCases := []struct {
		name    string
		gameID  string
		request map[string]any
		check   func(*http.Response)
	}{
		{
			name:    "success",
			gameID:  gameID.String(),
			request: map[string]any{"name": "Красные"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusCreated, resp.StatusCode)
			},
		},
		{
			name:    "empty name",
			gameID:  gameID.String(),
			request: map[string]any{"name": ""},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "game not found",
			gameID:  uuid.New().String(),
			request: map[string]any{"name": "Синие"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "invalid game id",
			gameID:  "abc",
			request: map[string]any{"name": "Синие"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Post(ctx, "/admin/games/"+tc.gameID+"/teams", nil, tc.request)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestListTeams() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	gameID := s.seedGame(ctx, packID, userID)

	testCases := []struct {
		name   string
		gameID string
		check  func(*http.Response)
	}{
		{
			name:   "success",
			gameID: gameID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:   "invalid game id",
			gameID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Get(ctx, "/admin/games/"+tc.gameID+"/teams", nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestRemoveTeam() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	gameID := s.seedGame(ctx, packID, userID)
	teamID := s.seedTeam(ctx, gameID)

	testCases := []struct {
		name   string
		teamID string
		check  func(*http.Response)
	}{
		{
			name:   "success",
			teamID: teamID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNoContent, resp.StatusCode)
			},
		},
		{
			name:   "invalid id",
			teamID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Delete(ctx, "/admin/teams/"+tc.teamID, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestGetBoard() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	gameID := s.seedGame(ctx, packID, userID)

	testCases := []struct {
		name   string
		gameID string
		check  func(*http.Response)
	}{
		{
			name:   "success",
			gameID: gameID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:   "invalid game id",
			gameID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Get(ctx, "/admin/games/"+tc.gameID+"/board", nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestAnswerQuestion() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	roundID := s.seedRound(ctx, packID)
	categoryID := s.seedCategory(ctx, roundID)
	questionID := s.seedQuestion(ctx, categoryID)
	gameID := s.seedGame(ctx, packID, userID)
	teamID := s.seedTeam(ctx, gameID)

	testCases := []struct {
		name       string
		gameID     string
		questionID string
		request    map[string]any
		check      func(*http.Response)
	}{
		{
			name:       "success with team",
			gameID:     gameID.String(),
			questionID: questionID.String(),
			request:    map[string]any{"team_id": teamID},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:       "success without team",
			gameID:     gameID.String(),
			questionID: questionID.String(),
			request:    map[string]any{"team_id": nil},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:       "game not found",
			gameID:     uuid.New().String(),
			questionID: questionID.String(),
			request:    map[string]any{"team_id": nil},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:       "invalid game id",
			gameID:     "abc",
			questionID: questionID.String(),
			request:    map[string]any{"team_id": nil},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:       "invalid question id",
			gameID:     gameID.String(),
			questionID: "abc",
			request:    map[string]any{"team_id": nil},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Post(ctx, "/admin/games/"+tc.gameID+"/questions/"+tc.questionID+"/answer", nil, tc.request)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}
