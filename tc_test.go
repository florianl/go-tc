//+build linux

package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nltest"
)

func testConn(t *testing.T) (*Tc, func()) {
	t.Helper()

	c := &Tc{
		con: nltest.Dial(func(req []netlink.Message) ([]netlink.Message, error) {
			if len(req) == 0 {
				// skip validation requests
				return []netlink.Message{}, nil
			}
			if diff := cmp.Diff(1, len(req)); diff != "" {
				t.Fatalf("unexpected number of request messages (-want +got):\n%s", diff)
			}

			var responses []response
			switch req[0].Header.Type {
			case rtmGetQdisc:
				responses = qdiscGetResponses(t)
			default:
			}
			// Return many messages in response to the single request.
			h := netlink.Header{
				Sequence: req[0].Header.Sequence,
				PID:      req[0].Header.PID,
			}
			msgs := make([]netlink.Message, 0, len(responses))
			for _, r := range responses {
				_ = r
				var data []byte
				tcmsg, err := tcmsgEncode(&r.Msg)
				if err != nil {
					t.Fatalf("could not encode %v: %v", r.Msg, err)
				}
				data = append(data, tcmsg...)
				data = append(data, r.data...)

				msgs = append(msgs, netlink.Message{
					Header: h,
					Data:   data,
				})
			}

			return msgs, nil
		}),
	}

	return c, func() {
		if err := c.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}
	}
}

type response struct {
	Msg
	data []byte
}
