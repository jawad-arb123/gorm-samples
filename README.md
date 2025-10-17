# GORM + Oracle Example (Go + React)

About
Minimal example showing how to use GORM with Oracle (github.com/oracle-samples/gorm-oracle).  
The app connects to Oracle, auto-migrates a simple Customer table, exposes a tiny HTTP API, and serves a small React page to create/list customers.

Quick Start
1) Set environment variables (use the provided script):
   source ./setenv.sh
   This sets ORA_USER, ORA_PASSWORD, ORA_CONNECT_STRING, and (optionally) ORA_LIB_DIR for the Instant Client.
   Note: setenv.sh is ignored by Git to avoid leaking credentials.

2) Run the server (uses the godror build tag):
   go run -tags godror .

3) Open the UI in your browser:
   http://localhost:8080
   - Use the form to create a customer (name + email)
   - See the list update live via the API

What It Does
- Builds a DSN from environment variables and connects via the Oracle GORM dialect
- Auto-creates a CUSTOMERS table
- HTTP API:
  - GET  /api/customers     -> list customers
  - POST /api/customers     -> create a customer (JSON: { "name": "...", "email": "..." })
- Frontend served at / (web/index.html) using React (CDN)

Environment Variables
- ORA_USER            (e.g., ADMIN)
- ORA_PASSWORD        (your password)
- ORA_CONNECT_STRING  (e.g., localhost:1521/XEPDB1 or an ADB connect descriptor)
- ORA_LIB_DIR         (optional; path to Instant Client for godror on macOS/Linux)

Links
- Oracle GORM Dialect: https://github.com/oracle-samples/gorm-oracle
- GORM: https://gorm.io

Security
- Do not commit credentials. setenv.sh is in .gitignore.
