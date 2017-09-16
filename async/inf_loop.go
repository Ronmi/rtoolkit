package async

import "context"

// InfiniteLoop loops your function and capable to cancel-on-demand
//
// There are few things you should take care of:
//
//    - It will not interrupt current loop.
//    - It will not wait any second between tasks.
//    - It will stop if task() returns non-nil.
func InfiniteLoop(task func() error) (cancel context.CancelFunc, err chan error) {
	ctx, cancel := context.WithCancel(context.Background())
	err = make(chan error)
	go doInfiniteLooping(ctx, err, task)

	return
}

func doInfiniteLooping(ctx context.Context, errchan chan error, task func() error) {
	var err error
	for err == nil {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
			err = task()
		}
	}

	errchan <- err
}
