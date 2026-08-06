# Go best-practice review: `tsk`

**Scope:** codebase structure, CLI design, flags handling, package layout.  
**Method:** inspect live sources against [go-best-practice](https://github.com/xhd2015) topics: `cli` (color, dry-run, streaming), `flags-parsing` (types, subcommand, cut), `cmd-exec`, plus notes on `kool-create` / `go-embed-assets` applicability.  
**Constraint:** review only — no implementation changes (docs-only file).

**Module:** `github.com/xhd2015/tsk` (Go 1.25.10)  
**Packages:** `cmd/tsk`, `tskcli` (+ `pipeline`, `storage`), `channel` (+ `file`), `script/check-channel-activity` (+ `signals`)

---

## Executive summary

`tsk` is a solid, test-heavy task/channel CLI. It already follows several go-best-practice patterns well: thin `main`, `less-flags` with `HelpNoExit`, nested `channel` dispatch via `StopOnFirstArg`, domain interface under `channel`, and a single dry-run pipeline in `check-channel-activity`.

The largest gaps vs the skill recipes are:

1. **Color policy** is partial (`--color` + TTY auto only; no `--no-color` / `NO_COLOR`).
2. **Nested subcommand `--help`** is incomplete for `label` / `topic` / `clarify` leaves (and a few empty-arg defaults).
3. **External notify command** is a shell string + `os/exec` + hand-scanned flags, not `Cut` + `xgo/support/cmd`.

`kool-create` and `go:embed` assets are **not applicable** to this repo today (existing CLI, no generated UI tree).

---

## What already looks good

| Area | Evidence | Topic |
|------|----------|--------|
| Thin entrypoint | `cmd/tsk/main.go` only calls `tskcli.Run` | package layout |
| Flag library | Widespread `github.com/xhd2015/less-flags` | `flags-parsing` |
| Help + library-friendly exit | `Help("-h,--help", …).HelpNoExit()` + `errors.Is(err, lessflags.ErrHelp)` on leaf commands | `flags-parsing` |
| Root / no-toplevel-flags dispatch | Manual `-h`/`--help` and switch in `dispatch` | `flags-parsing/subcommand` |
| Nested channel parent flags | `StopOnFirstArg()` + merge of parent/leaf `--channel-id` / `--user` | `flags-parsing/subcommand` |
| Channel help depth | Per-leaf help for create/list/send/… and participant add/remove | `flags-parsing/subcommand` |
| Domain split | `channel.Store` interface + `channel/file` impl | package layout |
| Dry-run shape | One `decideAndAct` path; dry-run gates exec + state write | `cli/dry-run` |
| Status auto-format | Documented precedence; extensive doctests | CLI UX |
| Test depth | Large doctest trees under `cmd/tsk/tests` and `channel/tests` | quality |

---

## Findings (by severity)

### High

#### H1. Color policy incomplete vs `cli/color`

**Where:** `tskcli/status.go` (`--color` / `--plain`, `isStdoutTTY`), `tskcli/channel.go` (ANSI green/gray on success/status when TTY).

**Current behavior:**

- `status`: force color with `--color`; force plain/no-ANSI with `--plain`; otherwise auto color iff stdout is a char device.
- `channel` list/archive/send/…: color when `isStdoutTTY()` only; no flags.
- No `--no-color`.
- No `NO_COLOR` env handling.
- TTY check is `os.Stdout.Stat()` / `ModeCharDevice`, not `golang.org/x/term` (acceptable alternative per recipe, but not shared).

**Why it matters (recipe):** `cli/color` requires a three-mode model:

| Mode | Selection | Color |
|------|-----------|-------|
| Auto | neither flag | TTY and `NO_COLOR` empty → on |
| Always | `--color` | always on |
| Never | `--no-color` | always off |

Conflict must error: `--color and --no-color cannot be specified together`.

**Concrete bugs / UX holes:**

1. Users cannot force color off on a TTY without piping (except `status --plain`, which also changes box drawing).
2. Industry `NO_COLOR=1` is ignored (channel status lines and status diagram auto-color still emit SGR when TTY).
3. With less-flags bools, `--color=false` leaves `colorFlag == false`, then auto TTY turns color **back on** — so “false” does not mean Never.

**Recommended change:**

- Introduce shared `ColorMode` + `ResolveColor(mode, stdoutIsTTY, noColorEnv)` (copy from `cli/color`).
- On `status` (and any human-facing channel output): parse `--color` and `--no-color`; reject both; resolve once.
- Keep `--plain` as a **format** concern (ASCII boxes), orthogonal to color mode where possible; document interaction.
- Document `NO_COLOR=1` in `statusHelp` / top help.
- Prefer one helper used by both `status` and `channel` instead of duplicated ANSI constants + ad-hoc TTY checks.

---

#### H2. Nested leaves missing real `--help` (`label`, `topic`, `clarify`)

**Where:** `tskcli/label.go`, `tskcli/topic.go`, `tskcli/clarify.go` (add/list only).

**Pattern today:** parent level handles empty / `-h` / `--help` and prints parent help. Leaves take positionals only — no `lessflags.Help`.

**Broken UX examples:**

| Invoked | Likely result today |
|---------|---------------------|
| `tsk label add --help` | `parseID("--help")` → invalid task id |
| `tsk topic set --help` | “task id and path required” / invalid id |
| `tsk clarify add --help` | treated as id + question tokens |
| `tsk clarify list --help` | invalid id |
| `tsk clarify confirm --help` | works (lessflags help), but help text is **parent** `clarifyHelp()`, not confirm-specific |

**Recipe:** `flags-parsing/subcommand` — every level that users can land on must answer `-h`/`--help` with **that** level’s usage.

**Recommended change:**

- For flag-less leaves, still run a minimal parse:

  ```go
  remaining, err := lessflags.Help("-h,--help", labelAddHelp()).HelpNoExit().Parse(args)
  ```

- Or accept `--help` as first positional before `parseID`.
- Split help strings: `clarifyConfirmHelp`, `labelAddHelp`, `topicSetHelp`, etc.
- Empty args at confirm already fail “task id required”; consider printing confirm help when args empty (optional friendliness).

---

#### H3. Notify exec: shell string + hand flag scan vs `Cut` / `cmd-exec`

**Where:** `script/check-channel-activity/main.go`, `run.go`.

**Current design:**

- `--exec-if-idle-1h LINE` is a single `String` (shell-ish line).
- Presence of the flag is detected by **manual argv scan** (`hasExecFlag`), not by less-flags results / `**string`.
- Line is re-tokenized with `go-shellwords`, then `exec.Command(argv[0], argv[1:]...)`.
- Stdout/stderr inherited (good), but no `xgo/support/cmd` builder.

**Why it matters:**

| Topic | Expectation | Gap |
|-------|-------------|-----|
| `flags-parsing/cut` | Opaque command tail as raw argv; no re-parse of child flags | Single string requires quoting; child `--flags` fragile |
| `flags-parsing/types` | Use `**string` / presence, not hand-scan | `hasFlag` / loop over `args` |
| `cmd-exec` | Prefer `github.com/xhd2015/xgo/support/cmd` (Debug, Dir, Env, inherit I/O) | Raw `os/exec` only |

**Also:** flag name `--exec-if-idle-1h` hard-codes “1h” while `--idle` defaults to 1h but is configurable — name and semantics diverge.

**Recommended change:**

1. Prefer:

   ```text
   check-channel-activity --channel-id ID [options] --exec <cmd> [args...]
   ```

   with `lessflags.Cut("--exec", &execArgs)` (see `flags-parsing/cut`).

2. If a single-line form must stay for cron/config compatibility, keep it as a secondary alias, but implement presence via `**string` (unset vs empty) — drop hand-scan.

3. Run with `cmd.Debug().Run(execArgs[0], execArgs[1:]...)` or non-debug `cmd.Run` when quiet is desired (`cmd-exec`).

4. Rename or alias to `--exec-if-idle` (idle duration is already `--idle`).

5. In dry-run, print the **resolved argv** (or planned command line) so dry-run answers “what this run would exec”.

---

### Medium

#### M1. `tsk channel` with no subcommand errors instead of help

**Where:** `runChannel` — `len(remaining) == 0` → `tsk channel: subcommand required`.

**Recipe:** empty args at a dispatch-only level should print that level’s help (friendly default). Root and `label`/`topic`/`clarify` already do this.

**Recommended change:** `fmt.Print(channelHelp()); return nil` when no subcommand (same for any other dispatch node that still errors).

---

#### M2. `tskcli/channel.go` is a god file (~800 lines)

**Where:** flags, merge helpers, all channel subcommands, ANSI helpers, and all channel help strings in one file.

**Impact:** harder review, high merge conflict risk, mixed concerns (CLI parse vs presentation vs domain calls).

**Recommended package layout (without changing public UX):**

```text
tskcli/
  channel/
    dispatch.go      # runChannel + StopOnFirstArg
    create.go
    list.go
    send.go
    …
    help.go
    color.go         # or share with status color helper
```

Keep `channel` domain package as-is (`channel`, `channel/file`). CLI package name can stay under `tskcli` to avoid import cycles with `storage`.

---

#### M3. Inconsistent error surface

**Where:**

- `channelFail` prefixes `Error: …`.
- Most task commands use `tsk <cmd>: …` or bare `invalid transition: …`.
- `fail(err)` is an identity function (`return err`) — dead abstraction.
- `main` always `fmt.Fprintln(os.Stderr, err)` once (good for “error once”; covered by doctest).

**Recommended change:**

- One policy: either always `Error: <message>` (matches check-channel-activity) or always `tsk: <message>` — document and apply in one helper.
- Delete or repurpose `fail` (e.g. wrap with prefix) so call sites are meaningful.
- Keep single print in `main` (already correct).

---

#### M4. Global `currentCtx` for event recording

**Where:** `tskcli/util.go` `var currentCtx *invocationContext`; every handler calls `setCommand(currentCtx, …)`.

**Issues:** package-level mutable state; hard to test in parallel; nil risk if `Run` is not used.

**Recommended change:** pass `*invocationContext` (or a small `Env` struct with home + event sink) into `dispatch` / handlers. Event append can stay in `defer` on `Run`.

---

#### M5. Secondary binary lives under `script/` as importable `main`

**Where:** `script/check-channel-activity` is `package main` and a `go list` package.

**Idiomatic Go:** `cmd/<name>/` for all mains (`cmd/tsk`, `cmd/check-channel-activity`). `script/` is fine for non-module shell helpers, less so for first-class Go tools.

**Recommended change:** move to `cmd/check-channel-activity` (or keep path with a thin main + library package if reused). Update doctest paths / CI only as needed.

---

#### M6. Help text location split

**Where:** task command help in `tskcli/help.go`; channel help functions at bottom of `channel.go`.

**Recommended change:** either all help next to handlers (subcommand recipe preference) or all under `help_*.go` — pick one convention and stick to it when splitting channel.

---

### Low

#### L1. Dry-run output polish (`cli/dry-run`)

**Where:** `decideAndAct` returns `"would notify (dry-run)"`; status printer formats it.

**Good:** same path as live; gates exec + state write; already-notified still shared.

**Improve:**

- Prefix planned action lines with `[dry-run]` per recipe.
- In dry-run, still resolve/parse exec argv and print planned command (even if not run) so dry-run validates the real plan.

---

#### L2. `--max-ticks` hidden test hook

**Where:** parsed in `main.go`, documented only in doctest DOCTEST.md as test hook, omitted from `helpText`.

**Options:** document as “testing”, hide behind `//go:build` / env only, or leave as-is with a one-line help note so operators are not confused if they discover it.

---

#### L3. No project README

No root README describing install, `TSK_HOME`, command map, or relationship of `tsk` vs `check-channel-activity`. Doctests are excellent internal docs but not onboarding.

**Recommended:** short README (install, env vars, command tree, link to doctest trees). Not a go-best-practice topic per se, but CLI UX hygiene.

---

#### L4. Streaming N/A for most commands

**Topic:** `cli/streaming`.

Most `tsk` commands print small, finite results — buffering is fine. No change required unless long-running channel watch or status of many tasks is added later (then stream or NDJSON).

---

#### L5. `parseCreatedAt` placement

Defined in `list.go` but used from `next.go`. Minor cohesion issue; move to `util.go` or `storage` timestamps helper.

---

### Informational / N/A topics

| Topic | Applicability |
|-------|----------------|
| `kool-create` | N/A — project already exists; no need to rescaffold with `kool create`. |
| `go-embed-assets` | N/A — no SPA/extension assets, no `//go:embed` of generated trees. |
| `cli/skill-cli` | N/A — product CLI, not a skill host. |
| `cli/inline-tui-mouse` | N/A — no inline TUI. |
| `flags-parsing/collect` | Optional later if parent flags should be reconstructed/forwarded to a child process; current merge helper is fine for in-process dual placement. |

---

## Package layout snapshot

```text
cmd/tsk/main.go              # thin main ✅
tskcli/                      # CLI + orchestration (large; channel.go oversized)
tskcli/storage/              # task FS home, index, transitions ✅
tskcli/pipeline/             # status rendering ✅
channel/                     # domain types + Store interface ✅
channel/file/                # FS implementation ✅
script/check-channel-activity/  # second main (prefer cmd/) ⚠️
```

**Strengths:** domain vs CLI vs storage separation for channels; pipeline isolated for status art.  
**Weaknesses:** all task CLI verbs flat in `tskcli` root; channel CLI megfile; tool binary under `script/`.

Suggested target shape (evolutionary, not big-bang):

```text
cmd/tsk/
cmd/check-channel-activity/
tskcli/                 # Run, dispatch, shared parse/color/errors
tskcli/task/            # create, list, status, advance, …
tskcli/channelcmd/      # channel CLI only (name avoids clash with domain)
channel/ + channel/file/
internal/…              # optional: event sink, color, exec helpers if unexported
```

---

## Flags & subcommand matrix (quick audit)

| Command | less-flags | Help at level | Notes |
|---------|------------|---------------|--------|
| (root) | manual | yes | no global flags — OK |
| create/list/show/status/advance/stage/next/followup/done | yes | yes | good |
| label / topic / clarify (parent) | manual | yes | good |
| label add/rm, topic set/mkdir, clarify add/list | no | **no** | **H2** |
| clarify confirm | yes | parent help only | tighten help text |
| channel (parent) | StopOnFirstArg | yes on -h; empty → error | **M1** |
| channel leaves + participant add/remove | yes | yes | good model to copy |
| check-channel-activity | yes | yes | exec/cut/color/docs gaps |

---

## Recommended change order

1. **H1 color** — shared `ResolveColor`, `--no-color`, `NO_COLOR`; fix status + channel human output; add doctests (force off on TTY, `NO_COLOR` auto, conflict).  
2. **H2 nested help** — mirror channel leaf pattern for label/topic/clarify; fix confirm help string.  
3. **M1** — empty `tsk channel` → help (one-liner + test).  
4. **H3 exec** — `Cut("--exec")` and/or `**string` presence; `cmd` package; rename flag; dry-run prints plan.  
5. **M3/M4** — error helper + pass context explicitly.  
6. **M2/M5** — split `channel.go`; move secondary binary under `cmd/`.  
7. **L1–L3** — dry-run polish, README, optional `--max-ticks` docs.

---

## Suggested acceptance checks (when implementing)

- `go test ./...` and existing doctest suites stay green.  
- New doctests:  
  - `tsk status --no-color <id>` on forced TTY path / or env `NO_COLOR` with auto.  
  - `tsk status --color --no-color` → conflict error.  
  - `tsk label add --help` / `tsk topic set --help` → usage, exit 0.  
  - `tsk channel` (no args) → channel help, exit 0.  
  - check-channel-activity dry-run prints planned argv when idle.  
- No double error lines on stderr (`ux/error-once` still holds).

---

## Summary table

| ID | Severity | Topic | One-line fix |
|----|----------|--------|--------------|
| H1 | High | `cli/color` | `--color` / `--no-color` / `NO_COLOR` three-mode resolve |
| H2 | High | `flags-parsing/subcommand` | leaf `--help` for label/topic/clarify |
| H3 | High | `flags-parsing/cut`, `cmd-exec` | Cut/argv + `xgo/support/cmd`; stop hand-scan |
| M1 | Medium | `flags-parsing/subcommand` | empty `channel` → help |
| M2 | Medium | package layout | split `channel.go` |
| M3 | Medium | CLI UX | unify error prefix; drop no-op `fail` |
| M4 | Medium | package layout | pass ctx, no package global |
| M5 | Medium | package layout | `cmd/check-channel-activity` |
| M6 | Medium | consistency | one home for help strings |
| L1 | Low | `cli/dry-run` | `[dry-run]` + planned command |
| L2 | Low | flags docs | document or hide `--max-ticks` |
| L3 | Low | onboarding | root README |
| — | N/A | `kool-create`, `go-embed-assets` | not applicable |

---

*Generated as a review-only deliverable. Implementation should be phased and doctest-backed per the order above.*
