package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaQfqUnspec = iota
	tcaQfqWeight
	tcaQfqLmax
)

// Qfq contains attributes of the qfq discipline
type Qfq struct {
	Weight uint32
	Lmax   uint32
}

func extractQfqOptions(data []byte, info *Qfq) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaQfqWeight:
			info.Weight = ad.Uint32()
		case tcaQfqLmax:
			info.Lmax = ad.Uint32()
		default:
			return fmt.Errorf("extractQfqOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

const (
	tcaHfscUnspec = iota
	tcaHfscRsc
	tcaHfscFsc
	tcaHfscUsc
)

// Hfsc contains attributes of the hfsc discipline
type Hfsc struct {
	Rsc *ServiceCurve
	Fsc *ServiceCurve
	Usc *ServiceCurve
}

func extractHfscOptions(data []byte, info *Hfsc) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaHfscRsc:
			curve := &ServiceCurve{}
			if err := extractServiceCurve(ad.Bytes(), curve); err != nil {
				return err
			}
			info.Rsc = curve
		case tcaHfscFsc:
			curve := &ServiceCurve{}
			if err := extractServiceCurve(ad.Bytes(), curve); err != nil {
				return err
			}
			info.Fsc = curve
		case tcaHfscUsc:
			curve := &ServiceCurve{}
			if err := extractServiceCurve(ad.Bytes(), curve); err != nil {
				return err
			}
			info.Usc = curve
		default:
			return fmt.Errorf("extractHfscOptions()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// ServiceCurve from include/uapi/linux/pkt_sched.h
type ServiceCurve struct {
	M1 uint32
	D  uint32
	M2 uint32
}

func extractServiceCurve(data []byte, info *ServiceCurve) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}
