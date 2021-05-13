package leaderboard

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/bigflood/leaderboard/api"
)

type User = api.User

type LeaderBoard struct {
	NowFunc func() time.Time

	Storage Storage
}

type Storage interface {
	Count(ctx context.Context) (int, error)
	GetData(ctx context.Context, keys ...string) ([][]byte, error)
	SetData(ctx context.Context, key string, data []byte, score int) error
	GetRanks(ctx context.Context, keys ...string) ([]int, error)
	GetSortedRange(ctx context.Context, rank, count int) ([]string, error)
}

func (lb *LeaderBoard) now() time.Time {
	if lb.NowFunc != nil {
		return lb.NowFunc()
	}
	return time.Now()
}

func (lb *LeaderBoard) UserCount(ctx context.Context) (int, error) {
	return lb.Storage.Count(ctx)
}

func (lb *LeaderBoard) GetUser(ctx context.Context, userId string) (User, error) {
	users, err := lb.Storage.GetData(ctx, userId)
	if err != nil {
		return User{}, err
	}

	if len(users[0]) == 0 {
		return User{}, api.ErrorWithStatusCode(errors.New("not found"), http.StatusNotFound)
	}

	user := User{}
	if err := json.Unmarshal(users[0], &user); err != nil {
		return User{}, err
	}

	ranks, err := lb.Storage.GetRanks(ctx, userId)
	if err != nil {
		return User{}, err
	}

	user.Rank = ranks[0]
	return user, nil
}

func (lb *LeaderBoard) SetUser(ctx context.Context, userId string, score int) error {
	users, err := lb.Storage.GetData(ctx, userId)
	if err != nil {
		return err
	}

	oldData := users[0]
	oldUser := User{}
	if len(oldData) != 0 {
		if err := json.Unmarshal(oldData, &oldUser); err != nil {
			return err
		}
	}

	if oldUser.Score == score {
		return nil
	}

	newUser := User{
		Id:        userId,
		Score:     score,
		UpdatedAt: lb.now(),
	}

	newData, err := json.Marshal(newUser)
	if err != nil {
		return err
	}

	return lb.Storage.SetData(ctx, userId, newData, score)
}

func (lb *LeaderBoard) GetRanks(ctx context.Context, rank, count int) ([]User, error) {
	if rank < 1 {
		return nil, api.ErrorWithStatusCode(errors.New("invalid rank"), http.StatusBadRequest)
	}

	if count <= 0 {
		return nil, api.ErrorWithStatusCode(errors.New("invalid count"), http.StatusBadRequest)
	}

	userIds, err := lb.Storage.GetSortedRange(ctx, rank, count)
	if err != nil {
		return nil, err
	}

	userDataList, err := lb.Storage.GetData(ctx, userIds...)
	if err != nil {
		return nil, err
	}

	returnUsers := make([]User, len(userDataList))
	for i, u := range userDataList {
		if err := json.Unmarshal(u, &returnUsers[i]); err != nil {
			return nil, err
		}
		returnUsers[i].Rank = rank + i
	}

	return returnUsers, nil
}
