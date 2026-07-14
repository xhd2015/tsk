## Expected

- Exit code 0; stdout `deleted temp-ch\n`.
- Tombstone exists; active dir and index gone.
- `list --all` does not show channel.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "deleted temp-ch")
	assertFileExists(t, channelAbs(req, "tombstones/temp-ch"))
	assertFileNotExists(t, activeChannelDir(req, "temp-ch"))
	assertFileNotExists(t, channelAbs(req, "index/temp-ch"))

	listResp := runTskOK(t, req, "channel", "list", "--all")
	assertNotContains(t, listResp.Stdout, "temp-ch")
}
```
