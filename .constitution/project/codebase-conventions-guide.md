---
status: Accepted
ratified_by: 5903b59
---

# conventions — codebase guide

**Loaded when:** writing or reviewing code.

## WDI Engineering Playbook Integration

Proyek pustaka Mailketing Go ini mengadopsi Single Source of Truth (SSOT) rekayasa terpusat WDI:
- **Konvensi Inti (`03-essential-conventions.md`):** Tiga lapis penegakan `[L1-Tool]`, `[L2-Guard]`, `[L3-Review]`.
- **Backend Go (`stack/go.md` & `03`):**
  - Pembungkusan error kontekstual `%w`.
  - Timeout eksplisit pada seluruh panggilan HTTP keluar (maks 10 detik).
  - Sanitasi logging dan larangan mencetak API token Mailketing ke stream log.
  - Zero panic di jalur eksekusi biasa.
