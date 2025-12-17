# Software Requirements Specification (SRS)

## 1. Introduction
- **Document purpose:** Define functional and non-functional requirements for ConfigBin service
- **Scope:** Remote encrypted configuration storage with web UI and API access
- **Audience:** Developers, operators, API consumers
- **Definitions and abbreviations:**
  - Bin: encrypted configuration container identified by UUID
  - AES-GCM: encryption algorithm (AES-256 in GCM mode)

## 2. General description
- **System context:** HTTP service storing encrypted configs. Users edit via web UI; apps fetch via API. Both use password auth.
- **Assumptions and constraints:**
  - Passwords not recoverable (not stored)
  - SQLite for persistence
  - Single binary deployment
- **Assumptions:**
  - Users manage passwords securely
  - Network transport secured (HTTPS recommended)

## 3. Functional requirements

### 3.1 FR-1: Bin creation
- **Description:** Create encrypted config bin with UUID and password
- **Use case scenario:** User visits `/`, gets generated UUID/password, submits form with content → bin created
- **Acceptance criteria:**
  - [x] UUID auto-generated
  - [x] Password auto-generated (8 chars, alphanumeric + specials)
  - [x] Content encrypted with AES-GCM
  - [x] Stored in SQLite
  - [x] Password cookie set (24h, HttpOnly, SameSite=Lax)
  - [x] Redirect to bin view page

### 3.2 FR-2: Bin viewing (web)
- **Description:** Display bin content via web UI after password auth
- **Use case scenario:** User visits `/{bid}` → if no cookie, show auth form → after auth, show content
- **Acceptance criteria:**
  - [x] Auth form shown if no valid cookie
  - [x] Password validation via cookie or form
  - [x] Decrypted content displayed
  - [x] Edit form available

### 3.3 FR-3: Bin authentication
- **Description:** Authenticate with password to access bin
- **Use case scenario:** User submits password at `/{bid}/auth` → cookie set → redirect to bin view
- **Acceptance criteria:**
  - [x] Password validated against stored encrypted data
  - [x] Cookie set on success
  - [x] Error on invalid password
  - [x] Redirect to bin view

### 3.4 FR-4: Bin update
- **Description:** Update bin content (creates new version)
- **Use case scenario:** User edits content, submits → new version created → redirect to view
- **Acceptance criteria:**
  - [x] Version incremented
  - [x] Previous versions preserved in history
  - [x] Content re-encrypted with same password
  - [x] Cookie-based auth required

### 3.5 FR-5: API access
- **Description:** Fetch bin content via HTTP API with Basic auth
- **Use case scenario:** App sends `GET /api/v1/{bid}` with `Authorization: Basic` header → returns JSON
- **Acceptance criteria:**
  - [x] Basic auth password extraction
  - [x] JSON response with bin data
  - [x] History included in response
  - [x] Error responses in JSON format

### 3.6 FR-6: Version history
- **Description:** Store and retrieve bin version history
- **Use case scenario:** On update, new version inserted; on get, all versions returned
- **Acceptance criteria:**
  - [x] Versions stored with incrementing numbers
  - [x] History array populated on retrieval
  - [x] Versions ordered DESC

## 4. Non-functional requirements
- **Performance:**
  - Read timeout: 5min
  - Write timeout: 10s
  - SQLite with index on (uuid, version)
- **Reliability:**
  - Liveness probe: `/liveness` (always OK)
  - Readiness probe: `/readiness` (checks DB ping)
- **Security:**
  - AES-256-GCM encryption
  - Password hashed via SHA-256 for key derivation
  - Passwords never stored
  - HttpOnly cookies
- **Scalability:**
  - Single-instance design
  - SQLite suitable for moderate load
- **Availability/UX:**
  - Prometheus metrics exposed
  - Structured logging (logrus)
  - Error pages for web UI
  - JSON errors for API

## 5. Interfaces
- **APIs and integrations:**
  - HTTP REST API: `GET /api/v1/{bid}` (Basic auth)
  - Web UI: HTML forms, cookie-based sessions
- **Protocols and data formats:**
  - HTTP/1.1
  - JSON for API responses
  - HTML for web UI
  - Base64-encoded encrypted data in DB
- **UI/UX constraints:**
  - Simple form-based interface
  - CodeMirror for content editing (planned)

## 6. Acceptance criteria
- System accepted if:
  - [x] Bins can be created, viewed, updated via web UI
  - [x] API returns encrypted bin data with correct password
  - [x] Versions preserved on updates
  - [x] Health checks functional
  - [x] Metrics collected
  - [ ] Code examples on edit page (TODO)
  - [ ] Version revert functionality (TODO)
  - [ ] Syntax highlighting/validation (TODO)

