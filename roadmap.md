Ideen für die **Roadmap** 

### 1. Cache Metrics (Observability)

Entwickler lieben es zu sehen, ob ihr Cache effizient arbeitet. Hinzufügen einer kleinen Struktur, die Treffer und Fehlversuche zählt.

* **Feature:** `Hits`, `Misses`, `HitRate()`.
* **Vorteil:** Nutzer können entscheiden, ob sie die Kapazität erhöhen müssen.

### 2. Sharded Cache (Performance)

Aktuell schützt ein einziger Mutex den gesamten Cache. Bei extrem vielen gleichzeitigen Zugriffen (High Concurrency) kann das zum Flaschenhals werden.

* **Feature:** Intern den Cache in z.B. 16 "Shards" (Teil-Caches) aufteilen, die jeweils ihren eigenen Mutex haben.
* **Vorteil:** Massiv höherer Durchsatz auf Mehrkern-Systemen.

### 3. Size-based Eviction

Manchmal ist nicht die Anzahl der Items das Problem, sondern der Speicherverbrauch (Bytes).

* **Feature:** Jedes Item bekommt ein "Weight" (Größe in Bytes), und der Cache begrenzt die Gesamt-Bytes statt der Anzahl der Einträge.
* **Vorteil:** Schützt besser vor `Out-of-Memory` Fehlern bei sehr unterschiedlichen Datengrößen.

### 4. Functional Options beim Initialisieren

Statt immer mehr Parameter in die `New()`-Funktion zu packen, könntest du das "Functional Options"-Pattern nutzen.

* **Konzept:** `nexcache.New(100, nexcache.WithTTL(10*time.Minute))`.
* **Vorteil:** Die API bleibt extrem sauber und erweiterbar, ohne bestehenden Code zu brechen.

