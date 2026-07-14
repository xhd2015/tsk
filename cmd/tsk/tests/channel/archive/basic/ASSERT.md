## Expected

- Exit code 0; stdout `archived eng-alerts\n`.
- `channels/archive/eng-alerts/` exists; `channels/active/eng-alerts/` gone.
- Index `archive/eng-alerts`; `channel.json` status `archived`.
- Default `list` excludes channel.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "archived eng-alerts")

	assertFileNotExists(t, activeChannelDir(req, "eng-alerts"))
	archDir := archiveChannelDir(req, "eng-alerts")
	assertDirExists(t, archDir)
	assertChannelIndexEquals(t, req, "eng-alerts", "archive/eng-alerts")

	ch := readChannelJSON(t, archDir)
	if ch.Status != "archived" {
		t.Fatalf("status: got %q want archived", ch.Status)
	}

	listResp := runTskOK(t, req, "channel", "list")
	assertNotContains(t, listResp.Stdout, "eng-alerts")
}
```
