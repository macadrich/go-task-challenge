package rediscache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCacheRaceConditions(t *testing.T) {
	cache := NewRedisCache(10)

	var wg sync.WaitGroup
	numRoutines := 5 // Reduced number of routines for faster testing

	// Context with a timeout to ensure the test doesn't run too long
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Testing concurrent SET operations
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				key := "key" + fmt.Sprintf("%d", i)
				cache.Set(key, i, 1*time.Second)
			}
		}(i)
	}

	// Testing concurrent GET operations
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				key := "key" + fmt.Sprintf("%d", i)
				cache.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// Further tests to ensure keys are properly handled after concurrent operations
	cache.Set("key100", "value100", 1*time.Second)
	cache.Set("key101", "value101", 0)

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(2 * time.Second)
			if value := cache.Get("key100"); value != nil {
				t.Errorf("Expected key100 to be expired, but got %v", value)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			if value := cache.Get("key101"); value != "value101" {
				t.Errorf("Expected key101 to have value 'value101', but got %v", value)
			}
		}
	}()

	wg.Wait()
}
