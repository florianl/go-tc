//+build integration,linux

package tc

import (
	"fmt"
	"testing"
)

// Get requests don't need special priviledges
func TestLinuxTcQueueGet(t *testing.T) {

	config := Config{}

	tcSocket, err := Open(&config)
	if err != nil {
		t.Fatalf("Could not open socket for TC: %v", err)
	}
	defer func() {
		if err := tcSocket.Close(); err != nil {
			t.Fatalf("Coult not close TC socket: %v", err)
		}
	}()

	queues, qErr := tcSocket.Qdisc().Get()
	if qErr != nil {
		t.Fatalf("Could not get queues from TC socket: %v", qErr)
	}
	for _, queue := range queues {
		fmt.Printf("%#v\n", queue)
	}
}
