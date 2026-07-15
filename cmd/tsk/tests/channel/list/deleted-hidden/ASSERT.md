## Expected

- Exit code 0.
- Stdout does not contain `gone-channel`.
- Tombstone still on disk.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertNotContains(t, resp.Stdout, "gone-channel")
	assertFileExists(t, channelAbs(req, "tombstones/gone-channel.json"))
}
```
