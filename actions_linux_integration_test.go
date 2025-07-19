//go:build integration && linux
// +build integration,linux

package tc

import (
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/jsimonetti/rtnetlink"
)

func TestTCActions(t *testing.T) {
	tcTestIface := "mirror"

	rtnl, err := setupDummyInterface(tcTestIface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not setup dummy interface: %v\n", err)
		return
	}
	defer rtnl.Close()

	devID, err := net.InterfaceByName(tcTestIface)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get interface ID: %v\n", err)
		return
	}
	defer func(devID uint32, rtnl *rtnetlink.Conn) {
		if err := rtnl.Link.Delete(devID); err != nil {
			fmt.Fprintf(os.Stderr, "could not delete interface: %v\n", err)
		}
	}(uint32(devID.Index), rtnl)

	config := Config{}

	tcSocket, err := Open(&config)
	if err != nil {
		t.Fatalf("could not open socket for TC: %v", err)
	}
	defer func() {
		if err := tcSocket.Close(); err != nil {
			t.Fatalf("could not close TC socket: %v", err)
		}
	}()

	existingMirredActions, err := tcSocket.Actions().Get("mirred")
	if err != nil {
		t.Fatalf("failed to get singe action: %v", err)
	}
	t.Logf("existing mirred actions: %d", len(existingMirredActions))

	mirrorIf, err := net.InterfaceByName(tcTestIface)
	if err != nil {
		t.Fatalf("failed t oget mirror interface: %v", err)
	}
	mirredActionIdx := uint32(42)
	gactActionIdx := uint32(1337)

	if err := tcSocket.Actions().Add([]*Action{
		&Action{
			Kind: "mirred",
			Mirred: &Mirred{
				Parms: &MirredParam{
					Index:   mirredActionIdx,
					Eaction: 4, /* mirror packet to INGRESS */
					IfIndex: uint32(mirrorIf.Index),
				},
			},
		},
		&Action{
			Kind: "gact",
			Gact: &Gact{
				Parms: &GactParms{
					Index:  gactActionIdx,
					Action: 2, /* drop */
				},
			},
		},
	}); err != nil {
		t.Fatalf("failed to add mirred action: %v", err)
	}

	defer func() {
		if err := tcSocket.Actions().Delete([]*Action{
			&Action{
				Kind:  "gact",
				Index: gactActionIdx,
			},
		}); err != nil {
			t.Fatalf("failed to delete action: %v", err)
		}
	}()

	mirredActions, err := tcSocket.Actions().Get("mirred")
	if err != nil {
		t.Fatalf("failed to get singe action: %v", err)
	}

	if len(mirredActions) != 1+len(existingMirredActions) {
		t.Fatalf("expected 1 mirred action but got %d", len(mirredActions))
	}

	foundInstalledMirredAction := false
	for _, a := range mirredActions {
		gotMirredActionIndex := a.Mirred.Parms.Index
		if gotMirredActionIndex == mirredActionIdx {
			foundInstalledMirredAction = true
		}
	}

	if !foundInstalledMirredAction {
		t.Fatalf("did not find installed mirred action index with index %d",
			mirredActionIdx)
	}

	if err := tcSocket.Actions().Delete([]*Action{
		&Action{
			Kind:  "mirred",
			Index: mirredActionIdx,
		},
	}); err != nil {
		t.Fatalf("failed to delete mirred action: %v", err)
	}

	mirredActions, err = tcSocket.Actions().Get("mirred")
	if err != nil {
		t.Fatalf("failed to get mirred action: %v", err)
	}

	if len(mirredActions) != len(existingMirredActions) {
		t.Fatalf("expected %d mirred action but got %d",
			len(existingMirredActions),
			len(mirredActions))
	}
}
