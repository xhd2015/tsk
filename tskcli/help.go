package tskcli

func topHelp() string {
	return `tsk — task workflow CLI

Usage:
  tsk <command> [arguments]

Commands:
  create     create a new task
  list       list task ids (optional filters)
  show       show task metadata
  status     show stage pipeline for a task
  advance    advance task to next stage
  stage      set task stage directly
  next       print oldest in_process task id
  label      add or remove labels
  topic      set topic path or mkdir topic tree
  clarify    manage clarification questions
  followup   add followup context from summary
  done       mark task done (use --force to bypass the workflow stage requirement)
  delete     permanently remove a task (use --recursive for nested sub-tasks)
  channel    manage conversational channels
  note       add or list timestamped notes on a task
  progress   record and list progress entries on a task
  search     search task titles and note/progress/topic text
  tree       print all tasks organized by topic tree
  skill      show/install embedded skill docs

Run tsk <command> --help for command-specific options.
Run tsk skill --help and tsk skill --install --help for skill docs.
`
}

func createHelp() string {
	return `Usage: tsk create [--label LABEL]... [--topic PATH | --parent ID] [--note TEXT]... <title>

Create a new task in inbox, under a topic path, or as a nested sub-task.
Optional --note flags append task notes after create (same store as note add).

Flags:
  --label LABEL   label to attach (repeatable)
  --topic PATH    topic path (e.g. eng/backend)
  --parent ID     create as nested sub-task under this task (any depth)
  --note TEXT     append a task note after create (repeatable)
  -h, --help      show this help
`
}

func listHelp() string {
	return `Usage: tsk list [--stage STAGE] [--label LABEL] [--topic PREFIX]

List task ids, optionally filtered.

Flags:
  --stage STAGE   filter by stage
  --label LABEL   filter by label
  --topic PREFIX  filter by topic path prefix
  -h, --help      show this help
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
  where <topic>             print absolute topic directory
  info <topic>              show topic metadata and counts
  note <topic> <text...>    append a timestamped topic note
  notes <topic>             list topic notes
  view <topic>              print topic / sub-topic / task tree
  alias add <topic> <alias> add an alias (creates topic.json)

  -h, --help                show this help
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
  --label LABEL   label for this note (repeatable)
  -h, --help      show this help
`
}

func topicNotesHelp() string {
	return `Usage: tsk topic notes [options] <topic>

List timestamped topic notes (oldest first).

  <topic>     topic path or alias
  --label LABEL   only notes that have this label (repeatable, AND)
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
  --label LABEL   label for this note (repeatable)
  -h, --help      show this help
`
}

func noteListHelp() string {
	return `Usage: tsk note list [options]

List notes for a task (oldest first).

Flags:
  --id ID          task id (required)
  --label LABEL   only notes that have this label (repeatable, AND)
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
  --label LABEL    filter by label (repeatable, AND)
  --status STATUS  replace status (optional; preserved if omitted)
  --append         append text to existing note text instead of replacing
  -h, --help       show this help
`
}

func treeHelp() string {
	return `Usage: tsk tree [options]

Print all tasks organized by topic tree, like the tree CLI.
Inbox tasks (no topic) appear at the root level alongside top-level topics.

With --id, print a pruned branch from root to one task, including its
notes and progress entries nested under the task leaf. Done task leaves and
progress entries with status done or archived are gray and struck through on a terminal.

Flags:
  --id ID     show only the branch for one task (with notes + progress)
  --color     force ANSI color and strikethrough
  --plain     force no ANSI color or strikethrough
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

  -h, --help         show this help
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

Mark task done from summary or user_followup stage.
Use --force to complete it from any non-terminal stage.

Flags:
  --force         bypass the normal workflow stage requirement
  -h, --help      show this help
`
}

func deleteHelp() string {
	return `Usage: tsk delete [--recursive] <id>

Permanently remove a task directory and its index entry.
Nested sub-tasks require --recursive (deletes the whole subtree).

Unlike tsk done, the task does not remain in the tree.

Flags:
  --recursive   also delete nested sub-tasks under this task
  -h, --help    show this help
`
}
