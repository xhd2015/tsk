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
                 ↘ archived   (alternate terminal; off-spine)
```

`done` and `archived` are terminal. Stage is stored in `task.json` only; the directory basename `[id]-<slug>` does not change on stage transitions.

## Commands

| Command | Role |
|---------|------|
| `tsk advance <id>` | Next allowed edge |
| `tsk stage <id> <stage>` | Set stage directly (invalid jumps error) |
| `tsk clarify …` | Questions; `confirm -y` → implementation |
| `tsk followup <id> <msg>` | From summary → `user_followup` + `context/followup-*.md` |
| `tsk done <id>` | From any non-terminal stage → `done` (`--force` accepted, no extra effect) |
| `tsk archive <id>` | From any non-terminal stage → `archived` (`--force` accepted, no extra effect) |
| `tsk delete <id>` | Permanently remove task dir + index; `--recursive` for nested sub-tasks |
| `tsk status <id>` | Pipeline art (`diagram` / `agent` formats) |
| `tsk next` | Oldest `in_process` id (or empty) |

`done` / `archived` keep the task in the tree (struck when colored). `delete` removes it; ids are never reused. Terminal ↔ terminal moves are refused.

## Agent guidance

- Prefer `advance` along the pipeline for real work; use `stage` sparingly.
- `status --format=agent` (or auto under agent hosts) for compact facts including `dir:`.
- Record decisions in notes/progress rather than only changing stage.
- Use `done` to finish or `archive` to shelve from any open stage; both are terminal.
- Use `delete` for mistakes/duplicates; refuse without `--recursive` when children exist.
- When the user asks to work on / finish a task end-to-end, follow topic
  `working-on-task` (tsk + wrk + kck playbook), not only these tips.
