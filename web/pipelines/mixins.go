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

func (mixin *PipelineAndStageMixin) Run() *RunInstance {
	runInstance := mixin.NewRunInstance()

	go func() {
		if mixin.Uri != "" {
			resp, err := mixin.Do()

			runInstance.PubResponse(resp)
			runInstance.PubError(err)
		}

		runInstancesChan := make(chan *RunInstance)

		for _, stage := range mixin.Stages {
			go func() {
				if mixin.Uri == "" {
					runInstancesChan <- stage.Run()
				} else {
					var parentErr interface{}

					parentErr = <-runInstance.SubError()
					parentErr = parentErr.(error)

					if parentErr == nil {
						runInstancesChan <- stage.Run()
					}
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
	ri.PubSub = pubsub.New(1)
	ri.RunInstances = make([]*RunInstance, 0)

	return ri
}

type RunInstance struct {
	TimestampUnixNano       int64
	TimestampUnixNanoString string
	PubSub                  *pubsub.PubSub
	RunInstances            []*RunInstance
}

func (ri *RunInstance) ResponseTopic() string {
	return ri.TimestampUnixNanoString + "-response"
}

func (ri *RunInstance) ErrorTopic() string {
	return ri.TimestampUnixNanoString + "-error"
}

func (ri *RunInstance) PubResponse(data interface{}) {
	ri.PubSub.Pub(data, ri.ResponseTopic())
}

func (ri *RunInstance) PubError(data error) {
	ri.PubSub.Pub(data, ri.ErrorTopic())
}

func (ri *RunInstance) SubResponse() chan interface{} {
	return ri.PubSub.Sub(ri.ResponseTopic())
}

func (ri *RunInstance) SubError() chan interface{} {
	return ri.PubSub.Sub(ri.ErrorTopic())
}
