package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// SyncMap is a concurrency-safe map guarded by a RWMutex.
// It can optionally be persisted to a JSON file (see NewSyncMap).
type SyncMap[K comparable, V any] struct {
	smap     map[K]V
	mu       sync.RWMutex
	datafile string
}

// NewSyncMap creates an empty SyncMap.
//
// filename is optional: when omitted, the map lives only in memory.
// When given, the map is loaded from that JSON file on creation and
// written back to it after every Set, Update, Remove or Clear.
// Only the first value is used; any extra values are ignored.
func NewSyncMap[K comparable, V any](filename ...string) *SyncMap[K, V] {
	sm := &SyncMap[K, V]{
		smap: map[K]V{},
		mu:   sync.RWMutex{},
	}
	if len(filename) > 0 {
		sm.datafile = filename[0]
		sm.loadFromFile()
	}
	return sm
}

func (m *SyncMap[K, V]) getSourceFile() (*os.File, error) {
	file, err := os.Open(m.datafile)

	if err != nil {
		if os.IsNotExist(err) {
			file, err = os.Create(m.datafile)

			if err != nil {
				return nil, fmt.Errorf("error creating file: %w", err)
			}

			return file, nil
		}
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return file, nil
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Println("Error closing file:", err)
	}
}

func (m *SyncMap[K, V]) loadFromFile() {
	file, err := m.getSourceFile()

	if err != nil {
		log.Fatal("[Error getting]", err)
	}

	defer closeFile(file)

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	if len(data) == 0 {
		data = []byte("{}")
	}

	err = json.Unmarshal(data, &m.smap)

	if err != nil {
		log.Fatal("Error unmarshalling JSON:", err)
	}
}

func (m *SyncMap[K, V]) saveToFile() {
	data, err := json.MarshalIndent(m.smap, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	err = os.WriteFile(m.datafile, data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
}

func (m *SyncMap[K, V]) Set(key K, val V) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.exists(key) {
		return fmt.Errorf("data already exists")
	}

	m.smap[key] = val
	if m.datafile != "" {
		m.saveToFile()
	}

	return nil
}

func (m *SyncMap[K, V]) Get(key K) (V, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var zero V

	if !m.exists(key) {
		return zero, ErrNotFound
	}

	return m.smap[key], nil
}

func (m *SyncMap[K, V]) GetAll() ([]V, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var vs []V

	for _, v := range m.smap {
		vs = append(vs, v)
	}

	return vs, nil
}

func (m *SyncMap[K, V]) Remove(key K) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.exists(key) {
		return fmt.Errorf("data doesn't exist")
	}

	delete(m.smap, key)
	if m.datafile != "" {
		m.saveToFile()
	}

	return nil
}

func (m *SyncMap[K, V]) Update(key K, newval V) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.exists(key) {
		return fmt.Errorf("data doesn't exist")
	}

	m.smap[key] = newval

	if m.datafile != "" {
		m.saveToFile()
	}

	return nil
}

func (m *SyncMap[K, V]) Count() int {
	return len(m.smap)
}

func (m *SyncMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for k := range m.smap {
		delete(m.smap, k)
	}
	if m.datafile != "" {
		m.saveToFile()
	}
}

func (m *SyncMap[K, V]) exists(key K) bool {
	_, ok := m.smap[key]
	return ok
}
