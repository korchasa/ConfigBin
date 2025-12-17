# Software Design Specification (SDS)

## 1. Introduction
- **Document purpose:** Describe system architecture, components, and implementation details
- **Relation to SRS:** Implements requirements from `requirements.md`

## 2. System Architecture

### Overview
```
┌─────────────┐
│   Client   │
│  (Browser) │
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────────────────────────┐
│         HTTP Server             │
│  ┌───────────────────────────┐  │
│  │   Router (Gorilla Mux)    │  │
│  └───────────┬───────────────┘  │
│              │                   │
│  ┌───────────▼───────────────┐ │
│  │   Handlers + Middleware    │ │
│  └───────────┬───────────────┘ │
└──────────────┼──────────────────┘
               │
       ┌───────┴────────┐
       │                │
       ▼                ▼
┌─────────────┐  ┌──────────────┐
│   Storage   │  │  Encryptor   │
│  (SQLite)   │  │   (AES-GCM)  │
└─────────────┘  └──────────────┘
```

### Main subsystems
1. **HTTP Server:** Request routing, handlers, middleware
2. **Storage:** SQLite persistence with encryption layer
3. **Encryption:** AES-256-GCM with password-derived keys
4. **Templates:** HTML rendering (Go templates)
5. **Metrics:** Prometheus instrumentation

## 3. Components

### 3.1 Server (`pkg/server`)
- **Purpose:** HTTP request handling, routing, middleware
- **Interfaces:**
  - `Storage`: bin CRUD operations
  - `Metrics`: request/metrics collection
  - `TemplatesProvider`: template access
  - `Responder`: response formatting
- **Dependencies:** gorilla/mux, logrus, prometheus
- **Key handlers:**
  - `handleRoot()`: generate UUID/password, show create form
  - `handleBinCreate()`: create bin, set cookie, redirect
  - `handleBinShow()`: display bin or auth form
  - `handleBinAuth()`: validate password, set cookie
  - `handleBinUpdate()`: create new version
  - `handleAPIGetBin()`: JSON API endpoint
  - `handleLiveness()`: health check
  - `handleReadiness()`: readiness check (DB ping)

### 3.2 Storage (`pkg/storage/sqlite`)
- **Purpose:** Persist encrypted bins with versioning
- **Interfaces:** `Storage` (CreateBin, GetBin, UpdateBin, IsReady, Close)
- **Dependencies:** mattn/go-sqlite3, `storage.Encryptor`
- **Schema:**
  ```sql
  CREATE TABLE bins (
      uuid TEXT,
      data TEXT,        -- base64 encrypted
      version INTEGER,
      created_at TIMESTAMP
  );
  CREATE INDEX bins_uuid_version_index ON bins (uuid, version);
  ```
- **Operations:**
  - `CreateBin`: insert version 0
  - `GetBin`: fetch latest + history (ORDER BY version DESC)
  - `UpdateBin`: insert new version (MAX(version)+1)

### 3.3 Encryptor (`pkg/encryptor/aes`)
- **Purpose:** AES-256-GCM encryption/decryption
- **Interfaces:** `storage.Encryptor`
- **Dependencies:** crypto/aes, crypto/cipher, crypto/sha256
- **Algorithm:**
  - Key derivation: SHA-256(password)
  - Mode: AES-GCM (authenticated encryption)
  - Nonce: random 12 bytes per encryption
  - Encoding: base64 output
- **Methods:**
  - `Encrypt(data, password)`: encrypt → base64
  - `Decrypt(data, password)`: base64 → decrypt

### 3.4 Responder (`pkg/server/responder`)
- **Purpose:** Format HTTP responses (JSON/HTML)
- **Interfaces:** `Responder` (JSONError, HTMLError, JSONSuccess)
- **Dependencies:** html/template
- **Response formats:**
  - JSON: `{success: bool, result: {...}, error: {code, message}}`
  - HTML: error template rendering

### 3.5 Templates (`pkg/server/templates`)
- **Purpose:** HTML template compilation and access
- **Dependencies:** html/template, sprig functions
- **Templates:**
  - `root.gohtml`: create form
  - `bin.gohtml`: bin view/edit
  - `auth.gohtml`: password form
  - `error.gohtml`: error page
  - `layouts/main.gohtml`: base layout

### 3.6 Metrics (`pkg/metrics/prometheus`)
- **Purpose:** Prometheus metrics collection
- **Dependencies:** prometheus/client_golang
- **Metrics:**
  - `http_producer_requests_total`: request counter
  - `http_producer_successful_requests_total{code}`: success counter
  - `http_producer_failed_requests_total{code}`: failure counter
  - `http_producer_successful_requests_duration`: latency histogram
  - (Kafka/event metrics unused in current implementation)

### 3.7 Utils (`pkg/server/utils`)
- **Purpose:** Helper functions
- **Functions:**
  - `ExtractBinIDFromPathVar()`: parse UUID from route vars
  - `PasswordCookie()`: create password cookie
  - `ReadPasswordCookie()`: read cookie value
  - `PasswordFromHeader()`: extract from Basic auth
  - `GeneratePassword()`: random password generation

## 4. Data and Storage

### Entities
- **Bin:**
  - `ID`: uuid.UUID
  - `Data`: string (decrypted content)
  - `Version`: int
  - `CreatedAt`: time.Time
  - `History`: []Bin (previous versions)

### ER diagram
```
bins
├── uuid (TEXT, indexed)
├── data (TEXT, encrypted)
├── version (INTEGER, indexed)
└── created_at (TIMESTAMP)
```

### Migration policies
- Schema auto-created on startup (`ApplySchema()`)
- No migrations; schema is idempotent (CREATE IF NOT EXISTS)

## 5. Algorithms and Logic

### Encryption flow
1. Password → SHA-256 → 32-byte key
2. Generate random 12-byte nonce
3. AES-GCM encrypt (nonce + plaintext)
4. Base64 encode result

### Decryption flow
1. Base64 decode ciphertext
2. Extract nonce (first 12 bytes)
3. Password → SHA-256 → key
4. AES-GCM decrypt

### Versioning
- Create: version = 0
- Update: version = MAX(version) + 1 (per UUID)
- Get: fetch all versions, order DESC, latest first

### Cookie management
- Name: bin UUID string
- Value: password
- MaxAge: 86400 (24h)
- HttpOnly: true
- SameSite: Lax

## 6. Non-functional Aspects

### Scalability
- Single-instance design
- SQLite suitable for <100k bins
- No horizontal scaling

### Fault tolerance
- DB connection errors handled
- Encryption errors return 400
- Template errors return 500
- Graceful shutdown (Close() on exit)

### Security
- Passwords never logged
- Encrypted at rest
- HttpOnly cookies prevent XSS
- SameSite=Lax reduces CSRF risk
- No password recovery (by design)

### Monitoring and logging
- Structured logs (logrus, fields: method, duration, URI)
- Prometheus metrics (counters, histograms)
- Health endpoints (liveness, readiness)

## 7. Constraints and Trade-offs

### Simplified
- No user accounts (password-only auth)
- No password reset/recovery
- No RBAC or permissions
- No rate limiting
- No request size limits (except timeouts)

### Deferred
- Code examples on edit page
- Version revert UI
- Syntax highlighting/validation
- Multi-format support (YAML, TOML, etc.)

## 8. Future Extensions
- Multi-format validation (JSON, YAML, TOML)
- Syntax highlighting in editor
- Version comparison/diff view
- API rate limiting
- Bin expiration/TTL
- Export/import functionality
- Multi-user support with accounts

