package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/f1monkey/search/pkg/errs"
)

type aofData[K comparable, V any] struct {
	Key       K    `json:"key"`
	Value     *V   `json:"value,omitempty"`
	IsDeleted bool `json:"isDeleted,omitempty"`
}

// AOF append-only file storage
type AOF[K comparable, V any] struct {
	mtx   sync.RWMutex
	file  *os.File
	items map[K]V
}

func NewAOF[K comparable, V any](f *os.File) *AOF[K, V] {
	return &AOF[K, V]{
		items: make(map[K]V),
		file:  f,
	}
}

func NewAOFFromPath[K comparable, V any](path string) (*AOF[K, V], error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return nil, errs.Errorf("err open file %q: %w", path, err)
	}

	// @todo close file on app stop

	return NewAOF[K, V](f), nil
}

// Create element
func (s *AOF[K, V]) Create(key K, value V) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.items[key]; ok {
		return ErrAlreadyExists
	}

	if err := s.writeData(aofData[K, V]{Key: key, Value: &value}); err != nil {
		return err
	}
	s.items[key] = value

	return nil
}

// Get element from storage
func (s *AOF[K, V]) Get(key K) (V, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if value, ok := s.items[key]; ok {
		return value, nil
	}

	var empty V

	return empty, ErrNotFound
}

// Delete element from storage
func (s *AOF[K, V]) Delete(key K) error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if _, ok := s.items[key]; ok {
		if err := s.writeData(aofData[K, V]{Key: key, IsDeleted: true}); err != nil {
			return err
		}
		delete(s.items, key)
		return nil
	}

	return ErrNotFound
}

func (s *AOF[K, V]) All() []V {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	result := make([]V, 0, len(s.items))
	for _, item := range s.items {
		result = append(result, item)
	}

	return result
}

func (s *AOF[K, V]) Init(ctx context.Context) error {
	scanner := bufio.NewScanner(s.file)

	// // @todo implement to process long lines
	// const maxCapacity int = 10000000000
	// buf := make([]byte, maxCapacity)
	// scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		v := aofData[K, V]{}
		if err := json.Unmarshal(scanner.Bytes(), &v); err != nil {
			// @todo handle corrupted data
			return errs.Errorf("data unmarshal err: %w", err)
		}
		if v.IsDeleted {
			delete(s.items, v.Key)
		} else if v.Value == nil {
			// @todo handle corrupted data
			return errs.Errorf("got nil value for key #%v", v.Key)
		} else {
			s.items[v.Key] = *v.Value
		}
	}

	if err := scanner.Err(); err != nil {
		return errs.Errorf("init err: %w", err)
	}

	return nil
}

func (s *AOF[K, V]) writeData(dat aofData[K, V]) error {
	// @todo make it async
	data, err := json.Marshal(dat)
	if err != nil {
		return errs.Errorf("element marshal err: %w", err)
	}
	data = append(data, "\n"...)
	if _, err := s.file.Write(data); err != nil {
		return errs.Errorf("element write err: %w", err)
	}

	return nil
}
