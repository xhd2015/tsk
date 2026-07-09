# tsk Test Cases

## Version
0.0.2

Decision tree covering the `tsk` CLI: task creation (inbox and topic placement),
listing and filtering, show/status display, stage transitions (advance, stage,
clarify, followup, done), topic management, label management, `next` selection,
and append-only `events.jsonl` auditing.

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
- **status** — pipeline view of a task; flags `--format=diagram|agent` (default `diagram`), `--color` (default on TTY for diagram), `--plain` (ASCII boxes for diagram, no ANSI). **diagram**: hand-made compact pipeline via `tskcli/pipeline` (~34 col max, 3-line boxes with labels inside `│ … │` rows); semantic ANSI overlay (current=green bold, visited=grey, edge-into-current=orange). **agent**: strict 2-row plain-text spine (`create -> … -> done` with `name[doing]` / `(name)` / bare marks) plus back line (`refine`, `questions`, `user_followup` — no `satisfied` on art) and facts block (`id`, `stage`, `terminal`, `advance`/`next`); no ANSI even with `--color`; no rectangle chrome; no 36-col cap. Invalid `--format` → exit 1, single stderr line. `context/pipeline.mmd` ignored (may remain on disk harmlessly).
- **Request.Args** — CLI arguments passed to `tsk` (subcommand + flags + positionals).
- **Request.TaskID** — task id for multi-step setups and assertions.
- **Session fixtures** — doctest injects `DOCTEST_SESSION_ID`; `getTskBin` builds once per session to `{cache}/bin/tsk` with file lock across leaf processes.

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
│   ├── diagram/                  # clarification + --color → compact art + green highlight
│   ├── at-create/                # create stage + │ create │ + green ANSI
│   ├── at-done/                  # done stage + │ done │ + green ANSI
│   ├── no-color-pipe/            # piped stdout → box chars, no ANSI
│   ├── plain-ascii/              # --plain → ASCII + boxes, no ANSI
│   ├── compact-width/            # every stdout line width ≤ 36
│   ├── box-format/               # each stage has │ <stage> │ box row
│   ├── arrows/                   # ▼ main flow, branch arrows, followup before ◉
│   ├── edge-labels/              # claim/research/confirmed/questions/satisfied order
│   ├── fork-semantics/           # no followup vs questions rows; satisfied ► into done rail
│   ├── agent/                    # --format=agent (2-row plain + facts)
│   │   ├── spine/                # create: spine order, create[doing], facts, no boxes
│   │   ├── two-rows/             # back line refine+questions; no satisfied on art
│   │   ├── marks-mid/            # implementation[doing]; past bare; future (name)
│   │   ├── at-clarification/     # blocked advance; next clarify confirm
│   │   ├── at-summary/           # next followup + done
│   │   ├── at-user-followup/     # user_followup[doing]; next refine + done
│   │   ├── at-done/              # terminal true; done[doing]; advance blocked
│   │   └── no-ansi/              # --format=agent --color → no ANSI
│   ├── format-invalid/           # --format=nope → exit 1; stderr once
│   └── help/                     # status --help documents --format
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
| 12 | status/diagram | at clarification + `--color` → compact box art, `│ clarification │`, width ≤ 36, edge labels `refine`/`confirmed`, green on clarification |
| 25 | status/at-create | create only + `status --color` → `│ create │` with green ANSI |
| 26 | status/at-done | at done + `status --color` → `│ done │` with green ANSI |
| 27 | status/no-color-pipe | clarification, piped → `│ clarification │`, box chars, no ANSI |
| 28 | status/plain-ascii | `status --plain` → `| create |` or `+` ASCII boxes, no ANSI |
| 29 | status/compact-width | full diagram → every stdout line rune width ≤ 36 |
| 30 | status/box-format | full diagram → each stage has `│ <stage> │` (or ascii `| <stage> |`) box row |
| 31 | status/arrows | full diagram → ≥6 `▼`, branch `►`/`──►`, `◄──` refine, followup before `◉` |
| 32 | status/edge-labels | full diagram → edge labels in correct order (claim, research, confirmed, questions, satisfied) |
| 33 | status/fork-semantics | full diagram → no followup on horizontal branch; questions separate; satisfied has ►; no ╰──▼ on done |
| 34 | status/agent/spine | `--format=agent` at create → spine order, `create[doing]`, facts, no rect chrome, no ANSI |
| 35 | status/agent/two-rows | agent art has `user_followup`/`refine`/`questions`; no `satisfied` on art |
| 36 | status/agent/marks-mid | at implementation → `implementation[doing]`; past bare; future `(…)` |
| 37 | status/agent/at-clarification | `clarification[doing]`; `advance: blocked`; next mentions clarify confirm |
| 38 | status/agent/at-summary | `summary[doing]`; next has followup + done |
| 39 | status/agent/at-user-followup | `user_followup[doing]`; advance→clarification; next advance + done |
| 40 | status/agent/at-done | `terminal: true`; `done[doing]`; advance blocked |
| 41 | status/agent/no-ansi | `--format=agent --color` → no `\x1b[` |
| 42 | status/format-invalid | `--format=nope` → exit 1; single stderr line |
| 43 | status/help | `status --help` documents `--format` |
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

# Run individual leaves
doctest test ./tests/create/no-topic
doctest test ./tests/advance/basic
doctest test ./tests/advance/invalid/stage-jump
doctest test ./tests/clarify/confirm
doctest test ./tests/topic/set-to-topic
doctest test ./tests/next/oldest
doctest test ./tests/done/from-summary
doctest test ./tests/followup/basic
doctest test ./tests/status/diagram
doctest test ./tests/status/at-create
doctest test ./tests/status/at-done
doctest test ./tests/status/no-color-pipe
doctest test ./tests/status/plain-ascii
doctest test ./tests/status/compact-width
doctest test ./tests/status/box-format
doctest test ./tests/status/arrows
doctest test ./tests/status/edge-labels
doctest test ./tests/status/agent
doctest test ./tests/status/agent/spine
doctest test ./tests/status/agent/two-rows
doctest test ./tests/status/agent/marks-mid
doctest test ./tests/status/agent/at-clarification
doctest test ./tests/status/agent/at-summary
doctest test ./tests/status/agent/at-user-followup
doctest test ./tests/status/agent/at-done
doctest test ./tests/status/agent/no-ansi
doctest test ./tests/status/format-invalid
doctest test ./tests/status/help
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
	Message  string // followup message body
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