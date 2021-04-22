package leaderboard

import (
	"context"
	"time"
)

type LeaderBoard struct {
}

func (lb *LeaderBoard) UserCount(ctx context.Context) (int, error) {
	return 100, nil
}

func (lb *LeaderBoard) GetUser(ctx context.Context, userId string) (User, error) {
	user := User{
		Id:        userId,
		Score:     100,
		Rank:      123,
		UpdatedAt: time.Now(),
	}
	return user, nil
}

func (lb *LeaderBoard) SetUser(ctx context.Context, userId string, score int) error {
	return nil
}

func (lb *LeaderBoard) GetRanks(ctx context.Context, rank, count int) ([]User, error) {
	users := []User{
		{
			Id:        "a",
			Score:     300,
			Rank:      1,
			UpdatedAt: time.Now(),
		},
		{
			Id:        "b",
			Score:     200,
			Rank:      2,
			UpdatedAt: time.Now(),
		},
		{
			Id:        "a",
			Score:     100,
			Rank:      3,
			UpdatedAt: time.Now(),
		},
	}
	return users, nil
}

type User struct {
	Id        string    `json:"id"`
	Score     int       `json:"score"`
	Rank      int       `json:"rank"`
	UpdatedAt time.Time `json:"updated_at"`
}
