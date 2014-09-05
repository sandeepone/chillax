package libnumber

import (
	"testing"
)

func TestLargestInt(t *testing.T) {
	if LargestInt([]int{44, 2, 21}) != 44 {
		t.Errorf("Failed to get largest int correctly")
	}
}

func TestFirstGapIntSlice(t *testing.T) {
	if FirstGapIntSlice([]int{4, 2, 1}) != 3 {
		t.Errorf("Failed to get first gap in slice: 4,2,1")
	}
	if FirstGapIntSlice([]int{2, 99, 1, 44}) != 98 {
		t.Errorf("Failed to get first gap in slice: 99,44,2,1")
	}
}
