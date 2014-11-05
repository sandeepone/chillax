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

func TestVmstatDarwinSerialization(t *testing.T) {
	data := NewVmstatDarwin()
	_, err := data.ToJson()
	if err != nil {
		t.Errorf("Serializing to JSON should not break. err: %v", err)
	}

	_, err = data.ToToml()
	if err != nil {
		t.Errorf("Serializing to TOML should not break. err: %v", err)
	}
}
