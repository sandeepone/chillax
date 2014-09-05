package pipelines

import (
	"fmt"
	"time"

	"github.com/franela/goreq"
	"github.com/tuxychandru/pubsub"
)

// Create new Stage struct
// Every Stage is capable to make HTTP call.
// By default, the HTTP verb is set to POST and timeout is set to 1 second.
func NewStage(uri string) *Stage {
	stage := &Stage{
		Request: goreq.Request{
			Uri:     uri,
			Method:  "POST",
			Timeout: 1 * time.Second,
		},
		Stages: make([]*Stage, 0),
	}
	return stage
}

type Stage struct {
	goreq.Request
	Stages []*Stage
}

func (s *Stage) NewStageRun() *StageRun {
	sr := &StageRun{}
	sr.TimestampUnixNano = time.Now().UnixNano()
	sr.TimestampUnixNanoString = fmt.Sprintf("%v", sr.TimestampUnixNano)
	sr.PubSub = pubsub.New(1)
	return sr
}

func (s *Stage) Run() (chan interface{}, chan interface{}) {
	sr := s.NewStageRun()

	responseChan := sr.PubSub.Sub(sr.TimestampUnixNanoString + "-response")
	errChan := sr.PubSub.Sub(sr.TimestampUnixNanoString + "-error")

	go func(responseChan chan interface{}, errChan chan interface{}) {
		resp, err := s.Do()

		responseChan <- resp
		errChan <- err
	}(responseChan, errChan)

	return responseChan, errChan
}

type StageRun struct {
	TimestampUnixNano       int64
	TimestampUnixNanoString string
	PubSub                  *pubsub.PubSub
}
