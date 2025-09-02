// lruCache.go
package hCache

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

// LRUCache ist die Hauptstruktur des Caches mit TTL und Auto-Cleanup
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
	mu       sync.Mutex
	ttl      time.Duration
	stopCh   chan struct{}
}

var (
	NeverTTL   time.Duration = 800000 * time.Hour
	DefaultTTL time.Duration
)

// NewLRUCache erstellt einen neuen LRU-Cache mit TTL und startet das Auto-Cleanup
func NewLRUCache(capacity int, ttl time.Duration, cleanupInterval time.Duration) *LRUCache {
	cache := &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
		ttl:      ttl,
		stopCh:   make(chan struct{}),
	}
	go cache.startCleanup(cleanupInterval) // Cleanup-Routine starten
	return cache
}

// Get holt einen Wert aus dem Cache oder gibt nil zurück, falls nicht vorhanden oder abgelaufen
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.cache[key]; found {
		entry := element.Value.(*CacheEntry)

		// Prüfen, ob der Eintrag abgelaufen ist
		if time.Now().After(entry.ExpiresAt) {
			c.removeElement(element) // Entfernen, wenn TTL überschritten ist
			return nil, false
		}

		// Nach vorne verschieben, da er gerade benutzt wurde
		c.list.MoveToFront(element)
		return entry.Value, true
	}
	return nil, false
}

// Set speichert einen Wert im Cache und entfernt das älteste Element, falls nötig
func (c *LRUCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Falls Schlüssel existiert, aktualisieren und nach vorne verschieben
	if element, found := c.cache[key]; found {
		entry := element.Value.(*CacheEntry)
		entry.Value = value
		entry.ExpiresAt = time.Now().Add(c.ttl)
		c.list.MoveToFront(element)
		return
	}

	// Falls der Cache voll ist, das älteste Element entfernen
	if c.list.Len() >= c.capacity {
		c.evictOldest()
	}

	// Neues Element erstellen und vorne einfügen
	entry := &CacheEntry{Key: key, Value: value, ExpiresAt: time.Now().Add(c.ttl)}
	element := c.list.PushFront(entry)
	c.cache[key] = element
}

// Delete löscht einen Wert aus dem Cache oder gibt nil zurück, falls nicht vorhanden oder abgelaufen
func (c *LRUCache) Delete(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, found := c.cache[key]; found {
		entry := element.Value.(*CacheEntry)
		c.removeElement(element)
		return entry.Value, true
	}
	return nil, false
}

// removeElement entfernt ein Element aus der Liste und Map
func (c *LRUCache) removeElement(element *list.Element) {
	entry := element.Value.(*CacheEntry)
	delete(c.cache, entry.Key)
	c.list.Remove(element)
}

// evictOldest entfernt das älteste Element aus dem Cache
func (c *LRUCache) evictOldest() {
	oldest := c.list.Back()
	if oldest != nil {
		c.removeElement(oldest)
	}
}

// startCleanup entfernt periodisch abgelaufene Einträge
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

// cleanupExpiredEntries entfernt alle abgelaufenen Einträge
func (c *LRUCache) cleanupExpiredEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for element := c.list.Back(); element != nil; {
		entry := element.Value.(*CacheEntry)
		prev := element.Prev() // Speichern des vorherigen Elements, da `Remove` den aktuellen Zeiger verändert

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

// SaveToFile speichert den Cache als JSON
func (c *LRUCache) SaveToFile(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Daten in eine Liste umwandeln
	var entries []CacheEntry
	for element := c.list.Front(); element != nil; element = element.Next() {
		entry := element.Value.(*CacheEntry)
		entries = append(entries, *entry)
	}

	return json.NewEncoder(file).Encode(entries)
}

// LoadFromFile lädt den Cache aus einer JSON-Datei
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

	// Alten Cache leeren
	c.cache = make(map[string]*list.Element)
	c.list = list.New()

	// Daten wiederherstellen
	for _, entry := range entries {
		if time.Now().Before(entry.ExpiresAt) { // Nur gültige Einträge übernehmen
			element := c.list.PushFront(&entry)
			c.cache[entry.Key] = element
		}
	}
	return nil
}
