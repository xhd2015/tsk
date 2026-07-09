## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Output is **agent** format even without CODEX/PI host env.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertAutoFormatOK(t, resp, err)
	assertAgentStatusFormat(t, resp.Stdout)
	assertContains(t, stripANSI(resp.Stdout), "title: auto tsk format agent")
}
```
