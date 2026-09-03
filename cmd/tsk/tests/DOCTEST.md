# tsk Test Cases

## Version
0.0.2

Decision tree covering the `tsk` CLI: task creation (inbox and topic placement),
listing and filtering, show/status display, stage transitions (advance, stage,
clarify, followup, done), permanent **delete** (optional `--recursive` for
nested sub-tasks), topic management, label management, `next` selection,
Slack-like **channel** spaces (create/list/archive/delete, send, messages,
participants), append-only `events.jsonl` auditing, task **progress**
entries with optional status, a **search** command over titles/notes/progress/topic
notes, and a global **tree** view of all tasks organized by topic.

# DSN (Domain Specific Notion)

- **tsk CLI** — standalone binary; subcommand dispatcher with `less-flags` per handler; no top-level flags except `-h`/`--help`; empty args or help flags print `topHelp` on stdout (exit 0); each handler uses `lessflags.ErrHelp` for command help; errors to stderr once (no duplicate from `fail()` + `main`), exit code 1 on failure; non-empty stdout ends with trailing `\n`, empty stdout has no bytes; `create` success prints task id + `\n` on stdout.
- **TSK_HOME** — storage root env var (default `~/.tsk`); tests isolate per run at `{WorkRoot}/.tsk`.
- **TSK_DATE** — optional env var (`YYYY-MM-DD`) for deterministic timestamps; all tests set `TSK_DATE=2026-07-09`.
- **Work root** — temp directory per leaf; holds isolated `TSK_HOME`.
- **counter** — plain-text monotonic integer at `{TSK_HOME}/counter`; flock on read-modify-write for ID allocation.
- **index/<id>** — UTF-8 single line: relative path from `TSK_HOME` to task directory; updated on create and topic move (not on stage change); atomic write via temp + rename.
- **events.jsonl** — append-only audit log; one JSON object per CLI invocation (success or failure).
- **Task directory** — name `[id]-<slug>/` under `inbox/` (no topic), `topics/<path>/` (topic tree), or nested under a parent task dir; contains `task.json`, `context/` (empty on create), and `clarify/` (during clarification with `batch.json`).
- **task.json** — metadata: `id`, `title`, `slug`, `labels` (sorted), `topic_path` (null in inbox), optional `parent_id` (nested sub-tasks), optional `cwd` (CLI recording directory), optional `project` `{id,name}` (canonical origin key + basename), `stage`, `created_at`, `updated_at`, `stage_history`.
- **project** — `tsk project add|tree|list|which|register|unregister`. Prefer `project.origin` on tasks when git remote exists; else registered `project.name`. Manual registry `projects.json`; auto ledger `projects-auto.json` (upsert on add; main-repo `location` only, tilde-form; local-TZ `first_seen_at`/`last_seen_at`; legacy ledger `cwd` migrated to `location` on read). `list` default/`--all` = union auto+registered with `tasks=` and `LOCATION` column; `--json` includes `location`; `--auto` / `--registered` select one source; `--active` filters. `show` prints task `cwd:` in tilde form and `project: <location>` when resolvable, else name, else origin. `tree` = task forest; default also scans nested git dirs (max depth 3) via `scan_repo` and streams root-first (`--streaming` default; `--no-streaming` buffers). No mark import.
- **add --parent** — `tsk add --parent <id> <title>` nests under the parent task directory (any depth); child inherits parent `topic_path`; mutually exclusive with `--topic`.
- **Slug** — lowercase, non-letter-digit → `-`, collapse, trim, max 64 runes; immutable after create.
- **Stage workflow** — `create → in_process → clarification → implementation → verification → summary → user_followup (loop to clarification) OR done`; `archived` is an alternate off-spine terminal. `done` and `archived` are terminal.
- **Transitions** — `advance` follows allowed edges; `stage` sets stage directly (invalid jumps error; `done` still only via stage from `summary`/`user_followup`); `clarify confirm -y` confirms all items and auto-advances to `implementation`; `followup` writes `context/followup-<ts>.md` and sets `user_followup`; `done` / `archive` from any non-terminal stage.
- **topic set** — moves entire task directory; `--inbox` or empty path → `inbox/`; updates `topic_path` and `index/<id>`. Alias refs resolve to the canonical path before the move.
- **topic mkdir** — creates topic directory tree under `topics/`. Errors if the path is an alias for another topic.
- **topic.json** — optional metadata in `topics/<path>/topic.json`: `path`, `title`, `aliases`, `notes`, `updated_at`. Created by `topic note` / `topic alias add`. Reserved filename; not a task dir.
- **topic where** — stdout is the absolute topic directory + `\n`; `Error: topic not found: <ref>` (exit 1) when the dir does not exist. Accepts path or alias.
- **topic info** — facts `path`, `title`, `aliases`, `dir`, `notes`, `tasks` (exact `topic_path` match, not descendants), `subtopics` (child dirs that are not `[id]-<slug>`). Missing `topic.json` uses defaults (title = last path segment). `--json` encodes one object, no ANSI.
- **topic alias** — `alias add <topic> <alias>` appends to `topic.json`; conflict if another topic already claims the alias or a different topic dir exists at that path.
- **topic note** — `note [--label LABEL]... <topic> <text...>` **appends** one `{ts,text,labels?}` line to `topics/<path>/notes.jsonl` (does not replace). Labels support bare or `key=value` (same rules as task notes). Migrates a legacy `topic.json` `"notes"` blob into the jsonl once, then clears the blob. Unlabeled lines omit `labels`.
- **topic notes** — `notes [--json] [--limit N] [--label LABEL]... <topic>` lists the journal (oldest first; `--limit` is last N; `--label` AND-filters with bare-key presence matching). Missing jsonl → `0 notes` / `[]`. `info` `notes:` is the **count**. Human lines are `ts  text` or `ts  [a, b]  text` when labeled.
- **task notes** — `note add [--label LABEL]... --id ID <text...>` appends the same `{ts,text,labels?,status?}` jsonl into the **task directory**. `--id` required. Success stdout `added note`. Labels are bare tokens or `key=value` (first `=` splits; empty key rejected; empty value `key=` allowed). `note list [--json] [--limit N] [--label LABEL]... [--show-index] --id ID` lists (oldest first; AND filter; bare `--label key` matches `key` or `key=*`; `key=value` is exact; `--show-index` prefixes `N.  ` 1-based). Missing jsonl → `0 notes` / `[]`. Duplicate notes allowed (journal). `note edit [--label LABEL]... --id ID --index N [--status STATUS] [--append] <text...>` edits a note in place: `--index N` (1-based, required) selects within the `--label`-filtered set; text replaces by default, `--append` concatenates; labels and `ts` preserved, `--status` replaces status if given. Success stdout `edited note`. `RewriteNotes` atomically rewrites `notes.jsonl` (temp + rename). `show` prints `notes: N` (count). Grok session tracking convention: `--label grok-session-id=<uuid>` with description as text (legacy `grok`/`session-id` + id-in-text may still exist in old journals).
- **label list** — `label list` prints sorted, deduplicated **label names** from all task `labels` and all `notes.jsonl` (task + topic). For `key=value` note labels the **key** is listed once (`grok-session-id=…` → `grok-session-id`). Footer `N label(s)`. Empty store → `0 labels`.
- **progress** — `progress add --status STATUS --id ID <text...>` appends a `{ts,text,labels:["progress"],status}` line to the same task `notes.jsonl`; status is required and one of `in-progress`, `blocked`, `done`, `archived`. Legacy `started` entries remain readable. `progress list [--json] [--limit N] [--status STATUS] [--show-index] --id ID` lists progress-labeled entries (oldest first; `--status` filters; `--show-index` prefixes 1-based indices). `progress edit --status STATUS --id ID --index N [text...]` updates one 1-based progress entry (and replaces text only when supplied). `progress archive --id ID --index N` sets that entry to `archived`. `progress show --id ID` prints the latest entry or `no progress`. Human format: `ts  [progress]  (status)  text`. `done`/`archived` entries are gray + struck through when output is an interactive terminal (or `tree --color`); JSON and `tree --plain` contain no ANSI. `show` prints `progress: <status> (N entries)` when the latest has a status, else `progress: N entries`; omitted when no progress entries.
- **search** — `search [--task|--note|--progress|--topic|--all]… [--color|--no-color] [--json] <query>` case-insensitive substring search under `TSK_HOME`. Kind flags use less-flags `Group(CollectParsedFlags)`; no kind flags or any `--all` → all surfaces (`--all` wins over other kinds, no error); otherwise OR of kinds. `--note` excludes progress-labeled entries; `--progress` is only those. Surfaces: task titles, task notes, progress, topic `notes.jsonl`. Color via `dot-pkgs/terminal/color` (`ModeFromFlags` / `EnabledFor`): auto (TTY + `NO_COLOR`), `--color` always, `--no-color` never; conflict errors; human colors green anchor / blue kind / gray meta+footer and bolds the first EqualFold query span; `--json` never ANSI. Human: hit header + indented text + `N match(es)`; zero hits → `0 matches` exit 0. `--json` array of `{kind,task_id?,topic?,text,ts?,labels?,status?}`. Empty query → `Error: tsk search: query required`.
- **done** — `done <id>` marks a task done from any non-terminal stage; `--force` is accepted for compatibility (no extra effect). Updates `task.json` stage and stage history (directory basename unchanged). Refuses when already `done` or `archived`.
- **archive** — `archive <id>` marks a task archived from any non-terminal stage (shelve without claiming completion); `--force` accepted as no-op. Same storage update as done. Refuses when already terminal; no `done`↔`archived` moves.
- **delete** — `delete <id>` permanently removes the task directory and `index/<id>`; nested sub-tasks require `--recursive` (whole subtree + all descendant indexes). Stdout `deleted <id>`. No tombstone (task ids never reuse). Unlike `done`/`archive`, the task does not remain in the tree.
- **topic view** — `view [--json] <topic>` prints the topic tree (sub-topics + task nodes `[id]-<slug>` with color stage styling, or `[id]-<slug>  (stage)` when color is off). Skips `topic.json` / `notes.jsonl`. Empty topic: header + `(empty)`. Unicode `├──`/`└──`. `--json` nested `{path,aliases,tasks,subtopics}`; task nodes may include nested `tasks`; child `path` is the directory name.
- **add --note** — `create … [--note TEXT]… <title>` appends each note to the new task’s `notes.jsonl` after create (same format as `note add`); stdout remains the id only; empty `--note` errors.
- **skill** — `skill --show|--install|--list` embeds Shape-3 docs (`docs/SKILL.md` + `docs/<topic>/TOPIC.md`) via `skillcmd.SingleSkill`. Actions are flags (both orders). `--list` prints `tsk` then sorted topic paths (includes `working-on-task`). `--help` appends Available topics. Install via `skillcmd.HandleInstall` (flags in `--install --help` only — not in SKILL.md). Root help lists `skill`. Topic `working-on-task` is the daily playbook as numbered subsections: pick/approve a task (`tsk show` if id given, else `tsk project tree --all` preferring current project via `tsk project which`) → clarify (project / intent kind / E2E acceptance) → reuse live linked worktree `cwd` from `tsk show` when the user gave an id, else `wrk <project-location> --no-config -t` (optional `--bring`) → `kck` worker (`new` → messages/snapshot → `send` → `wait`) → main verify → user gate → `wrk` land/peel/`--done` → `tsk done`, with must-notes and stage updates throughout.
- **tree** — `tree [--json] [--id ID] [--color|--plain]` prints topics (primary) and projects (secondary). Inbox tasks with `project` group under project nodes; ungrouped inbox stay root leaves. Markers: topic `▣`/`#`, project `◆`/`@` (TTY vs plain). Leaves: padded `[id]` + title (512-rune cap; on-disk dir stays `[id]-<slug>`). Color on → ANSI stage styling only; color off → `  (stage)`. Auto color respects `NO_COLOR`. Footer `N tasks, M topics, P projects`. `--json`: `inbox`, `inbox_projects`, `topics` (task nodes include `title`); `--id` pruned path with markers; notes/progress under task.
- **topic alias on create/list** — `create --topic` and `list --topic` resolve aliases to the canonical slash path so `知识库` does not fork `topic_path`.
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
├── add/                       # tsk add
│   ├── no-topic/                 # inbox placement, index, task.json
│   ├── with-topic/               # topics/<path>/ placement
│   └── with-labels/              # --label flags, sorted labels
├── advance/                      # tsk advance
│   ├── basic/                    # create → advance updates stage only
│   └── invalid/
│       └── stage-jump/           # create → stage implementation errors
├── clarify/                      # tsk clarify *
│   └── confirm/                  # add questions, confirm -y → implementation
├── topic/                        # tsk topic *
│   ├── set-to-topic/             # inbox → topic path, dir move
│   ├── set-to-inbox/             # topic → inbox, topic_path null
│   ├── where/                    # topic where
│   │   ├── basic/                # mkdir → abs dir
│   │   ├── alias/                # 知识库 → knowledge-base dir
│   │   └── missing/              # not found error
│   ├── info/                     # topic info
│   │   ├── empty/                # mkdir-only defaults, no topic.json
│   │   ├── with-task/            # exact-path task count
│   │   ├── subtopic/             # child topic names
│   │   └── json/                 # --json object
│   ├── alias/                    # topic alias add
│   │   ├── add/                  # writes topic.json
│   │   └── conflict/             # second owner errors
│   ├── note/set/                 # topic note → info count
│   ├── note/append/              # second note keeps the first
│   ├── note/labeled/             # topic note --label
│   ├── notes/                    # topic notes list
│   │   ├── empty/
│   │   ├── json/
│   │   ├── limit/
│   │   └── migrate/              # legacy topic.json blob
│   ├── help/notes/               # topic notes --help
│   ├── create-via-alias/         # create --topic alias uses canonical path
│   ├── list-alias/               # list --topic alias
│   ├── help/where/               # topic where --help
│   ├── view/                     # topic view tree
│   │   ├── empty/
│   │   ├── tree/
│   │   ├── json/
│   │   ├── missing/
│   │   └── alias/
│   └── help/view/
├── next/                         # tsk next
│   └── oldest/                   # two in_process → older id on stdout
├── done/                         # tsk done (any non-terminal → done)
│   ├── from-create/              # create → done without --force
│   ├── from-summary/             # at summary → done, terminal stage
│   ├── force-from-create/        # --force still works (compat no-op)
│   ├── force-already-done/       # already done → error
│   ├── already-archived/         # archived → done refused
│   └── help/                     # documents any non-terminal + --force compat
├── archive/                      # tsk archive (any non-terminal → archived)
│   ├── from-create/              # create → archived
│   ├── already-archived/         # second archive refused
│   ├── already-done/             # done → archive refused
│   └── help/                     # documents any non-terminal + --force compat
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
├── note/                         # tsk note add / list (task journal)
│   ├── add/
│   │   ├── basic/                # unlabeled add + list
│   │   ├── labeled/              # --label grok --label session-id
│   │   ├── missing-id/
│   │   ├── missing-text/
│   │   └── missing-task/
│   ├── list/
│   │   ├── empty/
│   │   ├── json/
│   │   ├── filter-label/
│   │   ├── limit/
│   │   └── show-index/          # 1-based index prefix
│   ├── edit/
│   │   ├── basic/                # replace text in place
│   │   ├── append/               # --append adds to existing text
│   │   ├── labeled/              # --label filter + --index select
│   │   ├── status/               # --status on progress entry
│   │   ├── missing-index/
│   │   ├── index-out-of-range/
│   │   └── help/
│   ├── help/
│   │   ├── top/
│   │   ├── add/
│   │   └── list/
│   └── show-count/               # show notes: N
├── progress/                    # tsk progress lifecycle
│   ├── add/
│   │   ├── basic/                # --status in-progress required
│   │   ├── status/               # add --status then list
│   │   ├── missing-id/
│   │   ├── missing-status/
│   │   ├── invalid-status/
│   │   └── missing-text/
│   ├── list/
│   │   ├── empty/                # 0 entries
│   │   ├── basic/                # multiple entries with status
│   │   ├── status-filter/        # --status filters
│   │   ├── show-index/           # 1-based entry prefix
│   │   └── json/                 # JSON array
│   ├── edit/
│   │   └── basic/                # indexed status/text update
│   ├── archive/
│   │   └── basic/                # indexed archive shortcut
│   ├── show/
│   │   ├── latest/               # latest entry
│   │   └── empty/                # no progress
│   ├── help/
│   │   └── top/                  # progress -h
│   └── show-integration/         # show prints progress: line
├── tree/                        # tsk tree (topic + project grouping)
│   ├── empty/
│   ├── inbox/
│   ├── inbox-project/           # inbox grouped under @ project
│   ├── topic-project/           # topic → @ project → task
│   ├── full/
│   ├── json/                    # inbox + inbox_projects + topics
│   ├── id/
│   │   ├── full/
│   │   ├── inbox/
│   │   ├── json/
│   │   ├── color/
│   │   └── missing/
│   └── help/
├── project/                     # tsk project (tasks + registry)
│   ├── help/                    # project / add / tree / list help
│   ├── add/                     # project add
│   │   ├── basic/               # HTTPS origin → project.origin + cwd
│   │   ├── with-note/           # --note after add
│   │   ├── scp-origin/          # gitlab@host:path.git → origin key
│   │   ├── by-name/             # --project registered non-git → name
│   │   └── no-git/              # unregistered non-git → Error + register hint
│   ├── register/                # idempotent register; optional --name
│   │   ├── basic/               # first register + conflicting cwd error
│   │   ├── idempotent/          # up-to-date / fill empty location
│   │   └── no-name/             # auto basename / match by cwd
│   ├── which/basic/             # origin/name/cwd probe
│   ├── tree/                    # project tree (task forest)
│   │   ├── current/
│   │   ├── dir/                 # --dir PATH resolves project (not cwd)
│   │   ├── dir-missing/         # --dir bad path → Error
│   │   ├── exclude-done/        # default active-only; --done/--archived/--all stage filters
│   │   ├── all/
│   │   ├── empty/
│   │   └── sub-dirs/            # nested git scan; default --streaming; --no-streaming
│   ├── list/                    # project list (all/auto/registered)
│   │   ├── empty/
│   │   ├── after-add/           # auto row with tasks=
│   │   └── union-registered/    # register-only appears in default list
│   └── show-integration/        # show: cwd + project location|name|origin
├── list/                         # tsk list
│   └── filter/                   # --stage create filters ids
├── events/                       # events.jsonl audit
│   └── append/                   # any command appends one line
├── help/                         # --help / -h at every level
│   ├── root-empty/               # no args → top help
│   ├── root-flag/                # --help → top help
│   ├── root-h/                   # -h → top help
│   ├── add/                   # create --help → flags
│   ├── topic/                    # topic --help → set, mkdir
│   ├── label/                    # label list (deduped names); help → add, rm, list
│   └── clarify/                  # clarify --help → add, list, confirm
├── channel/                      # tsk channel *
│   ├── add/                   # channel create
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
│   │   ├── add/               # create --help documents --channel-id
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
| 1 | add/no-topic | `tsk add "add dark mode"` → `inbox/[1]-add-dark-mode/`, index, task.json |
| 2 | add/with-topic | `tsk add --topic eng/backend "x"` → dir under `topics/eng/backend/` |
| 3 | add/with-labels | `tsk add --label bug --label urgent "x"` → sorted labels in task.json |
| 4 | advance/basic | create + `tsk advance` → dir renamed to `…-in_process-…`, index updated |
| 5 | advance/invalid/stage-jump | create + `tsk stage <id> implementation` → error, dir unchanged |
| 6 | clarify/confirm | add 2 questions, `clarify confirm -y` → implementation, batch confirmed |
| 7 | topic/set-to-topic | inbox task → `topic set <path>` → dir moved, index updated |
| 8 | topic/set-to-inbox | topic task → `topic set --inbox` → inbox, `topic_path` null |
| 58 | topic/where/basic | mkdir `knowledge-base` → `where` prints abs `…/topics/knowledge-base` |
| 59 | topic/where/alias | alias `知识库` → `where 知识库` same dir as knowledge-base |
| 60 | topic/where/missing | `where no-such` → exit 1, `Error: topic not found: no-such` |
| 61 | topic/info/empty | mkdir-only → defaults, `tasks: 0`, no `topic.json` |
| 62 | topic/info/with-task | create under topic → `tasks: 1` |
| 63 | topic/info/subtopic | mkdir child `reports` → `subtopics: 1` |
| 64 | topic/info/json | `--json` object, no ANSI |
| 65 | topic/alias/add | writes `topic.json` aliases |
| 66 | topic/alias/conflict | second topic cannot claim the same alias |
| 67 | topic/note/set | `note` then `info` shows `notes: 1` + jsonl |
| 71 | topic/note/append | two notes listed, footer `2 notes` |
| 83 | topic/note/labeled | `topic note --label grok` lists `[grok]` |
| 72 | topic/notes/empty | no jsonl → `0 notes` |
| 73 | topic/notes/json | `--json` array with `ts`/`text` |
| 74 | topic/notes/limit | `--limit 1` last line only |
| 75 | topic/notes/migrate | legacy `topic.json` notes blob → jsonl |
| 76 | topic/help/notes | `topic notes --help` |
| 68 | topic/create-via-alias | `create --topic 知识库` stores `["knowledge-base"]` |
| 69 | topic/list-alias | `list --topic 知识库` prints the task id |
| 70 | topic/help/where | `topic where --help` |
| 77 | topic/view/empty | mkdir-only → `(empty)` |
| 78 | topic/view/tree | task + nested sub-topic task, unicode branches |
| 79 | topic/view/json | `--json` nested object |
| 80 | topic/view/missing | not found |
| 81 | topic/view/alias | `view 知识库` canonical header |
| 82 | topic/help/view | `topic view --help` |
| 84 | note/add/basic | unlabeled `note add` then `note list` |
| 85 | note/add/labeled | `--label grok --label session-id` |
| 86 | note/add/missing-id | `--id` required |
| 87 | note/add/missing-text | text required |
| 88 | note/add/missing-task | `Error: task not found: 99` |
| 89 | note/list/empty | `0 notes` |
| 90 | note/list/json | JSON array with labels |
| 91 | note/list/filter-label | AND filter |
| 92 | note/list/limit | last N |
| 93 | note/help/top | `note -h` add+list |
| 94 | note/help/add | `--id` `--label` |
| 95 | note/help/list | `--json` `--limit` |
| 96 | note/show-count | `show` `notes: 2` |
| 118 | note/edit/basic | replace text in place, `ts` preserved |
| 119 | note/edit/append | `--append` concatenates to existing text |
| 120 | note/edit/labeled | `--label` filter + `--index` select among filtered |
| 121 | note/edit/status | `--status done` on a progress entry |
| 122 | note/edit/missing-index | `--index` required |
| 123 | note/edit/index-out-of-range | `Error: index 5 out of range` |
| 124 | note/edit/help | `note edit -h` shows `--index` `--label` `--status` `--append` |
| 125 | note/list/show-index | `1.  ` and `2.  ` prefix |
| 97 | tree/empty | fresh store → `.` + `(empty)` + `0 tasks, 0 topics` |
| 98 | tree/inbox | inbox-only task at root, footer `1 task, 0 topics` |
| 99 | tree/full | inbox + topic with aliases + nested subtopic, footer `4 tasks, 1 topic` |
| 100 | tree/json | `--json` `{inbox:[...],topics:[...]}` no ANSI |
| 101 | tree/help | `tree -h` usage with `--json` and root level |
| 114 | tree/id/full | pruned branch with notes + progress under task |
| 115 | tree/id/inbox | inbox task, no notes/progress sections |
| 116 | tree/id/json | `--id --json` `{task,notes,progress}` |
| 117 | tree/id/missing | `Error: task not found: 99` |
| 126 | tree/id/color | `--color` styles done/archived leaf content only |
| 102 | progress/add/basic | required `--status in-progress` → `added progress` |
| 103 | progress/add/status | add `--status in-progress` then list shows `(in-progress)` |
| 104 | progress/add/missing-id | `--id` required |
| 105 | progress/add/missing-text | text required after valid status |
| 127 | progress/add/missing-status | `--status` required |
| 128 | progress/add/invalid-status | legacy `started` rejected for new entries |
| 106 | progress/list/empty | no entries → `0 entries` |
| 107 | progress/list/basic | three entries with status display, footer `3 entries` |
| 108 | progress/list/status-filter | `--status blocked` filters to one entry |
| 109 | progress/list/json | JSON array with `status` field |
| 129 | progress/list/show-index | 1-based progress entry prefixes |
| 130 | progress/edit/basic | indexed update replaces status and optional text |
| 131 | progress/archive/basic | indexed archive shortcut sets `archived` |
| 110 | progress/show/latest | latest entry printed, no footer |
| 111 | progress/show/empty | no progress → `no progress` |
| 112 | progress/help/top | `progress -h` lists add/list/edit/archive/show |
| 113 | progress/show-integration | `show` prints `progress: in-progress (1 entry)` |
| 9 | next/oldest | two `in_process` tasks → stdout = older id |
| 10 | done/from-summary | at summary → `tsk done` → stage done, dir unchanged |
| 10a | done/from-create | create → `tsk done` → stage done |
| 10b | done/already-archived | archived → `tsk done` → already archived |
| 10c | archive/from-create | create → `tsk archive` → stage archived |
| 10d | archive/already-done | done → `tsk archive` → already done |
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
| 45 | status/agent/dir | create `"add dark mode"` → agent facts `dir: <abs path>` after `topic:`; absolute; contains `inbox/[id]-add-dark-mode`; no `path`/`path_rel` |
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
| 15 | events/append | `tsk add` → `events.jsonl` gains one audit line |
| 16 | help/root-empty | `tsk` (no args) → exit 0; stdout has `Usage:` + command list; stderr empty |
| 17 | help/root-flag | `tsk --help` → exit 0; top help on stdout; stderr empty |
| 18 | help/root-h | `tsk -h` → exit 0; stdout contains `Usage:` |
| 19 | help/add | `tsk add --help` → create usage with `--label` and `--topic` |
| 20 | help/topic | `tsk topic --help` → lists `set`, `mkdir` subcommands |
| 21 | help/label | `tsk label --help` → lists `add`, `rm`, `list` subcommands |
| 22 | help/clarify | `tsk clarify --help` → lists `add`, `list`, `confirm` |
| 23 | ux/error-once | `tsk advance` (no id) → exit 1; `task id required` on stderr exactly once |
| 24 | ux/add-prints-id | `tsk add "hello"` → stdout `1\n`; inbox dir created; stderr empty |
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
| 97 | channel/help/add | `channel create --help` documents `--channel-id` |
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
doctest test ./tests/add
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
doctest test ./tests/add/no-topic
doctest test ./tests/advance/basic
doctest test ./tests/advance/invalid/stage-jump
doctest test ./tests/clarify/confirm
doctest test ./tests/topic/set-to-topic
doctest test ./tests/topic/where
doctest test ./tests/topic/info
doctest test ./tests/topic/alias
doctest test ./tests/topic/note
doctest test ./tests/topic/create-via-alias
doctest test ./tests/topic/list-alias
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
doctest test ./tests/help/add
doctest test ./tests/help/topic
doctest test ./tests/help/label
doctest test ./tests/help/clarify
doctest test ./tests/ux/error-once
doctest test ./tests/ux/add-prints-id
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

	"github.com/xhd2015/doctest/session"
)

type Request struct {
	WorkRoot string
	TskHome  string
	Args     []string
	TaskID   int
	OpenID     int // secondary id for an open (non-terminal) task in the same leaf
	ArchivedID int // secondary id when a leaf also creates an archived task
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

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	bin := getTskBin(t)
	cmd := exec.Command(bin, req.Args...)
	cmd.Dir = req.WorkRoot
	cmd.Env = tskEnv(req)
	return captureCommandOutput(cmd)
}
```