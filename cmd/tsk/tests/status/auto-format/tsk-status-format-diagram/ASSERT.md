## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Output is **diagram** despite `CODEX_THREAD_ID` (env format override wins over detect).
- No agent `id:`/`title:` fact block.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertAutoFormatOK(t, resp, err)
	assertDiagramStatusFormat(t, resp.Stdout)
}
```
