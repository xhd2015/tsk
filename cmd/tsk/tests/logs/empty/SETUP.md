# Scenario

**Feature**: empty events.jsonl prints `0 logs`

```
tsk logs -> 0 logs
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"logs"}
	return nil
}
```
