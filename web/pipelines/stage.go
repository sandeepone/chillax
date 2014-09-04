package pipelines

import (
    "fmt"
    "time"
    "github.com/franela/goreq"
    "github.com/tuxychandru/pubsub"
)

func NewStage(uri string) *Stage {
    return &Stage{&goreq.Request{
            Uri:     uri,
            Method:  "POST",
            Timeout: 1 * time.Second,
        },
    }
}

type Stage struct {
    *goreq.Request
}

func (s *Stage) NewStageRun() *StageRun {
    sr := &StageRun{}
    sr.TimestampUnixNano       = time.Now().UnixNano()
    sr.TimestampUnixNanoString = fmt.Sprintf("%v", sr.TimestampUnixNano)
    sr.PubSub                  = pubsub.New(1)
    return sr
}

func (s *Stage) Run() (chan interface{}, chan interface{}) {
    sr := s.NewStageRun()

    responseChan := sr.PubSub.Sub(sr.TimestampUnixNanoString + "-response")
    errChan      := sr.PubSub.Sub(sr.TimestampUnixNanoString + "-error")

    go func(responseChan chan interface{}, errChan chan interface{}){
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
