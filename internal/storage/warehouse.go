package storage

import (
	"context"
	"sync"
)

const accessPoolSize = 10

type Warehouse struct {
	mu      sync.RWMutex
	storage map[uint64]*Product

	accessPool chan struct{}
}

func NewWarehouse() *Warehouse {
	return &Warehouse{
		storage:    make(map[uint64]*Product),
		accessPool: make(chan struct{}, accessPoolSize),
	}
}

func (w *Warehouse) Lock() {
	w.accessPool <- struct{}{}
	w.mu.Lock()
}

func (w *Warehouse) Unlock() {
	w.mu.Unlock()
	<-w.accessPool
}

func (w *Warehouse) RLock() {
	w.accessPool <- struct{}{}
	w.mu.RLock()
}

func (w *Warehouse) RUnlock() {
	w.mu.RUnlock()
	<-w.accessPool
}

func (w *Warehouse) LockWithContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case w.accessPool <- struct{}{}:
		w.mu.Lock()
		return nil
	}
}

func (w *Warehouse) RLockWithContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case w.accessPool <- struct{}{}:
		w.mu.RLock()
		return nil
	}
}
