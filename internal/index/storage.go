package index

import (
	"encoding/gob"
	"fmt"
	"io"
	"sync"
)

var ErrAlreadyExists = fmt.Errorf("index already exists")
var ErrNotFound = fmt.Errorf("index not found")

type Storage struct {
	mtx   sync.RWMutex
	items map[string]Index
}

func NewStorage() *Storage {
	return &Storage{
		items: make(map[string]Index),
	}
}

// Create add index to storage
func (s *Storage) Create(index Index) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, ok := s.items[index.Name]; ok {
		return ErrAlreadyExists
	}

	// @todo validate

	s.items[index.Name] = index

	return nil
}

// Get get index from storage
func (s *Storage) Get(name string) (Index, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if index, ok := s.items[name]; ok {
		return index, nil
	}

	return Index{}, ErrNotFound
}

// Delete delete index from storage
func (s *Storage) Delete(name string) error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	if _, ok := s.items[name]; ok {
		delete(s.items, name)
		return nil
	}

	return ErrNotFound
}

type indexData struct {
	Items map[string]Index
}

// Save encodes data and writes it to the provided writer
func (s *Storage) Save(w io.Writer) error {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	data := indexData{
		Items: s.items,
	}

	return gob.NewEncoder(w).Encode(data)
}

// Load reads data from the provided reader and decodes it
func Load(reader io.Reader) (*Storage, error) {
	data := indexData{}

	err := gob.NewDecoder(reader).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &Storage{
		items: data.Items,
	}, nil
}

var DefaultStorage = NewStorage()

// Create add index to default storage
func Create(index Index) error {
	return DefaultStorage.Create(index)
}

// Delete delete index from default storage
func Delete(name string) {
	DefaultStorage.Delete(name)
}

// Get get index from default storage
func Get(name string) (Index, error) {
	return DefaultStorage.Get(name)
}
