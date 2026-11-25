# Goption

![goption](goption.png "Goption")

A fast, zero-allocation, type-safe generic Optional library for Go.

`goption` provides specialized generic containers tailored to Go's type system, cleanly solving the challenge of distinguishing between **absence** and **valid zero values** (such as `0`, `""`, `false`, or empty structs).

---

## Features

- **`Optional[T]` (Nullability & Presence)**: Distinguishes between absence/null and any valid value, including zero values (`0`, `""`, `false`).
- **`Zero[T]` (Zero-Value Awareness)**: Treats Go zero values (`0`, `""`, `false`, `nil`) as absent, and non-zero values as present.
- **Zero-Allocation Accessors**: 0 heap allocations and sub-nanosecond presence checks.
- **JSON Serialization & PATCH Support**: Full RFC compliance, seamless JSON PATCH handling, and Go 1.24 `omitzero` tag compatibility.
- **SQL Driver Integration**: Built-in support for `database/sql.Scanner` and `database/sql/driver.Valuer`.
- **Form & Query Parsing**: `UnmarshalText` supports relaxed boolean parsing (`"on"`, `"off"`, `"yes"`, `"no"`, `"1"`, `"0"`).
- **Functional Combinators**: Monadic chaining and transformation with `Map`, `FlatMap`, `Filter`, `IfPresent`, `IfPresentOrElse`, and `Or`.

---

## Requirements

- Go 1.23 or later (Go 1.24+ recommended for `omitzero` tag support).

---

## Installation

```bash
go get github.com/manicar2093/goption
```

---

## When to Use `Optional[T]` vs `Zero[T]`

| Concept | `Optional[T]` | `Zero[T]` |
| :--- | :--- | :--- |
| **Primary Philosophy** | Value is valid unless explicitly `null` or `nil` | Value is present only if it is non-zero and non-nil |
| **Zero Values (`0`, `""`, `false`)** | **Valid (Present)** | **Absent (Empty)** |
| **`nil` Pointers / Slices / Maps** | **Invalid (Absent/NULL)** | **Absent (Empty)** |
| **Best For** | REST APIs, JSON PATCH, SQL Databases | Config merging, CLI flags, fallback cascades |

---

## Real-World Use Cases & Examples

### Use Case 1: REST API Partial Updates (JSON PATCH) &rarr; Use `Optional[T]`

In modern REST APIs, you must distinguish between:
1. **Omitted field** &rarr; Do not update the column in the database.
2. **Explicit `null`** &rarr; Set the database column to `NULL` (clear field).
3. **Explicit zero value (`0`, `""`, `false`)** &rarr; Update column to `0`, empty string, or `false`.

```go
package main

import (
    "encoding/json"
    "fmt"
    "github.com/manicar2093/goption"
)

type UpdateUserProfileRequest struct {
    DisplayName goption.Optional[string] `json:"display_name,omitzero"`
    Age         goption.Optional[int]    `json:"age,omitzero"`
    Bio         goption.Optional[string] `json:"bio,omitzero"`
    Newsletter  goption.Optional[bool]   `json:"newsletter,omitzero"`
}

func updateUser(payloadJSON []byte) {
    var req UpdateUserProfileRequest
    json.Unmarshal(payloadJSON, &req)

    // Scenario A: JSON Payload = `{"age": 0, "bio": null, "newsletter": false}`
    // - DisplayName was omitted:
    if req.DisplayName.IsPresent() {
        // Skipped: not included in update query
    }

    // - Age is a valid 0 (e.g. newborn profile):
    if req.Age.IsPresent() {
        fmt.Println("Updating age to:", req.Age.MustGet()) // Outputs: 0
    }

    // - Bio was explicitly set to null:
    if !req.Bio.IsPresent() {
        fmt.Println("Clearing bio to SQL NULL")
    }

    // - Newsletter is explicitly false:
    if req.Newsletter.IsPresent() {
        fmt.Println("Updating newsletter subscription to:", req.Newsletter.MustGet()) // Outputs: false
    }
}
```

---

### Use Case 2: SQL Database Nullable Columns &rarr; Use `Optional[T]`

In SQL databases (PostgreSQL, MySQL, SQLite), an empty string `""` or `0` is distinct from a `NULL` column value. `Optional[T]` implements `sql.Scanner` and `driver.Valuer` without reflection overhead.

```go
package main

import (
    "database/sql"
    "github.com/manicar2093/goption"
)

type Product struct {
    ID          int64
    SKU         string
    Discount    goption.Optional[float64] // NULL = no discount; 0.0 = 0% discount
    Notes       goption.Optional[string]  // NULL = no notes; "" = empty notes string
}

func saveProduct(db *sql.DB, p *Product) error {
    _, err := db.Exec(
        "INSERT INTO products (sku, discount, notes) VALUES (?, ?, ?)",
        p.SKU,
        p.Discount, // Serializes to SQL NULL if invalid, or float value (e.g. 0.0) if valid
        p.Notes,    // Serializes to SQL NULL if invalid, or string (e.g. "") if valid
    )
    return err
}

func loadProduct(db *sql.DB, id int64) (*Product, error) {
    row := db.QueryRow("SELECT id, sku, discount, notes FROM products WHERE id = ?", id)
    var p Product
    // Scans NULL into invalid Optional[T], and scans 0.0 or "" as valid present values
    err := row.Scan(&p.ID, &p.SKU, &p.Discount, &p.Notes)
    return &p, err
}
```

---

### Use Case 3: Configuration Merging & Fallback Cascades &rarr; Use `Zero[T]`

When combining settings from **CLI flags**, **environment variables**, and **default configuration**, zero values (`0`, `""`, `false`) indicate that the user didn't specify an override. `Zero[T]` makes chaining fallbacks clean and expressive using `.Or()` and `.OrElse()`:

```go
package main

import (
    "fmt"
    "os"
    "strconv"
    "github.com/manicar2093/goption"
)

type Config struct {
    Port    int
    Timeout int
    Host    string
}

func loadConfig(cliPort int, cliHost string) Config {
    // CLI flags: if flag not passed, value is 0 or "" (treated as absent by SomeZero)
    optCLIPort := goption.SomeZero(cliPort)
    optCLIHost := goption.SomeZero(cliHost)

    // Environment variables
    envPort, _ := strconv.Atoi(os.Getenv("APP_PORT"))
    optEnvPort := goption.SomeZero(envPort)

    optEnvHost := goption.SomeZero(os.Getenv("APP_HOST"))

    // Defaults
    defaultPort := 8080
    defaultHost := "localhost"

    return Config{
        // Precedence: CLI Flag -> ENV Var -> Default
        Port: optCLIPort.Or(optEnvPort).OrElse(defaultPort),
        Host: optCLIHost.Or(optEnvHost).OrElse(defaultHost),
    }
}
```

---

### Use Case 4: Presence Toggles & Active Filters &rarr; Use `Zero[T]`

When building search query filters or analytics counters, non-zero values represent active filters:

```go
package main

import (
    "fmt"
    "github.com/manicar2093/goption"
)

type SearchFilter struct {
    MinScore goption.Zero[int]
    Category goption.Zero[string]
}

func applyFilters(f SearchFilter) {
    // Only apply filter if min score is non-zero
    f.MinScore.IfPresent(func(score int) {
        fmt.Printf("Filtering by minimum score: %d\n", score)
    })

    // Only apply category filter if string is non-empty
    f.Category.IfPresent(func(cat string) {
        fmt.Printf("Filtering by category: %s\n", cat)
    })
}
```

---

## Go 1.24 `omitzero` Tag Compatibility

Both `Optional[T]` and `Zero[T]` implement `IsZero() bool` to integrate directly with Go 1.24's `omitzero` struct tag:

```go
type Payload struct {
    OptField  goption.Optional[string] `json:"opt_field,omitzero"`
    ZeroField goption.Zero[string]     `json:"zero_field,omitzero"`
}

// 1. When OptField is Some(""):
//    - IsZero() == false -> serializes as `{"opt_field": ""}`
// 2. When OptField is uninitialized / invalid:
//    - IsZero() == true  -> omitted: `{}`
// 3. When ZeroField is SomeZero(""):
//    - IsZero() == true  -> omitted: `{}`
// 4. When ZeroField is SomeZero("hello"):
//    - IsZero() == false -> serializes as `{"zero_field": "hello"}`
```

---

## HTML Form & URL Query String Parsing

Most Go HTTP form and query decoders (such as `gorilla/schema` or `go-playground/form`) decode values via `encoding.TextUnmarshaler`. `goption` implements `UnmarshalText` with relaxed boolean parsing (`"on"`, `"off"`, `"yes"`, `"no"`, `"1"`, `"0"`), keeping standard JSON strict.

Both wrappers support defined string, boolean, and numeric types (such as `type Name string`) in text encoding. Empty text and `null` are absent; `Zero[T]` also treats decoded zero values as absent. Custom text marshaling interfaces take precedence.

```go
type SubscriptionForm struct {
    Email  goption.Optional[string] `schema:"email"`
    Agree  goption.Optional[bool]   `schema:"agree"`  // Parses "on", "yes", "1", "true" as true
}
```

---

## Functional Combinators

```go
opt := goption.Some(21)

// Map: Transform value to another type
doubled := goption.Map(opt, func(v int) string {
    return fmt.Sprintf("%d", v*2)
}) // Some("42")

// Filter: Keep value only if predicate is met
filtered := opt.Filter(func(v int) bool {
    return v > 50
}) // Invalid Optional (filtered out)

// Chaining fallback with Or
finalOpt := filtered.Or(goption.Some(100)) // Some(100)

// Conditional actions
opt.IfPresent(func(v int) {
    fmt.Printf("Value is present: %d\n", v)
})

opt.IfPresentOrElse(
    func(v int) { fmt.Printf("Present: %d\n", v) },
    func() { fmt.Println("No value present") },
)
```

---

## Performance Benchmarks

All presence checks and extractions run in constant sub-nanosecond time with zero allocations:

```
BenchmarkOptional_Get-10          1000000000          0.63 ns/op        0 B/op        0 allocs/op
BenchmarkZero_Get-10              1000000000          0.63 ns/op        0 B/op        0 allocs/op
BenchmarkPointer_Get-10           1000000000          0.31 ns/op        0 B/op        0 allocs/op
BenchmarkSQLNull_Get-10           1000000000          0.31 ns/op        0 B/op        0 allocs/op
BenchmarkOptional_Map-10           313313709          3.79 ns/op        0 B/op        0 allocs/op
```

---

## License

MIT License. See [LICENSE](LICENSE) for details.
