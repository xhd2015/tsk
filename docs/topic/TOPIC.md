---
name: tsk/topic
description: >-
  Topic tree management: mkdir, set, view, alias, topic notes; rules for
  nested tasks.
---

# topic

```text
tsk topic mkdir <path>
tsk topic set <id> <path|--inbox>
tsk topic where|info|view|notes <topic>
tsk topic note <topic> <text…>
tsk topic alias add <topic> <alias>
```

## Classification

- `topic set <id> <path>` moves a **root-level** task into a topic (creates path dirs as needed via prior `mkdir` or as part of create).
- `topic set <id> --inbox` clears classification (`topic_path: null`).
- Moving a parent task moves its nested children and updates their `topic_path`.

## Nested tasks

`topic set` on a task with `parent_id` is **rejected**. Reparent / move to root first; do not silently escape a parent.

## Aliases

`topic alias add knowledge-base 知识库` — `create --topic` / `list --topic` resolve aliases to the canonical slash path so agents do not fork duplicate topic trees.

## View

`tsk topic view <topic>` prints that topic’s subtopics and tasks (including nested sub-tasks), similar to a subtree of `tsk tree`.
