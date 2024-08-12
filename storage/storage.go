package storage

import "github.com/pkg/errors"

type Storage struct {
	data map[int]interface{}
}

func New() *Storage {
	return &Storage{data: make(map[int]interface{})}
}

func (s *Storage) Store(id int, data interface{}) error {
	s.data[id] = data

	return nil
}

func (s *Storage) Get(id int) (interface{}, error) {
	data, ok := s.data[id]
	if !ok {
		return nil, errors.Errorf("unable to find record with id: %s", id)
	}

	return data, nil
}
