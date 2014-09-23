package pipelines

import (
	"encoding/json"
	"fmt"
	"github.com/franela/goreq"
	"io/ioutil"
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

func (mixin *PipelineAndStageMixin) Run() *RunInstance {
	runInstance := mixin.NewRunInstance()

	go func() {
		if mixin.Uri != "" {
			runInstance.Response, runInstance.Error = mixin.Do()

			if runInstance.Error == nil && runInstance.Response != nil {
				runInstance.ResponseBodyBytes, runInstance.Error = ioutil.ReadAll(runInstance.Response.Body)
			}
		}

		if runInstance.Error == nil && len(mixin.Stages) > 0 {
			runInstancesChan := make(chan *RunInstance)

			for _, stage := range mixin.Stages {
				go func() {
					// Merge the JSON body of previous stage to next stage.
					if runInstance.ResponseBodyBytes != nil {
						json.Unmarshal(runInstance.ResponseBodyBytes, stage.Body)
					}

					runInstancesChan <- stage.Run()
				}()
			}

			for childRunInstance := range runInstancesChan {
				runInstance.RunInstances = append(runInstance.RunInstances, childRunInstance)
			}
			close(runInstancesChan)
		}
	}()

	return runInstance
}

func (mixin *PipelineAndStageMixin) NewRunInstance() *RunInstance {
	ri := &RunInstance{}
	ri.TimestampUnixNano = time.Now().UnixNano()
	ri.TimestampUnixNanoString = fmt.Sprintf("%v", ri.TimestampUnixNano)
	ri.RunInstances = make([]*RunInstance, 0)

	return ri
}

type RunInstance struct {
	TimestampUnixNano       int64
	TimestampUnixNanoString string
	Response                *goreq.Response
	ResponseBodyBytes       []byte
	Error                   error
	RunInstances            []*RunInstance
}
