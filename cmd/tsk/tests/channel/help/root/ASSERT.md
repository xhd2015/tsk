## Expected

- Exit code 0; help on stdout; stderr empty.
- Documents subcommands: create, list, archive, delete, send, messages, participants, participant.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	assertHelpOK(t, resp)
	for _, sub := range []string{"create", "list", "archive", "delete", "send", "messages", "participants", "participant"} {
		assertContains(t, resp.Stdout, sub)
	}
}
```
