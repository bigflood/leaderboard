package tests

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/bigflood/leaderboard/api"
	"github.com/bigflood/leaderboard/pkg/leaderboard"
	. "github.com/onsi/gomega"
)

func TestLeaderBoard(t *testing.T) {
	now := time.Now()

	lb := &leaderboard.LeaderBoard{
		NowFunc: func() time.Time {
			return now
		},
	}

	testLeaderBoard(t, lb, now)
}

func testLeaderBoard(t *testing.T, lb api.LeaderBoard, now time.Time) {
	ctx := context.Background()

	g := NewWithT(t)

	count, err := lb.UserCount(ctx)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(0))

	err = lb.SetUser(ctx, "a", 100)
	g.Expect(err).NotTo(HaveOccurred())

	count, err = lb.UserCount(ctx)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(count).To(Equal(1))

	user, err := lb.GetUser(ctx, "a")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(user).To(Equal(leaderboard.User{
		Id:        "a",
		Score:     100,
		Rank:      1,
		UpdatedAt: now,
	}))

	err = lb.SetUser(ctx, "a", 10)
	g.Expect(err).NotTo(HaveOccurred())

	err = lb.SetUser(ctx, "b", 20)
	g.Expect(err).NotTo(HaveOccurred())

	err = lb.SetUser(ctx, "c", 30)
	g.Expect(err).NotTo(HaveOccurred())

	user, err = lb.GetUser(ctx, "a")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(user).To(Equal(leaderboard.User{
		Id:        "a",
		Score:     10,
		Rank:      3,
		UpdatedAt: now,
	}))

	users, err := lb.GetRanks(ctx, 1, 1000)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(users).To(Equal([]leaderboard.User{
		{
			Id:        "c",
			Score:     30,
			Rank:      1,
			UpdatedAt: now,
		},
		{
			Id:        "b",
			Score:     20,
			Rank:      2,
			UpdatedAt: now,
		},
		{
			Id:        "a",
			Score:     10,
			Rank:      3,
			UpdatedAt: now,
		},
	}))

}

func TestMultiGoroutines(t *testing.T) {
	lb := &leaderboard.LeaderBoard{}
	testMultiGoroutines(t, lb)
}

func testMultiGoroutines(t *testing.T, lb api.LeaderBoard) {
	g := NewWithT(t)

	ctx := context.Background()

	wg := sync.WaitGroup{}

	const n = 100
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			time.Sleep(time.Duration(rand.Intn(10)))

			userId := fmt.Sprint(i)
			err := lb.SetUser(ctx, userId, (n - i))
			g.Expect(err).NotTo(HaveOccurred())
		}(i)
	}

	wg.Wait()

	users, err := lb.GetRanks(ctx, 1, n)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(users)).To(Equal(n))

	for i, user := range users {
		userId := fmt.Sprint(i)
		g.Expect(user.Id).To(Equal(userId))
		g.Expect(user.Score).To(Equal(n - i))
		g.Expect(user.Rank).To(Equal(i + 1))
	}
}
