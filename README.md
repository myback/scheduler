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

	var group int
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			group++
		}

		n := i
		g := group

		j := func() {
			fmt.Printf("run jobFunc: testing/test%d/%d\n", g, n)
			time.Sleep(300 * time.Microsecond)
		}

		if err := sch.AddJob(
			ctx,
			fmt.Sprintf("testing/test%d", g),
			fmt.Sprintf("%d", n),
			n,
			j); err != nil {
			fmt.Println(err)
		}
	}

	time.Sleep(10 * time.Second)

	// Error ErrorJobAlreadyExists
	if err := sch.AddJob(ctx, "testing/test0", "1", 3, func() { fmt.Println("run updated jobFunc testing/test1/1") }); err != nil {
		fmt.Println(err)
	}

	// Error ErrorJobNotFound
	if err := sch.CancelJob("testing"); err != nil {
		fmt.Println(err)
	}

	if err := sch.UpdateJobInterval("testing/test5/10", 1); err != nil {
		fmt.Println("testing/test5/10", err)
	} else {
		fmt.Println("Interval for job testing/test5/10 updated")
	}

	fmt.Println("Running jobs:", sch.GetNumberJobs())
	time.Sleep(3 * time.Second)

	if err := sch.CancelJob("testing/test0/1"); err != nil {
		fmt.Println("testing/test0/1", err)
	}

	if err := sch.CancelJob("testing/test2"); err != nil {
		fmt.Println("testing/test2", err)
	}

	fmt.Println("Running jobs:", sch.GetNumberJobs())
	time.Sleep(3 * time.Second)

	sch.GracefulShutdown()
	sch.WaitWithTimeout(30)
}
```
