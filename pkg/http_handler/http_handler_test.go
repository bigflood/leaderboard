package http_handler_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bigflood/leaderboard/api"
	"github.com/bigflood/leaderboard/api/apifakes"
	"github.com/bigflood/leaderboard/pkg/http_handler"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/gomega"
)

func TestHttpHandler(t *testing.T) {
	g := NewWithT(t)

	type TestData struct {
		description        string
		httpMethod         string
		path               string
		setup, after       func(*apifakes.FakeLeaderBoard)
		expectedStatusCode int
		data               interface{}
		expectedData       interface{}
	}

	type UserCountData struct {
		Count int
	}

	type MessageData struct {
		Message string
	}

	now := time.Now().UTC()

	testDataList := []TestData{
		{
			description: "usercount request",
			httpMethod:  http.MethodGet,
			path:        "/usercount",
			setup: func(fake *apifakes.FakeLeaderBoard) {
				fake.UserCountReturns(123, nil)
			},
			expectedStatusCode: http.StatusOK,
			data:               &UserCountData{},
			expectedData:       &UserCountData{Count: 123},
		},
		{
			description: "usercount error",
			httpMethod:  http.MethodGet,
			path:        "/usercount",
			setup: func(fake *apifakes.FakeLeaderBoard) {
				fake.UserCountReturns(0, errors.New("test error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			data:               &MessageData{},
			expectedData:       &MessageData{"test error"},
		},
		{
			description: "usercount: http error",
			httpMethod:  http.MethodGet,
			path:        "/usercount",
			setup: func(fake *apifakes.FakeLeaderBoard) {
				fake.UserCountReturns(
					0,
					api.ErrorWithStatusCode(errors.New("test error"), http.StatusForbidden))
			},
			expectedStatusCode: http.StatusForbidden,
			data:               &MessageData{},
			expectedData:       &MessageData{"test error"},
		},
		{
			description: "get users",
			httpMethod:  http.MethodGet,
			path:        "/users/abc",
			setup: func(fake *apifakes.FakeLeaderBoard) {
				fake.GetUserReturns(api.User{Id: "abc", Score: 100, Rank: 5, UpdatedAt: now}, nil)
			},
			expectedStatusCode: http.StatusOK,
			data:               &api.User{},
			expectedData:       &api.User{Id: "abc", Score: 100, Rank: 5, UpdatedAt: now},
		},
		{
			description: "get users error",
			httpMethod:  http.MethodGet,
			path:        "/users/abc",
			setup: func(fake *apifakes.FakeLeaderBoard) {
				fake.GetUserReturns(api.User{}, errors.New("test error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			description: "get users: not found",
			httpMethod:  http.MethodGet,
			path:        "/users/unknown",
			setup: func(fake *apifakes.FakeLeaderBoard) {
				fake.GetUserReturns(
					api.User{},
					api.ErrorWithStatusCode(errors.New("unknown not found"), http.StatusNotFound))
			},
			expectedStatusCode: http.StatusNotFound,
			data:               &MessageData{},
			expectedData:       &MessageData{"unknown not found"},
		},
		{
			description: "set users",
			httpMethod:  http.MethodPut,
			path:        "/users/abc?score=300",
			after: func(fake *apifakes.FakeLeaderBoard) {
				_, userId, score := fake.SetUserArgsForCall(0)
				g.Expect(userId).To(Equal("abc"))
				g.Expect(score).To(Equal(300))
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			description:        "set users: empty user id",
			httpMethod:         http.MethodPut,
			path:               "/users/",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			description:        "set users: empty score",
			httpMethod:         http.MethodPut,
			path:               "/users/abc",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "set users: not found",
			httpMethod:  http.MethodPut,
			path:        "/users/xxx?score=123",
			setup: func(fake *apifakes.FakeLeaderBoard) {
				fake.SetUserReturns(
					api.ErrorWithStatusCode(errors.New("xxx not found"), http.StatusNotFound))
			},
			expectedStatusCode: http.StatusNotFound,
			data:               &MessageData{},
			expectedData:       &MessageData{"xxx not found"},
		},
		{
			description: "get ranks",
			httpMethod:  http.MethodGet,
			path:        "/ranks?rank=10&count=20",
			setup: func(fake *apifakes.FakeLeaderBoard) {
				fake.GetRanksReturns(
					[]api.User{
						{Id: "a10", Score: 109, Rank: 10},
						{Id: "a11", Score: 108, Rank: 11},
						{Id: "a12", Score: 107, Rank: 12},
					},
					nil)
			},
			after: func(fake *apifakes.FakeLeaderBoard) {
				_, rank, count := fake.GetRanksArgsForCall(0)
				g.Expect(rank).To(Equal(10))
				g.Expect(count).To(Equal(20))
			},
			expectedStatusCode: http.StatusOK,
			data:               &[]api.User{},
			expectedData: &[]api.User{
				{Id: "a10", Score: 109, Rank: 10},
				{Id: "a11", Score: 108, Rank: 11},
				{Id: "a12", Score: 107, Rank: 12},
			},
		},
		{
			description:        "get ranks: empty rank",
			httpMethod:         http.MethodGet,
			path:               "/ranks?rank=&count=20",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:        "get ranks: empty count",
			httpMethod:         http.MethodGet,
			path:               "/ranks?rank=11&count=",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "get ranks: forbidden",
			httpMethod:  http.MethodGet,
			path:        "/ranks?rank=1&count=20",
			setup: func(fake *apifakes.FakeLeaderBoard) {
				fake.GetRanksReturns(
					nil,
					api.ErrorWithStatusCode(errors.New("xyz message"), http.StatusForbidden))
			},
			expectedStatusCode: http.StatusForbidden,
			data:               &MessageData{},
			expectedData:       &MessageData{"xyz message"},
		},
	}

	for _, testData := range testDataList {
		fake := &apifakes.FakeLeaderBoard{}

		if testData.setup != nil {
			testData.setup(fake)
		}

		e := echo.New()
		http_handler.New(fake).Setup(e)
		rw := httptest.NewRecorder()

		urlPrefix := "http://leaderboard.xx"
		req, err := http.NewRequest(testData.httpMethod, urlPrefix+testData.path, nil)
		g.Expect(err).NotTo(HaveOccurred())

		e.ServeHTTP(rw, req)

		resp := rw.Result()
		body, err := io.ReadAll(resp.Body)
		g.Expect(err).NotTo(HaveOccurred())

		g.Expect(resp.StatusCode).To(Equal(testData.expectedStatusCode),
			"%s: %s, body=%s", testData.description, resp.Status, string(body))

		if testData.data != nil {
			err = json.Unmarshal(body, testData.data)
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(testData.data).To(Equal(testData.expectedData),
				"%s: body=%s", testData.description, string(body))
		}

		if testData.after != nil {
			testData.after(fake)
		}
	}
}
