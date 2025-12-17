# ConfigBin

**Secure, Encrypted Configuration Storage Service.**

ConfigBin is a lightweight web service for storing configuration files securely. It uses client-side password handling to ensure the server never stores the decryption key, making your data accessible only to those with the password.

## Features

- **🔒 Zero-Knowledge Security**: Content is encrypted with AES-256-GCM. Passwords are never stored in the database.
- **📜 Version Control**: Every update creates a new version. Access full history of changes.
- **🌐 Dual Access**:
  - **Web UI**: Clean interface for humans to create and edit configs.
  - **HTTP API**: Simple JSON API with Basic Auth for automated retrieval.
- **📦 Lightweight**: Built with Go 1.24, utilizing SQLite for storage. Single binary deployment.
- **🐳 Container Ready**: Optimized Docker image based on Alpine 3.23.
- **👀 Observable**: Built-in Prometheus metrics and health checks.

## Getting Started

### Prerequisites

- **Go**: 1.24 or later
- **Deno**: For running task scripts
- **Docker**: Optional, for containerized execution

### Installation

```bash
git clone https://github.com/korchasa/ConfigBin.git
cd ConfigBin
./run.ts init
```

### Running Locally

Start the development server (default port 8080):

```bash
./run.ts dev
```

### Running with Docker

```bash
docker build -t configbin .
docker run -p 8080:8080 configbin
```

## Usage

### 1. Create a Bin
Open `http://localhost:8080` in your browser. The system will generate a unique Bin ID and a secure Password.
> **⚠️ IMPORTANT:** Save the password immediately. It is not stored on the server and cannot be recovered if lost.

### 2. Add Content
Paste your configuration data into the text area and click **Save**.

### 3. API Access
Retrieve your configuration programmatically using the Bin ID and Password:

```bash
curl -u "user:<YOUR_PASSWORD>" http://localhost:8080/api/v1/<BIN_ID>
```

Response:
```json
{
  "success": true,
  "result": {
    "id": "...",
    "data": "...",
    "version": 1,
    "created_at": "...",
    "history": [...]
  }
}
```

## Development

The project uses a `run.ts` script for common tasks:

- **Check Code**: Run formatting, linting, and tests.
  ```bash
  ./run.ts check
  ```
- **Run Specific Test**:
  ```bash
  ./run.ts test-one pkg/encryptor/aes
  ```
- **Clean Artifacts**:
  ```bash
  ./run.ts cleanup
  ```

## License

MIT

