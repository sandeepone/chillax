package pipelines

import (
	"bufio"
	"github.com/didip/chillax/libtime"
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

func TestNestedRuns(t *testing.T) {
	pipeline := NewPipelineForTest()

	runInstance := pipeline.Run()

	libtime.SleepString("1s")

	if len(runInstance.RunInstances) != 2 {
		t.Errorf("pipeline.Run should have 2 runInstances. runInstance.RunInstances: %v", runInstance.RunInstances)
	}

	for _, childLvl1RunInstance := range runInstance.RunInstances {
		if childLvl1RunInstance == nil {
			t.Fatalf("Children RunInstances should not be nil.")
		}
	}

	for _, childLvl1RunInstance := range runInstance.RunInstances {
		errChan := childLvl1RunInstance.SubError()

		err := <-errChan

		if err == nil {
			t.Errorf("All stages are expected to be broken. Error: %v", err)
		}
	}
}
