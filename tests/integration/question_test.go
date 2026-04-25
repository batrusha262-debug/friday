//go:build integration

package integration_test

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (s *Suite) TestGetQuestion() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	roundID := s.seedRound(ctx, packID)
	categoryID := s.seedCategory(ctx, roundID)
	questionID := s.seedQuestion(ctx, categoryID)

	testCases := []struct {
		name       string
		questionID string
		check      func(*http.Response)
	}{
		{
			name:       "success",
			questionID: questionID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:       "not found",
			questionID: uuid.New().String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:       "invalid id",
			questionID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Get(ctx, "/admin/questions/"+tc.questionID, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestCreateQuestion() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	roundID := s.seedRound(ctx, packID)
	categoryID := s.seedCategory(ctx, roundID)

	validReq := func() map[string]any {
		return map[string]any{
			"price":    100,
			"type":     "standard",
			"question": "Столица России?",
			"answer":   "Москва",
		}
	}

	testCases := []struct {
		name       string
		categoryID string
		request    map[string]any
		check      func(*http.Response)
	}{
		{
			name:       "success",
			categoryID: categoryID.String(),
			request:    validReq(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusCreated, resp.StatusCode)
			},
		},
		{
			name:       "empty question",
			categoryID: categoryID.String(),
			request:    map[string]any{"price": 100, "type": "standard", "question": "", "answer": "Москва"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:       "empty answer",
			categoryID: categoryID.String(),
			request:    map[string]any{"price": 100, "type": "standard", "question": "Q?", "answer": ""},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:       "zero price",
			categoryID: categoryID.String(),
			request:    map[string]any{"price": 0, "type": "standard", "question": "Q?", "answer": "A"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:       "invalid type",
			categoryID: categoryID.String(),
			request:    map[string]any{"price": 100, "type": "unknown", "question": "Q?", "answer": "A"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:       "category not found",
			categoryID: uuid.New().String(),
			request:    validReq(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:       "invalid category id",
			categoryID: "abc",
			request:    validReq(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Post(ctx, "/admin/categories/"+tc.categoryID+"/questions", nil, tc.request)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestListQuestions() {
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
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:       "invalid category id",
			categoryID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Get(ctx, "/admin/categories/"+tc.categoryID+"/questions", nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestUpdateQuestion() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	roundID := s.seedRound(ctx, packID)
	categoryID := s.seedCategory(ctx, roundID)
	questionID := s.seedQuestion(ctx, categoryID)

	validReq := map[string]any{
		"price":    200,
		"type":     "auction",
		"question": "Обновлённый вопрос?",
		"answer":   "Обновлённый ответ",
	}

	testCases := []struct {
		name       string
		questionID string
		request    map[string]any
		check      func(*http.Response)
	}{
		{
			name:       "success",
			questionID: questionID.String(),
			request:    validReq,
			check: func(resp *http.Response) {
				rq.Equal(http.StatusOK, resp.StatusCode)
			},
		},
		{
			name:       "not found",
			questionID: uuid.New().String(),
			request:    validReq,
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNotFound, resp.StatusCode)
			},
		},
		{
			name:       "invalid id",
			questionID: "abc",
			request:    validReq,
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
		{
			name:       "zero price",
			questionID: questionID.String(),
			request:    map[string]any{"price": 0, "type": "standard", "question": "Q?", "answer": "A"},
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Put(ctx, "/admin/questions/"+tc.questionID, nil, tc.request)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}

func (s *Suite) TestDeleteQuestion() {
	ctx := context.Background()
	rq := s.Require()

	userID := s.seedUser(ctx)
	packID := s.seedPack(ctx, userID)
	roundID := s.seedRound(ctx, packID)
	categoryID := s.seedCategory(ctx, roundID)
	questionID := s.seedQuestion(ctx, categoryID)

	testCases := []struct {
		name       string
		questionID string
		check      func(*http.Response)
	}{
		{
			name:       "success",
			questionID: questionID.String(),
			check: func(resp *http.Response) {
				rq.Equal(http.StatusNoContent, resp.StatusCode)
			},
		},
		{
			name:       "invalid id",
			questionID: "abc",
			check: func(resp *http.Response) {
				rq.Equal(http.StatusBadRequest, resp.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			resp, err := s.http.Delete(ctx, "/admin/questions/"+tc.questionID, nil)
			rq.NoError(err)
			defer resp.Body.Close()
			tc.check(resp)
		})
	}
}
