## Expected Output

```
---
version: 3
---
tsk
add
channel
install
note
overview
project
topic
tree
workflow
working-on-task
```

## Exit Code

- 0

```go
import (
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	assertErrIsNil(t, err)
	if resp.ExitCode != 0 {
		t.Fatalf("exit %d stderr=%q", resp.ExitCode, resp.Stderr)
	}
	assert.Output(t, resp.Stdout, `---
version: 3
---
tsk
add
channel
install
note
overview
project
topic
tree
workflow
working-on-task
`)
}
```
