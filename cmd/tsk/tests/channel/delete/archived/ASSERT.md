## Expected

- Exit code 0; stdout `deleted old-ch\n`.
- Tombstone exists; archive dir gone.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertStdoutTrimmedEquals(t, resp.Stdout, "deleted old-ch")
	assertFileExists(t, channelAbs(req, "tombstones/old-ch"))
	assertFileNotExists(t, archiveChannelDir(req, "old-ch"))
}
```
