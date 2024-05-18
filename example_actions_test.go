//go:build linux
// +build linux

package tc_test

import (
	"fmt"
	"os"

	"github.com/florianl/go-tc"
)

func ExampleActions() {
	tcnl, err := tc.Open(&tc.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open rtnetlink socket: %v\n", err)
		return
	}
	defer func() {
		if err := tcnl.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "could not close rtnetlink socket: %v\n", err)
		}
	}()

	// Create a gact Actions.
	if err := tcnl.Actions().Add([]*tc.Action{
		{
			Kind: "gact",
			Gact: &tc.Gact{
				Parms: &tc.GactParms{
					Action: 2, // drop
				},
			},
		},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to add actions: %v\n", err)
		return
	}

	// Delete the gact Actions on Index 1.
	if err := tcnl.Actions().Delete([]*tc.Action{
		{
			Kind:  "gact",
			Index: 1,
		},
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete gact actions: %v\n", err)
		return
	}
}
