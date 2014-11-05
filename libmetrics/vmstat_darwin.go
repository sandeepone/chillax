// +build darwin

package libmetrics

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"os/exec"
	"strconv"
	"strings"
)

func NewVmstatDarwin() VmstatDarwin {
	v := make(VmstatDarwin)

	out, _ := exec.Command("vm_stat").Output()

	scanner := bufio.NewScanner(bytes.NewReader(out))

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "Mach Virtual Memory Statistics") {
			keyAndValue := strings.Split(line, ":")
			key := strings.TrimSpace(keyAndValue[0])
			value := strings.Replace(strings.TrimSpace(keyAndValue[1]), ".", "", -1)
			valueInt64, err := strconv.ParseInt(value, 10, 64)

			if err == nil {
				// Convert value from pages to bytes
				v[key] = valueInt64 * 4096
			}
		}
	}

	return v
}

type VmstatDarwin map[string]int64

func (v *VmstatDarwin) ToJson() ([]byte, error) {
	return json.Marshal(v)
}

func (v *VmstatDarwin) ToToml() ([]byte, error) {
	var buffer bytes.Buffer

	err := toml.NewEncoder(&buffer).Encode(v)

	return buffer.Bytes(), err
}
