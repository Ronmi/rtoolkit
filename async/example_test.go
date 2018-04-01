package async

import (
	"fmt"
	"time"
)

func Example() {
	// our function
	x := 0
	f := func() error {
		// do some time-consume task like calling external api
		x++
		fmt.Printf("%d ", x)

		return nil
	}

	// - run the function 3 times
	// - at least 500ms between two executions
	f1 := OnceAtMost(500*time.Millisecond, f)
	f1()
	f1()
	f1()

	// - run the function as many times as possible for 500ms
	// - at least 100ms between two executions
	//
	// WARNING: THIS IS NOT A GOOD APPROACH IF f() RUNS FAST.
	begin := time.Now()
	f2 := OnceWithin(100*time.Millisecond, f)
	for time.Now().Sub(begin) < 500*time.Millisecond {
		f2()
	}

	// function always spend at least 100ms
	RunAtLeast(100*time.Millisecond, f)()

	// Output: 1 2 3 4 5 6 7 8 9
}
