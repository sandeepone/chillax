package pipelines

import (
	"bufio"
	"io/ioutil"
	"os"
	"testing"
)

func NewPipelineForTest() *Pipeline {
	fileHandle, _ := os.Open("./example-pipeline.toml")
	bufReader := bufio.NewReader(fileHandle)
	definition, _ := ioutil.ReadAll(bufReader)
	pipeline := NewPipeline(string(definition))
	return pipeline
}

func TestNewPipeline(t *testing.T) {
	pipeline := NewPipelineForTest()

	if pipeline.Body == nil {
		t.Errorf("Failed to parse pipeline TOML definition. pipeline.Body: %v", pipeline.Body)
	}
	if pipeline.Body["AwsAccessKey"] != "abc" {
		t.Errorf("Failed to parse pipeline TOML definition. pipeline.Body[AwsAccessKey]: %v", pipeline.Body["AwsAccessKey"])
	}
	if pipeline.Body["AwsSecretKey"] != "xyz" {
		t.Errorf("Failed to parse pipeline TOML definition. pipeline.Body[AwsSecretKey]: %v", pipeline.Body["AwsSecretKey"])
	}

	for _, stage := range pipeline.Stages {
		if stage.Body == nil {
			t.Errorf("Pipeline body should be copied to stages if each stage does not have body. stage.Body: %v", stage.Body)
		}
		if stage.Body["AwsAccessKey"] != "abc" {
			t.Errorf("stage.Body is missing a key. stage.Body[AwsAccessKey]: %v", stage.Body["AwsAccessKey"])
		}
		if stage.Body["AwsSecretKey"] != "xyz" {
			t.Errorf("stage.Body is missing a key. stage.Body[AwsSecretKey]: %v", stage.Body["AwsSecretKey"])
		}
	}
}
