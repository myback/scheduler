# Yet another scheduler

### Example
```go
package main

import (
	"context"
	"fmt"
	"github.com/myback/scheduler"
	"time"
)

func main() {
	sch := scheduler.New()
	ctx := context.Background()

	for i := 1; i <= 10; i++ {
		n := i
		j := func() {
			fmt.Printf("run jobFunc: %d\n", n)
			time.Sleep(300 * time.Microsecond)
		}
		if err := sch.AddJob(
			ctx,
			fmt.Sprintf("test%d", i),
			n,
			j); err != nil {
			fmt.Println(err)
		}
	}

	time.Sleep(15 * time.Second)

	// Error ErrorJobAlreadyExists
	if err := sch.AddJob(ctx, "test1", 3, func() { fmt.Println("run updated jobFunc_1") }); err != nil {
		fmt.Println(err)
	}

	// Error ErrorJobNotFound
	if err := sch.CancelJob("test"); err != nil {
		fmt.Println(err)
	}

	sch.UpdateJobInterval("test10", 1)

	fmt.Println("Running jobs:", sch.GetNumberJobs())
	time.Sleep(5 * time.Second)

	sch.CancelJob("test1")
	sch.CancelJob("test4")

	fmt.Println("Running jobs:", sch.GetNumberJobs())
	time.Sleep(5 * time.Second)

	sch.GracefulShutdown()
	sch.WaitWithTimeout(30)
}
```
