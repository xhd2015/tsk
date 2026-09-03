---
name: tsk/working-on-task
description: >-
  When the user asks to work on a task: pick/approve a task; clarify project,
  intent kind, and E2E acceptance; reuse task cwd if it is a live linked
  worktree else wrk-create; drive a kck grok worker (new → watch → send →
  wait); update tsk stages and notes; verify; ship with wrk on user approve;
  tsk done after wrk --done.
---

# working-on-task

Playbook for the **main** agent to mimic daily work on a tsk task using three
CLIs: **`tsk`** (stages + notes), **`wrk`** (worktree + land), **`kck`**
(worker session). Self-contained — do not load other skill docs for this flow.
Extra flags: each tool’s `--help`. Procedure is the numbered sections below.

| Role | Owns |
|------|------|
| **Main** | Intake, `tsk` stages/notes/progress, `wrk` create/bring/ship, `kck` spawn/watch/send/wait, independent verify, present to user, `tsk done` |
| **Worker** | Clarify then implement in the worktree only; **never** `tsk done` / `wrk --done` |

## When to use

- User asks to **work on** / **drive** / **finish** a task (id or from tree)
- Goal: verified result, user-approved land, or a clean **blocked** stop with notes
- Not for ad-hoc create/list/note-only sessions

## Stages (main only)

Spine: `create → in_process → clarification → implementation → verification → summary → done`  
(alt terminal: `archived`). Prefer `tsk advance`; use `tsk stage` sparingly.

| Moment | Stage action |
|--------|----------------|
| §1, still `create` | `tsk advance <id>` → `in_process` |
| §2 answers recorded | `tsk clarify confirm -y <id>` → `implementation` |
| §3–§7 | stay `implementation`; progress notes |
| §8 verify starts | advance → `verification` |
| §8 fail / §9 followup | honest stage; progress `blocked` or back to implement work |
| §11 after `wrk --done` | summary note → **`tsk done <id>`** |
| User skips ship after approve | `tsk done` only if they accept verified-only |

## Must-notes (user-checkable)

Use `tsk note add --id <id> --label …` and `tsk progress add --id <id> --status …`.

| When | Label / kind | Content |
|------|--------------|---------|
| §2 | (text or `acceptance`) | project, intent kind, E2E acceptance |
| §3 | `worktree=<abs-path>` | path used (reused or created); note “reused existing worktree” when reusing |
| §3 | `brought=<basenames>` | if `--bring` used |
| §4 | `grok-session-id=<uuid>` | worker session id (codex: `codex-session-id=`) |
| §5–§7 | progress `in-progress` | milestones / kickoff |
| §8 | `verify` / progress | pass/fail + evidence paths |
| §8/§9 wait | note | presented to user / followup ask |
| §10–§11 | `ship` | land/tag/peel/`wrk --done` summary |

## 1. Pick a task

1. If the user **explicitly gave a task id** → use it. Run `tsk show <id>` and
   keep `cwd:` / `project:` for §3 (`project:` is the project filesystem path).
   Confirm with the user only if the id is ambiguous.
2. Else discover candidates:
   - Run `tsk project tree --all`.
   - If the current directory resolves to a project (`tsk project which`
     succeeds), **prefer suggesting** a non-terminal task under that project
     from the tree.
   - Otherwise propose from the full `--all` view (say briefly why).
3. **Ask the user for approval** of the chosen id before continuing.
4. Orient: `tsk status --format=agent <id>`. If still `create` →
   `tsk advance <id>` → `in_process`.

## 2. Clarify (with user)

Ask only what is not already clear on the task (`tsk show`):

| Field | Ask |
|-------|-----|
| **Project** | Correct checkout / path? (confirm `project:` from show) |
| **Intent kind** | `feature` / `doc change` / `issue report` / `bug repro` / other |
| **E2E acceptance** | Real, human-like checks that must pass for “done” |

Record answers in notes (and `tsk clarify add` as needed). Then
`tsk clarify confirm -y <id>` → `implementation`. Do **not** create a worktree
or spawn a worker until this is done.

## 3. Worktree (reuse or `wrk`)

**Default when the user explicitly gave the task id** — from `tsk show <id>`:

```text
cwd set AND path exists AND is a linked git worktree
  → reuse that path (do not wrk-create)
else
  → wrk <project-location> --no-config -t "task-[id]: <short description>"
```

- **`<project-location>`** = `project:` from `tsk show` (already a path; expand
  `~` as needed), or the path confirmed in §2 if show had none / user
  overrode it.
- **Linked worktree** = a secondary checkout (e.g. `.git` is a `gitdir:` file /
  listed under `git worktree list`) — **not** the project main repo alone.
  If `cwd` is missing, gone, or only the main checkout → create.
- If the agent picked the task (no explicit id from the user) → **create** with
  `wrk` unless the user later points at an existing worktree to reuse.

```bash
wrk <project-location> --no-config -t "task-[id]: <short description>"
```

Note `worktree=<abs-path>` (reused or created). On reuse, note text
**reused existing worktree**. If dependency projects are needed:

```bash
wrk --bring agent-pro dot-pkgs   # basenames; example only
```

Note `brought=`. All later edits/tests/ship use the worktree (and
`./external/…` for brought trees) — not the bare main checkout.

## 4. Spawn worker (`kck`)

Prefer **grok**. Do **not** pass `--submit` on `new` (kickoff is §6).

```bash
kck grok new --dir <worktree> "$(cat <<'EOF'
Work on tsk <id>: <title>

Project: …
Intent kind: …
E2E acceptance (main will re-run): …
…

Clarify anything unclear first. Wait for an explicit kickoff before implementing.
Do not run tsk done or wrk --done.
EOF
)"
```

`kck grok new` prepends `/brainstorm` — expected. Default new terminal (user-visible);
avoid `--here` unless asked. Capture `session-id:` → note `grok-session-id=`.
Codex if user asks: same shape with `kck codex new` / `codex-session-id=`
(no `codex wait` yet — use messages/snapshot).

## 5. Watch worker

While the worker clarifies (or after a kickoff turn), read progress from the
session — do not rely on memory alone:

| Action | Command |
|--------|---------|
| Chat | `kck grok messages <session-id>` |
| Pane | `kck grok snapshot <session-id>` |

Prefer `snapshot` when it may be waiting on input. Record meaningful progress
on the tsk task. Do not invent a monitor + hash watcher.

**Blocked** (stop; note unblock ask): prerequisite missing, decisions outside
agreed scope, pane stuck after snapshot + user check.

## 6. Kickoff implement

When worker clarification is finished, **`send`** the implement kickoff (plan
already agreed; implement only; leave evidence for E2E):

```bash
kck grok send --session-id <session-id> "…"
```

## 7. Wait for turn

```bash
kck grok wait <session-id>
```

Block until the turn completes (or timeout / not-running). Then continue to
verify (§8) or re-watch (§5) as needed.

## 8. Verify (main)

Do **not** trust the worker’s claim alone. Advance to `verification`. Re-run the
§2 **E2E acceptance** as a human would (real journeys, UI/runtime when listed).
Confirm evidence exists.

| Result | Next |
|--------|------|
| **FAIL** | Note evidence; `send` (§6) or re-watch (§5); do not present as done |
| **PASS** | Note verify pass; **present result to user and wait** for instruction |

## 9. User followup

If the user wants changes → loop **§5** / **§6** (watch / `send`), keep notes
honest. If they approve → §10.

## 10. Ship (`wrk`)

Only after **user approve**. cwd = feature worktree. **No PR** unless the user
explicitly asks for `--pr`.

```bash
wrk --add-all --commit -m "<summary>" \
  --merge-back --tag-next --push --sync --reinstall-local
```

If there are brought `./external/…` trees:

1. `wrk --status` in the feature worktree.
2. Land each brought dependency **least-deps-first** (same commit + merge-back
   + tag/push/sync/reinstall line in that tree).
3. `wrk --status` again on the feature worktree.

## 11. Cleanup and `tsk done`

When status is clean:

```bash
wrk --done
```

Note ship summary → `tsk done <id>`. Never `wrk --done` before land if unmerged
work remains only on the worktree.

## Report

| Outcome | Include |
|---------|---------|
| **Finished** | `done`; acceptance; verify evidence; worktree; session-id; ship/`wrk --done` |
| **Awaiting user** | verify pass notes; what to approve or follow up |
| **Blocked / verify fail** | stage; last messages/snapshot insight; exact unblock ask |

## Commands cheat sheet

```bash
# tsk
tsk show <id>
tsk project which
tsk project tree --all
tsk status --format=agent <id>
tsk clarify add <id> <question…>
tsk clarify confirm -y <id>
tsk advance <id>
tsk progress add --id <id> --status in-progress|done|blocked "…"
tsk note add --id <id> --label worktree=<path> "…"
tsk note add --id <id> --label grok-session-id=<uuid> "…"
tsk done <id>

# wrk (create when not reusing a live linked worktree cwd from tsk show)
wrk <project-location> --no-config -t "task-[id]: <short description>"
wrk --bring <basename…>
wrk --status
wrk --add-all --commit -m "…" --merge-back --tag-next --push --sync --reinstall-local
wrk --done

# kck (grok; codex same shape where available)
kck grok new --dir <worktree> "…"
kck grok messages <session-id>
kck grok snapshot <session-id>
kck grok send --session-id <session-id> "…"
kck grok wait <session-id>
```
