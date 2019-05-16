//+build linux

package tc

import (
	"bytes"
	"encoding/binary"
	"testing"

	"golang.org/x/sys/unix"
)

func TestQdisc(t *testing.T) {
	tcSocket, done := testConn(t)
	defer done()

	err := tcSocket.Qdisc().Add(nil)
	if err != ErrNoArg {
		t.Fatalf("expected ErrNoArg, received: %v", err)
	}

	testQdisc := Object{
		Msg{
			Family:  unix.AF_UNSPEC,
			Ifindex: 123,
			Handle:  BuildHandle(0xFFFF, 0x0000),
			Parent:  0xFFFFFFF1,
			Info:    0,
		},
		Attribute{
			Kind: "clsact",
		},
	}

	if err := tcSocket.Qdisc().Add(&testQdisc); err != nil {
		t.Fatalf("could not add new qdisc: %v", err)
	}

	qdiscs, err := tcSocket.Qdisc().Get()
	if err != nil {
		t.Fatalf("could not get qdiscs: %v\n", err)
		return
	}
	for _, qdisc := range qdiscs {
		t.Logf("%#v\n", qdisc)
	}

	if err := tcSocket.Qdisc().Delete(&testQdisc); err != nil {
		t.Fatalf("could not delete qdisc: %v", err)
	}
}

func qdiscGetResponses(t *testing.T) []response {
	t.Helper()
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
	noqueue, err := marshalAttributes([]tcOption{
		tcOption{Interpretation: vtString, Type: tcaKind, Data: "noqueue"},
		tcOption{Interpretation: vtUint8, Type: tcaHwOffload, Data: uint8(0)},
		tcOption{Interpretation: vtBytes, Type: tcaStats2, Data: stats2.Bytes()},
	})
	if err != nil {
		t.Fatalf("could not marshal attributes: %v", err)
	}

	fqcodel, err := marshalAttributes([]tcOption{
		tcOption{Interpretation: vtString, Type: tcaKind, Data: "fq_codel"},
		tcOption{Interpretation: vtBytes, Type: tcaOptions, Data: []byte{0x08, 0x00, 0x01, 0x00, 0x87, 0x13, 0x00, 0x00, 0x08, 0x00, 0x02, 0x00, 0x00, 0x28, 0x00, 0x00, 0x08, 0x00, 0x03, 0x00, 0x9f, 0x86, 0x01, 0x00, 0x08, 0x00, 0x04, 0x00, 0x01, 0x00, 0x00, 0x00, 0x08, 0x00, 0x06, 0x00, 0xea, 0x05, 0x00, 0x00, 0x08, 0x00, 0x08, 0x00, 0x40, 0x00, 0x00, 0x00, 0x08, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x02, 0x08, 0x00, 0x05, 0x00, 0x00, 0x04, 0x00, 0x00}},
	})
	if err != nil {
		t.Fatalf("could not marshal attributes: %v", err)
	}

	responses := []response{
		response{
			Msg: Msg{
				Family:  unix.AF_UNSPEC,
				Ifindex: 1,
				Handle:  321,
				Parent:  456,
				Info:    0,
			},
			data: noqueue,
		},
		response{
			Msg: Msg{
				Family:  unix.AF_UNSPEC,
				Ifindex: 2,
				Handle:  789,
				Parent:  987,
				Info:    2,
			},
			data: fqcodel,
		},
	}

	return responses
}
