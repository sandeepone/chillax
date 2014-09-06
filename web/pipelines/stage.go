package pipelines

import (
	"time"

	"github.com/franela/goreq"
)

// Create new Stage struct
// Every Stage is capable to make HTTP call.
// By default, the HTTP verb is set to POST and timeout is set to 1 second.
func NewStage(uri string) *Stage {
	stage := &Stage{
		RunMixin: RunMixin{
			goreq.Request{
				Uri:         uri,
				Method:      "POST",
				Timeout:     1 * time.Second,
				Accept:      "application/json",
				ContentType: "application/json",
			},
		},
		Stages: make([]*Stage, 0),
	}
	return stage
}

type Stage struct {
	RunMixin
	Body   map[string]interface{}
	Stages []*Stage
}
