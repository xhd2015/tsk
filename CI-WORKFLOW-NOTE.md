# CI workflow note

## Status

**updated** (pushed)

## Branch / remote / push

| Field | Value |
|-------|--------|
| Branch | `master-2026-08-06-use-go-best-practice-to-review-current-project` |
| Remote | `ssh://git@github.com/xhd2015/tsk.git` (`origin`) |
| CI commit SHA | `f6aedbc9eb26292c836f976020bab40e637ac543` (short: `f6aedbc`) — workflow + helper |
| Branch tip | `afdfa7e36855a06f90f27a2e45fb484177d429ff` (includes this note) |
| Push result | success — upstream set; CI push `2b08c3d..f6aedbc`, note push `f6aedbc..afdfa7e` |

## Paths changed (in `f6aedbc`)

- `.github/workflows/test.yml` — full doctest-pattern CI (was incomplete: container + bare `go test` / `--label-all`)
- `script/ci/coverage-package-table.py` — package coverage table for step summary (`github.com/xhd2015/tsk/…`)

## How to view Actions for this push

- Repo: https://github.com/xhd2015/tsk  
- Actions (all): https://github.com/xhd2015/tsk/actions  
- Branch filter: https://github.com/xhd2015/tsk/actions?query=branch%3Amaster-2026-08-06-use-go-best-practice-to-review-current-project  
- Workflow file on branch: https://github.com/xhd2015/tsk/blob/master-2026-08-06-use-go-best-practice-to-review-current-project/.github/workflows/test.yml  
- Commit: https://github.com/xhd2015/tsk/commit/f6aedbc9eb26292c836f976020bab40e637ac543  

Workflow name: **Test** (`on: push`, `on: pull_request`).

## How this differs from doctest’s workflow

| Aspect | doctest reference | this repo (`tsk`) |
|--------|-------------------|-------------------|
| Module / `COVERPKG` | `github.com/xhd2015/doctest/...` | `github.com/xhd2015/tsk/...` |
| Install doctest | `go install ./cmd/doctest` (checkout) | `go install github.com/xhd2015/doctest/cmd/doctest@latest` |
| Discovery / e2e | `!e2e` then `e2e` | same; no `e2e`-labeled leaves today |
| Package table | `script/ci/coverage-package-table.py` for doctest module | same helper shape; module prefix `github.com/xhd2015/tsk/` |
| Host | `ubuntu-latest` + `setup-go` from `go.mod` | same |
| xgo merge + artifacts | yes | yes |

## Workflow stages

1. `go test ./...` → `coverage-gotest.out`  
2. Doctest discovery (`--label '!e2e'`) → `coverage-doctest-discovery.out`  
3. Doctest e2e (`--label e2e`) → `coverage-doctest-e2e.out`  
4. Install xgo → merge profiles → `GITHUB_STEP_SUMMARY` + package table  
5. Upload coverage artifacts (`if: always()`)
