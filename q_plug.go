package tc

import "fmt"

const (
	PlugBuffer PlugAction = iota
	PlugReleaseOne
	PlugReleaseIndefinite
	PlugLimit
)

type Plug struct {
	Action PlugAction
	Limit  uint32
}

type PlugAction int32

func marshalPlug(info *Plug) ([]byte, error) {
	if info == nil {
		return []byte{}, fmt.Errorf("Plug: %w", ErrNoArg)
	}
	return marshalStruct(info)
}

func unmarshalPlug(data []byte, info *Plug) error {
	// So far the kernel does not implement this functionality.
	return ErrNotImplemented
}
