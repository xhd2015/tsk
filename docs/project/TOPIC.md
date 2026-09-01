---
name: tsk/project
description: >-
  Project-scoped tasks via tsk project: origin-preferred identity, optional
  name registry for non-git dirs, tree view, same ~/.tsk store as other todos.
---

# project

Non-blocking replacement for using `mark` tabs as project remarks. Tasks live
in the same `TSK_HOME` store; `mark` stays separate (no import). Shortcut:
`tsk install pmark` installs `~/.local/bin/pmark` → `tsk project add` (see
`install` topic).

```text
tsk project add [--dir PATH] [--project NAME] [--note TEXT]... <title>
tsk project tree [--name NAME | --project KEY] [--stage STAGE] [--all]
                 [--color|--plain] [--json]
tsk project list [--all|--auto|--registered] [--active] [--json]
tsk project which [--dir PATH]
tsk project register --name NAME [--cwd PATH] [--origin ORIGIN]
tsk project unregister <name>
```

## Identity (prefer origin)

| Situation | Stored on task |
|-----------|----------------|
| Has git `remote.origin.url` | `project: { "origin": "host/path" }` |
| No origin, registered name | `project: { "name": "…" }` |
| Neither | Error — see `tsk project register --help` |

Even if a friendly name is registered for an origin, the **task** still stores
`origin` only (more stable). Names live in `projects.json` for display/lookup.

## Registry — `TSK_HOME/projects.json`

Explicit register only. Origin-only projects need no row.

```json
{
  "projects": [
    { "origin": "github.com/xhd2015/dot-pkgs", "name": "dot-pkgs", "cwd": "~/Projects/xhd2015/dot-pkgs" },
    { "name": "seatalk-local-bot", "cwd": "~/seatalk-local-bot" }
  ]
}
```

- `name` unique when set
- `cwd` stored with `~` via pathfmt when under `$HOME`

## Auto ledger — `TSK_HOME/projects-auto.json`

Written on every successful `project add` (upsert-only). Identity is `origin` XOR
`name`. `cwd` is the **main repo** root (linked worktrees still record main),
tilde-form, set once; later adds only bump `last_seen_at` (local TZ offset).

`tsk project list` (default / `--all`): union of auto + registered as an aligned
table (`NAME` `ORIGIN` `CWD` `TASKS`; registered-only rows show `0` until add).
When `TASKS` is shown, rows sort by count descending (tie-break name, origin).
`--auto` / `--registered` select one source; `--active` filters `tasks>0`.
`--registered` omits the `TASKS` column unless `--active` (then name/origin order).

## Placement / lifecycle

Default inbox. Use `tsk done` / `delete` / `show` / `note`.

To attach or clear project on an **existing** task (without creating a new one):

```text
tsk update <id> --set-project REF    # name, origin, or unique basename
tsk update <id> --clear-project
```


## Tree

Like `tsk tree`. Default: current project, exclude `done`.
