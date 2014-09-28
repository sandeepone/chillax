package pipelines

import (
	"bufio"
	"fmt"
	"github.com/didip/chillax/libtime"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestBadNestedRuns(t *testing.T) {
	pipeline := NewPipelineForTest()

	runInstance := pipeline.Run()

	libtime.SleepString("100ms")

	if len(runInstance.RunInstances) != 2 {
		t.Errorf("pipeline.Run should have 2 runInstances. runInstance.RunInstances: %v", runInstance.RunInstances)
	}

	for _, childLvl1RunInstance := range runInstance.RunInstances {
		if childLvl1RunInstance == nil {
			t.Fatalf("Children RunInstances should not be nil.")
		}
		if childLvl1RunInstance.Error == nil {
			t.Errorf("All stages are expected to be broken. Error: %v", childLvl1RunInstance.Error)
		}
	}
}

func TestGoodNestedRuns(t *testing.T) {
	// Setup mock endpoints for stage[0], stage[1], and stage[1[0]]
	server0Body := `{"pick": "me"}`
	server0 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, server0Body)
	}))
	defer server0.Close()

	server1Body := `{"nextParam": "aaa", "result": "success!"}`
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, server1Body)
	}))
	defer server1.Close()

	server10Body := `{"result": "success again!"}`
	server10 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, server10Body)
	}))
	defer server10.Close()

	// Create a new pipeline
	pipeline := NewPipelineForTest()

	// Stub pipeline endpoint URLs with mock endpoints.
	stage0 := pipeline.Stages[0]
	stage0.Uri = server0.URL

	stage1 := pipeline.Stages[1]
	stage1.Uri = server1.URL

	stage10 := pipeline.Stages[1].Stages[0]
	stage10.Uri = server10.URL

	runInstance := pipeline.Run()

	libtime.SleepString("75ms")

	// Assert RunInstance of stage[0]
	stage0RunInstance := runInstance.RunInstances[0]

	if stage0RunInstance.Error != nil {
		t.Errorf("stage0RunInstance should complete successfully. Error: %v", stage0RunInstance.Error)
	}
	if !strings.Contains(string(stage0RunInstance.ResponseBodyBytes), server0Body) {
		t.Errorf("stage0RunInstance received wrong ResponseBodyBytes. ResponseBodyBytes: %v", string(stage0RunInstance.ResponseBodyBytes))
	}

	// Assert RunInstance of stage[1]
	stage1RunInstance := runInstance.RunInstances[1]

	if stage1RunInstance.Error != nil {
		t.Errorf("stage1RunInstance should complete successfully. Error: %v", stage1RunInstance.Error)
	}
	if !strings.Contains(string(stage1RunInstance.ResponseBodyBytes), server1Body) {
		t.Errorf("stage1RunInstance received wrong ResponseBodyBytes. ResponseBodyBytes: %v", string(stage1RunInstance.ResponseBodyBytes))
	}

	// Next stage, which is stage[1[0]], should contain params from stage[1]
	if stage10.Body["Token"] != "lolz" {
		t.Errorf("stage10.Body is incorrect. stage10.Body: %v", stage10.Body)
	}
	if stage10.Body["nextParam"] != "aaa" {
		t.Errorf("stage10.Body is incorrect. stage10.Body: %v", stage10.Body)
	}
	if stage10.Body["result"] != "success!" {
		t.Errorf("stage10.Body is incorrect. stage10.Body: %v", stage10.Body)
	}

	// Assert RunInstance of stage[1[0]]
	stage10RunInstance := runInstance.RunInstances[1].RunInstances[0]
	if stage10RunInstance.Error != nil {
		t.Errorf("stage10RunInstance should complete successfully. Error: %v", stage10RunInstance.Error)
	}
	if !strings.Contains(string(stage10RunInstance.ResponseBodyBytes), server10Body) {
		t.Errorf("stage10RunInstance received wrong ResponseBodyBytes. ResponseBodyBytes: %v", string(stage10RunInstance.ResponseBodyBytes))
	}
}
