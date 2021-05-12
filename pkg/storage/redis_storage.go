package storage

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisStorage struct {
	KeyPrefix string
	Client    *redis.Client
}

func (s *RedisStorage) Count(ctx context.Context) (int, error) {
	key := s.KeyPrefix + "_scores"

	count, err := s.Client.ZCard(ctx, key).Result()

	return int(count), err
}

func (s *RedisStorage) GetData(ctx context.Context, keys ...string) ([][]byte, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	var redisKeys = make([]string, len(keys))
	for i, k := range keys {
		redisKeys[i] = s.KeyPrefix + "_data_" + k
	}

	values, err := s.Client.MGet(ctx, redisKeys...).Result()
	if err != nil {
		return nil, err
	}

	returnArr := make([][]byte, len(values))
	for i, v := range values {
		if v != nil {
			returnArr[i] = []byte(fmt.Sprint(v))
		}
	}

	return returnArr, nil
}

func (s *RedisStorage) SetData(ctx context.Context, key string, data []byte, score int) error {
	dataKey := s.KeyPrefix + "_data_" + key
	scoresKey := s.KeyPrefix + "_scores"

	pipe := s.Client.TxPipeline()
	setCmd := pipe.Set(ctx, dataKey, data, 0)
	zaddCmd := pipe.ZAdd(ctx, scoresKey, &redis.Z{
		Score:  float64(score),
		Member: key,
	})

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	if _, err := setCmd.Result(); err != nil {
		return err
	}

	if _, err := zaddCmd.Result(); err != nil {
		return err
	}

	return nil
}

func (s *RedisStorage) GetRanks(ctx context.Context, keys ...string) ([]int, error) {
	if len(keys) == 0 {
		return nil, nil
	}

	scoresKey := s.KeyPrefix + "_scores"

	cmds := make([]*redis.IntCmd, len(keys))

	pipe := s.Client.TxPipeline()

	for i, k := range keys {
		cmds[i] = pipe.ZRevRank(ctx, scoresKey, k)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return nil, err
	}

	ranks := make([]int, len(cmds))
	for i, cmd := range cmds {
		rank, err := cmd.Result()
		if err != nil {
			return nil, err
		}
		ranks[i] = int(rank + 1)
	}

	return ranks, nil
}

func (s *RedisStorage) GetSortedRange(ctx context.Context, rank, count int) ([]string, error) {
	scoresKey := s.KeyPrefix + "_scores"

	return s.Client.ZRevRange(ctx, scoresKey, int64(rank-1), int64(count)).Result()
}
