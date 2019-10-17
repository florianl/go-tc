//+build gofuzz

package tc

import "bytes"

func Fuzz(data []byte) int {
	return fuzzFu32(data)
}

func fuzzFu32(data []byte) int {
	var orig U32
	err := unmarshalU32(data, &orig)
	if err != nil {
		return 0
	}

	var data2 []byte
	if data2, err = marshalU32(&orig); err != nil {
		panic(err)
	}

	if bytes.Compare(data, data2) != 0 {
		panic(err)
	}

	return 1
}
