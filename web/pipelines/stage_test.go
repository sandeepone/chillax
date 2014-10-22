package pipelines

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

func NewStageForTest() *Stage {
	fileHandle, _ := os.Open("./example-pipeline.toml")
	bufReader := bufio.NewReader(fileHandle)
	definition, _ := ioutil.ReadAll(bufReader)
	pipeline, _ := NewPipeline(string(definition))
	return pipeline.Stages[0]
}

func TestNewStage(t *testing.T) {
	stage := NewStageForTest()

	if stage.Method != "POST" {
		t.Error("Default method should be POST.")
	}
	if stage.TimeoutString != "1s" {
		t.Error("Default TimeoutString should be 1 second.")
	}
	if stage.FailMax != 3 {
		t.Error("Default FailMax should be 3.")
	}
	if stage.Accept != "application/json" {
		t.Error("Default Accept should be application/json.")
	}
	if stage.ContentType != "application/json" {
		t.Error("Default ContentType should be application/json.")
	}
}
