// +build darwin

package libmetrics

import (
	"testing"
)

func TestNewVmstatDarwin(t *testing.T) {
	data := NewVmstatDarwin()
	if len(data) <= 0 {
		t.Errorf("Getting value from vm_stat should work. data: %v", data)
	}
	for key, datum := range data {
		if datum < 0 {
			t.Errorf("Data inside vmstat should never be < 0. key: %v, datum: %v", key, datum)
		}
	}
}
