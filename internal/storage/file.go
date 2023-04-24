package storage

import (
	"encoding/gob"
	"fmt"
	"io"
	"sync"
)

var ErrAlreadyExists = fmt.Errorf("element already exists")
var ErrNotFound = fmt.Errorf("element not found")

type File[K comparable, V any] struct {
	mtx   sync.RWMutex
	items map[K]V
}

func NewFile[K comparable, V any]() *File[K, V] {
	return &File[K, V]{
		items: make(map[K]V),
	}
}

// Create element
func (s *File[K, V]) Create(key K, value V) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.items[key]; ok {
		return ErrAlreadyExists
	}

	s.items[key] = value

	return nil
}

// Get element from storage
func (s *File[K, V]) Get(key K) (V, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if value, ok := s.items[key]; ok {
		return value, nil
	}

	var empty V

	return empty, ErrNotFound
}

// Delete element from storage
func (s *File[K, V]) Delete(key K) error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if _, ok := s.items[key]; ok {
		delete(s.items, key)
		return nil
	}

	return ErrNotFound
}

type fileData[K comparable, V any] struct {
	Items map[K]V
}

// Save encodes data and writes it to the provided writer
func (s *File[K, V]) Save(w io.Writer) error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	data := fileData[K, V]{
		Items: s.items,
	}

	return gob.NewEncoder(w).Encode(data)
}

// Load reads data from the provided reader and decodes it
func Load[K comparable, V any](reader io.Reader) (*File[K, V], error) {
	data := fileData[K, V]{}

	err := gob.NewDecoder(reader).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &File[K, V]{
		items: data.Items,
	}, nil
}
