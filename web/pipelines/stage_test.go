package pipelines

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func NewStageForTest() *Stage {
	fileHandle, _ := os.Open("./example-pipeline.toml")
	bufReader := bufio.NewReader(fileHandle)
	definition, _ := ioutil.ReadAll(bufReader)
	pipeline := NewPipeline(string(definition))
	return pipeline.Stages[0]
}

func TestNewStage(t *testing.T) {
	stage := NewStageForTest()

	if stage.Method != "POST" {
		t.Error("Default method should be POST.")
	}
	if stage.Timeout != 1*time.Second {
		t.Error("Default timeout should be 1 second.")
	}
}

func TestStageRunBadRequest(t *testing.T) {
	stage := NewStageForTest()
	stage.Timeout = 1 * time.Millisecond

	runInstance := stage.Run()

	err := <-runInstance.errChan

	if err == nil {
		t.Error("Request should fail.")
	}
}
