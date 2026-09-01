---
name: tsk/add
description: >-
  Create inbox, topic-scoped, or nested sub-tasks with tsk add
  (--topic or --parent), optionally seeding notes with --note.
---

# add

```text
tsk add [--label LABEL]... [--topic PATH | --parent ID] [--note TEXT]... <title>
```

Success prints the new task id on stdout. Each `--note` appends a task note
(same `notes.jsonl` as `tsk note add`) after the task is created.

## Placement

| Flags | Result |
|-------|--------|
| (none) | `inbox/[id]-create-<slug>/`, `topic_path: null` |
| `--topic PATH` | Under `topics/<path>/…`; resolves aliases to canonical path |
| `--parent ID` | Nested under that task’s directory; inherits parent `topic_path`; sets `parent_id` |

`--topic` and `--parent` together → error (child inherits parent location).

## Sub-tasks

When the user wants a sub-task under task N:

```bash
tsk add --parent N "title"
```

Arbitrary depth: `--parent` may point at another nested task.

Do **not** create a sibling under the same topic and treat it as a sub-task.

## Examples

```bash
tsk add "inbox item"
tsk add --topic knowledge-base "report progress"
tsk add --topic 知识库 "via alias"          # stores canonical topic_path
tsk add --parent 6 "add more doctests to pricing"
tsk add --parent 9 "cover discount edge cases"
tsk add --label report --topic knowledge-base "weekly update"
tsk add --topic agent-pro \
  --note "grok session … track stall" \
  "flaky issue: …"
```

## After add

- Prefer `--note` on `add` when context is known up front (dirs, grok
  session ids). Use `tsk note add --id <id> …` only for later notes.
- `tsk show <id>` — metadata (includes `parent:` when nested; `notes: N`)
- `tsk tree` / `tsk tree --id <id>` — location and note/progress leaves
