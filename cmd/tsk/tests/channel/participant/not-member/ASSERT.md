## Expected

- Exit code 1; stderr `Error:`.
- Dave not added.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	ch := readChannelJSON(t, activeChannelDir(req, "closed-ch"))
	assertChannelParticipantsSorted(t, ch, []string{"agent", "alice"})
}
```
