## Expected

- Exit code ≠ 0; stderr has exactly one `Error:` prefix.
- Stderr mentions conflict and both channel-id values (`ch-a`, `ch-b`).
- No messages written under `ch-a`.

## Errors

- Conflicting parent vs leaf `--channel-id` (not silent leaf-wins).

## Exit Code

- non-zero

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode == 0 {
		t.Fatal("expected non-zero exit for conflicting --channel-id")
	}
	assertStderrErrorPrefix(t, resp.Stderr)
	low := strings.ToLower(resp.Stderr)
	if !strings.Contains(low, "conflict") {
		t.Fatalf("expected conflict in stderr, got %q", resp.Stderr)
	}
	assertContains(t, resp.Stderr, "ch-a")
	assertContains(t, resp.Stderr, "ch-b")

	msgs := readMessagesJSONL(t, activeChannelDir(req, "ch-a"))
	if len(msgs) != 0 {
		t.Fatalf("expected no messages on conflict, got %d", len(msgs))
	}
}
```
