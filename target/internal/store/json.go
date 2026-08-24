package store

import (
	"encoding/json"
	"errors"
	"os"
)

// ErrNotFound is returned when a requested record does not exist.
var ErrNotFound = errors.New("record not found")

// IsNotFound reports whether err represents a missing record.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// ReadJSON decodes a JSON record from the store, mapping missing files to
// ErrNotFound.
func ReadJSON[T any](s Store, name string) (T, error) {
	var out T
	data, err := s.Read(name)
	if err != nil {
		if os.IsNotExist(err) {
			return out, ErrNotFound
		}
		return out, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, err
	}
	return out, nil
}

// WriteJSON encodes a record and persists it atomically.
func WriteJSON[T any](s Store, name string, value T) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return s.Write(name, data)
}
