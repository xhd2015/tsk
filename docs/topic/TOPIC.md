---
name: tsk/topic
description: >-
  Topic tree management: mkdir, set, rm, view, alias, topic notes; rules for
  nested tasks. Task field updates via tsk update.
---

# topic

```text
tsk topic mkdir <path>
tsk topic set <id> <path|--inbox>
tsk topic rm <path>
tsk topic where|info|view|notes <topic>
tsk topic note <topic> <text…>
tsk topic alias add <topic> <alias>
tsk update <id> [--set-project REF|--clear-project] [--set-topic PATH|--clear-topic]
```

## Classification

- `topic set <id> <path>` moves a **root-level** task into a topic (creates path dirs as needed via prior `mkdir` or as part of create).
- `topic set <id> --inbox` clears classification (`topic_path: null`).
- `tsk update <id> --set-topic PATH` / `--clear-topic` does the same moves; also supports `--set-project` / `--clear-project` (project is metadata only).
- Moving a parent task moves its nested children and updates their `topic_path`.

## Remove topic

- `topic rm <path>` deletes the topic directory when it has **no tasks** at or under that path and **no subtopics**.
- No force flag. Leftover `topic.json` / `notes.jsonl` are removed with the directory when otherwise empty.

## Nested tasks

`topic set` / `update --set-topic` on a task with `parent_id` is **rejected**. Reparent / move to root first; do not silently escape a parent.

## Aliases

`topic alias add knowledge-base 知识库` — `create --topic` / `list --topic` resolve aliases to the canonical slash path so agents do not fork duplicate topic trees.

## View

`tsk topic view <topic>` prints that topic’s subtopics and tasks (including nested sub-tasks), similar to a subtree of `tsk tree`.
