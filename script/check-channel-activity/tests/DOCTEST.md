# check-channel-activity Test Cases

## Version
0.0.2

Decision tree for the `check-channel-activity` CLI: monitors a tsk channel's last
message activity, runs a notify command when idle ≥ threshold, persists anti-spam
state, and optionally loops with `--forever`.

# DSN (Domain Specific Notion)

- **check-channel-activity CLI** — standalone binary at `script/check-channel-activity/main.go`; reads channel storage under `TSK_HOME`, compares last activity to wall clock, optionally executes a notify command via `exec.Command`; prints four-line status on stdout; errors as single `Error: ...` line on stderr (exit 1); SIGINT/SIGTERM prints `stopped\n` on stderr (exit 0).
- **TSK_HOME** — storage root (default `~/.tsk`); tests isolate per leaf at `{WorkRoot}/.tsk`.
- **Channel index** — `$TSK_HOME/channels/index/<channel-id>` UTF-8 line `active/<id>` or `archive/<id>`; missing index → error exit 1.
- **Active channel** — `$TSK_HOME/channels/active/<channel-id>/` with `channel.json`, `messages.jsonl`, `msg-counter`.
- **Archived channel** — `$TSK_HOME/channels/archive/<channel-id>/`; readonly; checker errors exit 1 (does not notify).
- **Last activity** — last message `created_at` in `messages.jsonl` when any messages exist; else `channel.created_at`.
- **Idle threshold** — `--idle` duration (default `1h`); idle when `now - last_activity ≥ threshold`.
- **Notify command** — `--exec-if-idle-1h LINE` single shell command line (quote when invoking if LINE contains spaces); parsed with `github.com/mattn/go-shellwords` `Parse()` (no env expansion), then `exec.Command(argv[0], argv[1:]...)`; empty LINE → `Error: --exec-if-idle-1h requires a command`; parse error → `Error: parse exec command: ...`; runs only when idle and not already notified for current `last_activity_at`.
- **State file** — `$TSK_HOME/channels/state/<channel-id>.json` (override `--state-dir`); fields `channel_id`, `last_activity_at`, `last_notified_at`; prevents re-notify until `last_activity_at` advances.
- **Dry-run** — `--dry-run` prints `status: would notify (dry-run)` without exec or state write.
- **Forever loop** — `--forever` repeats check → sleep `--interval` (default `1m`) until signal; each tick prints status block.
- **Test hook** — `--max-ticks N` (undocumented in help) stops forever loop after N ticks; for doctest only.
- **Process-local binary** — `getCheckBin` builds once per process under an in-memory mutex into `os.MkdirTemp("", "check-channel-activity-doctest-bin-")` (one-process suite; no session disk flock).
- **Work root** — temp dir per leaf; holds isolated `TSK_HOME`, notify marker path, and touch script for exec verification.

## Tree Overview

```
check-channel-activity tests
├── active/                     # recent activity → no notify
│   └── recent-message/
├── idle/                       # stale activity → notify semantics
│   ├── notify/                 # first run executes command
│   ├── already-notified/       # state prevents re-exec
│   └── reset-on-message/       # new message resets anti-spam
├── empty-channel/              # no messages; last_activity = created_at
├── dry-run/
│   └── would-notify/           # idle but no exec/state
├── exec/
│   └── quoted-args/            # LINE with quoted spaces parses correctly
├── error/
│   ├── not-found/
│   └── archived/
├── forever/
│   └── max-ticks/              # --forever --max-ticks 2
├── signal/
│   └── sigint/                 # SIGINT graceful stop (slow)
└── help/
    └── root/
```

## Test Index

| # | Leaf | Scenario |
|---|------|----------|
| 1 | active/recent-message | recent message → `status: active`, no marker |
| 2 | idle/notify | old message → exec touches marker, state written |
| 3 | idle/already-notified | pre-seeded state → `already notified`, no marker |
| 4 | idle/reset-on-message | notify then new message → notify again |
| 5 | empty-channel | empty messages.jsonl, old `created_at` → notify |
| 6 | dry-run/would-notify | idle + `--dry-run` → would notify, no marker |
| 7 | exec/quoted-args | LINE with `"hello world"` → argv preserved |
| 8 | error/not-found | missing channel → exit 1 |
| 9 | error/archived | archive path → exit 1 |
| 10 | forever/max-ticks | two ticks: notified then already notified |
| 11 | signal/sigint | SIGINT → `stopped\n`, exit 0 |
| 12 | help/root | `-h` documents flags and LINE quoting |

## How to Run

```sh
cd script/check-channel-activity
doctest vet ./tests
doctest test ./tests

doctest test ./tests/active
doctest test ./tests/idle
doctest test ./tests/exec
doctest test ./tests/error
doctest test ./tests/forever
doctest test ./tests/signal
doctest test ./tests/help

# Skip slow leaves
doctest test ./tests --label=-slow
```

```go
import (
	"os/exec"
	"testing"
)

type Request struct {
	WorkRoot     string
	TskHome      string
	ChannelID    string
	Args         []string
	MarkerPath   string
	ExecScript   string
	ArgvPath     string   // exec/quoted-args: file recording parsed argv
	LastActivity string // RFC3339 expected in stdout
	ExtraEnv     []string
	SIGINTStop   bool   // run forever then send SIGINT (signal/sigint leaf)
}

type Response struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.SIGINTStop {
		return runWithSIGINT(t, req)
	}
	bin := getCheckBin(t)
	cmd := exec.Command(bin, req.Args...)
	cmd.Dir = req.WorkRoot
	cmd.Env = checkEnv(req)
	return captureCommandOutput(cmd)
}
```