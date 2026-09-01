---
name: tsk/install
description: >-
  Install convenience CLI wrappers into ~/.local/bin (e.g. pmark → tsk project
  add) and ensure that directory is on PATH.
---

# install

```text
tsk install [--dry-run] <name>
tsk install pmark
tsk install --dry-run pmark
```

Writes a `#!/bin/sh` forwarder under `~/.local/bin/<name>` and ensures
`~/.local/bin` is on PATH (bash/zsh rc marker block). `--dry-run` probes
current state and prints `[dry-run] would` / `skip` lines without writing.

| Name | Forwards to |
|------|-------------|
| `pmark` | `tsk project add` |

After install:

```text
pmark "fix flaky CI"   # ≡ tsk project add "fix flaky CI"
```

Re-running `tsk install pmark` rewrites the wrapper (idempotent). Open a new
terminal if PATH was just patched, or `export PATH="$HOME/.local/bin:$PATH"`.
