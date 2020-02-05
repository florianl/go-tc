package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaEtsUnspec = iota
	tcaEtsNBands
	tcaEtsNStrict
	tcaEtsQuanta
	tcaEtsQuantaBand
	tcaEtsPrioMap
	tcaEtsPrioMapBand
)

// Ets represents a struct for Enhanced Transmission Selection, a 802.1Qaz-based Qdisc.
// More info at https://lwn.net/Articles/805229/
type Ets struct {
	NBands  *uint8
	NStrict *uint8
	Quanta  *[]uint32
	PrioMap *[]uint8
}

// unmarshalEtsQuanta
func unmarshalEtsQuanta(data []byte, info *[]uint32) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaEtsQuantaBand:
			*info = append(*info, ad.Uint32())
		default:
			return fmt.Errorf("unmarshalEtsQuanta()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalEtsQuanta
func marshalEtsQuanta(info *[]uint32) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("marshalEtsQuanta: %w", ErrNoArg)
	}

	for _, band := range *info {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaEtsQuantaBand, Data: band})
	}

	return marshalAttributes(options)
}

// unmarshalEtsPrioMap
func unmarshalEtsPrioMap(data []byte, info *[]uint8) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaEtsPrioMapBand:
			*info = append(*info, ad.Uint8())
		default:
			return fmt.Errorf("unmarshalEtsPrioMap()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalEtsPrioMap
func marshalEtsPrioMap(info *[]uint8) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("marshalEtsPrioMap: %w", ErrNoArg)
	}

	for _, band := range *info {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaEtsPrioMapBand, Data: band})
	}

	return marshalAttributes(options)
}

// unmarshalEts parses the Ets-encoded data and stores the result in the value pointed to by info.
func unmarshalEts(data []byte, info *Ets) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaEtsNBands:
			tmp := ad.Uint8()
			info.NBands = &tmp
		case tcaEtsNStrict:
			tmp := ad.Uint8()
			info.NStrict = &tmp
		case tcaEtsQuanta:
			var tmp []uint32
			if err := unmarshalEtsQuanta(ad.Bytes(), &tmp); err != nil {
				return err
			}
			info.Quanta = &tmp
		case tcaEtsPrioMap:
			var tmp []uint8
			if err := unmarshalEtsPrioMap(ad.Bytes(), &tmp); err != nil {
				return err
			}
			info.PrioMap = &tmp
		default:
			return fmt.Errorf("unmarshalEts()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalEts returns the binary encoding of Ets
func marshalEts(info *Ets) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ets: %w", ErrNoArg)
	}

	if info.NBands != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaEtsNBands, Data: *info.NBands})
	}
	if info.NStrict != nil {
		options = append(options, tcOption{Interpretation: vtUint8, Type: tcaEtsNStrict, Data: *info.NStrict})
	}
	if info.Quanta != nil {
		data, err := marshalEtsQuanta(info.Quanta)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaEtsQuanta, Data: data})
	}
	if info.PrioMap != nil {
		data, err := marshalEtsPrioMap(info.PrioMap)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaEtsPrioMap, Data: data})
	}

	return marshalAttributes(options)
}
