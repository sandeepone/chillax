package libmetrics

import (
	"testing"
)

func TestParseUptime(t *testing.T) {
	data := ParseUptime()
	if len(data) != 3 {
		t.Errorf("ParseUptime should return 3 load averages value. data: %v", data)
	}
	for _, datum := range data {
		if datum <= 0 {
			t.Errorf("Load average should never be below 0. datum: %v", datum)
		}
	}
}

func TestNewCpuMetrics(t *testing.T) {
	c := NewCpuMetrics()

	for _, datum := range c.LoadAverages {
		if datum <= 0 {
			t.Errorf("Load average should never be below 0. datum: %v", datum)
		}
	}
	if c.NumCpu <= 0 {
		t.Errorf("NumCpu should never be below 0. c.NumCpu: %v", c.NumCpu)
	}
	for _, datum := range c.LoadAveragesPerCpu {
		if datum <= 0 {
			t.Errorf("Load average per CPU should never be below 0. datum: %v", datum)
		}
	}
}
