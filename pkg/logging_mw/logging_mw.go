package logging_mw

import (
	"context"
	"log"

	"github.com/bigflood/leaderboard/api"
)

type LoggingMiddleware struct {
	Logger   *log.Logger
	Receiver api.LeaderBoard
}

var _ api.LeaderBoard = (*LoggingMiddleware)(nil)

func (mw *LoggingMiddleware) UserCount(ctx context.Context) (int, error) {
	count, err := mw.Receiver.UserCount(ctx)
	mw.Logger.Printf("LeaderBoard.UserCount() -> %v, err=%v\n", count, err)
	return count, err
}

func (mw *LoggingMiddleware) GetUser(ctx context.Context, userId string) (api.User, error) {
	user, err := mw.Receiver.GetUser(ctx, userId)
	mw.Logger.Printf("LeaderBoard.GetUser(userId=%v) -> %+v, err=%v\n", userId, user, err)
	return user, err
}

func (mw *LoggingMiddleware) SetUser(ctx context.Context, userId string, score int) error {
	err := mw.Receiver.SetUser(ctx, userId, score)
	mw.Logger.Printf("LeaderBoard.SetUser(userId=%v, score=%v) -> err=%v\n", userId, score, err)
	return err
}

func (mw *LoggingMiddleware) GetRanks(ctx context.Context, rank, count int) ([]api.User, error) {
	users, err := mw.Receiver.GetRanks(ctx, rank, count)
	mw.Logger.Printf("LeaderBoard.GetRanks(rank=%v, count=%v) -> %+v, err=%v\n", rank, count, users, err)
	return users, err
}
