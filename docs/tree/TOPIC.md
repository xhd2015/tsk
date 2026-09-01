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
- Kind markers: topic `▣` / `#`, project `◆` / `@` (TTY vs `--plain`). With color, project origin (`github.com/…`) is grey; the short name stays plain.
- Leaf label: padded `[id]` + title (title capped at 512 runes). Color on → stage via ANSI only (`create` plain, mid-pipeline tinted, `done` gray+strike); color off → trailing `  (stage)`. On-disk dir stays `[id]-<slug>`.
- Nested sub-tasks stay under their parent leaf.
- Footer: `N tasks, M topics, P projects`.
- Color auto on TTY unless `NO_COLOR` is set; `--color` forces; `--plain` disables color and fancy markers.

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
