## Expected

- Exit code 0; stdout `my-room\n`.
- Directory `channels/active/my-room/` exists.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "my-room")
	assertDirExists(t, activeChannelDir(req, "my-room"))
	assertChannelIndexEquals(t, req, "my-room", "active/my-room")
}
```
