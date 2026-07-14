## Expected

- Exit code 1; stderr `Error:`.
- Agent still present.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	ch := readChannelJSON(t, activeChannelDir(req, "solo-ch"))
	if len(ch.Participants) != 1 || ch.Participants[0].Handle != "agent" {
		t.Fatalf("expected only agent left, got %+v", ch.Participants)
	}
}
```
