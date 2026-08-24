package store

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Store is the durable key/value file abstraction used by every component.
type Store interface {
	Read(name string) ([]byte, error)
	Write(name string, data []byte) error
	Exists(name string) bool
	List(prefix string) ([]string, error)
}

// FileStore persists records under a root directory using atomic writes.
type FileStore struct {
	root string
}

// NewFileStore returns a FileStore rooted at dir, creating it when needed.
func NewFileStore(dir string) (*FileStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create store root: %w", err)
	}
	return &FileStore{root: dir}, nil
}

// Root returns the absolute directory this store writes into.
func (s *FileStore) Root() string {
	return s.root
}

// Resolve turns a relative record name into an absolute path inside the root.
func (s *FileStore) Resolve(name string) (string, error) {
	clean := filepath.Clean(filepath.FromSlash(name))
	sep := string(filepath.Separator)
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+sep) || filepath.IsAbs(clean) {
		return "", fmt.Errorf("unsafe store name: %q", name)
	}
	return filepath.Join(s.root, clean), nil
}

// Read returns the raw bytes of one record.
func (s *FileStore) Read(name string) ([]byte, error) {
	path, err := s.Resolve(name)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Write persists one record atomically.
func (s *FileStore) Write(name string, data []byte) error {
	path, err := s.Resolve(name)
	if err != nil {
		return err
	}
	return writeAtomic(path, data)
}

// Exists reports whether a record is present.
func (s *FileStore) Exists(name string) bool {
	path, err := s.Resolve(name)
	if err != nil {
		return false
	}
	_, err = os.Stat(path)
	return err == nil
}

// List returns record names under a prefix, sorted lexically.
func (s *FileStore) List(prefix string) ([]string, error) {
	base, err := s.Resolve(prefix)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0)
	err = filepath.WalkDir(base, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			if os.IsNotExist(walkErr) {
				return nil
			}
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, relErr := filepath.Rel(s.root, path)
		if relErr != nil {
			return relErr
		}
		names = append(names, filepath.ToSlash(relative))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
