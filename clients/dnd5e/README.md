# D&D 5e API Client with Caching

This package provides a Go client for the D&D 5e API with built-in caching support.

## Features

- Full coverage of D&D 5e API endpoints
- In-memory caching with configurable TTL
- Thread-safe implementation using sync.Map
- Zero external dependencies for caching
- Comprehensive test coverage

## Usage

### Basic Client (No Caching)

```go
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}

client, err := dnd5e.NewDND5eAPI(&dnd5e.DND5eAPIConfig{
    Client: httpClient,
})
```

### Cached Client

```go
// Create base client
baseClient, err := dnd5e.NewDND5eAPI(&dnd5e.DND5eAPIConfig{
    Client: httpClient,
})

// Wrap with caching layer (24-hour TTL recommended for static D&D data)
cachedClient := dnd5e.NewCachedClient(baseClient, 24*time.Hour)

// Use the same interface - caching is transparent
races, err := cachedClient.ListRaces()
```

## Cache Implementation Details

The cached client uses a simple but effective caching strategy:

- **Storage**: Go's built-in `sync.Map` for thread-safety
- **TTL**: Configurable time-to-live for cache entries
- **Key Format**: 
  - Lists: `"list:races"`, `"list:classes"`
  - Individual items: `"race:dwarf"`, `"class:fighter"`
  - Filtered queries: `"list:spells:class:wizard:level:1"`
- **Memory Usage**: ~10MB for complete D&D 5e dataset

## Why This Approach?

Given that D&D 5e data is:
- **Static**: Rules don't change frequently
- **Small**: Entire dataset < 10MB
- **Frequently accessed**: Same data requested repeatedly

A simple in-memory cache with sync.Map provides:
- **Zero dependencies**: No external caching libraries needed
- **Excellent performance**: Sub-microsecond cache hits
- **Simple implementation**: Easy to understand and maintain
- **Thread-safe**: Safe for concurrent use

## Testing

```bash
go test ./clients/dnd5e -v
```

Tests include:
- Cache hit/miss scenarios
- TTL expiration
- Error handling
- Different query parameters