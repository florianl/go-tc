//+build linux

package tc

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

func TestQdisc(t *testing.T) {
	tcSocket, done := testConn(t)
	defer done()

	err := tcSocket.Qdisc().Add(nil)
	if err != ErrNoArg {
		t.Fatalf("expected ErrNoArg, received: %v", err)
	}

	fqCodelOptions := &FqCodel{
		Target: 42,
		Limit:  0xCAFE,
	}

	tests := map[string]struct {
		kind    string
		fqCodel *FqCodel
	}{
		"clsact":   {kind: "clsact"},
		"fq_codel": {kind: "fq_codel", fqCodel: fqCodelOptions},
	}

	tcMsg := Msg{
		Family:  unix.AF_UNSPEC,
		Ifindex: 123,
		Handle:  BuildHandle(0xFFFF, 0x0000),
		Parent:  0xFFFFFFF1,
		Info:    0,
	}
	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {

			testQdisc := Object{
				tcMsg,
				Attribute{
					Kind:    testcase.kind,
					FqCodel: testcase.fqCodel,
				},
			}

			if err := tcSocket.Qdisc().Add(&testQdisc); err != nil {
				t.Fatalf("could not add new qdisc: %v", err)
			}

			qdiscs, err := tcSocket.Qdisc().Get()
			if err != nil {
				t.Fatalf("could not get qdiscs: %v", err)
			}
			for _, qdisc := range qdiscs {
				t.Logf("%#v\n", qdisc)
			}

			if err := tcSocket.Qdisc().Delete(&testQdisc); err != nil {
				t.Fatalf("could not delete qdisc: %v", err)
			}

		})
	}

}

func qdiscAlterResponses(t *testing.T, cache *[]netlink.Message) []byte {
	t.Helper()
	var tmp []Object
	var dataStream []byte

	// Decode data from cache
	for _, msg := range *cache {
		var result Object
		if err := extractTcmsgAttributes(msg.Data[20:], &result.Attribute); err != nil {
			t.Fatalf("could not decode attributes: %v", err)
		}
		tmp = append(tmp, result)
	}

	var stats2 bytes.Buffer
	if err := binary.Write(&stats2, nativeEndian, &Stats2{
		Bytes:      42,
		Packets:    1,
		Qlen:       1,
		Backlog:    0,
		Drops:      0,
		Requeues:   0,
		Overlimits: 42,
	}); err != nil {
		t.Fatalf("could not encode stats2: %v", err)
	}

	var stats bytes.Buffer
	if err := binary.Write(&stats, nativeEndian, &Stats{
		Bytes:      32,
		Packets:    1,
		Drops:      0,
		Overlimits: 0,
		Bps:        1,
		Pps:        1,
		Qlen:       1,
		Backlog:    0,
	}); err != nil {
		t.Fatalf("could not encode stats: %v", err)
	}

	// Alter and marshal data
	for _, obj := range tmp {
		var data []byte
		var attrs []tcOption

		attrs = append(attrs, tcOption{Interpretation: vtString, Type: tcaKind, Data: obj.Kind})
		attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaStats2, Data: stats2.Bytes()})
		attrs = append(attrs, tcOption{Interpretation: vtBytes, Type: tcaStats, Data: stats.Bytes()})
		attrs = append(attrs, tcOption{Interpretation: vtUint8, Type: tcaHwOffload, Data: uint8(0)})

		marshaled, err := marshalAttributes(attrs)
		if err != nil {
			t.Fatalf("could not marshal attributes: %v", err)
		}
		data = append(data, marshaled...)

		dataStream = append(dataStream, data...)

	}
	return dataStream
}
