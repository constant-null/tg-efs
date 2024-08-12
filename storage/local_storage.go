package storage

import (
	"context"

	"github.com/pkg/errors"
)

type LocalStorage struct {
	data map[int64][]byte
}

func NewLocal() *LocalStorage {
	return &LocalStorage{data: make(map[int64][]byte)}
}

func (s *LocalStorage) Store(ctx context.Context, id int64, data []byte) error {
	s.data[id] = data

	return nil
}

func (s *LocalStorage) Get(ctx context.Context, id int64) ([]byte, error) {
	data, ok := s.data[id]
	if !ok {
		return nil, errors.Errorf("unable to find record with id: %d", id)
	}

	return data, nil
}
