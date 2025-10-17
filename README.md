# GORM + Oracle Example

About
This repository is a minimal, working example showing how to use GORM with Oracle via the official dialect at github.com/oracle-samples/gorm-oracle. It connects to Oracle, auto-migrates a model, and demonstrates basic CRUD.

Quick Start
- Prerequisites:
  - Go 1.20+ (CGO enabled)
  - Oracle database (XE or Autonomous DB)
  - Oracle Instant Client (when using the godror driver)

- Setup:
  1) Export environment variables (recommended: use the provided script):
     source ./setenv.sh
     This sets ORA_USER, ORA_PASSWORD, ORA_CONNECT_STRING, and optional ORA_LIB_DIR.
     Note: setenv.sh is excluded from Git via .gitignore to avoid leaking credentials.

  2) Run the example with the required build tag:
     go run -tags godror .

What It Does
- Builds a DSN from environment variables
- Opens a GORM DB using the Oracle dialect
- Auto-creates a CUSTOMERS table
- Executes Create, Read, Update, and List queries

Environment Variables
- ORA_USER          (e.g., ADMIN)
- ORA_PASSWORD      (your password)
- ORA_CONNECT_STRING (e.g., localhost:1521/XEPDB1 or an ADB connect descriptor)
- ORA_LIB_DIR       (optional; path to Instant Client for godror)

Links
- Oracle GORM Dialect: https://github.com/oracle-samples/gorm-oracle
- GORM: https://gorm.io
