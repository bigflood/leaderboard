package storage

import (
	"context"
	"errors"
	"sort"
	"sync"
)

type Storage struct {
	mutex        sync.Mutex
	values       map[string][]byte
	scores       map[string]Score
	sortedScores []Score
}

type Score struct {
	key   string
	score int
	rank  int
}

func (storage *Storage) Count(ctx context.Context) (int, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	return len(storage.sortedScores), nil
}

func (storage *Storage) GetData(ctx context.Context, keys ...string) ([][]byte, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	returnList := make([][]byte, len(keys))

	for i, key := range keys {
		returnList[i] = storage.values[key]
	}

	return returnList, nil
}

func (storage *Storage) SetData(ctx context.Context, key string, data []byte, score int) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	if storage.values == nil {
		storage.values = map[string][]byte{}
	}
	storage.values[key] = data

	if storage.scores == nil {
		storage.scores = map[string]Score{}
	}
	storage.scores[key] = Score{key: key, score: score}

	storage.sortedScores = storage.sortedScores[:0]

	for k, s := range storage.scores {
		storage.sortedScores = append(storage.sortedScores, Score{key: k, score: s.score})
	}

	sort.Slice(storage.sortedScores, func(i, j int) bool {
		return storage.sortedScores[i].score > storage.sortedScores[j].score
	})

	for i := range storage.sortedScores {
		storage.sortedScores[i].rank = i + 1
		s := storage.sortedScores[i]
		storage.scores[s.key] = s
	}

	return nil
}

func (storage *Storage) GetRanks(ctx context.Context, keys ...string) ([]int, error) {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	returnData := make([]int, len(keys))

	for i, key := range keys {
		s, ok := storage.scores[key]
		if ok {
			returnData[i] = s.rank
		}
	}

	return returnData, nil
}

func (storage *Storage) GetSortedRange(ctx context.Context, rank, count int) ([]string, error) {
	if rank < 1 {
		return nil, errors.New("invalid rank")
	}

	if count <= 0 {
		return nil, errors.New("invalid count")
	}

	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	baseIndex := rank - 1
	if baseIndex >= len(storage.sortedScores) {
		return nil, nil
	}

	if maxCount := len(storage.sortedScores) - baseIndex; count > maxCount {
		count = maxCount
	}

	returnData := make([]string, count)

	for i := range returnData {
		returnData[i] = storage.sortedScores[baseIndex+i].key
	}

	return returnData, nil
}
