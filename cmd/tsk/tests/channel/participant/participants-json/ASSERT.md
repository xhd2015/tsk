## Expected

- Exit code 0; valid JSON array; handles include `alice` and `agent`.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit code %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assertNoANSI(t, resp.Stdout)
	var arr []channelParticipant
	if err := json.Unmarshal([]byte(strings.TrimSpace(resp.Stdout)), &arr); err != nil {
		t.Fatalf("parse participants json: %v", err)
	}
	handles := participantHandles(channelJSON{Participants: arr})
	sort.Strings(handles)
	if len(handles) != 2 || handles[0] != "agent" || handles[1] != "alice" {
		t.Fatalf("participants: got %v", handles)
	}
}
```
