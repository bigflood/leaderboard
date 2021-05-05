package leaderboard

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/bigflood/leaderboard/api"
)

type User = api.User

type LeaderBoard struct {
	NowFunc func() time.Time

	mutex     sync.Mutex
	users     []*User
	userIdMap map[string]*User
}

func (lb *LeaderBoard) now() time.Time {
	if lb.NowFunc != nil {
		return lb.NowFunc()
	}
	return time.Now()
}

func (lb *LeaderBoard) UserCount(ctx context.Context) (int, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	return len(lb.users), nil
}

func (lb *LeaderBoard) GetUser(ctx context.Context, userId string) (User, error) {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	user := lb.userIdMap[userId]
	if user == nil {
		return User{}, api.ErrorWithStatusCode(errors.New("not found"), http.StatusNotFound)
	}

	return *user, nil
}

func (lb *LeaderBoard) SetUser(ctx context.Context, userId string, score int) error {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	user := lb.userIdMap[userId]
	if user == nil {
		user := &User{
			Id:        userId,
			Score:     score,
			UpdatedAt: lb.now(),
		}
		lb.users = append(lb.users, user)

		if lb.userIdMap == nil {
			lb.userIdMap = map[string]*User{}
		}
		lb.userIdMap[userId] = user
	} else {
		user.Score = score
		user.UpdatedAt = lb.now()
	}

	sort.Slice(lb.users, func(i, j int) bool {
		return lb.users[i].Score > lb.users[j].Score
	})

	for i, user := range lb.users {
		user.Rank = i + 1
	}

	return nil
}

func (lb *LeaderBoard) GetRanks(ctx context.Context, rank, count int) ([]User, error) {
	if rank < 1 {
		return nil, api.ErrorWithStatusCode(errors.New("invalid rank"), http.StatusBadRequest)
	}

	if count <= 0 {
		return nil, api.ErrorWithStatusCode(errors.New("invalid count"), http.StatusBadRequest)
	}

	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	baseIndex := rank - 1
	if baseIndex >= len(lb.users) {
		return nil, nil
	}

	if maxCount := len(lb.users) - baseIndex; count > maxCount {
		count = maxCount
	}

	returnUsers := make([]User, count)

	for i := range returnUsers {
		returnUsers[i] = *lb.users[baseIndex+i]
	}

	return returnUsers, nil
}
