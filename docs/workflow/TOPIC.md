---
name: tsk/workflow
description: >-
  Stage pipeline: advance, stage, clarify, followup, done, delete, status, and next.
---

# workflow

## Stages

```text
create → in_process → clarification → implementation → verification
       → summary → user_followup ⇄ clarification
                 ↘ done
```

`done` is terminal. Stage is stored in `task.json` only; the directory basename `[id]-<slug>` does not change on stage transitions.

## Commands

| Command | Role |
|---------|------|
| `tsk advance <id>` | Next allowed edge |
| `tsk stage <id> <stage>` | Set stage directly (invalid jumps error) |
| `tsk clarify …` | Questions; `confirm -y` → implementation |
| `tsk followup <id> <msg>` | From summary → `user_followup` + `context/followup-*.md` |
| `tsk done <id>` | From summary or user_followup; `--force` bypasses stage gate |
| `tsk delete <id>` | Permanently remove task dir + index; `--recursive` for nested sub-tasks |
| `tsk status <id>` | Pipeline art (`diagram` / `agent` formats) |
| `tsk next` | Oldest `in_process` id (or empty) |

`done` keeps the task in the tree (struck when done). `delete` removes it; ids are never reused.

## Agent guidance

- Prefer `advance` along the pipeline for real work; use `stage` sparingly.
- `status --format=agent` (or auto under agent hosts) for compact facts including `dir:`.
- Record decisions in notes/progress rather than only changing stage.
- Use `delete` for mistakes/duplicates; refuse without `--recursive` when children exist.
