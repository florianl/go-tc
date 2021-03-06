package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaActUnspec = iota
	tcaActKind
	tcaActOptions
	tcaActIndex
	tcaActStats
	tcaActPad
	tcaActCookie
)

const (
	ActBind      = 1
	ActNoBind    = 0
	ActUnbind    = 1
	ActNoUnbind  = 0
	ActReplace   = 1
	ActNoReplace = 0
)

const (
	ActOk         = 0
	ActReclassify = 1
	ActShot       = 2
	ActPipe       = 3
	ActStolen     = 4
	ActQueued     = 5
	ActRepeat     = 6
	ActRedirect   = 7
	ActTrap       = 8
)

// Action represents action attributes of various filters and classes
type Action struct {
	Kind      string
	Index     uint32
	Stats     *GenStats
	Cookie    *[]byte
	Bpf       *ActBpf
	ConnMark  *Connmark
	CSum      *Csum
	Defact    *Defact
	Ife       *Ife
	Ipt       *Ipt
	Mirred    *Mirred
	Nat       *Nat
	Sample    *Sample
	VLan      *VLan
	Police    *Police
	TunnelKey *TunnelKey
}

func unmarshalActions(data []byte, actions *[]*Action) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		action := &Action{}
		if err := unmarshalAction(ad.Bytes(), action); err != nil {
			return err
		}
		*actions = append(*actions, action)
	}
	return nil
}

// unmarshalAction parses the Action-encoded data and stores the result in the value pointed to by info.
func unmarshalAction(data []byte, info *Action) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	var actOptions []byte
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaActKind:
			info.Kind = ad.String()
		case tcaActIndex:
			info.Index = ad.Uint32()
		case tcaActOptions:
			actOptions = ad.Bytes()
		case tcaActCookie:
			tmp := ad.Bytes()
			info.Cookie = &tmp
		case tcaActStats:
			stats := &GenStats{}
			if err := unmarshalGenStats(ad.Bytes(), stats); err != nil {
				return err
			}
			info.Stats = stats
		case tcaActPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("unmarshalAction()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	if len(actOptions) > 0 {
		if err := extractActOptions(actOptions, info, info.Kind); err != nil {
			return err
		}
	}

	return nil
}

func marshalActions(info []*Action) ([]byte, error) {
	options := []tcOption{}

	for i, action := range info {
		data, err := marshalAction(action)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: uint16(i + 1), Data: data})
	}

	return marshalAttributes(options)
}

// marshalAction returns the binary encoding of Action
func marshalAction(info *Action) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Action: %w", ErrNoArg)
	}

	if len(info.Kind) == 0 {
		return []byte{}, fmt.Errorf("kind is missing")
	}

	// TODO: improve logic and check combinations
	switch info.Kind {
	case "bpf":
		data, err := marshalActBpf(info.Bpf)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "connmark":
		data, err := marshalConnmark(info.ConnMark)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "csum":
		data, err := marshalCsum(info.CSum)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "defact":
		data, err := marshalDefact(info.Defact)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "ife":
		data, err := marshalIfe(info.Ife)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "ipt":
		data, err := marshalIpt(info.Ipt)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "mirred":
		data, err := marshalMirred(info.Mirred)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "nat":
		data, err := marshalNat(info.Nat)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "sample":
		data, err := marshalSample(info.Sample)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "vlan":
		data, err := marshalVlan(info.VLan)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "police":
		data, err := marshalPolice(info.Police)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	case "tunnel_key":
		data, err := marshalTunnelKey(info.TunnelKey)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActOptions, Data: data})
	default:
		return []byte{}, fmt.Errorf("unknown kind '%s'", info.Kind)
	}
	options = append(options, tcOption{Interpretation: vtString, Type: tcaActKind, Data: info.Kind})

	if info.Index != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaActIndex, Data: info.Index})
	}
	if info.Stats != nil {
		data, err := marshalGenStats(info.Stats)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActStats, Data: data})
	}
	if info.Cookie != nil {
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaActCookie, Data: bytesValue(info.Cookie)})
	}
	return marshalAttributes(options)
}

func extractActOptions(data []byte, act *Action, kind string) error {
	switch kind {
	case "bpf":
		info := &ActBpf{}
		if err := unmarshalActBpf(data, info); err != nil {
			return err
		}
		act.Bpf = info
	case "connmark":
		info := &Connmark{}
		if err := unmarshalConnmark(data, info); err != nil {
			return err
		}
		act.ConnMark = info
	case "csum":
		info := &Csum{}
		if err := unmarshalCsum(data, info); err != nil {
			return err
		}
		act.CSum = info
	case "defact":
		info := &Defact{}
		if err := unmarshalDefact(data, info); err != nil {
			return err
		}
		act.Defact = info
	case "ife":
		info := &Ife{}
		if err := unmarshalIfe(data, info); err != nil {
			return err
		}
		act.Ife = info
	case "ipt":
		info := &Ipt{}
		if err := unmarshalIpt(data, info); err != nil {
			return err
		}
		act.Ipt = info
	case "mirred":
		info := &Mirred{}
		if err := unmarshalMirred(data, info); err != nil {
			return err
		}
		act.Mirred = info
	case "nat":
		info := &Nat{}
		if err := unmarshalNat(data, info); err != nil {
			return err
		}
		act.Nat = info
	case "sample":
		info := &Sample{}
		if err := unmarshalSample(data, info); err != nil {
			return err
		}
		act.Sample = info
	case "vlan":
		info := &VLan{}
		if err := unmarshalVLan(data, info); err != nil {
			return err
		}
		act.VLan = info
	case "police":
		info := &Police{}
		if err := unmarshalPolice(data, info); err != nil {
			return err
		}
		act.Police = info
	case "tunnel_key":
		info := &TunnelKey{}
		if err := unmarshalTunnelKey(data, info); err != nil {
			return err
		}
		act.TunnelKey = info
	default:
		return fmt.Errorf("extractActOptions(): unsupported kind: %s", kind)

	}
	return nil
}
