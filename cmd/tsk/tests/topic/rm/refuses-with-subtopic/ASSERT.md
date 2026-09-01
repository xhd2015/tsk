## Expected

- Non-zero exit; stderr mentions subtopics.

## Exit Code

- 1

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatalf("expected failure")
	}
	assertContains(t, resp.Stderr, "subtopics")
}
```
