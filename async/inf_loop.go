package async

import "context"

// InfiniteLoop loops your function and capable to cancel-on-demand
//
// There are few things you should take care of:
//
//    - It will not interrupt current loop.
//    - It will not wait any second between tasks.
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
	close(errchan)
}

// HookedInfiniteLoop is identical with InfiniteLoop, excepts it uses callback instead of channel
func HookedInfiniteLoop(task func() error, cb func(error)) (cancel context.CancelFunc) {
	cancel, errchan := InfiniteLoop(task)

	go func() {
		if cb != nil {
			cb(<-errchan)
		}
	}()

	return cancel
}
