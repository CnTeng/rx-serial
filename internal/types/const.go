package types

import "fmt"

type ConstType struct {
	Value byte
	*TypeMeta
}

func (c *ConstType) Parse(data []byte) error {
	b := data[:c.Len()]
	if c.Value != b[0] {
		return fmt.Errorf("value mismatch: %d != %d", c.Value, b[0])
	}
	return nil
}

func (c *ConstType) Serialize() []byte {
	return []byte{c.Value}
}

func (c *ConstType) Equal(other Type) bool {
	if oc, ok := other.(*ConstType); !ok {
		return false
	} else {
		return c.Value == oc.Value && c.TypeMeta.Equal(oc.TypeMeta)
	}
}
