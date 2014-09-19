package pipelines

import "github.com/BurntSushi/toml"

// Create a new Pipeline struct.
// It receive TOML definition as string.
func NewPipeline(definition string) *Pipeline {
	p := &Pipeline{}

	toml.Decode(definition, p)

	p.SetDefaults()

	return p
}

type Pipeline struct {
	PipelineAndStageMixin
}

func (p *Pipeline) SetDefaults() {
	if p.TimeoutString == "" {
		p.TimeoutString = "1s"
	}

	p.SetStagesDefaults()
	p.MergeBodyToStagesBody()
}

func (p *Pipeline) SetStagesDefaults() {
	for _, stage := range p.Stages {
		stage.SetDefaults()
	}
}

// Run all stages under 1 pipeline.
// Returns slice of tuple: outChan and errChan
func (p *Pipeline) RunStages() [][]chan interface{} {
	stagesChans := make([][]chan interface{}, len(p.Stages))

	for i, stage := range p.Stages {
		outChan, errChan := stage.Run()

		stagesChans[i] = make([]chan interface{}, 2)
		stagesChans[i][0] = outChan
		stagesChans[i][1] = errChan
	}

	return stagesChans
}
