---
name: tsk/tree
description: >-
  Print the topic/task forest with tsk tree; prune to one task with --id;
  JSON and color/plain modes.
---

# tree

```text
tsk tree [--json] [--color|--plain]
tsk tree --id ID [--json] [--color|--plain]
```

## Full tree

- Root `.` then inbox task dirs and top-level topics (sorted).
- Topics may nest; task dirs may nest sub-tasks.
- Leaf label: `[id]-<stage>-<slug>  task <id>  <stage>`
- Done tasks (and done/archived progress under `--id`) use gray + strikethrough on a TTY (`--color` forces; `--plain` disables).

## `--id`

Pruned branch from root to that task. Under the task leaf:

- `notes` — non-progress notes
- `progress` — notes labeled progress

Does not list sibling tasks outside the path to `--id`.

## `--json`

Machine-readable forest (`inbox` + `topics` with nested `tasks` / `subtopics`).
No ANSI. Task nodes may include nested `tasks` for sub-tasks.

## When to use

| Goal | Command |
|------|---------|
| Orient / whole board | `tsk tree` |
| One task + its notes | `tsk tree --id N` |
| Scripts | `tsk tree --json` |
