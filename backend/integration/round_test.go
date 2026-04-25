//go:build integration

package integration_test

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (s *Suite) TestCreateRound() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)

	testCases := []struct {
		name    string
		packID  string
		request map[string]any
		check   func(*http.Response)
	}{
		{
			name:    "success standard",
			packID:  packID.String(),
			request: map[string]any{"name": "Первый раунд", "type": "standard"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusCreated, resp.StatusCode)
			},
		},
		{
			name:    "success double",
			packID:  packID.String(),
			request: map[string]any{"name": "Двойная", "type": "double"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusCreated, resp.StatusCode)
			},
		},
		{
			name:    "empty name",
			packID:  packID.String(),
			request: map[string]any{"name": "", "type": "standard"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "invalid type",
			packID:  packID.String(),
			request: map[string]any{"name": "Round", "type": "unknown"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "pack not found",
			packID:  uuid.New().String(),
			request: map[string]any{"name": "Round", "type": "standard"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "invalid pack id",
			packID:  "abc",
			request: map[string]any{"name": "Round", "type": "standard"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Post(ctx, "/admin/packs/"+tc.packID+"/rounds", nil, tc.request)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestListRounds() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)

	testCases := []struct {
		name   string
		packID string
		check  func(*http.Response)
	}{
		{
			name:   "success",
			packID: packID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:   "invalid pack id",
			packID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Get(ctx, "/admin/packs/"+tc.packID+"/rounds", nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestDeleteRound() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	roundID := s.seedRound(ctx, packID)

	testCases := []struct {
		name    string
		roundID string
		check   func(*http.Response)
	}{
		{
			name:    "success",
			roundID: roundID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNoContent, resp.StatusCode)
			},
		},
		{
			name:    "invalid id",
			roundID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Delete(ctx, "/admin/rounds/"+tc.roundID, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}
