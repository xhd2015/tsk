## Expected

- Exit code 0; valid JSON array; handles include `alice` only.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	var arr []channelParticipant
	if err := json.Unmarshal([]byte(resp.Stdout), &arr); err != nil {
		t.Fatalf("parse participants json: %v", err)
	}
	handles := participantHandles(arr)
	if len(handles) != 1 || handles[0] != "alice" {
		t.Fatalf("participants: got %v want [alice]", handles)
	}
}
```