---
status: Accepted
ratified_by: 5903b59
playbook:
  repo: wiradeltaid/ops
  path: research/wdi-ecosystem-strategy/coding-playbook/
  local: D:\Developer\wiradeltaid\ops\research\wdi-ecosystem-strategy\coding-playbook\
  rev: 5903b59
reads:
  - 01-principles.md
  - 02-architecture-and-structure.md
  - 03-essential-conventions.md
  - 04-file-size-and-cohesion.md
  - 06-tooling-and-ratchet.md
  - stack/go.md
excludes:
  - 07-ui-architecture-and-design-system.md
  - 05-realtime-and-sync-protocols.md
  - stack/rust.md
  - stack/slint.md
  - stack/react-typescript.md
  - stack/kotlin.md
  - stack/python.md
---

# stack — codebase guide

**Loaded when:** writing or reviewing code.

## 1. Toolchains & Runtimes

- **Bahasa:** Go 1.22+ (Headless API Client & Mail Service).
- **Networking:** Standard library `net/http` dengan timeout wajib.

## 2. Command Verifikasi

```powershell
go test -v -race ./...
```
