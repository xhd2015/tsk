package tskcli

func topHelp() string {
	return `tsk — task workflow CLI

Usage:
  tsk <command> [arguments]

Commands:
  add        add a new task
  list       list task ids (optional filters)
  show       show task metadata
  status     show stage pipeline for a task
  advance    advance task to next stage
  stage      set task stage directly
  next       print oldest in_process task id
  label      add, remove, or list labels
  topic      set topic path or mkdir topic tree
  update     update task project and/or topic
  clarify    manage clarification questions
  followup   add followup context from summary
  done       mark task done from any non-terminal stage
  archive    mark task archived from any non-terminal stage
  delete     permanently remove a task (use --recursive for nested sub-tasks)
  channel    manage conversational channels
  note       add or list timestamped notes on a task
  progress   record and list progress entries on a task
  search     search task titles and note/progress/topic text
  tree       print all tasks organized by topic tree
  project    project-scoped tasks and registry (add/tree/list/register)
  install    install convenience CLI wrappers (e.g. pmark)
  skill      show/install embedded skill docs

Run tsk <command> --help for command-specific options.
Run tsk skill --help and tsk skill --install --help for skill docs.
`
}

func addHelp() string {
	return `Usage: tsk add [--label LABEL]... [--topic PATH | --parent ID] [--note TEXT]... <title>

Add a new task in inbox, under a topic path, or as a nested sub-task.
Optional --note flags append task notes after add (same store as note add).

Flags:
  --label LABEL   label to attach (repeatable)
  --topic PATH    topic path (e.g. eng/backend)
  --parent ID     add as nested sub-task under this task (any depth)
  --note TEXT     append a task note after add (repeatable)
  -h, --help      show this help
`
}

func listHelp() string {
	return `Usage: tsk list [--stage STAGE] [--label LABEL] [--topic PREFIX] [--project KEY] [--name NAME]

List task ids, optionally filtered.

Flags:
  --stage STAGE   filter by stage
  --label LABEL   filter by label
  --topic PREFIX  filter by topic path prefix
  --project KEY   filter by project.origin (or registered name)
  --name NAME     filter by project.name or origin basename / registry alias
  -h, --help      show this help
`
}

func projectHelp() string {
	return `Usage: tsk project <command> [arguments]

Subcommands:
  add          create a project-scoped task (non-blocking)
  tree         list project tasks as a tree (like tsk tree)
  list         list registered projects (projects.json)
  which        print resolved origin/name/cwd for cwd/--dir
  register     register a project (optional --name; idempotent re-register)
  unregister   remove a registered project name

  -h, --help  show this help
`
}

func projectAddHelp() string {
	return `Usage: tsk project add [--dir PATH] [--project NAME] [--note TEXT]... <title>

Create a task tagged with git origin when available, else a registered name.
Prints the new task id on stdout and exits (does not block).

Without --project: requires git remote.origin.url, or cwd matching a registered
project cwd. Otherwise errors (see tsk project register --help).

Flags:
  --dir PATH       resolve from PATH instead of cwd
  --project NAME   use a registered project name (override)
  --note TEXT      append a task note after create (repeatable)
  -h, --help       show this help
`
}

func projectTreeHelp() string {
	return `Usage: tsk project tree [--dir PATH | --name NAME | --project KEY] [--stage STAGE | --done] [--archived] [--all]
                  [--no-sub-dirs | --sub-dirs-depth N] [--streaming|--no-streaming]
                  [--color|--plain] [--json]

List project-scoped tasks as a tree (like tsk tree).
Default: current cwd's project plus projects from nested git repos under the
scan root (git toplevel or cwd/--dir, max depth 3), non-terminal stages only.
Human default mode streams: root project first, then each discovered project
as found. --no-streaming buffers then prints (root first, rest label-sorted).

Flags:
  --dir PATH           resolve project from PATH instead of cwd (conflicts with --name/--project/--all)
  --name NAME          filter by registry name / origin basename
  --project KEY        filter by origin key or registered name
  --stage STAGE        filter to one stage (conflicts with --done/--archived)
  --done               show only done tasks
  --archived           show only archived tasks (--done --archived shows both)
  --all                all projects and all stages (narrow with --done/--archived)
  --no-sub-dirs        do not scan nested git dirs for other projects
  --sub-dirs-depth N   max scan depth under scan root (default 3; must be >= 1)
  --streaming          stream output (default for human scan mode)
  --no-streaming       buffer the full tree, then print once
  --color              force ANSI stage styling (no trailing stage text)
  --plain              force no ANSI; show  (stage); ASCII markers
  --json               emit structured JSON instead of the tree
  -h, --help           show this help
`
}

func projectRegistryListHelp() string {
	return `Usage: tsk project list [--all|--auto|--registered] [--active] [--json]

List projects as an aligned table (NAME ORIGIN LOCATION TASKS). Default and --all:
union of projects-auto.json and projects.json (with live task counts).
When TASKS is shown, rows are sorted by count descending. Registered-only
names appear with tasks=0 until add.

Flags:
  --all           union of auto + registered (default)
  --auto          projects-auto.json only
  --registered    projects.json only
  --active        only projects with tasks > 0
  --json          emit JSON
  -h, --help      show this help
`
}

func projectRegisterHelp() string {
	return `Usage: tsk project register [--name NAME] [--cwd PATH] [--origin ORIGIN]

Register a project in projects.json (location tilde-form under $HOME).
Without --name: match by location (vs probe, then vs main location), else
basename(location) as name. Re-register is idempotent when fields match;
empty location/origin may be filled. Conflicting non-empty values error.
If --origin is omitted, origin is taken from git in --cwd/cwd when present.

Flags:
  --name NAME       project name (optional; default: basename of location)
  --cwd PATH        probe directory (default: process cwd; stored as location)
  --origin ORIGIN   optional origin key or git URL
  -h, --help        show this help
`
}

func projectUnregisterHelp() string {
	return `Usage: tsk project unregister <name>

Remove a registered project name from projects.json.

  -h, --help  show this help
`
}

func projectWhichHelp() string {
	return `Usage: tsk project which [--dir PATH]

Print resolved origin and/or registry name and cwd for the probe directory.

Flags:
  --dir PATH    resolve from PATH instead of cwd
  -h, --help    show this help
`
}

func showHelp() string {
	return `Usage: tsk show <id>

Show task metadata.

  -h, --help      show this help
`
}

func statusHelp() string {
	return `Usage: tsk status [--format=diagram|agent] [--color] [--plain] <id>

Show stage pipeline for a task.

Default format (when --format, --color, and --plain are all omitted):
  agent           if host agent detected (CODEX_THREAD_ID, PI_CODING_AGENT, or parent process)
  diagram         otherwise
  Override with TSK_STATUS_FORMAT=agent|diagram (debug/testing).

Formats:
  diagram         compact hand-made pipeline art; --color/--plain apply
  agent           2-row plain spine + facts; no ANSI, no boxes

Flags:
  --format FORMAT output format: diagram or agent (disables auto-detect)
  --color         force diagram + ANSI highlight (default on TTY for diagram; ignored for agent)
  --plain         force diagram; ASCII boxes, no ANSI
  -h, --help      show this help
`
}

func advanceHelp() string {
	return `Usage: tsk advance [--note NOTE] <id>

Advance task to the next allowed stage.

Flags:
  --note NOTE     optional note for stage history
  -h, --help      show this help
`
}

func stageHelp() string {
	return `Usage: tsk stage [--note NOTE] <id> <stage>

Set task stage directly (invalid transitions error).

Flags:
  --note NOTE     optional note for stage history
  -h, --help      show this help
`
}

func nextHelp() string {
	return `Usage: tsk next

Print id of oldest in_process task, or empty stdout when none.

  -h, --help      show this help
`
}

func topicHelp() string {
	return `Usage: tsk topic <command> [arguments]

Subcommands:
  set <id> <path|--inbox>   move task to topic path or inbox
  mkdir <path>              create topic directory tree
  rm <path>                 remove a topic directory (no tasks/subtopics)
  where <topic>             print absolute topic directory
  info <topic>              show topic metadata and counts
  note <topic> <text...>    append a timestamped topic note
  notes <topic>             list topic notes
  view <topic>              print topic / sub-topic / task tree
  alias add <topic> <alias> add an alias (creates topic.json)

  -h, --help                show this help
`
}

func topicRmHelp() string {
	return `Usage: tsk topic rm <path>

Remove a topic directory. Refuses if any task is classified at or under
this path, or if subtopic directories remain. No force flag.

  <path>      topic path or alias
  -h, --help  show this help
`
}

func updateHelp() string {
	return `Usage: tsk update <id> [--set-project REF | --clear-project]
                   [--set-topic PATH | --clear-topic]

Update an existing task's project and/or topic.

Flags:
  --set-project REF   set project (registered name, origin, or unique basename)
  --clear-project     remove project from the task
  --set-topic PATH    move task to topic path (same as topic set)
  --clear-topic       move task to inbox (clear topic_path)
  -h, --help          show this help
`
}

func topicWhereHelp() string {
	return `Usage: tsk topic where <topic>

Print the absolute directory of a topic.

  <topic>     topic path or alias
  -h, --help  show this help
`
}

func topicInfoHelp() string {
	return `Usage: tsk topic info [options] <topic>

Show topic details: path, title, aliases, dir, note count, task count, subtopics.

  <topic>     topic path or alias
  --json      machine-readable object (no ANSI)
  -h, --help  show this help
`
}

func topicNoteHelp() string {
	return `Usage: tsk topic note [options] <topic> <text...>

Append a timestamped note to the topic journal (notes.jsonl).
Does not replace earlier notes.

Flags:
  --label LABEL   label for this note (repeatable); bare name or key=value
  -h, --help      show this help
`
}

func topicNotesHelp() string {
	return `Usage: tsk topic notes [options] <topic>

List timestamped topic notes (oldest first).

  <topic>     topic path or alias
  --label LABEL   only notes that have this label (repeatable, AND);
                  bare name matches name or name=*; key=value is exact
  --json          JSON array (no ANSI)
  --limit N       last N notes only
  -h, --help      show this help
`
}

func noteHelp() string {
	return `Usage: tsk note <command> [arguments]

Subcommands:
  add [--label LABEL]... --id ID <text...>   append a timestamped note to a task
  list [--label LABEL]... --id ID            list notes (oldest first)
  edit [options] --id ID --index N <text...> edit an existing note in place

  -h, --help                                 show this help
`
}

func noteAddHelp() string {
	return `Usage: tsk note add [options] <text...>

Append a timestamped note to a task.

Flags:
  --id ID         task id (required)
  --label LABEL   label for this note (repeatable); bare name or key=value
  -h, --help      show this help
`
}

func noteListHelp() string {
	return `Usage: tsk note list [options]

List notes for a task (oldest first).

Flags:
  --id ID          task id (required)
  --label LABEL   only notes that have this label (repeatable, AND);
                  bare name matches name or name=*; key=value is exact
  --show-index     prefix each note with 1-based index
  --json          JSON array (no ANSI)
  --limit N       last N notes only
  -h, --help      show this help
`
}

func noteEditHelp() string {
	return `Usage: tsk note edit [options] --id ID --index N <text...>

Edit the text of an existing note in place. Labels and status are
preserved unless --status is given.

The note is identified by --index N (1-based) within the filtered
set. Use --label to narrow the candidate notes (AND filter), same
as ` + "`note list`" + `. Use ` + "`note list --show-index`" + ` to see indices.

Flags:
  --id ID          task id (required)
  --index N        1-based index into filtered notes (required)
  --label LABEL    filter by label (repeatable, AND); bare name matches
                   name or name=*; key=value is exact
  --status STATUS  replace status (optional; preserved if omitted)
  --append         append text to existing note text instead of replacing
  -h, --help       show this help
`
}

func treeHelp() string {
	return `Usage: tsk tree [options]

Print all tasks organized by topic (primary) and project (secondary).
Inbox tasks with a project are grouped under project nodes; ungrouped inbox
tasks stay at the root. Under a topic, project is a secondary level.

Kind markers: topic ▣ / # , project ◆ / @ (TTY vs --plain).
Task leaves: padded [id] + title (512-rune cap). On color, stage is ANSI only
(create plain; mid-pipeline tinted; done/archived gray+strike); without color,  (stage).

With --id, print a pruned branch from root to one task, including its
notes and progress entries nested under the task leaf. Done/archived progress
entries are gray and struck through when color is on.

Flags:
  --id ID     show only the branch for one task (with notes + progress)
  --color     force ANSI color (ignores TTY and NO_COLOR)
  --plain     force no ANSI; ASCII # / @ markers; show  (stage)
  --json      emit structured JSON instead of the tree
  -h, --help  show this help
`
}

func progressHelp() string {
	return `Usage: tsk progress <command> [arguments]

Subcommands:
  add --status STATUS --id ID <text...>       append a progress entry
  list [--status STATUS] [--show-index] --id ID
                                                list progress entries
  edit --status STATUS --id ID --index N [text...]
                                                update one progress entry
  archive --id ID --index N                    mark one entry archived
  show --id ID                                 show latest progress entry

Statuses: in-progress, blocked, done, archived

  -h, --help                                    show this help
`
}

func progressAddHelp() string {
	return `Usage: tsk progress add --status STATUS --id ID <text...>

Append a timestamped progress entry to a task journal.

Statuses: in-progress, blocked, done, archived

Flags:
  --id ID         task id (required)
  --status STATUS  progress status (required)
  -h, --help       show this help
`
}

func progressListHelp() string {
	return `Usage: tsk progress list [options] --id ID

List progress entries for a task (oldest first).

Flags:
  --id ID          task id (required)
  --status STATUS  only entries with this status
  --show-index     prefix each entry with 1-based index
  --json           JSON array (no ANSI)
  --limit N        last N entries only
  -h, --help       show this help
`
}

func progressEditHelp() string {
	return `Usage: tsk progress edit --status STATUS --id ID --index N [text...]

Update one progress entry. The index is 1-based among progress entries;
use ` + "`tsk progress list --show-index`" + ` to see indices. When text is
provided, it replaces the existing entry text.

Statuses: in-progress, blocked, done, archived

Flags:
  --id ID          task id (required)
  --index N        1-based progress entry index (required)
  --status STATUS  replacement status (required)
  -h, --help       show this help
`
}

func progressArchiveHelp() string {
	return `Usage: tsk progress archive --id ID --index N

Mark one progress entry archived. Equivalent to:
  tsk progress edit --id ID --index N --status archived

Flags:
  --id ID      task id (required)
  --index N    1-based progress entry index (required)
  -h, --help   show this help
`
}

func progressShowHelp() string {
	return `Usage: tsk progress show [options] --id ID

Print the latest progress entry for a task.

Flags:
  --id ID      task id (required)
  -h, --help   show this help
`
}

func searchHelp() string {
	return `Usage: tsk search [options] <query>

Search task titles and note/progress/topic note text under TSK_HOME.
Case-insensitive substring match. With no kind flags, searches all surfaces
(same as --all). When --all appears with other kind flags, --all wins.

Flags:
  --task       search task titles
  --note       search task notes (excludes progress entries)
  --progress   search progress entries
  --topic      search topic notes
  --all        search all surfaces (default when none given)
  --color      force ANSI color on (even when stdout is not a TTY)
  --no-color   force ANSI color off
  --json       JSON array (no ANSI)
  -h, --help   show this help
`
}

func topicViewHelp() string {
	return `Usage: tsk topic view [options] <topic>

Print the topic tree: sub-topics and tasks. Skips topic.json and notes.jsonl.

  <topic>     topic path or alias
  --json      nested object (no ANSI)
  -h, --help  show this help
`
}

func topicAliasHelp() string {
	return `Usage: tsk topic alias <command> [arguments]

Subcommands:
  add <topic> <alias>   bind an alias to an existing topic

  -h, --help            show this help
`
}

func topicAliasAddHelp() string {
	return `Usage: tsk topic alias add <topic> <alias>

Add an alias that resolves to the topic's canonical path.

  -h, --help  show this help
`
}

func labelHelp() string {
	return `Usage: tsk label <command> [arguments]

Subcommands:
  add <id> <label>   add label to task
  rm <id> <label>    remove label from task
  list               list deduped label names in use (tasks + notes)

  -h, --help         show this help
`
}

func labelListHelp() string {
	return `Usage: tsk label list

List deduplicated label names from all tasks and notes (including topic
notes). For key=value note labels, the key is listed once.

  -h, --help  show this help
`
}

func clarifyHelp() string {
	return `Usage: tsk clarify <command> [arguments]

Subcommands:
  add <id> <question...>       add clarification question
  list <id>                    list clarification items
  confirm [-y] <id>            confirm all items and advance

  -h, --help                   show this help
`
}

func followupHelp() string {
	return `Usage: tsk followup <id> <message...>

Add followup context from summary stage.

  -h, --help      show this help
`
}

func doneHelp() string {
	return `Usage: tsk done [--force] <id>

Mark task done from any non-terminal stage.
--force is accepted for compatibility and has no extra effect.

Flags:
  --force         accepted for compatibility (no extra effect)
  -h, --help      show this help
`
}

func archiveHelp() string {
	return `Usage: tsk archive [--force] <id>

Mark task archived from any non-terminal stage.
--force is accepted for compatibility and has no extra effect.

Flags:
  --force         accepted for compatibility (no extra effect)
  -h, --help      show this help
`
}

func deleteHelp() string {
	return `Usage: tsk delete [--recursive] <id>

Permanently remove a task directory and its index entry.
Nested sub-tasks require --recursive (deletes the whole subtree).

Unlike tsk done / tsk archive, the task does not remain in the tree.

Flags:
  --recursive   also delete nested sub-tasks under this task
  -h, --help    show this help
`
}
