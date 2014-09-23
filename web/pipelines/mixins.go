package pipelines

import (
	"fmt"
	"time"

	"github.com/franela/goreq"
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
		}

		runInstancesChan := make(chan *RunInstance)

		for _, stage := range mixin.Stages {
			go func() {
				if runInstance.Error == nil {
					runInstancesChan <- stage.Run()
				}
			}()
		}

		for childRunInstance := range runInstancesChan {
			runInstance.RunInstances = append(runInstance.RunInstances, childRunInstance)
		}
		close(runInstancesChan)
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
	Error                   error
	RunInstances            []*RunInstance
}
