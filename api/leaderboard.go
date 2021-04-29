package api

import (
	"context"
	"time"
)

type LeaderBoard interface {
	UserCount(ctx context.Context) (int, error)
	GetUser(ctx context.Context, userId string) (User, error)
	SetUser(ctx context.Context, userId string, score int) error
	GetRanks(ctx context.Context, rank, count int) ([]User, error)
}

type User struct {
	Id    string `json:"id"`
	Score int    `json:"score"`
	// 1부터 시작하는 순위
	Rank      int       `json:"rank"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ErrorWithStatusCode(err error, statusCode int) error {
	return Error{
		origin:     err,
		statusCode: statusCode,
	}
}

type Error struct {
	origin     error
	statusCode int
}

func (e Error) Error() string {
	return e.origin.Error()
}

func (e Error) Unwrap() error {
	return e.origin
}

func (e Error) StatusCode() int {
	return e.statusCode
}
