package types

import (
	"slices"
	"testing"
)

func TestStringType_Parse(t *testing.T) {
	d := []byte{104, 101, 108, 108, 111}
	s := StringType{TypeMeta: NewTypeMeta().WithLength(NewConstInt(5))}
	want := "hello"

	if err := s.Parse(d); err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}

	if s.Value != want {
		t.Errorf("expected %v, got %v", want, s.Value)
	}
}

func TestStringType_Serialize(t *testing.T) {
	s := StringType{Value: "hello", TypeMeta: NewTypeMeta().WithLength(NewConstInt(5))}
	result := s.Serialize()
	want := []byte{104, 101, 108, 108, 111}
	if !slices.Equal(result, want) {
		t.Errorf("expected %v, got %v", want, result)
	}
}
