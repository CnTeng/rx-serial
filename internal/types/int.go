package types

type IntType struct {
	Value int
	*TypeMeta
}

func NewConstInt(v int) *IntType {
	return &IntType{Value: v}
}

func (i *IntType) Parse(data []byte) error {
	i.Value = parseInt(data[:i.Len()], i.Config.ByteOrder)
	return nil
}

func (i *IntType) Serialize() []byte {
	return serializeInt(i.Value, i.Config.ByteOrder, i.Len())
}

func (i *IntType) Equal(other Type) bool {
	if oi, ok := other.(*IntType); !ok {
		return false
	} else {
		return i.Value == oi.Value && i.TypeMeta.Equal(oi.TypeMeta)
	}
}

func parseInt(buf []byte, order ByteOrder) int {
	var v int
	switch order {
	case LittleEndian:
		for i, b := range buf {
			v |= int(b) << (8 * i)
		}
	case BigEndian:
		for _, b := range buf {
			v = v<<8 | int(b)
		}
	}
	return v
}

func serializeInt(value int, order ByteOrder, len int) []byte {
	buf := make([]byte, len)
	switch order {
	case LittleEndian:
		for i := range buf {
			buf[i] = byte(value & 0xFF)
			value >>= 8
		}
	case BigEndian:
		for i := range buf {
			buf[len-1-i] = byte(value & 0xFF)
			value >>= 8
		}
	}
	return buf
}
