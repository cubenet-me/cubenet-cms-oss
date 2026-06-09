package cache

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// simulateDB adds artificial latency and returns a result
func simulateDB(key string) string {
	time.Sleep(time.Duration(5+rand.Int63n(5)) * time.Millisecond)
	return fmt.Sprintf("result:%s", key)
}

// bench a service pattern: N concurrent goroutines all fetching the same key
func benchConcurrent(b *testing.B, concurrency int, useCache bool) {
	var c *Cache[string]
	if useCache {
		c = New[string](time.Second, 0)
	}
	key := "servers"

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if useCache {
				c.GetOrSet(key, func() string {
					return simulateDB(key)
				})
			} else {
				simulateDB(key)
			}
		}
	})
}

func BenchmarkNoCache_1(b *testing.B)    { benchConcurrent(b, 1, false) }
func BenchmarkNoCache_10(b *testing.B)   { benchConcurrent(b, 10, false) }
func BenchmarkNoCache_100(b *testing.B)  { benchConcurrent(b, 100, false) }
func BenchmarkNoCache_500(b *testing.B)  { benchConcurrent(b, 500, false) }

func BenchmarkCache_1(b *testing.B)      { benchConcurrent(b, 1, true) }
func BenchmarkCache_10(b *testing.B)     { benchConcurrent(b, 10, true) }
func BenchmarkCache_100(b *testing.B)    { benchConcurrent(b, 100, true) }
func BenchmarkCache_500(b *testing.B)    { benchConcurrent(b, 500, true) }

// measure multi-key access (simulating servers list + get by slug)
func benchMixed(b *testing.B, concurrency int, useCache bool) {
	var c *Cache[string]
	if useCache {
		c = New[string](time.Second, 0)
	}
	keys := []string{"list", "slug:lobby", "slug:survival", "slug:creative", "manifest", "news"}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			key := keys[rand.Intn(len(keys))]
			if useCache {
				c.GetOrSet(key, func() string {
					return simulateDB(key)
				})
			} else {
				simulateDB(key)
			}
		}
	})
}

func BenchmarkMixedNoCache_500(b *testing.B) { benchMixed(b, 500, false) }
func BenchmarkMixedCache_500(b *testing.B)   { benchMixed(b, 500, true) }

// direct cache ops benchmark (no DB)
func BenchmarkCacheDirect(b *testing.B) {
	c := New[string](time.Second, 0)
	key := "test"

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set(key, "value")
			c.Get(key)
		}
	})
}

// verify singleflight dedup under concurrent GetOrSet
func TestGetOrSetDedup(t *testing.T) {
	c := New[string](time.Minute, 0)
	var mu sync.Mutex
	calls := 0

	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.GetOrSet("key", func() string {
				mu.Lock()
				calls++
				mu.Unlock()
				time.Sleep(50 * time.Millisecond)
				return "val"
			})
		}()
	}
	wg.Wait()

	if calls != 1 {
		t.Fatalf("expected 1 call to fn, got %d", calls)
	}
	got, ok := c.Get("key")
	if !ok || got != "val" {
		t.Fatalf("expected val, got %v %v", got, ok)
	}
}

// verify expiry works
func TestExpiry(t *testing.T) {
	c := New[string](10*time.Millisecond, 4)
	c.Set("k", "v")

	v, ok := c.Get("k")
	if !ok || v != "v" {
		t.Fatal("should be available immediately")
	}

	time.Sleep(50 * time.Millisecond)
	_, ok = c.Get("k")
	if ok {
		t.Fatal("should have expired")
	}
}

// verify singleflight doesn't block Close
func TestClose(t *testing.T) {
	c := New[string](time.Minute, 4)
	c.Close()
	c.Close() // double close must not panic
}
