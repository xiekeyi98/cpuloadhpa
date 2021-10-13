
```go
package main

import (
	"context"

	"git.code.oa.com/cpuload/cl"
)

func main() {
	l := cl.NewPayloadPercent(context.Background(), 50)
	l.Run()
	select {}
}

```