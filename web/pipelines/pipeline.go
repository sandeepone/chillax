package pipelines

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/franela/goreq"
)

type Pipeline struct {
	RunMixin
	Stages []*Stage
}

func NewPipeline(definition string) *Pipeline {
	p := &Pipeline{
		RunMixin: RunMixin{
			goreq.Request{
				Method:  "POST",
				Timeout: 1 * time.Second,
			},
		},
	}

	toml.Decode(definition, p)

	return p
}
