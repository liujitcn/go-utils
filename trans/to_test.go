package trans

import "testing"

func TestPtr(t *testing.T) {
}

func TestSliceOfPtrs(t *testing.T) {
	arr := SliceOfPtrs[int]()
	if len(arr) != 0 {
		t.Fatal("expected zero length")
	}
	arr = SliceOfPtrs(1, 2, 3, 4, 5)
	for i, v := range arr {
		if v != i+1 {
			t.Fatal("values don't match")
		}
	}
}
