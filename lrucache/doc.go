// Copyright 2026 Georg Hagn
// SPDX-License-Identifier: Apache-2.0

/*
Package lrucache bietet eine performante, thread-sichere LRU (Least Recently Used)
Cache-Implementierung mit TTL-Unterstützung und optionaler Persistenz.

Der Cache kombiniert eine Map für O(1) Zugriffe mit einer doppelt verketteten Liste,
um die Nutzungsreihenfolge zu verwalten. Abgelaufene Einträge werden entweder
beim Zugriff (Lazy) oder durch eine konfigurierbare Hintergrund-Routine (Eager) entfernt.

Ein besonderes Merkmal sind die Helper-Methoden für das "Lazy Loading" (GetOrLoad),
die den Boilerplate-Code für Cache-Miss-Szenarien drastisch reduzieren.

Beispiel für die Erstellung:

	cache := lrucache.New(1000, 15*time.Minute, 1*time.Minute)
	defer cache.StopCleanup()

Die Speicherung erfolgt über leere Interfaces (interface{}), was den Cache
flexibel für beliebige Datentypen macht.
*/
package lrucache
