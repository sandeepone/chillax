package libmetrics

import (
	"testing"
)

func TestParseLoadAverages(t *testing.T) {
	data := ParseLoadAverages()
	if len(data) != 3 {
		t.Errorf("ParseLoadAverages should return 3 load averages value. data: %v", data)
	}
	for _, datum := range data {
		if datum <= 0 {
			t.Errorf("Load average should never be below 0. datum: %v", datum)
		}
	}
}

func TestNewCpuMetrics(t *testing.T) {
	c, err := NewCpuMetrics()
	if err != nil {
		t.Errorf("Should be able to get CPU metrics. err: %v", err)
	}

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

func TestCpuMetricsSerialization(t *testing.T) {
	data, err := NewCpuMetrics()
	if err != nil {
		t.Errorf("Should be able to get CPU metrics. err: %v", err)
	}

	_, err = data.ToJson()
	if err != nil {
		t.Errorf("Serializing to JSON should not break. err: %v", err)
	}

	_, err = data.ToToml()
	if err != nil {
		t.Errorf("Serializing to TOML should not break. err: %v", err)
	}
}
