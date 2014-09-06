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
	RunMixin
	Body   map[string]interface{}
	Stages []*Stage
}

func (p *Pipeline) MergeBodyToStagesBody() {
	for _, stage := range p.Stages {
		if stage.Body == nil {
			stage.Body = p.Body
		} else {
			for pipelineKey, pipelineValue := range p.Body {
				if stage.Body[pipelineKey] == nil {
					stage.Body[pipelineKey] = pipelineValue
				}
			}
		}
	}

}
