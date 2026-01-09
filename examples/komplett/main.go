package main

import (
	"fmt"
	"time"

	"github.com/georghagn/nexcache/lrucache"
)

func main() {
	// Cache mit Kapazität 3, TTL 5s, Cleanup alle 2s
	cache := lrucache.New(3, 5*time.Second, 2*time.Second)
	defer cache.StopCleanup()

	// -------- Einfaches Set/Get --------
	cache.Set("foo", "bar")
	if val, ok := cache.Get("foo"); ok {
		fmt.Println("Get foo:", val) // → "bar"
	}

	// -------- TTL Ablauf testen --------
	cache.Set("temp", "value")
	fmt.Println("Set temp: value")
	time.Sleep(6 * time.Second) // länger als TTL
	if _, ok := cache.Get("temp"); !ok {
		fmt.Println("temp abgelaufen!")
	}

	// -------- GetOrLoad mit Loader --------
	val, err := cache.GetOrLoad("user:1", func() (interface{}, error) {
		fmt.Println("Loader aufgerufen für user:1")
		return "Alice", nil
	})
	fmt.Println("user:1 =", val, "err:", err)

	// Nächster Zugriff holt aus Cache, Loader wird NICHT aufgerufen
	val, _ = cache.GetOrLoad("user:1", func() (interface{}, error) {
		fmt.Println("Dieser Loader sollte nicht laufen!")
		return "Bob", nil
	})
	fmt.Println("user:1 =", val)

	// -------- GetOrLoadWithFallback --------
	val, err = cache.GetOrLoadWithFallback("user:2", func() (interface{}, error) {
		fmt.Println("Loader schlägt fehl für user:2")
		return nil, fmt.Errorf("DB down")
	}, "FallbackUser")
	fmt.Println("user:2 =", val, "err:", err)

	// -------- Persistenz: Save/Load --------
	cache.Set("session", "abc123")
	if err := cache.SaveToFile("cache.json"); err != nil {
		fmt.Println("Fehler beim Speichern:", err)
	} else {
		fmt.Println("Cache in cache.json gespeichert")
	}

	// Neuen Cache laden
	newCache := lrucache.New(3, 5*time.Second, 2*time.Second)
	defer newCache.StopCleanup()

	if err := newCache.LoadFromFile("cache.json"); err != nil {
		fmt.Println("Fehler beim Laden:", err)
	} else if val, ok := newCache.Get("session"); ok {
		fmt.Println("Geladener Wert session:", val)
	}
}
