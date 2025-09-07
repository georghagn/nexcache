// hCache_test
package hCache

import (
	"testing"
	"time"
)

type CacheRecord struct {
	id        int
	name      string
	generated time.Time
}

var fileName string = "cache.json"
var jetzt = time.Now()

func TestHCache(t *testing.T) {
	cache := NewLRUCache(3, 10*time.Second, 2*time.Second)

	recordA := CacheRecord{
		id: 1, name: "AAA", generated: jetzt,
	}
	recordB := &CacheRecord{
		id: 2, name: "BBB", generated: time.Now(),
	}
	recordC := &CacheRecord{
		id: 3, name: "CCC", generated: time.Now(),
	}

	// Daten setzen
	cache.Put("A", recordA)
	cache.Put("B", recordB)
	cache.Put("C", recordC)

	// Cache speichern
	t.Log("Save to file:", fileName)
	if err := cache.SaveToFile(fileName); err != nil {
		t.Error("Fehler beim Speichern:", err)
	}

	t.Log("Read from file:", fileName)
	newCache := NewLRUCache(3, 10*time.Second, 2*time.Second)
	if err := newCache.LoadFromFile("cache.json"); err != nil {
		t.Error("Fehler beim Laden", err)
	}

	if ergA, found := newCache.Get("A"); found {
		t.Log("A:", ergA) // 1
	} else {
		t.Error("Record A not found")
	}

	if ergB, found := newCache.Get("B"); found {
		t.Log("B:", ergB) // 2
	} else {
		t.Error("Record B not found")
	}

	if ergC, found := newCache.Get("C"); found {
		t.Log("C:", ergC) // 2
	} else {
		t.Error("Record C not found")
	}

	// Test auf delete entry
	if ergB, found := newCache.Delete("B"); found {
		t.Log("B deleted", ergB)
	} else {
		t.Error("B not found")
	}

	if ergB, found := newCache.Get("B"); found {
		t.Error("B found, should be 'NOT FOUND'", ergB)
	} else {
		t.Log("B not found")
	}
}
