## Expected

- Exit code 0; stderr empty.
- Stdout is a JSON object with `inbox` and `topics` keys.
- No ANSI.

## Exit Code

- 0

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	if resp.Stderr != "" {
		t.Fatalf("stderr should be empty, got %q", resp.Stderr)
	}
	assertNoANSI(t, resp.Stdout)
	want := "{\"inbox\":[{\"id\":1,\"stage\":\"done\",\"slug\":\"solo\",\"dir\":\"[1]-done-solo\"}],\"topics\":[{\"path\":\"kb\",\"aliases\":[],\"tasks\":[{\"id\":2,\"stage\":\"create\",\"slug\":\"report\",\"dir\":\"[2]-create-report\"}],\"subtopics\":[]}]}\n"
	if resp.Stdout != want {
		t.Fatalf("stdout=%q want %q", resp.Stdout, want)
	}
}
```
