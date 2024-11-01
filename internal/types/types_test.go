package types

import "testing"

func TestTypeMeta_Equal(t *testing.T) {
	tm := NewTypeMeta().WithLength(NewConstInt(5))
	tm2 := NewTypeMeta().WithLength(NewConstInt(5))
	if !tm.Equal(tm2) {
		t.Errorf("expected %v, got %v", tm, tm2)
	}
}
