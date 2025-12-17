# Software Requirements Specification (SRS)

## 1. Introduction
- **Purpose:** Functional/non-functional requirements for ConfigBin.
- **Scope:** Remote encrypted config storage; Web UI + API.
- **Definitions:**
  - Bin: Encrypted config container (UUID).
  - AES-GCM: AES-256 encryption (GCM mode).

## 2. General Description
- **Context:** HTTP service; SQLite storage; Password auth (not stored).
- **Constraints:**
  - Passwords irrecoverable.
  - Single binary; SQLite.
  - HTTPS recommended.

## 3. Functional Requirements
### 3.1 FR-1: Creation
- **Action:** Create bin with auto-generated UUID/Password.
- **Criteria:**
  - [x] UUID/Password (8 char alphanumeric) generated.
  - [x] Content AES-GCM encrypted.
  - [x] Stored in SQLite.
  - [x] Cookie set (24h, HttpOnly, Lax).
  - [x] Redirect to view.

### 3.2 FR-2: Viewing (Web)
- **Action:** Display content via Web UI.
- **Criteria:**
  - [x] Auth form if no cookie.
  - [x] Password validation.
  - [x] Content decrypted/displayed.
  - [x] Edit form accessible.

### 3.3 FR-3: Auth
- **Action:** Authenticate via password.
- **Criteria:**
  - [x] Password validation.
  - [x] Cookie set.
  - [x] Redirect to view.

### 3.4 FR-4: Update
- **Action:** Update content (new version).
- **Criteria:**
  - [x] Version incremented.
  - [x] History preserved.
  - [x] Re-encrypted with same password.

### 3.5 FR-5: API Access
- **Action:** `GET /api/v1/{bid}` via Basic Auth.
- **Criteria:**
  - [x] Basic auth handling.
  - [x] JSON response (data + history).

### 3.6 FR-6: History
- **Action:** Store/retrieve versions.
- **Criteria:**
  - [x] Versions stored (incrementing).
  - [x] DESC order retrieval.

## 4. Non-Functional
- **Performance:** Read <5min; Write <10s; DB indexed.
- **Reliability:** Liveness/Readiness probes.
- **Security:** AES-256-GCM; SHA-256 key derivation; No stored passwords; HttpOnly cookies.
- **Scalability:** Single-instance; SQLite.
- **Observability:** Prometheus metrics; Logrus logging.

## 5. Interfaces
- **API:** HTTP REST `GET /api/v1/{bid}`.
- **UI:** HTML forms.
- **Format:** JSON (API), HTML (UI), Base64 (Storage).

## 6. Acceptance
- [x] Web UI: Create, View, Update.
- [x] API: Encrypted retrieval.
- [x] Versioning: Preserved.
- [x] Ops: Health checks, Metrics.
- [ ] TODO: Code examples, Revert UI, Syntax highlight.
