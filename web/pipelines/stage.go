package pipelines

import (
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

type StageRun struct {
    TimestampUnixNano int64
}

type Stage struct {
    *goreq.Request
}

func (s *Stage) NewStageRun() *StageRun {
    sr := &StageRun{}
    sr.TimestampUnixNano = time.Now().UnixNano()
    return sr
}

func (s *Stage) Run() (*StageRun, error) {
    sr := s.NewStageRun()

    _, err := s.Do()
    if err == nil {
    } else {
    }
    return sr, err
}

func (s *Stage) NewPubSub(bufferSize int) *pubsub.PubSub {
    return pubsub.New(bufferSize)
}
