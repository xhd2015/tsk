## Expected

- Help exits 0.
- Root lists `add`, `tree`, `list`, `register`.
- `add` documents `--project`.
- `tree` documents `--dir` / `--all` / `--done` / `--archived` / `--json`.
- `list` is registry list.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	assertContains(t, resp.Stdout, "add")
	assertContains(t, resp.Stdout, "tree")
	assertContains(t, resp.Stdout, "list")
	assertContains(t, resp.Stdout, "register")

	rAdd := runTskCmd(t, req, "project", "add", "-h")
	assertHelpOK(t, rAdd)
	assertContains(t, rAdd.Stdout, "--project")
	assertContains(t, rAdd.Stdout, "--dir")

	rTree := runTskCmd(t, req, "project", "tree", "-h")
	assertHelpOK(t, rTree)
	assertContains(t, rTree.Stdout, "--dir")
	assertContains(t, rTree.Stdout, "--all")
	assertContains(t, rTree.Stdout, "--done")
	assertContains(t, rTree.Stdout, "--archived")
	assertContains(t, rTree.Stdout, "--json")
	assertContains(t, rTree.Stdout, "--no-sub-dirs")
	assertContains(t, rTree.Stdout, "--sub-dirs-depth")

	rList := runTskCmd(t, req, "project", "list", "-h")
	assertHelpOK(t, rList)
	assertContains(t, rList.Stdout, "--all")
	assertContains(t, rList.Stdout, "--auto")
	assertContains(t, rList.Stdout, "--registered")
	assertContains(t, rList.Stdout, "--active")
}
```


