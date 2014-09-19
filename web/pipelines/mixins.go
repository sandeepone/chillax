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

func (mixin *PipelineAndStageMixin) NewRunInstance() *RunInstance {
	ri := &RunInstance{}
	ri.TimestampUnixNano = time.Now().UnixNano()
	ri.TimestampUnixNanoString = fmt.Sprintf("%v", ri.TimestampUnixNano)
	ri.PubSub = pubsub.New(1)
	ri.responseChan = ri.PubSub.Sub(ri.TimestampUnixNanoString + "-response")
	ri.errChan = ri.PubSub.Sub(ri.TimestampUnixNanoString + "-error")
	ri.RunInstances = make([]*RunInstance, len(mixin.Stages))

	return ri
}

func (mixin *PipelineAndStageMixin) Run() *RunInstance {
	runInstance := mixin.NewRunInstance()

	go func(runInstance *RunInstance) {
		resp, err := mixin.Do()

		runInstance.responseChan <- resp
		runInstance.errChan <- err

	}(runInstance)

	for i, stage := range mixin.Stages {
		go func(runInstance *RunInstance) {
			var parentErr interface{}

			parentErr = <-runInstance.errChan

			if parentErr == nil {
				runInstance.RunInstances[i] = stage.Run()
			}
		}(runInstance)
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
