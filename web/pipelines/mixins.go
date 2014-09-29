package pipelines

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	chillax_storage "github.com/didip/chillax/storage"
	"github.com/franela/goreq"
	"io/ioutil"
	"sync"
	"time"
)

type PipelineAndStageMixin struct {
	goreq.Request

	// Default is "1s"
	TimeoutString string
	Body          map[string]interface{}
	Stages        []*Stage
}

func (mixin *PipelineAndStageMixin) MergeBodyToChildrenBody() {
	for _, stage := range mixin.Stages {
		if stage.Body == nil {
			stage.Body = mixin.Body
		} else {
			for pipelineKey, pipelineValue := range mixin.Body {
				if stage.Body[pipelineKey] == nil {
					stage.Body[pipelineKey] = pipelineValue
				}
			}
		}
	}
}

func (mixin *PipelineAndStageMixin) Run() RunInstance {
	var err error

	runInstance := mixin.NewRunInstance()

	if mixin.Uri != "" {
		response, err := mixin.Do()

		if err == nil && response != nil && response.Body != nil {
			responseBytes, err := ioutil.ReadAll(response.Body)

			if err == nil && len(responseBytes) > 0 {
				runInstance.ResponseBody = string(responseBytes)
			}
		}
		if err != nil {
			runInstance.ErrorMessage = err.Error()
		}
	}

	if err == nil && len(mixin.Stages) > 0 {
		runInstance.RunInstances = make([]RunInstance, len(mixin.Stages))

		var wg sync.WaitGroup

		for i, stage := range mixin.Stages {
			wg.Add(1)

			go func(runInstance RunInstance, i int, stage *Stage) {
				defer wg.Done()

				// Merge the JSON body of previous stage to next stage.
				if len(runInstance.ResponseBody) > 0 {
					json.Unmarshal([]byte(runInstance.ResponseBody), &stage.Body)
				}

				runInstance.RunInstances[i] = stage.Run()
			}(runInstance, i, stage)
		}
		wg.Wait()
	}

	runInstance.Save()

	return runInstance
}

func (mixin *PipelineAndStageMixin) NewRunInstance() RunInstance {
	ri := RunInstance{}
	ri.TimestampUnixNano = time.Now().UnixNano()

	return ri
}

type RunInstance struct {
	TimestampUnixNano int64
	ResponseBody      string
	ErrorMessage      string
	RunInstances      []RunInstance
}

func (ri *RunInstance) Error() error {
	if ri.ErrorMessage != "" {
		return errors.New(ri.ErrorMessage)
	}
	return nil
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
	return chillax_storage.NewStorage().Create(fmt.Sprintf("/pipelines/run-instances/%v", ri.TimestampUnixNano), inBytes)
}
