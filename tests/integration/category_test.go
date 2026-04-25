//go:build integration

package integration_test

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (s *Suite) TestCreateCategory() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	roundID := s.seedRound(ctx, packID)

	testCases := []struct {
		name    string
		roundID string
		request map[string]any
		check   func(*http.Response)
	}{
		{
			name:    "success",
			roundID: roundID.String(),
			request: map[string]any{"name": "История"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusCreated, resp.StatusCode)
			},
		},
		{
			name:    "empty name",
			roundID: roundID.String(),
			request: map[string]any{"name": ""},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "round not found",
			roundID: uuid.New().String(),
			request: map[string]any{"name": "История"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "invalid round id",
			roundID: "abc",
			request: map[string]any{"name": "История"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Post(ctx, "/admin/rounds/"+tc.roundID+"/categories", nil, tc.request)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestListCategories() {
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
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:    "invalid round id",
			roundID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Get(ctx, "/admin/rounds/"+tc.roundID+"/categories", nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestDeleteCategory() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	roundID := s.seedRound(ctx, packID)
	categoryID := s.seedCategory(ctx, roundID)

	testCases := []struct {
		name       string
		categoryID string
		check      func(*http.Response)
	}{
		{
			name:       "success",
			categoryID: categoryID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNoContent, resp.StatusCode)
			},
		},
		{
			name:       "invalid id",
			categoryID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Delete(ctx, "/admin/categories/"+tc.categoryID, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}
