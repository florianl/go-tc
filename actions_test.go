package tc

import "testing"

func TestActions(t *testing.T) {
	tcSocket, done := testConn(t)
	defer done()

	t.Run("ErrNoArg", func(t *testing.T) {
		err := tcSocket.Actions().Add(nil)
		if err != ErrNoArg {
			t.Fatalf("expected ErrNoArg, received: %v", err)
		}

		err = tcSocket.Actions().Replace(nil)
		if err != ErrNoArg {
			t.Fatalf("expected ErrNoArg, received: %v", err)
		}

		err = tcSocket.Actions().Delete(nil)
		if err != ErrNoArg {
			t.Fatalf("expected ErrNoArg, received: %v", err)
		}
	})

	t.Run("simple example", func(t *testing.T) {
		newActions := []*Action{
			{
				Kind: "mirred",
				Mirred: &Mirred{
					Parms: &MirredParam{
						Index:  42,
						Action: 1337,
					},
				},
			},
			{
				Kind: "gact",
				Gact: &Gact{
					Parms: &GactParms{
						Index:  1337,
						Action: 42,
					},
				},
			},
		}
		if err := tcSocket.Actions().Add(newActions); err != nil {
			t.Fatal(err)
		}

		if err := tcSocket.Actions().Replace(newActions); err != nil {
			t.Fatal(err)
		}

		if err := tcSocket.Actions().Delete(newActions); err != nil {
			t.Fatal(err)
		}
	})
}
