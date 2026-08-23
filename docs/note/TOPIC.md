---
name: tsk/note
description: >-
  Timestamped notes and progress entries on a task (session ids, dirs, status).
---

# note

Notes and progress share `notes.jsonl` on the task directory. Progress entries
are notes that include the `progress` label.

When creating a **new** task with initial context, prefer
`tsk create --topic … --note "…" "title"` (see topic `create`). Use the
commands below on an **existing** task id.

## Notes

```text
tsk note add --id ID [--label L]… <text…>
tsk note list --id ID [--label L]…
tsk note edit --id ID --index N [--append] <text…>
```

Use notes for durable pointers agents and humans need later, for example:

```bash
tsk note add --id 9 "dir W0/knowledge-workspace/credit-pricing-center-doctest-backfill"
tsk note add --id 9 "grok session 01a02283-…: backfill doctests for critical biz paths"
```

`tree --id N` shows notes and progress under the task leaf.

## Progress

```text
tsk progress add --id ID [--status STATUS] <text…>
tsk progress list --id ID
```

Typical statuses: `in-progress`, `done`, `archived`. Done/archived progress is
styled like done tasks in colored `tree --id` output.

## show

`tsk show <id>` prints `notes: <count>` and a short progress summary when
progress entries exist.
