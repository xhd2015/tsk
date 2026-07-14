## Expected

- Exit code 0; valid JSON array; no ANSI.
- One message with `body: json body`, `sender: alice`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertNoANSI(t, resp.Stdout)
	var arr []channelMessage
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Stdout)), &arr); err != nil {
		t.Fatalf("parse messages json: %v", err)
	}
	if len(arr) != 1 {
		t.Fatalf("messages: got %d want 1", len(arr))
	}
	if arr[0].Body != "json body" || arr[0].Sender != "alice" {
		t.Fatalf("message: got %+v", arr[0])
	}
}
```
