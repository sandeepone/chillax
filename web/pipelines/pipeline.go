package pipelines

import "github.com/BurntSushi/toml"

// Create a new Pipeline struct.
// It receive TOML definition as string.
func NewPipeline(definition string) *Pipeline {
	p := &Pipeline{}

	toml.Decode(definition, p)

	p.MergeBodyToStagesBody()

	return p
}

type Pipeline struct {
	PipelineAndStageMixin
}
