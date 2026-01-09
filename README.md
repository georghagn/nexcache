|<sub>🇩🇪 [German translation →](README.de.md)</sub>|
|----:|

||[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/georghagn/nexcache)](https://goreportcard.com/report/github.com/georghagn/nexcache)|
|----|----|
|![GSF-Suite-Logo](logo-gsf.png)| ***GSF-nexCache***<br>A lightweight, thread-safe **Least Recently Used (LRU)** cache fo go. Member of the **Go Small Frameworks (GSF)** family.|

<sup>***GSF*** stands for ***Go Small Frameworks*** — minimalist tools for robust applications.</sup>

# LRU Cache

A lightweight, thread-safe **Least Recently Used (LRU)** cache for Go, featuring Time-To-Live (TTL) expiration, background cleanup, and JSON persistence.

## Features

* **LRU Strategy:** Automatically evicts the least recently used items when capacity is reached.
* **TTL Support:** Entries expire automatically after a defined duration.
* **Thread-Safe:** Safe for concurrent use via `sync.Mutex`.
* **Loader Pattern:** Simplifies data fetching with `GetOrLoad` and fallback options.
* **Persistence:** Save and restore your cache state to/from JSON files.
* **Background Cleanup:** Active goroutine to prune expired entries.

## Installation

```bash
go get github.com/georghagn/nexcache"

```

## Quick Start

```go
package main

import (
	"fmt"
	"time"
    "github.com/georghagn/nexcache/lrucache"
)

func main() {
	// Initialize: capacity 100, 10m TTL, cleanup every 1m
	cache := lrucache.New(100, 10*time.Minute, 1*time.Minute)
	defer cache.StopCleanup()

	// Set a value
	cache.Set("user_1", "Georg")

	// Get a value
	if val, found := cache.Get("user_1"); found {
		fmt.Printf("Found: %v\n", val)
	}
}

```

## Extensions

### Lazy Loading (GetOrLoad)

Instead of checking for existence manually, provide a loader function. The cache handles the logic:

```go
val, err := cache.GetOrLoad("api_data", func() (interface{}, error) {
    return fetchDataFromRemoteAPI()
})

```

### Persistence

Easily persist your cache to disk to survive application restarts:

```go
// Save to file
cache.SaveToFile("backup.json")

// Load from file (only non-expired items are restored)
cache.LoadFromFile("backup.json")

```

---

## API Reference

| Method | Description |
| --- | --- |
| `New(cap, ttl, interval)` | Creates a new cache with capacity, TTL, and cleanup interval. |
| `Get(key)` | Returns the value. Updates the LRU position. |
| `Set(key, value)` | Saves a value and resets the TTL. |
| `GetOrLoad(key, loader)` | Retrieves the value or loads it if it's missing using the `loader` function. |
| `SaveToFile(path)` | Exports the cache contents as JSON. |
| `LoadFromFile(path)` | Imports cache contents (only non-expired files). |
| `StopCleanup()` | Stops the background cleanup goroutine. |

---

## Examples

The `examples/` directory contains four ready-to-run implementations:

1. **Basic:** Standard `Get` and `Set` operations.
2. **Lazy Loading:** Using `GetOrLoad` to fetch missing data.
3. **Persistence:** Demonstrating `SaveToFile` and `LoadFromFile`.
4. **TTL & Cleanup:** Showcasing how the background cleaner works.

To run an example:

```bash
go run examples/persistence/example_withPersistence.go

```
---

## How it works (LRU & TTL)

The cache combines a **hash map** for fast access with a **double linked list** to track the usage order.

* **Read access:** An item is moved to the top of the list.
* **Write access:** New items are moved to the front; when the cache capacity is exceeded, the oldest item is removed.
* **Expiration:** The background routine checks for expired timestamps at defined intervals, starting from the end of the list, to efficiently free up memory.

---

## Best Practices

### Choosing a Cleanup Interval

* **Frequent (e.g., 10s):** Use for small caches with high turnover where memory footprint is critical.
* **Balanced (e.g., 1m - 5m):** Recommended for most use cases.
* **Passive (e.g., 1h):** Sufficient if the cache is large and expired items are likely to be evicted by the LRU logic anyway.

### Type Assertions

Since the cache stores `interface{}`, always use type assertions when retrieving values:

```go
if val, found := cache.Get("myKey"); found {
    data := val.(string) // Assert to your expected type
}

```

---


### Organizational & Standards

* **Copyright:** © 2026 Georg Hagn.
* **Namespace:** `github.com/georghagn/nexcache/lrucache`
* **License:** Apache License, Version 2.0.


GSF-nexCache is an independent open-source project and is not affiliated with any corporation of a similar name.

---

## Contributing & Security

Contributions are welcome! Please use GitHub Issues for bug reports or feature ideas.
**Security-related topics** should not be discussed publicly; please refer to `SECURITY.md`.


## Contact

If you have any questions or are interested in this project, you can reach me at
📧 *georghagn [at] tiny-frameworks.io*

