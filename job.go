package scheduler

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type job struct {
	cancelFunc     context.CancelFunc
	jobFunc        func()
	namespacedName string
	ticker         *time.Ticker
}

func (j *job) cancel() {
	j.cancelFunc()
}

func (j *job) start(ctx context.Context) {
	ctx, j.cancelFunc = context.WithCancel(ctx)
	wg.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("%s: ERROR: %s\n", j.namespacedName, r)
			}
		}()

		runtime.Gosched()

		for {
			select {
			case <-j.ticker.C:
				j.jobFunc()
			case <-ctx.Done():
				j.ticker.Stop()
				wg.Done()
				return
			}
		}
	}()
}

func (j *job) updateInterval(intervalSec time.Duration) {
	j.ticker.Reset(intervalSec)
}
