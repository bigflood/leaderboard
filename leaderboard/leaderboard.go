package leaderboard

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"
)

type LeaderBoard struct {
	NowFunc func() time.Time

	mutext    sync.Mutex
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
	lb.mutext.Lock()
	defer lb.mutext.Unlock()

	return len(lb.users), nil
}

func (lb *LeaderBoard) GetUser(ctx context.Context, userId string) (User, error) {
	lb.mutext.Lock()
	defer lb.mutext.Unlock()

	user, ok := lb.userIdMap[userId]
	if !ok {
		return User{}, errors.New("user not found")
	}

	return *user, nil
}

func (lb *LeaderBoard) SetUser(ctx context.Context, userId string, score int) error {
	lb.mutext.Lock()
	defer lb.mutext.Unlock()

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
		return nil, errors.New("invalid rank")
	}

	if count <= 0 {
		return nil, errors.New("invalid count")
	}

	lb.mutext.Lock()
	defer lb.mutext.Unlock()

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

type User struct {
	Id    string `json:"id"`
	Score int    `json:"score"`
	// 1부터 시작하는 순위
	Rank      int       `json:"rank"`
	UpdatedAt time.Time `json:"updated_at"`
}
