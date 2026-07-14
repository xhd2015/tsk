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
  done       mark task done from summary or user_followup
  channel    manage conversational channels

Run tsk <command> --help for command-specific options.
`
}

func createHelp() string {
	return `Usage: tsk create [--label LABEL]... [--topic PATH] <title>

Create a new task in inbox or under a topic path.

Flags:
  --label LABEL   label to attach (repeatable)
  --topic PATH    topic path (e.g. eng/backend)
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

  -h, --help                show this help
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
	return `Usage: tsk done <id>

Mark task done from summary or user_followup stage.

  -h, --help      show this help
`
}