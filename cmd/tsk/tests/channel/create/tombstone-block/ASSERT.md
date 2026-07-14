## Expected

- Exit code 1; stderr `Error:`.
- Tombstone `channels/tombstones/eng-alerts` exists.
- No `channels/active/eng-alerts` or index entry.

## Exit Code

- 1

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	assertFileExists(t, channelAbs(req, "tombstones/eng-alerts"))
	assertFileNotExists(t, activeChannelDir(req, "eng-alerts"))
	assertFileNotExists(t, channelAbs(req, "index/eng-alerts"))
	ts := readTombstone(t, req, "eng-alerts")
	if ts.ID != "eng-alerts" {
		t.Fatalf("tombstone id: got %q want eng-alerts", ts.ID)
	}
}
```
