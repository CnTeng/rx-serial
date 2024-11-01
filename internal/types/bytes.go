package types

import "slices"

type BytesType struct {
	Value []byte
	*TypeMeta
}

func (bs *BytesType) Parse(data []byte) error {
	bs.Value = data[:bs.Len()]
	return nil
}

func (bs *BytesType) Serialize() []byte {
	return bs.Value
}

func (bs *BytesType) Equal(other Type) bool {
	if ob, ok := other.(*BytesType); !ok {
		return false
	} else {
		return slices.Equal(bs.Value, ob.Value) && bs.TypeMeta.Equal(ob.TypeMeta)
	}
}
