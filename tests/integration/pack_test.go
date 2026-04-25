//go:build integration

package integration_test

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (s *Suite) TestCreatePack() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)

	testCases := []struct {
		name    string
		request map[string]any
		check   func(*http.Response)
	}{
		{
			name:    "success",
			request: map[string]any{"title": "Моя игра", "author_id": userID},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusCreated, resp.StatusCode)
			},
		},
		{
			name:    "empty title",
			request: map[string]any{"title": "", "author_id": userID},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "missing author_id",
			request: map[string]any{"title": "Test", "author_id": uuid.Nil},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:    "author not found",
			request: map[string]any{"title": "Test", "author_id": uuid.New()},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Post(ctx, "/admin/packs", nil, tc.request)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestListPacks() {
	ctx := context.Background()
	rq := s.Require()

	testCases := []struct {
		name  string
		setup func()
		check func(*http.Response)
	}{
		{
			name:  "empty list",
			setup: func() {},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name: "with packs",
			setup: func() {
				userID := s.seedUser(ctx)
				s.seedPack(ctx, userID)
			},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.SetupTest()
			tc.setup()
			resp, err := s.http.Get(ctx, "/admin/packs", nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestGetPack() {
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
			name:   "not found",
			packID: uuid.New().String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:   "invalid id",
			packID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Get(ctx, "/admin/packs/"+tc.packID, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestDeletePack() {
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
				rq.Equal(http.StatusNoContent, resp.StatusCode)
			},
		},
		{
			name:   "invalid id",
			packID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Delete(ctx, "/admin/packs/"+tc.packID, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}
