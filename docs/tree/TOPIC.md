---
name: tsk/tree
description: >-
  Print the topic/task forest with tsk tree; project grouping; prune with --id;
  JSON and color/plain modes.
---

# tree

```text
tsk tree [--json] [--color|--plain]
tsk tree --id ID [--json] [--color|--plain]
```

## Full tree

- Root `.` then topics, project groups, and ungrouped inbox tasks (sorted).
- **Topic** is primary; **project** is secondary (under a topic, or at root for inbox).
- Kind markers: topic `▣` / `#`, project `◆` / `@` (TTY vs `--plain`); tasks use `[id]-…` only.
- Leaf label: `[id]-<stage>-<slug>  task <id>  <stage>`
- Nested sub-tasks stay under their parent leaf.
- Footer: `N tasks, M topics, P projects`.
- Done tasks (and done/archived progress under `--id`) use gray + strikethrough on a TTY (`--color` forces; `--plain` disables).

## `--id`

Pruned branch from root to that task (topic markers + optional project node).
Under the task leaf:

- `notes` — non-progress notes
- `progress` — notes labeled progress

## `--json`

`inbox` (ungrouped), `inbox_projects` `[{origin,name,label,tasks}]`, `topics`
(each topic may include `projects` buckets plus ungrouped `tasks` / `subtopics`).
No ANSI / no glyphs.

## When to use

| Goal | Command |
|------|---------|
| Orient / whole board | `tsk tree` |
| One task + its notes | `tsk tree --id N` |
| Scripts | `tsk tree --json` |
