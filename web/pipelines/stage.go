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
		PipelineAndStageMixin: PipelineAndStageMixin{
			Request: goreq.Request{
				Uri:         uri,
				Method:      "POST",
				Timeout:     1 * time.Second,
				Accept:      "application/json",
				ContentType: "application/json",
			},
			Body: make(map[string]interface{}),
		},
	}

	stage.MergeBodyToStagesBody()

	return stage
}

type Stage struct {
	PipelineAndStageMixin
}
