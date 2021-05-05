package leaderboard_test

import (
	"context"
	"github.com/benbjohnson/clock"
	. "github.com/bigflood/leaderboard/pkg/leaderboard"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

// SetUser함수가 score를 변경하는 경우 updatedAt을 변경 시간으로 갱신해야함
func TestLeaderBoard_UpdatedAt(t *testing.T) {
	g := NewWithT(t)

	ctx := context.Background()
	timeMock := clock.NewMock()

	lb := LeaderBoard{
		NowFunc: timeMock.Now,
	}

	t1 := time.Now().UTC().Truncate(time.Hour)
	timeMock.Set(t1)

	{
		const score = 100
		err := lb.SetUser(ctx, "user1",  score)
		g.Expect(err).NotTo(HaveOccurred())

		user, err := lb.GetUser(ctx, "user1")
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(user.UpdatedAt).To(Equal(t1))
		g.Expect(user.Score).To(Equal(score))
	}

	t2 := t1.Add(time.Second)
	timeMock.Set(t2)

	{
		const score = 200
		err := lb.SetUser(ctx, "user1",  score)
		g.Expect(err).NotTo(HaveOccurred())

		user, err := lb.GetUser(ctx, "user1")
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(user.UpdatedAt).To(Equal(t2))
		g.Expect(user.Score).To(Equal(score))
	}
}
