package journal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Store struct {
	mu      sync.Mutex
	dir     string
	logPath string
}

func NewStore(directory string) (*Store, error) {
	if directory == "" {
		return nil, fmt.Errorf("journal directory is required")
	}
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, fmt.Errorf("create journal directory: %w", err)
	}
	return &Store{dir: directory, logPath: filepath.Join(directory, "events.jsonl")}, nil
}

func (store *Store) Append(record Record) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	file, err := os.OpenFile(store.logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("open journal: %w", err)
	}
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(record); err != nil {
		file.Close()
		return fmt.Errorf("append journal: %w", err)
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return fmt.Errorf("sync journal: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close journal: %w", err)
	}
	return nil
}

func (store *Store) ReadAll() ([]Record, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	file, err := os.Open(store.logPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open journal: %w", err)
	}
	defer file.Close()
	var records []Record
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record Record
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			return nil, fmt.Errorf("truncated record: %w", err)
		}
		records = append(records, record)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan journal: %w", err)
	}
	return records, nil
}

func (store *Store) Path() string {
	return store.logPath
}
