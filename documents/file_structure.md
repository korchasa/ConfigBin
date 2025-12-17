# File Structure Map

## Root
```
ConfigBin/
├── bin/run/main.go      # Entry point
├── pkg/                 # Source code
├── vendor/              # Dependencies
├── documents/           # Documentation
├── Dockerfile           # Build config
├── go.mod/sum           # Go modules
├── README.md            # Overview
└── run.ts               # Task runner
```

## Sources (`pkg/`)
- **`server/`**: HTTP logic.
  - `handlers*.go`: Endpoints.
  - `routes.go`: Routing.
  - `templates/`: HTML templates.
  - `responder/`: Response formatting.
- **`storage/`**: Persistence interfaces.
  - `sqlite/`: SQLite implementation.
- **`encryptor/`**: Encryption interfaces.
  - `aes/`: AES-GCM implementation.
- **`metrics/`**: Instrumentation.
  - `prometheus/`: Prometheus impl.
- **`types.go`**: Core domain types.

## Patterns
- **Architecture:** Interface-based dependency injection.
- **Testing:** Co-located `*_test.go` files.
- **Naming:** `snake_case` files, `CamelCase` Go types.
