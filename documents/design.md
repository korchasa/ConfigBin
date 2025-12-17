# Software Design Specification (SDS)

## 1. Introduction
- **Purpose:** Architecture and implementation details.
- **Source:** Implements `requirements.md`.

## 2. Architecture
- **Type:** Monolithic HTTP Service.
- **Flow:** Client -> Router -> Handler -> Encryptor/Storage.
- **Stack:** Go 1.24, SQLite, Docker (Alpine 3.23).

### Components
1. **HTTP Server:** Gorilla Mux, Handlers, Middleware.
2. **Storage:** SQLite (encrypted data).
3. **Encryption:** AES-256-GCM (SHA-256 key derivation).
4. **Templates:** Go HTML templates.
5. **Observability:** Prometheus metrics.

## 3. Implementation
### 3.1 Server (`pkg/server`)
- **Deps:** gorilla/mux, logrus, prometheus.
- **Handlers:**
  - `handleRoot`: Gen UUID/Pass, Show Form.
  - `handleBinCreate/Update`: Store encrypted.
  - `handleBinShow/Auth`: Validate, Decrypt, Show.
  - `handleAPIGetBin`: JSON output.
  - `handleHealth`: Liveness/Readiness.

### 3.2 Storage (`pkg/storage/sqlite`)
- **Deps:** mattn/go-sqlite3.
- **Schema:**
  ```sql
  CREATE TABLE bins (uuid TEXT, data TEXT, version INTEGER, created_at TIMESTAMP);
  CREATE INDEX bins_idx ON bins (uuid, version);
  ```
- **Ops:** Create (v0), Update (MAX+1), Get (DESC).

### 3.3 Encryptor (`pkg/encryptor/aes`)
- **Algo:** AES-256-GCM.
- **Key:** SHA-256(password).
- **Nonce:** Random 12 bytes.
- **Format:** Base64(Encrypted).

### 3.4 Utils
- **Responder:** JSON/HTML formatting.
- **Templates:** `root`, `bin`, `auth`, `error`.
- **Auth:** Cookie (HttpOnly, Lax, 24h).

## 4. Data
- **Entity:** Bin (UUID, Data, Version, CreatedAt).
- **Storage:** Encrypted text, indexed by UUID+Version.
- **Migration:** Auto-init (`CREATE IF NOT EXISTS`).

## 5. Logic
- **Encryption:** Pass -> Key; Nonce + Data -> Encrypt -> Base64.
- **Decryption:** Base64 -> Nonce + Cipher -> Decrypt.
- **Versioning:** Linear increment per UUID.

## 6. Non-Functional
- **Security:** Zero-knowledge (server lacks password persistence).
- **Scalability:** Vertical only (SQLite limit).
- **Reliability:** Graceful shutdown; DB ping check.

## 7. Status
- **Simplified:** No accounts, RBAC, or recovery.
- **Deferred:** Syntax highlight, Revert UI, Multi-format.
