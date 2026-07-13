## Expected

- Exit code 0.
- Each workflow stage has a middle box row with the stage label between box/tee borders
  (`│`/`┤` … `│`/`├`, or ASCII `|`/`+` variants; padding spaces allowed e.g. `│  done  │`).
- Stages: `create`, `in_process`, `clarification`, `implementation`, `verification`, `summary`, `user_followup`, `done`.
- Stderr empty.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertStatusOK(t, resp)
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}

	assertAllStagesBoxed(t, resp.Stdout)
}
```