package pipelines

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	chillax_storage "github.com/didip/chillax/storage"
	"time"
)

type RunInstance struct {
	Id           int64
	ParentId     int64
	ResponseBody string
	ErrorMessage string
	RunInstances []RunInstance
}

func (ri *RunInstance) Error() error {
	if ri.ErrorMessage != "" {
		return errors.New(ri.ErrorMessage)
	}
	return nil
}

func (ri *RunInstance) HasErrorsRecursively() bool {
	if ri.ErrorMessage != "" {
		return true
	} else {
		for _, child := range ri.RunInstances {
			if child.HasErrorsRecursively() {
				return true
			}
		}
	}
	return false
}

func (ri *RunInstance) Serialize() ([]byte, error) {
	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(ri)

	return buffer.Bytes(), err
}

func (ri *RunInstance) Save() error {
	inBytes, err := ri.Serialize()
	if err != nil {
		return err
	}

	datetime := time.Unix(0, ri.Id)

	dataPath := fmt.Sprintf(
		"/logs/pipelines/run-instances/%v/%d/%v/%v/%v/%v",
		datetime.Year(), datetime.Month(), datetime.Day(), datetime.Hour(), datetime.Minute(), ri.Id)

	return chillax_storage.NewStorage().Create(dataPath, inBytes)
}
