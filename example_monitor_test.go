package tc_test

import (
	"context"
	"fmt"
	"time"

	"github.com/florianl/go-tc"
)

// This example demonstrates how Monitor() can be used
func ExampleTc_Monitor() {
	tcSocket, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Printf("could not open socket for TC: %v", err)
		return
	}
	defer func() {
		if err := tcSocket.Close(); err != nil {
			fmt.Printf("coult not close TC socket: %v", err)
			return
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Hook function mon, which is called every time,
	// something is received by the kernel on this socket
	mon := func(action uint16, m tc.Object) int {
		fmt.Printf("Action:\t%d\nObject: \t%#v\n", action, m)
		return 0
	}

	tcSocket.Monitor(ctx, 10*time.Millisecond, mon)

	<-ctx.Done()
}
