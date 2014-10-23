package pipelines

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	chillax_storage "github.com/didip/chillax/storage"
	"os"
	"sync"
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

func NewPipelineGivenJsonBytes(definition []byte) (*Pipeline, error) {
	p := &Pipeline{}

	err := json.Unmarshal(definition, p)

	p.SetDefaults()

	return p, err
}

func PipelineById(id string) (*Pipeline, error) {
	definition, err := chillax_storage.NewStorage().Get("/pipelines/" + id)
	if err != nil {
		return nil, err
	}

	return NewPipeline(string(definition))
}

func AllPipelinesByIds(pipelineIds []string) ([]*Pipeline, error) {
	pipelines := make([]*Pipeline, len(pipelineIds))

	for index, pipelineId := range pipelineIds {
		pipeline, err := PipelineById(pipelineId)
		if err != nil {
			return make([]*Pipeline, 0), err
		}
		pipelines[index] = pipeline
	}

	return pipelines, nil
}

func AllPipelines() ([]*Pipeline, error) {
	pipelineIds, err := chillax_storage.NewStorage().List("/pipelines")
	if err != nil {
		return make([]*Pipeline, 0), err
	}

	return AllPipelinesByIds(pipelineIds)
}

func AllInProgressPipelines() ([]*Pipeline, error) {
	hostname, _ := os.Hostname()

	pipelineIds, err := chillax_storage.NewStorage().List(fmt.Sprintf("/pipelines-inprogress/%v", hostname))
	if err != nil {
		return make([]*Pipeline, 0), err
	}

	return AllPipelinesByIds(pipelineIds)
}

func RunAllInProgressPipelinesAsync(numGoroutines int) {
	go func() {
		pipelines, err := AllInProgressPipelines()

		if err == nil && len(pipelines) > 0 {
			pipelinesChan := make(chan *Pipeline, len(pipelines))
			wg := new(sync.WaitGroup)

			for i := 0; i < numGoroutines; i++ {
				wg.Add(1)

				go func(pipelinesChan chan *Pipeline, wg *sync.WaitGroup) {
					defer wg.Done()

					for pipeline := range pipelinesChan {
						pipeline.RunWithCrashRequeue()
					}
				}(pipelinesChan, wg)
			}

			for _, pipeline := range pipelines {
				pipelinesChan <- pipeline
			}
		}
	}()
}

type Pipeline struct {
	PipelineAndStageMixin
	Id          int64
	Description string
}

type PipelineSerializable struct {
	PipelineAndStageSerializableMixin
	Id          int64
	Description string
}

func (p *Pipeline) SetDefaults() {
	p.Id = time.Now().UnixNano()

	p.SetStagesDefaults()
	p.SetCommonDefaults()
}

func (p *Pipeline) SetStagesDefaults() {
	for _, stage := range p.Stages {
		stage.SetCommonDefaults()
	}
}

func (p *Pipeline) Serialize() ([]byte, error) {
	serializable := &PipelineSerializable{}
	serializable.Id = p.Id
	serializable.Description = p.Description
	serializable.TimeoutString = p.TimeoutString
	serializable.FailCount = p.FailCount
	serializable.FailMax = p.FailMax
	serializable.Body = p.Body
	serializable.Stages = make([]*StageSerializable, len(p.Stages))

	for i, stage := range p.Stages {
		stageSerializable := &StageSerializable{}
		stageSerializable.TimeoutString = stage.TimeoutString
		stageSerializable.FailCount = stage.FailCount
		stageSerializable.FailMax = stage.FailMax
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

func (p *Pipeline) SaveInProgress() error {
	inBytes, err := p.Serialize()
	if err != nil {
		return err
	}

	hostname, _ := os.Hostname()

	return chillax_storage.NewStorage().Update(fmt.Sprintf("/pipelines-inprogress/%v/%v", hostname, p.Id), inBytes)
}

func (p *Pipeline) DeleteInProgress() error {
	hostname, _ := os.Hostname()

	return chillax_storage.NewStorage().Delete(fmt.Sprintf("/pipelines-inprogress/%v/%v", hostname, p.Id))
}

func (p *Pipeline) RunWithCrashRequeue() RunInstance {
	go func(p *Pipeline) {
		p.SaveInProgress()
	}(p)

	runInstance := p.Run()

	if runInstance.HasPerformedRecursively() {
		p.DeleteInProgress()
	}

	return runInstance
}
