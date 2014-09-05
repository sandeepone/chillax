package pipelines

import "github.com/BurntSushi/toml"

type Pipeline struct {
	RunMixin
	Stages []*Stage
}

func NewPipeline(definition string) *Pipeline {
	p := &Pipeline{}

	toml.Decode(definition, p)

	return p
}
