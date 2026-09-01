# Scenario

**Feature**: installed `pmark` forwards to `tsk project add` (e2e)

```
label: e2e
tsk install pmark; PATH=… pmark "note from wrapper"
  -> project task id on stdout
```

## Preconditions

- Git repo with origin in WorkRoot (project add needs origin).
- `HOME=WorkRoot`; Assert runs the wrapper with PATH including the tsk binary.

## Steps

1. Init git with origin.
2. Install pmark (Run also re-installs for the default harness).
3. Assert executes `~/.local/bin/pmark` with the title.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	installInitGitRepo(t, req.WorkRoot, "https://github.com/example/pmark-fwd.git")
	req.Title = "note from wrapper"
	req.Args = []string{"install", "pmark"}
	return nil
}
```
