---
name: tsk
description: >-
  Task workflow CLI: add/list tasks under topics or nested parents,
  stage pipeline, notes/progress, tree view, channels. Use when the user
  manages work with tsk, asks for task trees, sub-tasks, or /tsk.
  Load topics: tsk skill --show <topic>
---

# tsk — multi-topic task workflow skill

This skill is an **index**. Load detailed guidance with
`tsk skill --show <topic>` (or `tsk skill <topic> --show`).

**tsk** stores tasks under `TSK_HOME` (default `~/.tsk`). Topics classify work;
tasks are directories that can nest as sub-tasks. Domain commands remain
`tsk add`, `tsk tree`, `tsk note`, and so on — this skill documents the
product model and agent workflows.

## Core model (agents)

- Task directory name: `[id]-<stage>-<slug>` (brackets are literal on disk).
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
4. Advance stages when doing real work; use `tsk done` (or `--force`) only
   when finishing. Use `tsk delete` to permanently remove a mistaken or
   duplicate task (`--recursive` if it has nested sub-tasks).
5. Do not `tsk topic set` on a nested child — reparent first.

## Topics

- `overview` — storage layout, dirname rules, env vars
- `add` — inbox, topic, `--parent`, and one-shot `--note`
- `tree` — full tree, `--id`, JSON/color/plain
- `project` — project-scoped tasks + registry (`add`/`tree`/`list`/`register`)
- `install` — convenience wrappers (`tsk install pmark` → `~/.local/bin/pmark`)
- `topic` — mkdir/set/view/alias; nested-task restrictions
- `workflow` — stages, advance, clarify, followup, done, delete, status, next
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
tsk skill --show note

# YAML frontmatter only
tsk skill --show --header
tsk skill --show add --header
```

## Related CLI

Domain usage lives in `tsk <command> --help`. Use this skill for when/why and
nested-topic recipes, not as a full flag encyclopedia.
