package pipelines

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	chillax_storage "github.com/didip/chillax/storage"
	"time"
)

// Create a new Pipeline struct.
// It receive TOML definition as string.
func NewPipeline(definition string) (*Pipeline, error) {
	p := &Pipeline{}

	_, err := toml.Decode(definition, p)

	p.SetDefaults()

	return p, err
}

type Pipeline struct {
	PipelineAndStageMixin
	Id int64
}

type PipelineSerializable struct {
	PipelineAndStageSerializableMixin
	Id int64
}

func (p *Pipeline) SetDefaults() {
	if p.TimeoutString == "" {
		p.TimeoutString = "1s"
	}
	p.Id = time.Now().UnixNano()
	p.SetStagesDefaults()
	p.MergeBodyToChildrenBody()
}

func (p *Pipeline) SetStagesDefaults() {
	for _, stage := range p.Stages {
		stage.SetDefaults()
	}
}

func (p *Pipeline) Serialize() ([]byte, error) {
	serializable := &PipelineSerializable{}
	serializable.Id = p.Id
	serializable.TimeoutString = p.TimeoutString
	serializable.Body = p.Body
	serializable.Stages = make([]*StageSerializable, len(p.Stages))

	for i, stage := range p.Stages {
		stageSerializable := &StageSerializable{}
		stageSerializable.TimeoutString = stage.TimeoutString
		stageSerializable.Body = stage.Body

		serializable.Stages[i] = stageSerializable
	}

	var buffer bytes.Buffer
	err := toml.NewEncoder(&buffer).Encode(serializable)

	return buffer.Bytes(), err
}

func (p *Pipeline) Save() error {
	inBytes, err := p.Serialize()
	if err != nil {
		return err
	}
	return chillax_storage.NewStorage().Update(fmt.Sprintf("/pipelines/%v", p.Id), inBytes)
}
