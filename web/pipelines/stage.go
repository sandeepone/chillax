package pipelines

import (
    "time"
    "github.com/franela/goreq"
)

type Stage struct {
    LastSuccessUnixNano int64
}

func (s *Stage) NewRequest(uri string) *goreq.Request {
    return &goreq.Request{
        Uri:     uri,
        Method:  "POST",
        Timeout: 1 * time.Second,
    }
}

func (s *Stage) Run(req *goreq.Request) {
    _, err := req.Do()
    if err == nil {
        s.LastSuccessUnixNano = time.Now().UnixNano()
    } else {

    }
}