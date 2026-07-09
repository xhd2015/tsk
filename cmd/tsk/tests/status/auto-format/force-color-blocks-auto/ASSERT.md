## Expected

- Exit code 0; stderr empty; stdout ends with `\n`.
- Output is **diagram** path (box art / stage box), not agent fact block.
- May include ANSI on the current stage; must not look like agent facts/spine-only output.

## Exit Code

- 0

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	assertAutoFormatOK(t, resp, err)
	assertNotAgentFactBlock(t, resp.Stdout)
	assertCompactBoxArt(t, resp.Stdout)
	assertBoxLineForStage(t, resp.Stdout, "create")
}
```
