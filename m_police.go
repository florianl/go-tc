package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaPoliceUnspec = iota
	tcaPoliceTbf
	tcaPoliceRate
	tcaPolicePeakRate
	tcaPoliceAvRate
	tcaPoliceResult
	tcaPoliceTm
	tcaPolicePad
	tcaPoliceRate64
	tcaPolicePeakRate64
)

// PolicyAction defines the action that is applied by Policy.
type PolicyAction uint32

// Default Policy actions.
// PolicyUnspec - skipped as it is -1
const (
	PolicyOk PolicyAction = iota
	PolicyReclassify
	PolicyShot
	PolicyPipe
)

// Police represents policing attributes of various filters and classes
type Police struct {
	Tbf        *Policy
	Rate       *RateSpec
	PeakRate   *RateSpec
	AvRate     *uint32
	Result     *uint32
	Tm         *Tcft
	Rate64     *uint64
	PeakRate64 *uint64
}

// unmarshalPolice parses the Police-encoded data and stores the result in the value pointed to by info.
func unmarshalPolice(data []byte, info *Police) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaPoliceTbf:
			policy := &Policy{}
			if err := unmarshalStruct(ad.Bytes(), policy); err != nil {
				return err
			}
			info.Tbf = policy
		case tcaPoliceRate:
			rate := &RateSpec{}
			if err := unmarshalStruct(ad.Bytes(), rate); err != nil {
				return err
			}
			info.Rate = rate
		case tcaPolicePeakRate:
			rate := &RateSpec{}
			if err := unmarshalStruct(ad.Bytes(), rate); err != nil {
				return err
			}
			info.PeakRate = rate
		case tcaPoliceAvRate:
			info.AvRate = uint32Ptr(ad.Uint32())
		case tcaPoliceResult:
			info.Result = uint32Ptr(ad.Uint32())
		case tcaPoliceTm:
			tm := &Tcft{}
			if err := unmarshalStruct(ad.Bytes(), tm); err != nil {
				return err
			}
			info.Tm = tm
		case tcaPolicePad:
			// padding does not contain data, we just skip it
		case tcaPoliceRate64:
			info.Rate64 = uint64Ptr(ad.Uint64())
			return ErrNotImplemented
		case tcaPolicePeakRate64:
			info.PeakRate64 = uint64Ptr(ad.Uint64())
			return ErrNotImplemented
		default:
			return fmt.Errorf("UnmarshalPolice()\t%d\n\t%v", ad.Type(), ad.Bytes())

		}
	}
	return nil
}

// marshalPolice returns the binary encoding of Police
func marshalPolice(info *Police) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Police: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Rate != nil {
		data, err := marshalStruct(info.Rate)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaPoliceRate, Data: data})
	}
	if info.PeakRate != nil {
		data, err := marshalStruct(info.PeakRate)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaPolicePeakRate, Data: data})
	}
	if info.Tbf != nil {
		data, err := marshalStruct(info.Tbf)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaPoliceTbf, Data: data})
	}
	if info.AvRate != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPoliceAvRate, Data: uint32Value(info.AvRate)})
	}
	if info.Result != nil {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaPoliceResult, Data: uint32Value(info.Result)})
	}
	if info.Rate64 != nil {
		return []byte{}, fmt.Errorf("police: rate64: %w", ErrNotImplemented)
	}
	if info.PeakRate64 != nil {
		return []byte{}, fmt.Errorf("police: peakrate64: %w", ErrNotImplemented)
	}
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	return marshalAttributes(options)
}
