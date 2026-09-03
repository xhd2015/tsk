---
name: tsk/working-on-task
description: >-
  When the user asks to work on a task: clarify intent, project,
  prerequisites, definition of done, verification steps, and artifacts;
  spawn a kck worker after plan approval; main polls via messages/snapshot
  and verifies independently before done.
---

# working-on-task

Use this topic when the user asks to **work on** / **drive** / **finish** a
tsk task (by id or from `tsk tree` / `tsk project tree`).

Stage mechanics live in `workflow`. This topic is the **agent playbook**:
intake → approve → spawn worker (`kck`) → poll → independent verify → report.

```text
Main agent (this session)              Worker (kck grok|codex pane)
  Intake + plan with user
  User approves plan
  kck … new "brief…"  ───────────────►  implements in task cwd
  ← session-id; note on tsk task
  loop: messages + snapshot  ◄────────  progress / claims
  worker claims finished
  Independent verify (human perspective)
  tsk done | blocked report
```

**Roles**

| Role | Owns |
|------|------|
| **Main** | Intake, plan approval, `kck` spawn/poll, tsk stages, independent verify, `tsk done` / report |
| **Worker** | Implementation only against approved plan + DoD; **never** `tsk done` |

## When to use

- User picks a task id or says “work on this todo”
- Goal is **`done`** or a clean **blocked** stop with a report
- **Not** for ad-hoc create/list/note-only sessions — use Default agent
  workflow in the root skill index

## Phase 1 — Intake (required before spawn)

Clarify with the user (prefer `tsk clarify add <id> …`; record durable answers
in notes). Skip a field only when it is already clear on the task (e.g.
`project` / `cwd` from `tsk show`).

| Field | Ask | Pass criteria |
|-------|-----|----------------|
| **Intent** | What should be true when this is finished? (one sentence) | Shared understanding; no silent reinterpretation |
| **Project** | Which checkout / origin? | Matches task `project` / `cwd`, or user confirms a different target |
| **Prerequisites** | Tools, auth, deps (e.g. browser-agent, `kck`)? | Each item: ready, optional, or **BLOCKER** with unblock step |
| **Definition of done** | Pass/fail outcome criteria — not “looks good” | Clear success bar |
| **Verification steps** | Ordered human-style steps **main** will run after the worker claims finished | Main can re-run without asking “is this done?” |
| **Artifacts needed** | Evidence to collect (logs, screenshots, transcripts, URLs, note paths, session ids) | Worker leaves them; main checks they exist |

**Definition of done** = *what* success means.  
**Verification steps** = *how* main will check it.  
**Artifacts** = *what proof* must exist.

Also produce a short **plan** for user approval. Do **not** spawn a worker
while intake is incomplete or the plan is unapproved.

Then:

1. Orient: `tsk status --format=agent <id>`.
2. Claim if still `create`: `tsk advance <id>` → `in_process`.
3. After intake answers: `tsk clarify confirm -y <id>` → `implementation`
   (main may advance here when spawning the worker).

## Phase 2 — Spawn worker (`kck`)

After the user **approves the plan**, start a dedicated agent session in the
task workspace (`cwd` from `tsk show`). Prefer **grok** unless the user asks
for Codex.

```bash
kck grok new --dir <task-cwd> --submit "$(cat <<'EOF'
Work on tsk <id>: <title>

Approved plan:
…

Intent: …
Definition of done: …
Verification steps (main will re-run these — leave artifacts):
…
Artifacts needed: …

Implement only. Do not run tsk done. When finished, say so clearly with evidence paths.
EOF
)"
```

Capture **`session-id:`** from stdout and store it on the task:

```bash
tsk note add --id <id> --label grok-session-id=<uuid> \
  'kck worker for implementation'
```

Use `kck codex new` the same way when the runner is Codex (`codex-session-id=`
or the project’s usual label).

**Footgun:** `kck grok new` currently **prepends `/brainstorm`** to the prompt.
After plan approval the brief must state clearly: plan already approved —
**implement only** (do not re-brainstorm). If a no-brainstorm launch path
appears later, prefer that.

Do not use `--here` for this flow unless the user asks; default new terminal
keeps the worker pane **visible to the user**.

## Phase 3 — Poll worker (main ↔ `kck`)

Main interacts with the spawned session; do not rely on memory alone.

| Action | Command | Purpose |
|--------|---------|---------|
| Spawn | `kck grok new "…"` | Prints `session-id:`; open user-visible pane |
| Chat / claims | `kck grok messages --session-id <id>` | Message updates from the worker |
| Live pane | `kck grok snapshot <id>` | Terminal text the user also sees (stuck / waiting / idle) |

Codex: `kck codex messages` / `kck codex snapshot` with the same session id.

**Poll loop** until the worker clearly claims finished, asks a blocking
question, or appears stuck:

1. `messages --session-id` for progress and claims.
2. `snapshot` when status is unclear or the pane may be waiting on input
   (prefer snapshot for “is it waiting on the user?”).
3. Optional: `kck grok send --session-id <id> "…"` only for small nudges
   already covered by the approved plan; otherwise escalate to the user.
4. Record meaningful progress on the tsk task (`progress` / notes).

**Blocked** during poll (stop; report — do not invent answers):

- Prerequisite BLOCKER
- Worker needs decisions outside the approved plan
- Pane stuck with no path forward after snapshot + user check

## Phase 4 — Independent verify (main)

When the worker **claims finished**, main **does not** trust that claim alone.

1. Advance / stay in `verification` as appropriate.
2. Re-run the intake **verification steps** from a **human perspective**.
3. Confirm every **artifact** exists and matches DoD.
4. Prefer **verify-on-behalf-of-user** (scenario / human-shaped verify +
   transcript) when the surface needs real journeys, UI, or bring-up — not
   only the worker’s self-reported tests.

| Verify result | Next |
|---------------|------|
| **PASS** | Summary note → `tsk done <id>` |
| **FAIL** | Evidence back to worker (`send` or new brief) or user; do not `done` |
| **BLOCKED** | Report unblock ask; leave stage honest |

Worker self-check / doctest green is **not** sufficient for `tsk done` when
intake listed human verification steps or UI/runtime artifacts.

## Phase 5 — Report

Always end with a user-facing report:

| Outcome | Include |
|---------|---------|
| **Finished** | Stage `done`; intent; verify steps run; artifacts (paths); worker `session-id` |
| **Blocked** | Current stage; last `messages` / `snapshot` insight; exact unblock ask |
| **Verify FAIL** | What failed vs DoD; whether worker was re-nudged |

Do not claim finished without verify evidence. Do not leave a long run in
`implementation` with no progress note.

## tsk stage ownership (main)

| Stage | Main does |
|-------|-----------|
| `implementation` | Spawn/poll worker; progress notes |
| `verification` | Independent verify (steps + artifacts) |
| `summary` | Short summary note |
| `done` | Only after verify PASS |
| `archived` | Shelve without claiming success |

Prefer `tsk advance` along the spine; use `tsk stage` sparingly.

## Commands cheat sheet

```bash
# tsk
tsk show <id>
tsk status --format=agent <id>
tsk clarify add <id> <question…>
tsk clarify confirm -y <id>
tsk advance <id>
tsk progress add --id <id> --status in-progress|done|blocked "…"
tsk note add --id <id> --label grok-session-id=<uuid> "…"
tsk done <id>
tsk archive <id>

# kck (grok; same shape for codex)
kck grok new --dir <cwd> --submit "…"
kck grok messages --session-id <session-id>
kck grok snapshot <session-id>
kck grok send --session-id <session-id> "…"   # optional nudge
```

See also: `workflow` (stages), `note` (notes/progress), `project` (project tree).
External: `kck` CLI; **verify-on-behalf-of-user** for human-shaped verify.
