---
name: tsk/workflow
description: >-
  Stage pipeline: advance, stage, clarify, followup, done, status, and next.
---

# workflow

## Stages

```text
create → in_process → clarification → implementation → verification
       → summary → user_followup ⇄ clarification
                 ↘ done
```

`done` is terminal. Directory basename updates with stage (`[id]-create-…` → `[id]-in_process-…`); nested children stay under the renamed parent (indexes cascade).

## Commands

| Command | Role |
|---------|------|
| `tsk advance <id>` | Next allowed edge |
| `tsk stage <id> <stage>` | Set stage directly (invalid jumps error) |
| `tsk clarify …` | Questions; `confirm -y` → implementation |
| `tsk followup <id> <msg>` | From summary → `user_followup` + `context/followup-*.md` |
| `tsk done <id>` | From summary or user_followup; `--force` bypasses stage gate |
| `tsk status <id>` | Pipeline art (`diagram` / `agent` formats) |
| `tsk next` | Oldest `in_process` id (or empty) |

## Agent guidance

- Prefer `advance` along the pipeline for real work; use `stage` sparingly.
- `status --format=agent` (or auto under agent hosts) for compact facts including `dir:`.
- Record decisions in notes/progress rather than only changing stage.
