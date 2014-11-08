package libmetrics

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func ParseUptime() []float64 {
	out, _ := exec.Command("uptime").Output()

	chunks := strings.Split(string(out[:]), ",")
	loadAvgs := chunks[len(chunks)-1]
	loadAvgChunks := strings.Split(loadAvgs, ":")
	loadAvgValStr := strings.Fields(loadAvgChunks[len(loadAvgChunks)-1])

	loadAvg1, _ := strconv.ParseFloat(loadAvgValStr[0], 64)
	loadAvg2, _ := strconv.ParseFloat(loadAvgValStr[1], 64)
	loadAvg3, _ := strconv.ParseFloat(loadAvgValStr[2], 64)

	return []float64{loadAvg1, loadAvg2, loadAvg3}
}

func NewCpuMetrics() *CpuMetrics {
	c := &CpuMetrics{}
	c.LoadAverages = ParseUptime()
	c.NumCpu = runtime.NumCPU()
	c.LoadAveragesPerCpu = c.GetLoadAveragesPerCpu()
	return c
}

type CpuMetrics struct {
	LoadAverages       []float64
	NumCpu             int
	LoadAveragesPerCpu []float64
}

func (c *CpuMetrics) GetLoadAveragesPerCpu() []float64 {
	lapc := make([]float64, len(c.LoadAverages))

	for i, load := range c.LoadAverages {
		lapc[i] = load / float64(c.NumCpu)
	}
	return lapc
}

func (c *CpuMetrics) ToJson() ([]byte, error) {
	return json.Marshal(c)
}

func (c *CpuMetrics) ToToml() ([]byte, error) {
	var buffer bytes.Buffer

	err := toml.NewEncoder(&buffer).Encode(c)

	return buffer.Bytes(), err
}