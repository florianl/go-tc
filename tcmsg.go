package rtnetlink

import (
	"bytes"
	"encoding/binary"
)

// Tcmsg contains basic traffic controll elements
type Tcmsg struct {
	Family  uint32
	Ifindex uint32
	Handle  uint32
	Parent  uint32
	Info    uint32
}

func tcmsgEncode(i *Tcmsg) ([]byte, error) {
	var tcm bytes.Buffer
	err := binary.Write(&tcm, nativeEndian, *i)
	return tcm.Bytes(), err
}

func tcmsgDecode(data []byte, tc *Tcmsg) error {
	b := bytes.NewReader(data)
	if err := binary.Read(b, nativeEndian, tc); err != nil {
		return err
	}
	return nil
}
