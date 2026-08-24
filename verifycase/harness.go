package verifycase

import (
	"errors"
	"os"
	"sort"
	"strings"

	"bms/internal/store"
)

// memStore is an in-memory store.Store used to drive scenarios deterministically.
type memStore struct {
	files map[string][]byte
}

func newMemStore() *memStore {
	return &memStore{files: make(map[string][]byte)}
}

func (m *memStore) Read(name string) ([]byte, error) {
	data, ok := m.files[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return append([]byte(nil), data...), nil
}

func (m *memStore) Write(name string, data []byte) error {
	m.files[name] = append([]byte(nil), data...)
	return nil
}

func (m *memStore) Exists(name string) bool {
	_, ok := m.files[name]
	return ok
}

func (m *memStore) List(prefix string) ([]string, error) {
	names := make([]string, 0)
	for name := range m.files {
		if strings.HasPrefix(name, prefix) {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}

// failStore fails writes whose record name contains marker, simulating a
// crash between two durable steps.
type failStore struct {
	inner  store.Store
	marker string
}

func (f *failStore) Read(name string) ([]byte, error) {
	return f.inner.Read(name)
}

func (f *failStore) Write(name string, data []byte) error {
	if strings.Contains(name, f.marker) {
		return errors.New("simulated durable write failure")
	}
	return f.inner.Write(name, data)
}

func (f *failStore) Exists(name string) bool {
	return f.inner.Exists(name)
}

func (f *failStore) List(prefix string) ([]string, error) {
	return f.inner.List(prefix)
}
