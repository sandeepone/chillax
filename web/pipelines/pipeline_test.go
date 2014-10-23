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
	pipeline, _ := NewPipeline(string(definition))
	return pipeline
}

func TestNewPipeline(t *testing.T) {
	pipeline := NewPipelineForTest()

	if pipeline.Method != "POST" {
		t.Error("Default method should be POST.")
	}
	if pipeline.TimeoutString != "1s" {
		t.Error("Default TimeoutString should be 1 second.")
	}
	if pipeline.FailMax != 3 {
		t.Error("Default FailMax should be 3.")
	}
	if pipeline.Accept != "application/json" {
		t.Error("Default Accept should be application/json.")
	}
	if pipeline.ContentType != "application/json" {
		t.Error("Default ContentType should be application/json.")
	}

	if pipeline.Body == nil {
		t.Errorf("Failed to parse pipeline TOML definition. pipeline.Body: %v", pipeline.Body)
	}
	if pipeline.Body["AwsAccessKey"] != "abc" {
		t.Errorf("Failed to parse pipeline TOML definition. pipeline.Body[AwsAccessKey]: %v", pipeline.Body["AwsAccessKey"])
	}
	if pipeline.Body["AwsSecretKey"] != "xyz" {
		t.Errorf("Failed to parse pipeline TOML definition. pipeline.Body[AwsSecretKey]: %v", pipeline.Body["AwsSecretKey"])
	}

	for i, stage := range pipeline.Stages {
		if stage.Body == nil {
			t.Errorf("Pipeline body should be copied to stages if each stage does not have body. stage.Body: %v", stage.Body)
		}
		if stage.Body["AwsAccessKey"] != "abc" {
			t.Errorf("stage.Body is missing a key. stage.Body[AwsAccessKey]: %v", stage.Body["AwsAccessKey"])
		}
		if stage.Body["AwsSecretKey"] != "xyz" {
			t.Errorf("stage.Body is missing a key. stage.Body[AwsSecretKey]: %v", stage.Body["AwsSecretKey"])
		}

		// First stage
		if i == 0 {
			if stage.Body["Token"] != "hahaha" {
				t.Errorf("stage.Body[Token] should not be overriden. stage.Body[Token]: %v", stage.Body["Token"])
			}
		}

		// Last stage
		if i == 1 {
			if stage.Body["Token"] != "lolz" {
				t.Errorf("stage.Body[Token] should not be overriden. stage.Body[Token]: %v", stage.Body["Token"])
			}
		}
	}
}

func TestPipelineSave(t *testing.T) {
	pipeline := NewPipelineForTest()
	err := pipeline.Save()

	if err != nil {
		t.Errorf("Unable to save pipeline. Error: %v", err)
	}
}

func TestNestingStages(t *testing.T) {
	pipeline := NewPipelineForTest()

	stage := pipeline.Stages[1]

	if len(stage.Stages) != 1 {
		t.Errorf("final stage should contain 1 substage. stage.Stages: %v", stage.Stages)
	}

	substage := stage.Stages[0]

	if substage.Uri != "http://localhost:3000/work/step-02-01" {
		t.Errorf("substage.Uri must be correct. substage.Uri: %v", substage.Uri)
	}
}
