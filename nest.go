package rtnetlink

import (
	"github.com/mdlayher/netlink"
)

type RtNlOption struct {
	Type uint16
	Data interface{}
}

var optionEncoderType = map[uint16]uint8{
	TCA_KIND: 5,
}

func nestAttributes(options []RtNlOption) ([]byte, error) {
	ad := netlink.NewAttributeEncoder()

	for _, option := range options {
		switch optionEncoderType[option.Type] {
		case 1:
			ad.Uint8(option.Type, (option.Data).(uint8))
		case 2:
			ad.Uint16(option.Type, (option.Data).(uint16))
		case 3:
			ad.Uint32(option.Type, (option.Data).(uint32))
		case 4:
			ad.Uint64(option.Type, (option.Data).(uint64))
		case 5:
			ad.String(option.Type, (option.Data).(string))
		case 6:
			ad.Bytes(option.Type, (option.Data).([]byte))
		}
	}

	return ad.Encode()
}
