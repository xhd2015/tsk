# tsk Test Cases

## Version
0.0.2

Decision tree covering the `tsk` CLI: task creation (inbox and topic placement),
listing and filtering, show/status display, stage transitions (advance, stage,
clarify, followup, done), topic management, label management, `next` selection,
and append-only `events.jsonl` auditing.

# DSN (Domain Specific Notion)

- **tsk CLI** — standalone binary; subcommand dispatcher with `less-flags` per handler; no top-level flags; errors to stderr, exit code 1 on failure; non-empty stdout ends with trailing `\n`, empty stdout has no bytes.
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
- **status** — ASCII pipeline diagram with marker on current stage (not a list filter).
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
├── status/                       # tsk status (pipeline diagram)
│   └── diagram/                  # at clarification → ASCII + marker
├── show/                         # tsk show
│   └── basic/                    # metadata block for id
├── list/                         # tsk list
│   └── filter/                   # --stage create filters ids
└── events/                       # events.jsonl audit
    └── append/                   # any command appends one line
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
| 12 | status/diagram | at clarification → stdout has pipeline ASCII + current stage marker |
| 13 | show/basic | `tsk show <id>` → metadata block with title, stage, labels |
| 14 | list/filter | `tsk list --stage create` → matching ids one per line |
| 15 | events/append | `tsk create` → `events.jsonl` gains one audit line |

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
doctest test ./tests/show/basic
doctest test ./tests/list/filter
doctest test ./tests/events/append
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