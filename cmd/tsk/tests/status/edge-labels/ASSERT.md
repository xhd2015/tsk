## Expected

- Exit code 0.
- Edge labels appear in correct vertical order relative to stage boxes:
  - `claim` follows `│ create │` box row.
  - `research` appears after `│ in_process │` and before `│ clarification │`.
  - `confirmed` appears after `│ clarification │` and before `│ implementation │`.
  - `questions` follows `│ summary │` box row.
  - `satisfied` appears near the `user_followup`→`done` merge (not before `│ verification │`).
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

	assertEdgeLabelOrder(t, resp.Stdout)
}
```