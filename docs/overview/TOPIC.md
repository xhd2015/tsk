---
name: tsk/overview
description: >-
  tsk storage layout under TSK_HOME: inbox, topics, index, bracketed task
  dirs, parent_id nesting, and env vars.
---

# overview

## Home layout

Default root: `~/.tsk` (override with `TSK_HOME`).

| Path | Role |
|------|------|
| `counter` | Monotonic task id allocator |
| `index/<id>` | One-line relative path to the task directory |
| `inbox/` | Tasks with `topic_path: null` |
| `topics/<seg>/…` | Topic tree; may hold `topic.json`, notes, subtopics, tasks |
| `events.jsonl` | Append-only CLI audit log |
| `channels/` | Channel spaces (see `channel` topic) |

## Task directory

Name: `[id]-<slug>/` under inbox, a topic, or a parent task dir.

Contains:

- `task.json` — id, title, slug, labels, topic_path, optional parent_id, optional cwd, optional project `{id,name}`, stage, timestamps, stage_history
- `context/` — followup markdown and similar artifacts
- `clarify/` — present during clarification (`batch.json`)
- nested `[child]-…/` — sub-tasks (any depth)

Optional `project` is set by `tsk project add`: `{origin}` when git remote exists,
or `{name}` for a registered non-git project (prefer origin). Optional `cwd` is
the absolute CLI recording directory. Registry: `projects.json`.

Reserved non-task children: `context`, `clarify` (and note/progress jsonl files on the task or topic).

## Nesting vs topics

- **Topic** classifies a root-level task (`topic_path` segments).
- **Parent task** nests a child on disk under the parent directory; child
  inherits the parent's `topic_path` (or null).
- Stage lives in `task.json` only; the directory basename does not change on
  stage transitions. Topic moves still relocate the dir and cascade `index/`
  for descendants.

## Env

- `TSK_HOME` — storage root
- `TSK_DATE=YYYY-MM-DD` — deterministic timestamps (tests)
- `TSK_USER` — channel identity fallback (see `channel`)
