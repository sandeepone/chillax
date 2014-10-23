package pipelines

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Note:
// NewPipelineForTest is defined in pipeline_test.go

func TestRunInstanceParentId(t *testing.T) {
	pipeline := NewPipelineForTest()

	runInstance := pipeline.Run()

	if runInstance.ParentId > 0 {
		t.Errorf("Top most RunInstance must not have ParentId. runInstance.ParentId: %v", runInstance.ParentId)
	}

	for _, childLvl1RunInstance := range runInstance.RunInstances {
		if childLvl1RunInstance.ParentId <= 0 {
			t.Errorf("All RunInstances children must have ParentId")
		}
	}
}

func TestBadRunShouldRecordErrorOnRunInstances(t *testing.T) {
	pipeline := NewPipelineForTest()

	runInstance := pipeline.Run()

	if len(runInstance.RunInstances) != 2 {
		t.Errorf("pipeline.Run should have 2 runInstances. runInstance.RunInstances: %v", runInstance.RunInstances)
	}

	for _, childLvl1RunInstance := range runInstance.RunInstances {
		if childLvl1RunInstance.ErrorMessage == "" {
			t.Errorf("All stages are expected to be broken. Error: %v", childLvl1RunInstance.ErrorMessage)
		}
	}
}

func TestGoodNestedRunsShouldRecordResultsOnRunInstances(t *testing.T) {
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

	// Assert RunInstance of stage[0]
	stage0RunInstance := runInstance.RunInstances[0]

	if stage0RunInstance.Error() != nil {
		t.Errorf("stage0RunInstance should complete successfully. Error: %v", stage0RunInstance.Error)
	}
	if !strings.Contains(stage0RunInstance.ResponseBody, server0Body) {
		t.Errorf("stage0RunInstance received wrong ResponseBodyBytes. ResponseBodyBytes: %v", stage0RunInstance.ResponseBody)
	}

	// Assert RunInstance of stage[1]
	stage1RunInstance := runInstance.RunInstances[1]

	if stage1RunInstance.Error() != nil {
		t.Errorf("stage1RunInstance should complete successfully. Error: %v", stage1RunInstance.Error)
	}
	if !strings.Contains(stage1RunInstance.ResponseBody, server1Body) {
		t.Errorf("stage1RunInstance received wrong ResponseBodyBytes. ResponseBodyBytes: %v", stage1RunInstance.ResponseBody)
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
	if stage10RunInstance.Error() != nil {
		t.Errorf("stage10RunInstance should complete successfully. Error: %v", stage10RunInstance.Error)
	}
	if !strings.Contains(stage10RunInstance.ResponseBody, server10Body) {
		t.Errorf("stage10RunInstance received wrong ResponseBodyBytes. ResponseBodyBytes: %v", stage10RunInstance.ResponseBody)
	}
}
