package types

import (
	"slices"
	"testing"
)

func TestConstType_Parse(t *testing.T) {
	d := []byte{42}
	c := ConstType{
		Value:    42,
		TypeMeta: NewTypeMeta().WithLength(NewConstInt(1)),
	}

	if err := c.Parse(d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if c.Value != 42 {
		t.Errorf("expected value 42, got %d", c.Value)
	}
}

func TestConstType_Serialize(t *testing.T) {
	c := ConstType{Value: 42, TypeMeta: NewTypeMeta().WithLength(NewConstInt(1))}
	want := []byte{42}
	result := c.Serialize()
	if !slices.Equal(result, want) {
		t.Errorf("expected %v, got %v", want, result)
	}
}
