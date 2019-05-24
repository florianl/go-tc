//+build linux

package tc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mdlayher/netlink"
	"github.com/mdlayher/netlink/nltest"
	"golang.org/x/sys/unix"
)

func testConn(t *testing.T) (*Tc, func()) {
	t.Helper()

	var reqCache []netlink.Message

	c := &Tc{
		con: nltest.Dial(func(req []netlink.Message) ([]netlink.Message, error) {
			if len(req) == 0 {
				// skip validation requests
				return []netlink.Message{}, nil
			}
			if diff := cmp.Diff(1, len(req)); diff != "" {
				t.Fatalf("unexpected number of request messages (-want +got):\n%s", diff)
			}

			var altered []byte
			switch req[0].Header.Type {
			case rtmNewQdisc:
				reqCache = req
			case rtmGetQdisc:
				altered = qdiscAlterResponses(t, &reqCache)
			case rtmDelQdisc:
				reqCache = []netlink.Message{}
			default:
			}
			emptyMsg := make([]netlink.Message, 0, 1)
			var data []byte
			tcmsg, err := tcmsgEncode(&Msg{
				Family:  unix.AF_UNSPEC,
				Ifindex: 0,
				Handle:  0xC001,
				Parent:  0xCAFE,
				Info:    0,
			})
			if err != nil {
				t.Fatalf("could not encode dummy Msg{}: %v", err)
			}
			data = append(data, tcmsg...)
			data = append(data, altered...)

			emptyMsg = append(emptyMsg, netlink.Message{
				Header: netlink.Header{
					Sequence: req[0].Header.Sequence,
					PID:      req[0].Header.PID,
				},
				Data: data,
			})

			return emptyMsg, nil
		}),
	}

	return c, func() {
		if err := c.Close(); err != nil {
			t.Fatalf("failed to close: %v", err)
		}
	}
}
