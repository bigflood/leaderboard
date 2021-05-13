package storage

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	. "github.com/onsi/gomega"
	"sync"
	"testing"
)

func TestRedisStorage_SetData(t *testing.T) {
	g := NewWithT(t)

	ctx := context.Background()

	s, err := miniredis.Run()
	g.Expect(err).NotTo(HaveOccurred())

	defer s.Close()

	client := redis.NewClient(&redis.Options{Addr: s.Addr()})

	hook := &redisHook{}
	client.AddHook(hook)

	storage := &RedisStorage{
		KeyPrefix: "test",
		Client:    client,
	}

	err = storage.SetData(ctx, "user1", []byte("data1"), 123)
	g.Expect(err).NotTo(HaveOccurred())

	// SetData 함수에서 multi...exec 를 사용해서 redis 명령들을 드랜잭션으로 처리했는지 확인
	lastPipeCmds := hook.processPipeCmdsList[0]
	g.Expect(lastPipeCmds[0]).To(Equal("multi"))
	g.Expect(lastPipeCmds[len(lastPipeCmds)-1]).To(Equal("exec"))

	data, err := storage.GetData(ctx, "user1")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(string(data[0])).To(Equal("data1"))

	rank, err := storage.GetRanks(ctx, "user1")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(rank[0]).To(Equal(1))
}

type redisHook struct {
	mutex               sync.Mutex
	processCmdList      []string
	processPipeCmdsList [][]string
}

func (h *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.processCmdList = append(h.processCmdList, cmd.FullName())
	return nil
}

func (h *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (h *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cmdNames := make([]string, len(cmds))
	for i, cmd := range cmds {
		cmdNames[i] = cmd.FullName()
	}
	h.processPipeCmdsList = append(h.processPipeCmdsList, cmdNames)
	return nil
}
