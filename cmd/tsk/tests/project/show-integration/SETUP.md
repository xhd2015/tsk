# Scenario

**Feature**: `tsk show` project line prefers ledger `location` (tilde), else name, else origin; task `cwd:` stays as recorded

```
project add | update --set-project
  -> show: cwd: <task cwd>; project: <location|name|origin>
```

## Intentional exclusions

- Task.json still has create-time `cwd` (shown tilde-form); ledger uses `location` only.
