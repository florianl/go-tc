//+build linux

package tc

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"
)

// Tc represents a RTNETLINK wrapper
type Tc struct {
	con *netlink.Conn
}

// for detailes see https://github.com/tensorflow/tensorflow/blob/master/tensorflow/go/tensor.go#L488-L505
var nativeEndian binary.ByteOrder

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		nativeEndian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		nativeEndian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}

// Open establishes a RTNETLINK socket for traffic control
func Open(config *Config) (*Tc, error) {
	var tc Tc

	if config == nil {
		config = &Config{}
	}

	con, err := netlink.Dial(unix.NETLINK_ROUTE, &netlink.Config{NetNS: config.NetNS})
	if err != nil {
		return nil, err
	}
	tc.con = con

	return &tc, nil
}

// Close the connection
func (tc *Tc) Close() error {
	return tc.con.Close()
}

func (tc *Tc) query(req netlink.Message) ([]netlink.Message, error) {
	verify, err := tc.con.Send(req)
	if err != nil {
		return nil, err
	}

	if err := netlink.Validate(req, []netlink.Message{verify}); err != nil {
		return nil, err
	}

	return tc.con.Receive()
}

func (tc *Tc) action(action int, flags netlink.HeaderFlags, msg *Msg, opts []tcOption) error {
	tcminfo, err := tcmsgEncode(msg)
	if err != nil {
		return err
	}

	var data []byte
	data = append(data, tcminfo...)

	attrs, err := marshalAttributes(opts)
	if err != nil {
		return err
	}
	data = append(data, attrs...)
	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(action),
			Flags: netlink.Request | netlink.Acknowledge | flags,
		},
		Data: data,
	}

	msgs, err := tc.query(req)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		if msg.Header.Type == netlink.Error {
			// TODO: validate NLMSG_ERROR - see https://www.infradead.org/~tgr/libnl/doc/core.html#core_errmsg
		}
	}

	return nil
}

func (tc *Tc) get(action int, i *Msg) ([]Object, error) {
	var results []Object

	tcminfo, err := tcmsgEncode(i)
	if err != nil {
		return results, err
	}

	var data []byte
	data = append(data, tcminfo...)

	req := netlink.Message{
		Header: netlink.Header{
			Type:  netlink.HeaderType(action),
			Flags: netlink.Request | netlink.Dump,
		},
		Data: data,
	}

	msgs, err := tc.query(req)
	if err != nil {
		return results, err
	}

	for _, msg := range msgs {
		var result Object
		if err := tcmsgDecode(msg.Data[:20], &result.Msg); err != nil {
			return results, err
		}
		if err := extractTcmsgAttributes(msg.Data[20:], &result.Attribute); err != nil {
			return results, err
		}
		results = append(results, result)
	}

	return results, nil
}

// Object represents a generic traffic control object
type Object struct {
	Msg
	Attribute
}

// Attribute contains various elements for traffic control
type Attribute struct {
	Kind         string
	EgressBlock  uint32
	IngressBlock uint32
	HwOffload    uint8
	Chain        uint32
	Stats        *Stats
	XStats       *XStats
	Stats2       *Stats2

	// Filters
	BPF    *BPF
	U32    *U32
	Rsvp   *Rsvp
	Route4 *Route4
	Fw     *Fw
	Flow   *Flow

	// Classless qdiscs
	FqCodel *FqCodel
	Codel   *Codel
	Fq      *Fq
	Pie     *Pie
	Hhf     *Hhf
	Tbf     *Tbf
	Sfb     *Sfb
	Red     *Red
	MqPrio  *MqPrio
	Pfifo   *FifoOpt
	Bfifo   *FifoOpt

	// Classful qdiscs
	Htb    *Htb
	Hfsc   *Hfsc
	Dsmark *Dsmark
	Drr    *Drr
	Cbq    *Cbq
	Atm    *Atm
	Qfq    *Qfq
}

// BuildHandle is a simple helper function to construct the handle for the Tcmsg struct
func BuildHandle(major, minor uint16) uint32 {
	return uint32(major)<<16 | uint32(minor)
}

// XStats contains further statistics to the TCA_KIND
type XStats struct {
	Sfb     *SfbXStats
	Sfq     *SfqXStats
	Red     *RedXStats
	Choke   *ChokeXStats
	Htb     *HtbXStats
	Cbq     *CbqXStats
	Codel   *CodelXStats
	Hhf     *HhfXStats
	Pie     *PieXStats
	FqCodel *FqCodelXStats
}

const (
	tcaPoliceUnspec = iota
	tcaPoliceTbf
	tcaPoliceRate
	tcaPolicePeakRate
	tcaPoliceAvRate
	tcaPoliceResult
	tcaPoliceTm
	tcaPolicePad
)

// Police represents policing attributes of various filters and classes
type Police struct {
	Tbf      *Policy
	Rate     *RateSpec
	PeakRage *RateSpec
	AvRate   uint32
	Result   uint32
	Tm       *Tcft
}

func extractPoliceOptions(data []byte, info *Police) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaPoliceTbf:
			policy := &Policy{}
			if err := extractPolicy(ad.Bytes(), policy); err != nil {
				return err
			}
			info.Tbf = policy
		case tcaPoliceRate:
			rate := &RateSpec{}
			if err := extractRateSpec(ad.Bytes(), rate); err != nil {
				return err
			}
			info.Rate = rate
		case tcaPolicePeakRate:
			rate := &RateSpec{}
			if err := extractRateSpec(ad.Bytes(), rate); err != nil {
				return err
			}
			info.PeakRage = rate
		case tcaPoliceAvRate:
			info.AvRate = ad.Uint32()
		case tcaPoliceResult:
			info.Result = ad.Uint32()
		case tcaPoliceTm:
			tm := &Tcft{}
			if err := extractTcft(ad.Bytes(), tm); err != nil {
				return err
			}
			info.Tm = tm
		case tcaPolicePad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("extractPoliceOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())

		}
	}
	return nil
}
