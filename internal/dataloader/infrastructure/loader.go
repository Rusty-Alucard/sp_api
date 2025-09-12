// Package dataloader provides DataLoader infrastructure implementations.
package dataloader

import (
	"context"
	"fmt"
	"sync"
	"time"

	dataloader "github.com/Rusty-Alucard/sp_api/internal/dataloader/domain"
)

// loader implements the DataLoader interface.
type loader[K comparable, V any] struct {
	batchFunc dataloader.BatchFunc[K, V]
	config    dataloader.Config
	cache     dataloader.Cache[K, V]

	// Batch management
	pending   map[K][]chan result[V]
	pendingMu sync.Mutex

	// Timing
	timer   *time.Timer
	timerMu sync.Mutex
}

type result[V any] struct {
	value V
	err   error
}

// NewLoader creates a new DataLoader.
func NewLoader[K comparable, V any](
	batchFunc dataloader.BatchFunc[K, V],
	config dataloader.Config,
) dataloader.Loader[K, V] {
	if config.MaxBatch == 0 {
		config.MaxBatch = 100
	}
	if config.Wait == 0 {
		config.Wait = dataloader.DefaultWait
	}
	if config.CacheExpiry == 0 {
		config.CacheExpiry = dataloader.DefaultCacheExpiry
	}

	l := &loader[K, V]{
		batchFunc: batchFunc,
		config:    config,
		pending:   make(map[K][]chan result[V]),
	}

	if config.Cache {
		l.cache = dataloader.NewMemoryCache[K, V]()
	}

	return l
}

// Load loads a single value by key.
func (l *loader[K, V]) Load(ctx context.Context, key K) (V, error) {
	values, err := l.LoadMany(ctx, []K{key})
	if err != nil {
		var zero V
		return zero, err
	}

	if len(values) == 0 {
		var zero V
		return zero, fmt.Errorf("no value found for key: %v", key)
	}

	return values[0], nil
}

// LoadMany loads multiple values by keys.
func (l *loader[K, V]) LoadMany(ctx context.Context, keys []K) ([]V, error) {
	if len(keys) == 0 {
		return []V{}, nil
	}

	// Check cache first
	if l.cache != nil {
		cachedValues := make([]V, 0, len(keys))
		uncachedKeys := make([]K, 0, len(keys))
		uncachedIndices := make([]int, 0, len(keys))

		for i, key := range keys {
			if value, found := l.cache.Get(key); found {
				cachedValues = append(cachedValues, value)
			} else {
				uncachedKeys = append(uncachedKeys, key)
				uncachedIndices = append(uncachedIndices, i)
			}
		}

		// If all values are cached, return them
		if len(uncachedKeys) == 0 {
			return cachedValues, nil
		}

		// Load uncached values
		uncachedValues, err := l.loadUncached(ctx, uncachedKeys)
		if err != nil {
			return nil, err
		}

		// Merge cached and uncached values in correct order
		result := make([]V, len(keys))
		cachedIndex := 0
		uncachedIndex := 0

		for i, key := range keys {
			if cachedIndex < len(cachedValues) && l.cache != nil {
				if _, found := l.cache.Get(key); found {
					result[i] = cachedValues[cachedIndex]
					cachedIndex++
					continue
				}
			}

			if uncachedIndex < len(uncachedValues) {
				result[i] = uncachedValues[uncachedIndex]
				uncachedIndex++
			}
		}

		return result, nil
	}

	// No cache, load all values
	return l.loadUncached(ctx, keys)
}

// loadUncached loads values that are not in cache.
func (l *loader[K, V]) loadUncached(ctx context.Context, keys []K) ([]V, error) {
	if len(keys) == 0 {
		return []V{}, nil
	}

	// Create result channels for each key
	resultChans := make([]chan result[V], len(keys))
	for i := range resultChans {
		resultChans[i] = make(chan result[V], 1)
	}

	l.pendingMu.Lock()

	// Add to pending batch
	for i, key := range keys {
		l.pending[key] = append(l.pending[key], resultChans[i])
	}

	// Check if we should dispatch the batch
	shouldDispatch := len(l.pending) >= l.config.MaxBatch

	l.pendingMu.Unlock()

	if shouldDispatch {
		l.dispatchBatch(ctx)
	} else {
		// Set timer for next dispatch
		l.setTimer(ctx)
	}

	// Wait for results
	results := make([]V, len(keys))
	for i, resultChan := range resultChans {
		select {
		case res := <-resultChan:
			if res.err != nil {
				return nil, res.err
			}
			results[i] = res.value
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	return results, nil
}

// setTimer sets a timer to dispatch the batch after the wait period.
func (l *loader[K, V]) setTimer(ctx context.Context) {
	l.timerMu.Lock()
	defer l.timerMu.Unlock()

	if l.timer != nil {
		l.timer.Stop()
	}

	l.timer = time.AfterFunc(l.config.Wait, func() {
		l.dispatchBatch(ctx)
	})
}

// dispatchBatch dispatches the current batch of pending requests.
func (l *loader[K, V]) dispatchBatch(ctx context.Context) {
	l.pendingMu.Lock()

	if len(l.pending) == 0 {
		l.pendingMu.Unlock()
		return
	}

	// Collect all pending keys
	keys := make([]K, 0, len(l.pending))
	keyToChans := make(map[K][]chan result[V])

	for key, chans := range l.pending {
		keys = append(keys, key)
		keyToChans[key] = chans
	}

	// Clear pending
	l.pending = make(map[K][]chan result[V])

	l.pendingMu.Unlock()

	// Execute batch function
	values, err := l.batchFunc(ctx, keys)

	// Distribute results
	if err != nil {
		// Send error to all pending channels
		for _, chans := range keyToChans {
			for _, ch := range chans {
				ch <- result[V]{err: err}
			}
		}
		return
	}

	// Create a map of key to value for easy lookup
	valueMap := make(map[K]V)
	for i, key := range keys {
		if i < len(values) {
			valueMap[key] = values[i]

			// Cache the value
			if l.cache != nil {
				l.cache.Set(key, values[i], l.config.CacheExpiry)
			}
		}
	}

	// Send results to pending channels
	for key, chans := range keyToChans {
		value, exists := valueMap[key]
		for _, ch := range chans {
			if exists {
				ch <- result[V]{value: value}
			} else {
				ch <- result[V]{err: fmt.Errorf("no value found for key: %v", key)}
			}
		}
	}
}

// Clear clears a specific key from the cache.
func (l *loader[K, V]) Clear(key K) {
	if l.cache != nil {
		l.cache.Delete(key)
	}
}

// ClearAll clears all keys from the cache.
func (l *loader[K, V]) ClearAll() {
	if l.cache != nil {
		l.cache.Clear()
	}
}

// Prime primes the cache with a key-value pair.
func (l *loader[K, V]) Prime(key K, value V) {
	if l.cache != nil {
		l.cache.Set(key, value, l.config.CacheExpiry)
	}
}
