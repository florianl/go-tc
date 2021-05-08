package tc

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestIPSetMatch(t *testing.T) {
	tests := map[string]struct {
		ID  uint16
		Dir []IPSetDir
	}{
		"src":     {ID: 13, Dir: []IPSetDir{IPSetSrc}},
		"src,src": {ID: 1337, Dir: []IPSetDir{IPSetSrc, IPSetSrc}},
		"src,dst": {ID: 42, Dir: []IPSetDir{IPSetSrc, IPSetDst}},
	}

	for name, test := range tests {
		name := name
		test := test
		t.Run(name, func(t *testing.T) {
			in := IPSetMatch{
				IPSetID: test.ID,
				Dir:     test.Dir,
			}

			data, err := marshalIPSetMatch(&in)
			if err != nil {
				t.Fatal(err)
			}
			out := IPSetMatch{}
			if err := unmarshalIPSetMatch(data, &out); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(in, out); diff != "" {
				t.Fatalf("IPSetMatch missmatch (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("invalid direction", func(t *testing.T) {
		in := IPSetMatch{
			IPSetID: 3,
		}
		if _, err := marshalIPSetMatch(&in); !errors.Is(err, ErrInvalidArg) {
			t.Fatalf("Expected ErrInvalidArg but got '%v'", err)
		}
	})
}
