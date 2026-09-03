# Scenario

**Feature**: `tsk project register` without `--name` matches by path or auto-basename

```
register --cwd DIR (no --name) -> name = basename(DIR)
re-register no --name -> already up to date / match by cwd
```
