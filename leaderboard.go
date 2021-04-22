package main

import "time"

type LeaderBoard struct {
}

func (lb *LeaderBoard) UserCount() int {
	return 100
}

func (lb *LeaderBoard) GetUser(userId string) User {
	return User{
		Id:        userId,
		Score:     100,
		Rank:      123,
		UpdatedAt: time.Now(),
	}
}

func (lb *LeaderBoard) SetUser(userId string, score int) {
}

func (lb *LeaderBoard) GetRanks(rank, count int) []User {
	return []User{
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
}

type User struct {
	Id        string    `json:"id"`
	Score     int       `json:"score"`
	Rank      int       `json:"rank"`
	UpdatedAt time.Time `json:"updated_at"`
}
