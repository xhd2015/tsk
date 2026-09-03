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
tsk project tree [--dir PATH | --name NAME | --project KEY] [--stage STAGE] [--all]
                 [--no-sub-dirs | --sub-dirs-depth N] [--streaming|--no-streaming]
                 [--color|--plain] [--json]
tsk project list [--all|--auto|--registered] [--active] [--json]
tsk project which [--dir PATH]
tsk project register [--name NAME] [--cwd PATH] [--origin ORIGIN]
tsk project unregister <name>
```

`register` is idempotent: same name/location/origin → already up to date; empty
`location`/`origin` may be filled; conflicting non-empty values error.
Without `--name`, match location→basename(location) as name.


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
    { "origin": "github.com/xhd2015/dot-pkgs", "name": "dot-pkgs", "location": "~/Projects/xhd2015/dot-pkgs" },
    { "name": "seatalk-local-bot", "location": "~/seatalk-local-bot" }
  ]
}
```

- `name` unique when set
- `location` is the **main checkout** (git main worktree, else probe dir), tilde-form
- legacy `cwd` on read is migrated into `location` and dropped on write

## Auto ledger — `TSK_HOME/projects-auto.json`

Written on every successful `project add` (upsert-only). Identity is `origin` XOR
`name`. `location` is the **main repo** root (linked worktrees still record main),
tilde-form, set once; later adds only bump `last_seen_at` (and fill `location`
when empty). `tsk show` prints `project: <location>` when resolvable (else name,
else origin), and task `cwd:` in tilde form when recorded.

`tsk project list` (default / `--all`): union of auto + registered as an aligned
table (`NAME` `ORIGIN` `LOCATION` `TASKS`; registered-only rows show `0` until add).
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

Like `tsk tree`. Default: current project **plus** projects discovered by a
git-repo scan under the scan root (git toplevel of cwd/`--dir` when available,
else that probe dir) at max depth **3**, non-terminal stages only. Nested
checkouts (e.g. `external/dot-pkgs`) map via origin or registered cwd. Empty
discovered projects are omitted; the current project branch is kept even when
empty. Human default mode **streams**: root project prints immediately, then
each discovered project as `scan_repo` finds it (deduped by project key).
`--streaming` makes that explicit; `--no-streaming` buffers then prints (root
first, rest label-sorted). `--streaming` conflicts with `--no-streaming`,
`--json`, and `--all` / `--name` / `--project`. `--dir PATH` resolves the
project from PATH instead of cwd (conflicts with `--name`/`--project`/`--all`).
`--no-sub-dirs` disables the scan; `--sub-dirs-depth N` overrides depth
(`N >= 1`). Those two flags conflict with each other and with `--all` /
`--name` / `--project`.

`--done` / `--archived` filter to those stages (both ⇒ union). `--all` shows
every project and every stage (narrow with `--done`/`--archived`). `--stage`
filters to one stage and conflicts with `--done`/`--archived`.
