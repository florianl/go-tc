package tc

import (
	"fmt"
)

type U32Match struct {
	Mask    uint32 // big endian
	Value   uint32 // big endian
	Off     int32
	OffMask uint32
}

func unmarshalU32Match(data []byte, info *U32Match) error {
	return unmarshalStruct(data, info)
}

func marshalU32Match(info *U32Match) ([]byte, error) {
	if info == nil {
		return []byte{}, fmt.Errorf("marshalU32Match: %w", ErrNoArg)
	}
	return marshalStruct(info)
}
