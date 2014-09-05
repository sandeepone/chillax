package pipelines

import (
	"fmt"
	"time"

	"github.com/franela/goreq"
	"github.com/tuxychandru/pubsub"
)

type RunMixin struct {
	goreq.Request
}

func (rm *RunMixin) NewRunInstance() *RunInstance {
	ri := &RunInstance{}
	ri.TimestampUnixNano = time.Now().UnixNano()
	ri.TimestampUnixNanoString = fmt.Sprintf("%v", ri.TimestampUnixNano)
	ri.PubSub = pubsub.New(1)
	return ri
}

func (rm *RunMixin) Run() (chan interface{}, chan interface{}) {
	sr := rm.NewRunInstance()

	responseChan := sr.PubSub.Sub(sr.TimestampUnixNanoString + "-response")
	errChan := sr.PubSub.Sub(sr.TimestampUnixNanoString + "-error")

	go func(responseChan chan interface{}, errChan chan interface{}) {
		resp, err := rm.Do()

		responseChan <- resp
		errChan <- err
	}(responseChan, errChan)

	return responseChan, errChan
}

type RunInstance struct {
	TimestampUnixNano       int64
	TimestampUnixNanoString string
	PubSub                  *pubsub.PubSub
}
