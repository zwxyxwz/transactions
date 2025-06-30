package mvcc

import "sync"

// GCManager handles garbage collection of obsolete versions.
type GCManager struct {
	storage *MemoryStorage
	mu      sync.Mutex
}

func NewGCManager(storage *MemoryStorage) *GCManager {
	return &GCManager{storage: storage}
}

func (gc *GCManager) Run() {
	gc.mu.Lock()
	defer gc.mu.Unlock()
	for key, versions := range gc.storage.data {
		newVersions := make([]*VersionedValue, 0, len(versions))
		for _, v := range versions {
			if v.Committed {
				newVersions = append(newVersions, v)
			}
		}
		gc.storage.data[key] = newVersions
	}
}
