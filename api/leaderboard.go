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
