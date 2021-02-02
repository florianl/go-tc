package tc

import (
	"fmt"
	"net"

	"github.com/mdlayher/netlink"
)

const (
	tcaTunnelUnspec = iota
	tcaTunnelKeyTm
	tcaTunnelKeyParms
	tcaTunnelKeyEncIPv4Src // be32
	tcaTunnelKeyEncIPv4Dst // be32
	tcaTunnelKeyEncIPv6Src 
	tcaTunnelKeyEncIPv6Dst 
	tcaTunnelKeyEncKeyID   // be64
	tcaTunnelKeyPad
	tcaTunnelKeyEncDstPort
	tcaTunnelKeyNoCSUM
	tcaTunnelKeyEncOpts
	tcaTunnelKeyEncTOS
	tcaTunnelKeyEncTTL
)

// Tunnel contains attribute of the Tunnel discipline
type TunnelKey struct {
	Parms           *TunnelParms
	KeyEncIPv4Src   *net.IP
	KeyEncIPv4Dst   *net.IP
	KeyEncKeyID     *uint64
}

// TunnelParms from from include/uapi/linux/tc_act/tc_tunnel_key.h
type TunnelParms struct {
	Index        uint32
	Capab        uint32
	Action       uint32
	RefCnt       uint32
	BindCnt      uint32
	TunnelKeyAction uint32
}

// marshalVLan returns the binary encoding of Vlan
func marshalTunnelKey(info *TunnelKey) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Tunnel: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations

	if info.Parms != nil {
		data, err := marshalStruct(info.Parms)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaTunnelKeyParms, Data: data})
	}
	if info.KeyEncIPv4Src != nil {
		tmp, err := ipToUint32(*info.KeyEncIPv4Src)
		if err != nil {
			return []byte{}, fmt.Errorf("Tunnel - KeyEncIPv4Src: %w", err)
		}
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTunnelKeyEncIPv4Src, Data: tmp})
	}
	if info.KeyEncIPv4Dst != nil {
		tmp, err := ipToUint32(*info.KeyEncIPv4Dst)
		if err != nil {
			return []byte{}, fmt.Errorf("Tunnel - KeyEncIPv4Dst: %w", err)
		}
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaTunnelKeyEncIPv4Dst, Data: tmp})
	}
	if info.KeyEncKeyID != nil {
		options = append(options, tcOption{Interpretation: vtUint64Be, Type: tcaTunnelKeyEncKeyID, Data: *info.KeyEncKeyID})
	}
	return marshalAttributes(options)
}

// unmarshalTunnelKey parses the Tunnel-encoded data and stores the result in the value pointed to by info.
func unmarshalTunnelKey(data []byte, info *TunnelKey) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaTunnelKeyParms:
			parms := &TunnelParms{}
			if err := unmarshalStruct(ad.Bytes(), parms); err != nil {
				return err
			}
			info.Parms = parms
		case tcaTunnelKeyEncIPv4Src:
			tmp := uint32ToIP(ad.Uint32())
			info.KeyEncIPv4Src = &tmp
		case tcaTunnelKeyEncIPv4Dst:
			tmp := uint32ToIP(ad.Uint32())
			info.KeyEncIPv4Dst = &tmp
		case tcaTunnelKeyEncKeyID:
			tmp := ad.Uint64()
			info.KeyEncKeyID = &tmp
		default:
			return fmt.Errorf("unmarshalVLan()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}
