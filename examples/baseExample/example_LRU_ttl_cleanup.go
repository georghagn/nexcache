package main

import (
	"fmt"
	"time"

	"hConsult.biz/hCache/pkg/lrucache"
)

func main() {
	cache := lrucache.NewLRUCache(3, 3*time.Second, 2*time.Second) // TTL=3s, Cleanup alle 2s

	cache.Set("A", 1)
	cache.Set("B", 2)
	cache.Set("C", 3)

	// Direkt nach dem Setzen → Werte abrufen
	fmt.Println("Direkt nach Setzen:")

	item, found := cache.Get("A")
	fmt.Println("A:", item, found) // 1

	item, found = cache.Get("B")
	fmt.Println("B:", item, found) // 2

	item, found = cache.Get("C")
	fmt.Println("C:", item, found) // 3

	// 4 Sekunden warten, damit die Werte ablaufen
	time.Sleep(4 * time.Second)

	// Nach Ablauf der TTL → Werte abrufen (Cleanup sollte alte Einträge entfernt haben)
	fmt.Println("\nNach Ablauf der TTL + Cleanup:")

	item, found = cache.Get("A")
	fmt.Println("A:", item, found) // nil, weil abgelaufen

	item, found = cache.Get("B")
	fmt.Println("B:", item, found) // nil, weil abgelaufen

	item, found = cache.Get("C")
	fmt.Println("C:", item, found) // nil, weil abgelaufen

	// Neues Element hinzufügen → Cache leert sich automatisch
	cache.Set("D", 4)
	fmt.Println("\nNach dem Hinzufügen eines neuen Elements:")

	item, found = cache.Get("D")
	fmt.Println("D:", item, found) // ❌ nil, weil abgelaufen

	// Cache Cleanup beenden, falls gewünscht
	cache.StopCleanup()
}
