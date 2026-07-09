## Expected

- Exit code 0.
- Stdout contains `│ done │` box line.
- Line containing `done` includes green ANSI `\x1b[32m`.
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

	assertBoxLineForStage(t, resp.Stdout, "done")
	assertStageLineHasGreen(t, resp.Stdout, "done")
}
```