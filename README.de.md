
|<sub>🇬🇧 [English translation →](README.md)</sub>|
|----:|
|    |

||[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE) [![Dependencies](https://img.shields.io/badge/dependencies-zero-brightgreen.svg)](#)|
|----|----|
|![GSF-Suite-Logo](logo-gsf.png)| ***GSF-nexCache***<br>Din leichtgewichtiger, thread-sicherer **Least Recently Used (LRU)** Cache für Go. Teil der **Go Small Frameworks Suite**|

<sup>***GSF*** steht für ***Go Small Frameworks*** — eine Sammlung von minimalistischen Tools für robuste Anwendungen.</sup>

---

# LRU Cache

Ein leichtgewichtiger, thread-sicherer **Least Recently Used (LRU)** Cache für Go mit Time-To-Live (TTL) Unterstützung, automatischer Bereinigung und optionaler Persistenz.

## Features

* **LRU-Strategie:** Verdrängt automatisch die am längsten nicht verwendeten Elemente.
* **TTL (Time-To-Live):** Elemente laufen nach einer definierten Zeitspanne ab.
* **Thread-Safe:** Sicherer Zugriff aus mehreren Goroutinen durch `sync.Mutex`.
* **Loader-Pattern:** Bequemes Laden von Daten via `GetOrLoad` (Lazy Loading).
* **Persistenz:** Einfaches Speichern und Laden des Cache-Zustands als JSON.
* **Hintergrund-Cleanup:** Automatisches Entfernen abgelaufener Einträge.

## Installation

```bash
go get github.com/georghagn/nexcache/lrucache

```

## Schnelleinstieg

```go
package main

import (
	"fmt"
	"time"
	"github.com/georghagn/nexcache/lrucache"
)

func main() {
	// Cache erstellen: Kapazität 100, TTL 10 Min, Cleanup alle 1 Min
	cache := lrucache.New(100, 10*time.Minute, 1*time.Minute)
	defer cache.StopCleanup()

	// Wert setzen
	cache.Set("user_1", "Georg")

	// Wert abrufen
	if val, found := cache.Get("user_1"); found {
		fmt.Println("Gefunden:", val)
	}
}

```

## Fortgeschrittene Funktionen

### GetOrLoad (Lazy Loading)

Vermeidet Cache-Miss-Logik im restlichen Code. Wenn der Key fehlt oder abgelaufen ist, wird der Loader ausgeführt.

```go
val, err := cache.GetOrLoad("db_query", func() (interface{}, error) {
    // Hier Logik zum Laden aus der DB
    return "Ergebnis von DB", nil
})

```

### Persistenz

Speichere den aktuellen Zustand des Caches in einer Datei, um ihn nach einem Neustart wiederherzustellen.

```go
// Speichern
err := cache.SaveToFile("cache_dump.json")

// Laden
err := cache.LoadFromFile("cache_dump.json")

```

## API Referenz

| Methode | Beschreibung |
| --- | --- |
| `New(cap, ttl, interval)` | Erstellt einen neuen Cache mit Kapazität, TTL und Cleanup-Intervall. |
| `Get(key)` | Liefert den Wert. Aktualisiert die LRU-Position. |
| `Set(key, value)` | Speichert einen Wert und setzt die TTL zurück. |
| `GetOrLoad(key, loader)` | Holt den Wert oder lädt ihn bei Fehlen über die Funktion `loader`. |
| `SaveToFile(path)` | Exportiert den Cache-Inhalt als JSON. |
| `LoadFromFile(path)` | Importiert Cache-Inhalte (nur nicht-abgelaufene). |
| `StopCleanup()` | Beendet die Hintergrund-Goroutine für den Cleanup. |


---


## Beispiele

Im Ordner `examples/` findest du vier verschiedene Implementierungen, die die Vielseitigkeit des Caches zeigen:

1. **Basic Usage:** Standard `Get` und `Set` Operationen.
2. **Lazy Loading:** Komplexere Datenabfragen mit `GetOrLoad`.
3. **Persistence:** Strategien zum Speichern und Laden des Cache-States.
4. **TTL & Cleanup:** Demonstration der automatischen Speicherbereinigung.

Um ein Beispiel zu starten, wechsle einfach in das Verzeichnis und führe es aus:

```bash
go run examples/persistence/example_withPersistence.go

```

---

## Funktionsweise (LRU & TTL)

Der Cache kombiniert eine **Hash-Map** für schnellen Zugriff () mit einer **doppelten verketteten Liste**, um die Nutzungsreihenfolge zu tracken.

* **Lese-Zugriff:** Ein Element wird an die Spitze der Liste verschoben.
* **Schreib-Zugriff:** Neue Elemente kommen nach vorne; bei Überschreiten der Kapazität wird das letzte Element (Oldest) entfernt.
* **Ablauf:** Die Hintergrund-Routine prüft im definierten Intervall vom Ende der Liste her auf abgelaufene Zeitstempel, um den Speicher effizient freizugeben.

---

## Best Practices
### Cleanup-Intervall wählen

Das cleanupInterval bestimmt, wie oft eine Hintergrund-Goroutine den Speicher scannt.
    Aggressiv (z.B. 10s): Gut für sehr kleine Caches mit extrem flüchtigen Daten. Erhöht die CPU-Last minimal.
    Ausgewogen (z.B. 1m - 5m): Der Standard für die meisten Anwendungen.
    Passiv (z.B. 1h): Reicht aus, wenn der Cache sehr groß ist und abgelaufene Daten meistens sowieso durch die LRU-Logik (Verdrängung bei Kapazitätsgrenze) entfernt werden.

### Interface-Konvertierung

Da der Cache interface{} speichert, solltest du beim Abrufen den Type-Assertion-Check nutzen:
```Go

if val, found := cache.Get("key"); found {
    data := val.(MeinTyp) // oder sicher: data, ok := val.(MeinTyp)
}

```

---

### Organisation & Standards

* **Copyright:** © 2026 Georg Hagn.
* **Namespace:** `github.com/georghagn/nexcache/lrucache`
* **Lizenz:** Apache License, Version 2.0.

GSF-nexCache ist ein unabhängiges open-source project und ist mit keinem Unternehmen ähnlichen Namens verbunden.

---

## Mitwirken & Sicherheit

Beiträge sind willkommen! Bitte nutzen Sie GitHub Issues für Fehlerberichte oder Feature-Ideen.
**Sicherheitsrelevante Themen** sollten nicht öffentlich diskutiert werden; bitte beachten Sie hierzu die `SECURITY.de.md`.

---

## Kontakt

Bei Fragen oder Interesse an diesem Projekt erreichen Sie mich unter:
📧 *georghagn [at] tiny-frameworks.io*

<sup>*(Bitte keine Anfragen an die privaten GitHub-Account-Adressen)*</sup>


