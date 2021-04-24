package tc

type U32Match struct {
	Mask    uint32 // big endian
	Value   uint32 // big endian
	Off     int32
	OffMask int32
}

func unmarshalU32Match(data []byte, info *U32Match) error {
	return unmarshalStruct(data, info)
}

func marshalU32Match(info *U32Match) ([]byte, error) {
	return marshalStruct(info)
}
