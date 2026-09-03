---
name: tsk
description: >-
  Task workflow CLI: add/list tasks under topics or nested parents,
  stage pipeline, notes/progress, tree view, channels. Use when the user
  manages work with tsk, asks for task trees, sub-tasks, or /tsk — or asks
  to work on / finish a task via tsk + wrk + kck (see working-on-task).
  Load topics: tsk skill --show <topic>
---

# tsk — multi-topic task workflow skill

This skill is an **index**. Load detailed guidance with
`tsk skill --show <topic>` (or `tsk skill <topic> --show`).

**tsk** stores tasks under `TSK_HOME` (default `~/.tsk`). Topics classify work;
tasks are directories that can nest as sub-tasks. Domain commands remain
`tsk add`, `tsk tree`, `tsk note`, and so on — this skill documents the
product model and agent workflows (manual use and autonomous working-on-task).

## Core model (agents)

- Task directory name: `[id]-<slug>` (brackets are literal on disk; stage is in `task.json`).
- Placement: inbox, under a topic path, or nested under a parent task via
  `--parent`.
- Sub-tasks: use `tsk add --parent <id>`; do not create a sibling under the
  same topic and call it a sub-task.
- `topic_path` is topic segments only (or null in inbox); nesting is physical
  under the parent task dir (`parent_id` in `task.json`).
- Orient with `tsk tree` / `tsk tree --id N`; attach context with notes and
  progress (prefer `add --note` when creating; `note add` / `progress` later).

## Default agent workflow

1. `tsk tree` (or `tsk tree --id N`) to see where work lives.
2. Create with inbox / `--topic` / `--parent` as the user intends.
3. When session ids, dirs, or other pointers are known at create time, pass
   them as `tsk add … --note "…" "title"` (repeatable). Use
   `tsk note add --id N …` only for notes after the task already exists.
4. Advance stages when doing real work; use `tsk done` when finishing or
   `tsk archive` to shelve without claiming completion. Use `tsk delete` to
   permanently remove a mistaken or duplicate task (`--recursive` if it has
   nested sub-tasks).
5. Do not `tsk topic set` / `tsk update --set-topic` on a nested child — reparent first.
6. To attach an existing task to a project or move/clear its topic: `tsk update <id> …`.
   Remove an empty topic with `tsk topic rm <path>` (refuses if any task remains).

To run a task to finished without hand-holding each stage, see
**working-on-task** below (`tsk skill --show working-on-task`).

## working-on-task

When the user asks to **work on** / finish a task: **pick/approve** a task
(`tsk show` or `tsk project tree --all`, prefer current project); clarify
**project**, **intent kind**, and **E2E acceptance**; reuse a live linked
worktree from `tsk show` or **`wrk`**-create; drive a **`kck`** worker; keep
**`tsk`** stages and must-notes updated; verify; on user approve land with
`wrk` then `tsk done`. Full playbook: `tsk skill --show working-on-task`.

## Topics

- `overview` — storage layout, dirname rules, env vars
- `add` — inbox, topic, `--parent`, and one-shot `--note`
- `tree` — full tree, `--id`, JSON/color/plain
- `project` — project-scoped tasks + registry (`add`/`tree`/`list`/`register`)
- `install` — convenience wrappers (`tsk install pmark` → `~/.local/bin/pmark`)
- `topic` — mkdir/set/rm/view/alias; `tsk update` for project/topic fields
- `workflow` — stages, advance, clarify, followup, done, delete, status, next
- `working-on-task` — tsk + wrk + kck daily playbook (clarify → worktree → worker → verify → ship)
- `note` — notes and progress on existing tasks (see also `add --note`)
- `channel` — Slack-like channels

## Retrieve topics

```bash
# list skill name + every nested topic path
tsk skill --list

# root skill index (this document)
tsk skill --show

# topic (both flag orders)
tsk skill --show add
tsk skill add --show
tsk skill --show workflow
tsk skill --show working-on-task
tsk skill --show note

# YAML frontmatter only
tsk skill --show --header
tsk skill --show add --header
```

## Related CLI

Domain usage lives in `tsk <command> --help`. Use this skill for when/why and
nested-topic recipes, not as a full flag encyclopedia.
