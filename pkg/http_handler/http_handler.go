package http_handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/bigflood/leaderboard/api"
	"github.com/labstack/echo/v4"
)

type HttpHandler struct {
	lb api.LeaderBoard
}

func New(lb api.LeaderBoard) *HttpHandler {
	return &HttpHandler{lb: lb}
}

func errorJson(c echo.Context, err error) error {
	if err == nil {
		return nil
	}

	type messageData struct {
		Message string `json:"message"`
	}

	statusCode := http.StatusInternalServerError
	if s, ok := err.(interface{ StatusCode() int }); ok {
		statusCode = s.StatusCode()
	}

	return c.JSON(statusCode, messageData{err.Error()})
}

func (handler *HttpHandler) Setup(e *echo.Echo) {
	e.GET("/usercount", handler.HandleGetUserCount)
	e.GET("/users/:id", handler.HandleGetUsers)
	e.PUT("/users/:id", handler.HandlePutUsers)
	e.GET("/ranks", handler.HandleGetRanks)
}

func (handler *HttpHandler) HandleGetUserCount(c echo.Context) error {
	ctx := context.Background()
	count, err := handler.lb.UserCount(ctx)
	if err != nil {
		return errorJson(c, err)
	}

	type UserCountData struct {
		Count int `json:"count"`
	}

	return c.JSON(http.StatusOK, UserCountData{Count: count})
}

func (handler *HttpHandler) HandleGetUsers(c echo.Context) error {
	ctx := context.Background()
	userId := c.Param("id")
	user, err := handler.lb.GetUser(ctx, userId)
	if err != nil {
		return errorJson(c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (handler *HttpHandler) HandlePutUsers(c echo.Context) error {
	ctx := context.Background()
	userId := c.Param("id")
	score, err := strconv.Atoi(c.QueryParam("score"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, messageData{"score is empty or invalid format"})
	}

	if err := handler.lb.SetUser(ctx, userId, score); err != nil {
		return errorJson(c, err)
	}

	user, err := handler.lb.GetUser(ctx, userId)
	if err != nil {
		return errorJson(c, err)
	}

	return c.JSON(http.StatusOK, user)
}

func (handler *HttpHandler) HandleGetRanks(c echo.Context) error {
	ctx := context.Background()
	rank, err := strconv.Atoi(c.QueryParam("rank"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, messageData{"rank is empty or invalid format"})
	}

	count, err := strconv.Atoi(c.QueryParam("count"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, messageData{"count is empty or invalid format"})
	}

	users, err := handler.lb.GetRanks(ctx, rank, count)
	if err != nil {
		return errorJson(c, err)
	}

	return c.JSON(http.StatusOK, users)
}

type messageData struct {
	Message string `json:"message"`
}
