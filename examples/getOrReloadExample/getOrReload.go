package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/georghagn/nexCache/lrucache"
)

func main() {
	cache := lrucache.New(3, 5*time.Second, 2*time.Second)

	// Loader, der manchmal Fehler zurückgibt
	loader := func() (interface{}, error) {
		now := time.Now().Unix()
		if now%2 == 0 {
			fmt.Println(now)
			return nil, errors.New("DB nicht erreichbar")
		}
		return fmt.Sprintf("Wert geladen um %d", now), nil
	}

	// Mehrere Aufrufe
	for i := 0; i < 4; i++ {
		val, err := cache.GetOrLoad("user:42", loader)
		if err != nil {
			fmt.Println("Fehler:", err)
		} else {
			fmt.Println("Ergebnis:", val)
		}
		time.Sleep(2 * time.Second)
	}
}
