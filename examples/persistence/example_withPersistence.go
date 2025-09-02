// exampleWithPersistence
package main

import (
	"fmt"
	"hCache"
	"time"
)

func main() {
	cache := hCache.NewLRUCache(3, 10*time.Second, 2*time.Second)

	// Daten setzen
	cache.Put("A", 1)
	cache.Put("B", 2)
	cache.Put("C", 3)

	// Cache speichern
	if err := cache.SaveToFile("cache.json"); err != nil {
		fmt.Println("Fehler beim Speichern:", err)
	}

	// Neuen Cache laden
	newCache := hCache.NewLRUCache(3, 10*time.Second, 2*time.Second)
	if err := newCache.LoadFromFile("cache.json"); err != nil {
		fmt.Println("Fehler beim Laden:", err)
	}

	// Werte abrufen
	item, found := cache.Get("A")
	fmt.Println("A:", item, found) // 1
	item, found = cache.Get("B")
	fmt.Println("B:", item, found) // 2
	item, found = cache.Get("C")
	fmt.Println("C:", item, found) // 3
}
