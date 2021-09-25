//go:build gofuzz
// +build gofuzz

package tc

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
)

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

	if diff := cmp.Diff(data2, data); diff != "" {
		panic(fmt.Sprintf("Missmatch (-want +got):\n%s", diff))
	}

	return 1
}
