# tsk Test Cases

## Version
0.0.2

Decision tree covering the `tsk` CLI: task creation (inbox and topic placement),
listing and filtering, show/status display, stage transitions (advance, stage,
clarify, followup, done), topic management, label management, `next` selection,
Slack-like **channel** spaces (create/list/archive/delete, send, messages,
participants), and append-only `events.jsonl` auditing.

# DSN (Domain Specific Notion)

- **tsk CLI** — standalone binary; subcommand dispatcher with `less-flags` per handler; no top-level flags except `-h`/`--help`; empty args or help flags print `topHelp` on stdout (exit 0); each handler uses `lessflags.ErrHelp` for command help; errors to stderr once (no duplicate from `fail()` + `main`), exit code 1 on failure; non-empty stdout ends with trailing `\n`, empty stdout has no bytes; `create` success prints task id + `\n` on stdout.
- **TSK_HOME** — storage root env var (default `~/.tsk`); tests isolate per run at `{WorkRoot}/.tsk`.
- **TSK_DATE** — optional env var (`YYYY-MM-DD`) for deterministic timestamps; all tests set `TSK_DATE=2026-07-09`.
- **Work root** — temp directory per leaf; holds isolated `TSK_HOME`.
- **counter** — plain-text monotonic integer at `{TSK_HOME}/counter`; flock on read-modify-write for ID allocation.
- **index/<id>** — UTF-8 single line: relative path from `TSK_HOME` to task directory; updated on create, stage rename, topic move; atomic write via temp + rename.
- **events.jsonl** — append-only audit log; one JSON object per CLI invocation (success or failure).
- **Task directory** — name `<id>-<stage>-<slug>/` under `inbox/` (no topic) or `topics/<path>/` (topic tree); contains `task.json`, `context/` (empty on create), and `clarify/` (during clarification with `batch.json`).
- **task.json** — metadata: `id`, `title`, `slug`, `labels` (sorted), `topic_path` (null in inbox), `stage`, `created_at`, `updated_at`, `stage_history`.
- **Slug** — lowercase, non-letter-digit → `-`, collapse, trim, max 64 runes; immutable after create.
- **Stage workflow** — `create → in_process → clarification → implementation → verification → summary → user_followup (loop to clarification) OR done`; `done` is terminal.
- **Transitions** — `advance` follows allowed edges; `stage` sets stage directly (invalid jumps error); `clarify confirm -y` confirms all items and auto-advances to `implementation`; `followup` writes `context/followup-<ts>.md` and sets `user_followup`; `done` only from `summary` or `user_followup`.
- **topic set** — moves entire task directory; `--inbox` or empty path → `inbox/`; updates `topic_path` and `index/<id>`.
- **topic mkdir** — creates topic directory tree under `topics/`.
- **next** — stdout prints id of oldest `in_process` task by `created_at`, or empty stdout when none.
- **status** — pipeline view of a task; flags `--format=diagram|agent`, `--color` (default on TTY for diagram), `--plain` (ASCII boxes for diagram, no ANSI). **Default format** when `--format` is absent and neither `--color` nor `--plain` is present: if `TSK_STATUS_FORMAT=agent|diagram` is set use that; else if an agent host is detected (`CODEX_THREAD_ID`, `PI_CODING_AGENT`, or parent/grandparent process name via lean `agentrunner.Detect`) use `agent`; else `diagram`. Precedence (highest first): `--format` present → that value; `--color` or `--plain` present → diagram; `TSK_STATUS_FORMAT`; detect → agent; else diagram. **diagram**: hand-made compact pipeline via `tskcli/pipeline` (~40 col, 3-line boxes with labels inside mid-rows; tee borders `├`/`┤` OK on summary/user_followup); geometry: ●/create center-aligned on spine; **refine** left-rail from left mid of `user_followup` to left mid of `clarification` (no rail under done/◉); **no followup** right-rail from right mid of `summary` to right mid of `done`; **satisfied** vertical spine label under `user_followup` (no `satisfied►`); **done→◉** dead end; semantic ANSI overlay when colored (current=green bold, visited=grey, edge-into-current=orange). Exact art sealed by `status/diagram-golden` + `status/plain-golden` `expected.txt`. **agent**: strict 2-row plain-text spine (`create -> … -> done` with `name[doing]` / `(name)` / bare marks) plus back line (`refine`, `questions`, `user_followup` — no `satisfied` on art) and facts block (`id`, `title`, `stage`, `terminal`, `topic`, `dir` in that order, then after art `advance`/`next`); `title` is exact `task.json` create title (same key as `tsk show`); `topic` is always present above `dir:` — slash-joined `topic_path` segments (e.g. `eng/backend`) when set, or exactly `(not classified yet)` for inbox/null `topic_path` (differs from `tsk show`, which prints `topic: inbox`); `dir` is the absolute task directory path (from index + `TSK_HOME`; key `dir:` only — no `path`/`path_rel`); no ANSI even with `--color`; no rectangle chrome; no width cap. Invalid `--format` → exit 1, single stderr line. `context/pipeline.mmd` ignored (may remain on disk harmlessly).
- **Request.Args** — CLI arguments passed to `tsk` (subcommand + flags + positionals).
- **Request.TaskID** — task id for multi-step setups and assertions.
- **Request.ExtraEnv** — optional `KEY=value` strings appended to the child `tsk` process env (after `tskEnv` strips host agent / format-override vars for stable defaults).
- **Process-local binary** — `getTskBin` builds `tsk` once per process under an in-memory mutex into `os.MkdirTemp("", "tsk-doctest-bin-")` (one-process suite; no session disk flock).
- **channels/** — under `TSK_HOME`; layout `index/<channel-id>` (line `active/<id>` or `archive/<id>`), `active/<id>/` and `archive/<id>/` each with `channel.json` (metadata only), `participants.jsonl`, `messages.jsonl`, `msg-counter`; `tombstones/<id>.json` blocks id reuse after delete.
- **channel.json** — metadata only: `id`, `name`, `status` (`active`|`archived`), `created_at`, `updated_at` (no embedded `participants`).
- **participants.jsonl** — one `{"handle","joined_at"}` per line, sorted by `handle` on write; on create, creator handle only (no `agent` auto-join).
- **Channel message** — JSONL line `{"id", "sender", "body", "created_at"}`; monotonic ids via `msg-counter` (flock).
- **Channel identity** — precedence `--user <handle>` > `TSK_USER` env > `$USER`; empty `$USER` errors; handle format `^[a-z0-9][a-z0-9_-]{0,63}$` lowercase; channel id same format, default `Slugify(name)` when `--channel-id` omitted; `--user` on create, send, messages, participants, participant add/remove (not list/archive/delete).
- **Channel parent flags** — `--channel-id` and `--user` may appear directly after `channel` before the action subcommand (parent peel); merge with leaf flags (same value OK; different values → conflict error). `list` hard-rejects parent `--channel-id` / `--user`. Nested forms work: `channel --channel-id X participant add bob`.
- **Channel membership gate** — non-participants cannot `send`, `messages`, `participants`, `participant add`, or `participant remove`; archived channels are readonly for mutations but `messages`/`participants`/`list --all` still work.
- **Channel CLI output** — create prints `channel-id\n`; archive `archived <id>\n`; delete `deleted <id>\n`; send `sent message <id>\n`; participant add `added <handle>\n`; remove self `left <channel-id>\n`; remove other `removed <handle>\n`; list human table + gray count footer (TTY); `--json` arrays without ANSI; errors single stderr line `Error:` prefix, exit 1; every channel command appends `events.jsonl` with `command: channel`.
- **Request.ChannelID** / **Request.ChannelName** — channel id and display name for multi-step channel setups and assertions.

## Tree Overview

```
tsk tests
├── create/                       # tsk create
│   ├── no-topic/                 # inbox placement, index, task.json
│   ├── with-topic/               # topics/<path>/ placement
│   └── with-labels/              # --label flags, sorted labels
├── advance/                      # tsk advance
│   ├── basic/                    # create → advance renames dir + index
│   └── invalid/
│       └── stage-jump/           # create → stage implementation errors
├── clarify/                      # tsk clarify *
│   └── confirm/                  # add questions, confirm -y → implementation
├── topic/                        # tsk topic *
│   ├── set-to-topic/             # inbox → topic path, dir move
│   └── set-to-inbox/             # topic → inbox, topic_path null
├── next/                         # tsk next
│   └── oldest/                   # two in_process → older id on stdout
├── done/                         # tsk done
│   └── from-summary/             # at summary → done, terminal stage
├── followup/                     # tsk followup
│   └── basic/                    # at summary → user_followup + context file
├── status/                       # tsk status (diagram | agent formats)
│   ├── diagram-golden/           # --format=diagram exact stdout == expected.txt (unicode; no-followup rail aligned)
│   ├── plain-golden/             # --plain exact stdout == expected.txt (ASCII; no-followup rail aligned)
│   ├── color-box-only/           # --color at implementation: green on box; left refine │ outside box SGR
│   ├── diagram/                  # clarification + --color → compact art + green highlight
│   ├── at-create/                # create stage + │ create │ + green ANSI
│   ├── at-done/                  # done stage + │ done │ + green ANSI
│   ├── no-color-pipe/            # piped stdout → box chars, no ANSI
│   ├── plain-ascii/              # --plain → ASCII + boxes, no ANSI (soft; see plain-golden)
│   ├── compact-width/            # every stdout line width ≤ 42 (~40 geometry)
│   ├── box-format/               # each stage has box mid-row (tee borders OK)
│   ├── arrows/                   # ▼ spine; left refine ►│ clarification; ◄ into done
│   ├── edge-labels/              # claim/research/confirmed/questions/satisfied order
│   ├── fork-semantics/           # no followup vs questions; vertical satisfied; left refine
│   ├── agent/                    # --format=agent (2-row plain + facts)
│   │   ├── spine/                # create: spine order, create[doing], facts, no boxes
│   │   ├── title/                # facts title: exact create title; order id→…→topic→dir
│   │   ├── dir/                  # facts dir: absolute task path after topic; no path_rel
│   │   ├── topic/                # create --topic eng/backend → topic: eng/backend above dir
│   │   ├── two-rows/             # back line refine+questions; no satisfied on art
│   │   ├── marks-mid/            # implementation[doing]; past bare; future (name)
│   │   ├── at-clarification/     # blocked advance; next clarify confirm
│   │   ├── at-summary/           # next followup + done
│   │   ├── at-user-followup/     # user_followup[doing]; next refine + done
│   │   ├── at-done/              # terminal true; done[doing]; advance blocked
│   │   └── no-ansi/              # --format=agent --color → no ANSI
│   ├── format-invalid/           # --format=nope → exit 1; stderr once
│   ├── help/                     # status --help documents --format
│   └── auto-format/              # bare status format auto-select (detect / TSK_STATUS_FORMAT / flags)
│       ├── bare-human/           # no agent env → diagram (not agent facts)
│       ├── env-codex/            # CODEX_THREAD_ID → agent
│       ├── env-pi/               # PI_CODING_AGENT → agent
│       ├── tsk-status-format-agent/    # TSK_STATUS_FORMAT=agent → agent
│       ├── tsk-status-format-diagram/  # TSK_STATUS_FORMAT=diagram overrides CODEX → diagram
│       ├── force-diagram-flag/   # CODEX + --format=diagram → diagram
│       ├── force-plain-blocks-auto/    # CODEX + --plain → diagram, not agent
│       └── force-color-blocks-auto/    # CODEX + --color → diagram, not agent facts
├── show/                         # tsk show
│   └── basic/                    # metadata block for id
├── list/                         # tsk list
│   └── filter/                   # --stage create filters ids
├── events/                       # events.jsonl audit
│   └── append/                   # any command appends one line
├── help/                         # --help / -h at every level
│   ├── root-empty/               # no args → top help
│   ├── root-flag/                # --help → top help
│   ├── root-h/                   # -h → top help
│   ├── create/                   # create --help → flags
│   ├── topic/                    # topic --help → set, mkdir
│   ├── label/                    # label --help → add, rm
│   └── clarify/                  # clarify --help → add, list, confirm
├── channel/                      # tsk channel *
│   ├── create/                   # channel create
│   │   ├── basic/                # slug id, creator-only participants.jsonl, index, empty messages.jsonl
│   │   ├── custom-id/            # --channel-id
│   │   ├── user-flag/            # --user carol sets creator participant
│   │   ├── duplicate/            # same id → error
│   │   ├── tombstone-block/      # delete then recreate → error
│   │   └── invalid-id/           # bad id format → error
│   ├── list/                     # channel list
│   │   ├── empty/                # no channels
│   │   ├── active-only/          # archived hidden by default
│   │   ├── all/                  # --all shows archived
│   │   ├── json/                 # --json valid, no ANSI
│   │   ├── deleted-hidden/       # tombstoned absent from --all
│   │   ├── reject-parent-channel-id/  # parent --channel-id list → error
│   │   └── reject-parent-user/   # parent --user list → error
│   ├── archive/                  # channel archive
│   │   ├── basic/                # dir move, status archived, excluded from default list
│   │   ├── parent-channel-id/    # channel --channel-id X archive
│   │   ├── readonly/             # send blocked
│   │   ├── not-found/            # missing id → error
│   │   └── already-archived/     # double archive → error
│   ├── delete/                   # channel delete
│   │   ├── active/               # tombstone; not in list --all
│   │   ├── archived/             # delete from archive/
│   │   └── not-found/            # missing id → error
│   ├── send/                     # channel send
│   │   ├── basic/                # participant sends; jsonl + counter
│   │   ├── parent-channel-id/    # channel --channel-id X send
│   │   ├── parent-user/          # channel --channel-id X --user bob send
│   │   ├── same-channel-id/      # parent+leaf same --channel-id OK
│   │   ├── conflict-channel-id/  # parent+leaf different --channel-id → conflict
│   │   ├── conflict-user/        # parent+leaf different --user → conflict
│   │   ├── not-member/           # non-participant → error
│   │   ├── archived/             # archived → error
│   │   ├── not-found/            # missing channel → error
│   │   ├── user-flag/            # --user overrides identity
│   │   └── user-over-env/        # --user wins over TSK_USER
│   ├── messages/                 # channel messages
│   │   ├── human/                # chronological transcript
│   │   ├── parent-channel-id/    # channel --channel-id X messages
│   │   ├── json/                 # --json array
│   │   ├── limit/                # --limit 1 last message
│   │   ├── empty/                # no messages yet
│   │   ├── not-member/           # non-participant → error
│   │   ├── archived/             # archived channel readable
│   │   └── not-found/            # missing channel → error
│   ├── participant/              # channel participant *
│   │   ├── add/                  # add bob
│   │   ├── parent-channel-id-add/ # channel --channel-id X participant add bob
│   │   ├── add-dup/              # idempotent add
│   │   ├── remove-self/          # leave without handle
│   │   ├── remove-other/         # remove bob
│   │   ├── not-member/           # non-participant cannot add
│   │   ├── last-participant/     # cannot remove last member
│   │   ├── archived-readonly/    # add/remove blocked when archived
│   │   ├── not-found/            # missing channel → error
│   │   └── participants-json/    # participants --json roster
│   ├── help/                     # channel help
│   │   ├── root/                 # channel --help lists subcommands
│   │   ├── create/               # create --help documents --channel-id
│   │   ├── top/                  # tsk --help lists channel
│   │   └── send/                 # send --help documents --user
│   └── events/                   # channel events.jsonl
│       └── append/               # channel create appends audit line
└── ux/                           # CLI UX conventions
    ├── error-once/               # advance missing id → single stderr line
    └── create-prints-id/         # create prints id\n on stdout
```

## Test Case Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | create/no-topic | `tsk create "add dark mode"` → `inbox/1-create-add-dark-mode/`, index, task.json |
| 2 | create/with-topic | `tsk create --topic eng/backend "x"` → dir under `topics/eng/backend/` |
| 3 | create/with-labels | `tsk create --label bug --label urgent "x"` → sorted labels in task.json |
| 4 | advance/basic | create + `tsk advance` → dir renamed to `…-in_process-…`, index updated |
| 5 | advance/invalid/stage-jump | create + `tsk stage <id> implementation` → error, dir unchanged |
| 6 | clarify/confirm | add 2 questions, `clarify confirm -y` → implementation, batch confirmed |
| 7 | topic/set-to-topic | inbox task → `topic set <path>` → dir moved, index updated |
| 8 | topic/set-to-inbox | topic task → `topic set --inbox` → inbox, `topic_path` null |
| 9 | next/oldest | two `in_process` tasks → stdout = older id |
| 10 | done/from-summary | at summary → `tsk done` → stage done, dir renamed |
| 11 | followup/basic | at summary → `tsk followup` → `user_followup` + `context/followup-*.md` |
| 12 | status/diagram | at clarification + `--color` → compact box art, `│ clarification │`, edge labels `refine`/`confirmed`, green on clarification (geometry sealed by diagram-golden) |
| 55 | status/diagram-golden | `--format=diagram` (no color) → stdout byte-equal to unicode `expected.txt`; no-followup `┐`/`│`/`┘` same column |
| 56 | status/plain-golden | `--plain` → stdout byte-equal to ASCII `expected.txt`; no-followup `+`/`|`/`+` same column |
| 57 | status/color-box-only | at implementation + `--color` → green on box; leading left-rail `│` outside box SGR |
| 25 | status/at-create | create only + `status --color` → `│ create │` with green ANSI |
| 26 | status/at-done | at done + `status --color` → `│ done │` with green ANSI |
| 27 | status/no-color-pipe | clarification, piped → `│ clarification │`, box chars, no ANSI |
| 28 | status/plain-ascii | `status --plain` → `| create |` or `+` ASCII boxes, no ANSI (soft; plain-golden exact) |
| 29 | status/compact-width | full diagram → every stdout line rune width ≤ 42 |
| 30 | status/box-format | full diagram → each stage has box mid-row (tee borders/padding OK) |
| 31 | status/arrows | full diagram → ≥6 `▼`, `►│ clarification` + `└─refine`, `◄` into done, followup before `◉` |
| 32 | status/edge-labels | full diagram → edge labels in correct order (claim, research, confirmed, questions, vertical satisfied) |
| 33 | status/fork-semantics | full diagram → no followup vs questions; vertical satisfied (no satisfied►); left refine; done dead end |
| 34 | status/agent/spine | `--format=agent` at create → spine order, `create[doing]`, core facts (id/title/stage/terminal/topic/dir; inbox topic `(not classified yet)`), no rect chrome, no ANSI |
| 44 | status/agent/title | create `"add dark mode"` → agent facts `title: add dark mode` after `id:` before `stage:`; order locked through `topic` → `dir` |
| 45 | status/agent/dir | create `"add dark mode"` → agent facts `dir: <abs path>` after `topic:`; absolute; contains `inbox/<id>-create-add-dark-mode`; no `path`/`path_rel` |
| 46 | status/agent/topic | `create --topic eng/backend "…"` → agent facts `topic: eng/backend` after `terminal:` before `dir:`; `dir` contains `topics/eng/backend/` |
| 35 | status/agent/two-rows | agent art has `user_followup`/`refine`/`questions`; no `satisfied` on art |
| 36 | status/agent/marks-mid | at implementation → `implementation[doing]`; past bare; future `(…)` |
| 37 | status/agent/at-clarification | `clarification[doing]`; `advance: blocked`; next mentions clarify confirm |
| 38 | status/agent/at-summary | `summary[doing]`; next has followup + done |
| 39 | status/agent/at-user-followup | `user_followup[doing]`; advance→clarification; next advance + done |
| 40 | status/agent/at-done | `terminal: true`; `done[doing]`; advance blocked |
| 41 | status/agent/no-ansi | `--format=agent --color` → no `\x1b[` |
| 42 | status/format-invalid | `--format=nope` → exit 1; single stderr line |
| 43 | status/help | `status --help` documents `--format` |
| 47 | status/auto-format/bare-human | bare `status` + cleared agent env → diagram (box art; no agent facts) |
| 48 | status/auto-format/env-codex | `CODEX_THREAD_ID=t1` + bare `status` → agent (`id:`/`title:`/`topic:`/`dir:`) |
| 49 | status/auto-format/env-pi | `PI_CODING_AGENT=1` + bare `status` → agent |
| 50 | status/auto-format/tsk-status-format-agent | `TSK_STATUS_FORMAT=agent` + cleared host → agent |
| 51 | status/auto-format/tsk-status-format-diagram | `TSK_STATUS_FORMAT=diagram` + CODEX → diagram (env overrides detect) |
| 52 | status/auto-format/force-diagram-flag | CODEX + `--format=diagram` → diagram |
| 53 | status/auto-format/force-plain-blocks-auto | CODEX + `--plain` → diagram/plain, not agent facts |
| 54 | status/auto-format/force-color-blocks-auto | CODEX + `--color` → diagram (may ANSI), not agent facts |
| 13 | show/basic | `tsk show <id>` → metadata block with title, stage, labels |
| 14 | list/filter | `tsk list --stage create` → matching ids one per line |
| 15 | events/append | `tsk create` → `events.jsonl` gains one audit line |
| 16 | help/root-empty | `tsk` (no args) → exit 0; stdout has `Usage:` + command list; stderr empty |
| 17 | help/root-flag | `tsk --help` → exit 0; top help on stdout; stderr empty |
| 18 | help/root-h | `tsk -h` → exit 0; stdout contains `Usage:` |
| 19 | help/create | `tsk create --help` → create usage with `--label` and `--topic` |
| 20 | help/topic | `tsk topic --help` → lists `set`, `mkdir` subcommands |
| 21 | help/label | `tsk label --help` → lists `add`, `rm` subcommands |
| 22 | help/clarify | `tsk clarify --help` → lists `add`, `list`, `confirm` |
| 23 | ux/error-once | `tsk advance` (no id) → exit 1; `task id required` on stderr exactly once |
| 24 | ux/create-prints-id | `tsk create "hello"` → stdout `1\n`; inbox dir created; stderr empty |
| 58 | channel/create/basic | `tsk channel create "Eng Alerts"` → `eng-alerts\n`, active dir, alice-only participants.jsonl, metadata-only channel.json, empty messages.jsonl |
| 59 | channel/create/custom-id | `--channel-id my-room` → `my-room\n`, `channels/active/my-room/` |
| 59a | channel/create/user-flag | `create --user carol` → carol-only participants.jsonl (not alice) |
| 60 | channel/create/duplicate | second create same id → exit 1, `Error:` on stderr |
| 61 | channel/create/tombstone-block | delete then recreate same id → error; tombstone remains |
| 62 | channel/create/invalid-id | `--channel-id "BAD ID"` → exit 1 |
| 63 | channel/list/empty | no channels → empty or zero-count list |
| 64 | channel/list/active-only | archived hidden from default list |
| 65 | channel/list/all | `--all` shows archived channels |
| 66 | channel/list/json | `--json` valid array, no ANSI |
| 67 | channel/list/deleted-hidden | tombstoned channel absent from `list --all` |
| 68 | channel/archive/basic | move to `archive/`, status archived, excluded from default list |
| 69 | channel/archive/readonly | archived channel rejects send |
| 70 | channel/archive/not-found | archive missing id → error |
| 71 | channel/archive/already-archived | double archive → error |
| 72 | channel/delete/active | active delete → tombstone, `deleted <id>\n`, not in list |
| 73 | channel/delete/archived | delete from archive/ → tombstone |
| 74 | channel/delete/not-found | delete missing id → error |
| 75 | channel/send/basic | participant send → `sent message 1\n`, jsonl line |
| 76 | channel/send/not-member | non-participant send → error |
| 77 | channel/send/archived | send on archived → error |
| 78 | channel/send/not-found | send missing channel → error |
| 79 | channel/send/user-flag | `--user bob` sets message sender |
| 79a | channel/send/user-over-env | `TSK_USER=alice` + `--user bob` → sender bob |
| 80 | channel/messages/human | chronological human transcript |
| 81 | channel/messages/json | `--json` message array |
| 82 | channel/messages/limit | `--limit 1` returns last message only |
| 83 | channel/messages/empty | no messages → success, empty transcript |
| 84 | channel/messages/not-member | non-participant read → error |
| 85 | channel/messages/archived | archived channel messages readable |
| 86 | channel/messages/not-found | messages missing channel → error |
| 87 | channel/participant/add | `added bob\n`, bob in roster |
| 88 | channel/participant/add-dup | idempotent re-add bob |
| 89 | channel/participant/remove-self | `left <id>\n` when no handle |
| 90 | channel/participant/remove-other | `removed bob\n` |
| 91 | channel/participant/not-member | non-participant cannot add |
| 92 | channel/participant/last-participant | cannot remove last member |
| 93 | channel/participant/archived-readonly | add/remove blocked when archived |
| 94 | channel/participant/not-found | participant add on missing channel → error |
| 95 | channel/participant/participants-json | `participants --json` roster array |
| 96 | channel/help/root | `channel --help` lists subcommands |
| 97 | channel/help/create | `channel create --help` documents `--channel-id` |
| 98 | channel/help/top | `tsk --help` lists `channel` |
| 99 | channel/help/send | `channel send --help` documents `--user` (not `--as`) |
| 100 | channel/events/append | channel create appends `events.jsonl` with `command: channel` |

## How to Run

```sh
# Verify tree structure (no test execution)
doctest vet ./tests

# Run all leaves (expect RED until tsk CLI is implemented)
doctest test ./tests

# Run by command family
doctest test ./tests/create
doctest test ./tests/advance
doctest test ./tests/clarify
doctest test ./tests/topic
doctest test ./tests/next
doctest test ./tests/done
doctest test ./tests/followup
doctest test ./tests/status
doctest test ./tests/show
doctest test ./tests/list
doctest test ./tests/events
doctest test ./tests/help
doctest test ./tests/ux
doctest test ./tests/channel
doctest test ./tests/channel/create
doctest test ./tests/channel/list
doctest test ./tests/channel/archive
doctest test ./tests/channel/delete
doctest test ./tests/channel/send
doctest test ./tests/channel/messages
doctest test ./tests/channel/participant
doctest test ./tests/channel/help
doctest test ./tests/channel/events

# Run individual leaves
doctest test ./tests/create/no-topic
doctest test ./tests/advance/basic
doctest test ./tests/advance/invalid/stage-jump
doctest test ./tests/clarify/confirm
doctest test ./tests/topic/set-to-topic
doctest test ./tests/next/oldest
doctest test ./tests/done/from-summary
doctest test ./tests/followup/basic
doctest test ./tests/status/diagram-golden
doctest test ./tests/status/plain-golden
doctest test ./tests/status/color-box-only
doctest test ./tests/status/diagram
doctest test ./tests/status/at-create
doctest test ./tests/status/at-done
doctest test ./tests/status/no-color-pipe
doctest test ./tests/status/plain-ascii
doctest test ./tests/status/compact-width
doctest test ./tests/status/box-format
doctest test ./tests/status/arrows
doctest test ./tests/status/edge-labels
doctest test ./tests/status/fork-semantics
doctest test ./tests/status/agent
doctest test ./tests/status/agent/spine
doctest test ./tests/status/agent/title
doctest test ./tests/status/agent/dir
doctest test ./tests/status/agent/topic
doctest test ./tests/status/agent/two-rows
doctest test ./tests/status/agent/marks-mid
doctest test ./tests/status/agent/at-clarification
doctest test ./tests/status/agent/at-summary
doctest test ./tests/status/agent/at-user-followup
doctest test ./tests/status/agent/at-done
doctest test ./tests/status/agent/no-ansi
doctest test ./tests/status/format-invalid
doctest test ./tests/status/help
doctest test ./tests/status/auto-format
doctest test ./tests/status/auto-format/bare-human
doctest test ./tests/status/auto-format/env-codex
doctest test ./tests/status/auto-format/env-pi
doctest test ./tests/status/auto-format/tsk-status-format-agent
doctest test ./tests/status/auto-format/tsk-status-format-diagram
doctest test ./tests/status/auto-format/force-diagram-flag
doctest test ./tests/status/auto-format/force-plain-blocks-auto
doctest test ./tests/status/auto-format/force-color-blocks-auto
doctest test ./tests/show/basic
doctest test ./tests/list/filter
doctest test ./tests/events/append
doctest test ./tests/help/root-empty
doctest test ./tests/help/root-flag
doctest test ./tests/help/root-h
doctest test ./tests/help/create
doctest test ./tests/help/topic
doctest test ./tests/help/label
doctest test ./tests/help/clarify
doctest test ./tests/ux/error-once
doctest test ./tests/ux/create-prints-id
doctest test ./tests/channel/create/basic
doctest test ./tests/channel/list/json
doctest test ./tests/channel/send/basic
doctest test ./tests/channel/send/user-flag
doctest test ./tests/channel/send/user-over-env
doctest test ./tests/channel/create/user-flag
doctest test ./tests/channel/participant/last-participant
doctest test ./tests/channel/help/top
doctest test ./tests/channel/events/append
```

```go
import (
	"os/exec"
	"testing"
)

type Request struct {
	WorkRoot string
	TskHome  string
	Args     []string
	TaskID   int
	Title    string
	Topic    string
	Labels   []string
	Stage    string
	Message     string   // followup message body
	ChannelID   string   // channel id for multi-step channel setups
	ChannelName string   // channel display name
	ExtraEnv    []string // KEY=value injected into child tsk env (after tskEnv base strip)
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Run(t *testing.T, req *Request) (*Response, error) {
	bin := getTskBin(t)
	cmd := exec.Command(bin, req.Args...)
	cmd.Dir = req.WorkRoot
	cmd.Env = tskEnv(req)
	return captureCommandOutput(cmd)
}
```