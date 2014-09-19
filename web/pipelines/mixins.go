package pipelines

import (
	"fmt"
	"time"

	"github.com/franela/goreq"
	"github.com/tuxychandru/pubsub"
)

type PipelineAndStageMixin struct {
	goreq.Request

	// Default is "1s"
	TimeoutString string
	Body          map[string]interface{}
	Stages        []*Stage
}

func (rm *PipelineAndStageMixin) MergeBodyToChildrenBody() {
	for _, stage := range rm.Stages {
		if stage.Body == nil {
			stage.Body = rm.Body
		} else {
			for pipelineKey, pipelineValue := range rm.Body {
				if stage.Body[pipelineKey] == nil {
					stage.Body[pipelineKey] = pipelineValue
				}
			}
		}
	}
}

func (rm *PipelineAndStageMixin) NewRunInstance() *RunInstance {
	ri := &RunInstance{}
	ri.TimestampUnixNano = time.Now().UnixNano()
	ri.TimestampUnixNanoString = fmt.Sprintf("%v", ri.TimestampUnixNano)
	ri.PubSub = pubsub.New(1)
	ri.responseChan = ri.PubSub.Sub(ri.TimestampUnixNanoString + "-response")
	ri.errChan = ri.PubSub.Sub(ri.TimestampUnixNanoString + "-error")
	ri.RunInstances = make([]*RunInstance, len(rm.Stages))

	return ri
}

func (rm *PipelineAndStageMixin) Run() *RunInstance {
	runInstance := rm.NewRunInstance()

	go func(responseChan chan interface{}, errChan chan interface{}) {
		resp, err := rm.Do()

		responseChan <- resp
		errChan <- err

	}(runInstance.responseChan, runInstance.errChan)

	for i, stage := range rm.Stages {
		runInstance.RunInstances[i] = stage.Run()
	}

	return runInstance
}

type RunInstance struct {
	TimestampUnixNano       int64
	TimestampUnixNanoString string
	PubSub                  *pubsub.PubSub
	responseChan            chan interface{}
	errChan                 chan interface{}
	RunInstances            []*RunInstance
}
