# Scenario

**Feature**: `tsk project notes` — project-scoped journal under `projects/<id>/`

```
resolve project (origin or registered name)
  -> ensure shared project id
  -> notes.jsonl at TSK_HOME/projects/<id>/
```

## Preconditions

- Parent `project/SETUP.md` helpers available.
- Prefer `--dir` / `--project` so cwd need not be the project root.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	ensureProjectHelpersUsed()
	return nil
}
```
