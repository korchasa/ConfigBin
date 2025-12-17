# File Structure Map

## Root Directory
```
ConfigBin/
├── bin/run/main.go          # Application entry point
├── pkg/                     # Main package code
├── vendor/                  # Go dependencies (vendored)
├── documents/               # Project documentation
├── Dockerfile               # Container build definition
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── README.md                # Project overview
└── run.ts                   # Project automation script (Deno)
```

## Sources (`pkg/`)

### Core Types
```
pkg/
├── types.go                 # Bin struct definition
└── storage/
    └── encryptor.go         # Encryptor interface
```

### Server Package (`pkg/server/`)
```
pkg/server/
├── server.go                # Server struct, interfaces, main logic
├── routes.go                # Route definitions
├── middlewares.go           # Request logging middleware
│
├── handle_root.go           # GET / - create form
├── handle_bin_create.go     # POST /create - create bin
├── handle_bin_show.go       # GET /{bid} - view bin
├── handle_bin_auth.go       # POST /{bid}/auth - authenticate
├── handle_bin_update.go     # POST /{bid}/update - update bin
├── handle_api_bin_get.go    # GET /api/v1/{bid} - API endpoint
├── handle_liveness.go       # GET /liveness - health check
├── handle_readiness.go      # GET /readiness - readiness check
├── handle_not_found.go      # 404 handler
│
├── responder/
│   └── responder.go         # JSON/HTML response formatting
│
├── templates/
│   ├── templates.go         # Template compilation
│   ├── root.gohtml          # Create form template
│   ├── bin.gohtml           # Bin view/edit template
│   ├── auth.gohtml          # Password auth form
│   ├── error.gohtml         # Error page template
│   └── layouts/
│       └── main.gohtml      # Base layout
│
└── utils/
    ├── id.go                # UUID extraction from path
    └── password.go          # Cookie/password utilities
```

### Storage Package (`pkg/storage/sqlite/`)
```
pkg/storage/sqlite/
└── sqlite.go                # SQLite implementation of Storage interface
```

### Encryption Package (`pkg/encryptor/aes/`)
```
pkg/encryptor/aes/
├── aes_encryptor.go         # AES-GCM encryption implementation
└── aes_encryptor_test.go    # Encryption tests
```

### Metrics Package (`pkg/metrics/`)
```
pkg/metrics/
├── prometheus/
│   └── prometheus.go        # Prometheus metrics implementation
└── fake/
    └── fake.go              # Fake metrics for testing
```

### Entry Point (`bin/run/`)
```
bin/run/
└── main.go                  # Application initialization and startup
```

## File Organization Patterns

### Package Structure
- **Flat structure:** Most code in `pkg/` root or single-level subdirectories
- **Interface segregation:** Storage, Metrics, TemplatesProvider, Responder interfaces
- **Dependency injection:** Interfaces passed to constructors

### Naming Conventions
- **Handlers:** `handle_<action>.go` (snake_case)
- **Tests:** `*_test.go` alongside source files
- **Templates:** `*.gohtml` (Go HTML templates)
- **Packages:** lowercase, single word or short phrases

### Relationships
- `main.go` → `server.New()` → handlers
- Handlers → `Storage` interface → `sqlite.Storage`
- `sqlite.Storage` → `storage.Encryptor` → `aes.Encryptor`
- Handlers → `Responder` → JSON/HTML formatting
- Handlers → `TemplatesProvider` → HTML rendering
- Middleware → `Metrics` → Prometheus collection

### Test Files
- Tests co-located with source files
- Test files: `*_test.go`
- Test packages: same package (can test private methods)

