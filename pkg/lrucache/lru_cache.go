// lru_cache
package lrucache

import (
	"container/list"
	"encoding/json"
	"os"
	"sync"
	"time"
)

// CacheEntry speichert Schlüssel, Wert und Ablaufzeit
type CacheEntry struct {
	Key       string
	Value     interface{}
	ExpiresAt time.Time
}

// LRUCache ist die Hauptstruktur
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	mu       sync.Mutex
	ttl      time.Duration
	stopCh   chan struct{}
}

// NewLRUCache erstellt einen neuen Cache
func NewLRUCache(capacity int, ttl time.Duration, cleanupInterval time.Duration) *LRUCache {
	cache := &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
		ttl:      ttl,
		stopCh:   make(chan struct{}),
	}
	go cache.startCleanup(cleanupInterval)
	return cache
}

// ---------------------- Basisoperationen ----------------------

// Get holt einen Wert oder false, wenn nicht gefunden oder abgelaufen
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.cache[key]; found {
		entry := element.Value.(*CacheEntry)
		if time.Now().After(entry.ExpiresAt) {
			c.removeElement(element)
			return nil, false
		}
		c.list.MoveToFront(element)
		return entry.Value, true
	}
	return nil, false
}

// Set speichert einen Wert im Cache
func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.cache[key]; found {
		entry := element.Value.(*CacheEntry)
		entry.Value = value
		entry.ExpiresAt = time.Now().Add(c.ttl)
		c.list.MoveToFront(element)
		return
	}

	if c.list.Len() >= c.capacity {
		c.evictOldest()
	}

	entry := &CacheEntry{Key: key, Value: value, ExpiresAt: time.Now().Add(c.ttl)}
	element := c.list.PushFront(entry)
	c.cache[key] = element
}

// ---------------------- Erweiterungen ----------------------

// GetOrLoad: Holt Wert aus Cache oder ruft Loader auf.
// Nur erfolgreiche Loader-Ergebnisse werden gespeichert.
func (c *LRUCache) GetOrLoad(key string, loader func() (interface{}, error)) (interface{}, error) {
	c.mu.Lock()
	if element, found := c.cache[key]; found {
		entry := element.Value.(*CacheEntry)
		if time.Now().After(entry.ExpiresAt) {
			c.removeElement(element)
		} else {
			c.list.MoveToFront(element)
			val := entry.Value
			c.mu.Unlock()
			return val, nil
		}
	}
	c.mu.Unlock()

	val, err := loader()
	if err != nil {
		return nil, err
	}

	c.Set(key, val)
	return val, nil
}

// GetOrLoadWithFallback: wie GetOrLoad, aber liefert Fallback bei Fehler
func (c *LRUCache) GetOrLoadWithFallback(
	key string,
	loader func() (interface{}, error),
	fallback interface{},
) (interface{}, error) {

	c.mu.Lock()
	if element, found := c.cache[key]; found {
		entry := element.Value.(*CacheEntry)
		if time.Now().After(entry.ExpiresAt) {
			c.removeElement(element)
		} else {
			c.list.MoveToFront(element)
			val := entry.Value
			c.mu.Unlock()
			return val, nil
		}
	}
	c.mu.Unlock()

	val, err := loader()
	if err != nil {
		return fallback, err
	}

	c.Set(key, val)
	return val, nil
}

// ---------------------- Persistenz ----------------------

// SaveToFile speichert den Cache als JSON
func (c *LRUCache) SaveToFile(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var entries []CacheEntry
	for element := c.list.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*CacheEntry)
		entries = append(entries, *entry)
	}

	return json.NewEncoder(file).Encode(entries)
}

// LoadFromFile lädt Cache-Inhalt aus JSON-Datei
func (c *LRUCache) LoadFromFile(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var entries []CacheEntry
	if err := json.NewDecoder(file).Decode(&entries); err != nil {
		return err
	}

	c.cache = make(map[string]*list.Element)
	c.list = list.New()

	for _, entry := range entries {
		if time.Now().Before(entry.ExpiresAt) {
			element := c.list.PushFront(&entry)
			c.cache[entry.Key] = element
		}
	}
	return nil
}

// ---------------------- Hintergrund-Cleanup ----------------------

func (c *LRUCache) startCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.cleanupExpiredEntries()
		case <-c.stopCh:
			return
		}
	}
}

func (c *LRUCache) cleanupExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for element := c.list.Back(); element != nil; {
		entry := element.Value.(*CacheEntry)
		prev := element.Prev()
		if time.Now().After(entry.ExpiresAt) {
			c.removeElement(element)
		}
		element = prev
	}
}

// StopCleanup beendet die Cleanup-Routine
func (c *LRUCache) StopCleanup() {
	close(c.stopCh)
}

// ---------------------- Hilfsfunktionen ----------------------

func (c *LRUCache) removeElement(element *list.Element) {
	entry := element.Value.(*CacheEntry)
	delete(c.cache, entry.Key)
	c.list.Remove(element)
}

func (c *LRUCache) evictOldest() {
	oldest := c.list.Back()
	if oldest != nil {
		c.removeElement(oldest)
	}
}
