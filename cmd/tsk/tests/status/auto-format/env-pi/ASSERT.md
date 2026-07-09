## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Output is **agent** format (facts + spine, no box chrome).

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertAutoFormatOK(t, resp, err)
	assertAgentStatusFormat(t, resp.Stdout)
	assertContains(t, stripANSI(resp.Stdout), "title: auto env pi")
}
```
