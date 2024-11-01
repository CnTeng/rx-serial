package types

import (
	"slices"
	"testing"
)

func TestIntType_Parse(t *testing.T) {
	d := []byte{0x03, 0x45, 0x2F, 0x1B}
	i := IntType{
		TypeMeta: NewTypeMeta().
			WithLength(NewConstInt(4)).
			WithConfig(TypeConfig{ByteOrder: BigEndian}),
	}
	bigWant := 0x03452F1B
	littleWant := 0x1B2F4503

	if err := i.Parse(d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if i.Value != bigWant {
		t.Errorf("expected value 0x%x, got 0x%x", bigWant, i.Value)
	}

	i.Config.ByteOrder = LittleEndian
	if err := i.Parse(d); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if i.Value != littleWant {
		t.Errorf("expected value 0x%x, got 0x%x", littleWant, i.Value)
	}
}

func TestIntType_Serialize(t *testing.T) {
	i := IntType{
		Value: 0x03452F1B,
		TypeMeta: NewTypeMeta().
			WithLength(NewConstInt(4)).
			WithConfig(TypeConfig{ByteOrder: BigEndian}),
	}
	bigWant := []byte{0x03, 0x45, 0x2F, 0x1B}
	littleWant := []byte{0x1B, 0x2F, 0x45, 0x03}

	bigResult := i.Serialize()
	if !slices.Equal(bigResult, bigWant) {
		t.Errorf("expected %v, got %v", bigWant, bigResult)
	}

	i.Config.ByteOrder = LittleEndian
	littleResult := i.Serialize()
	if !slices.Equal(littleResult, littleWant) {
		t.Errorf("expected %v, got %v", littleWant, littleResult)
	}
}
