package types

import (
	"slices"
	"testing"
	"time"
)

func TestTimeStamp_Parse(t *testing.T) {
	d := []byte{0x20, 0x24, 0x10, 0x31, 0x17, 0x32, 0x15}
	ts := TimeStamp{TypeMeta: NewTypeMeta().WithLength(NewConstInt(7))}
	want := time.Date(2024, 10, 31, 17, 32, 15, 0, time.UTC)

	if err := ts.Parse(d); err != nil {
		t.Errorf("expected %v, got %v", nil, err)
	}

	if ts.Value != want {
		t.Errorf("expected %v, got %v", want, ts.Value)
	}
}

func TestTimeStamp_Serialize(t *testing.T) {
	ts := TimeStamp{
		Value:    time.Date(2024, 10, 31, 17, 32, 15, 0, time.UTC),
		TypeMeta: NewTypeMeta().WithLength(NewConstInt(7)),
	}
	want := []byte{0x20, 0x24, 0x10, 0x31, 0x17, 0x32, 0x15}
	result := ts.Serialize()
	if !slices.Equal(result, want) {
		t.Errorf("expected %v, got %v", want, result)
	}
}
