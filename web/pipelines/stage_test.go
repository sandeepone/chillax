package pipelines

import (
	"testing"
	"time"
)

func TestNewStage(t *testing.T) {
	stage := NewStage("http://localhost:3000")

	if stage.Uri != "http://localhost:3000" {
		t.Error("URI was not set correctly")
	}
	if stage.Method != "POST" {
		t.Error("Default method should be POST.")
	}
	if stage.Timeout != 1*time.Second {
		t.Error("Default timeout should be 1 second.")
	}
}

func TestStageRunBadRequest(t *testing.T) {
	stage := NewStage("http://localhost:3000")
	stage.Timeout = 1 * time.Millisecond

	_, errChan := stage.Run()

	err := <-errChan

	if err == nil {
		t.Error("Request should fail.")
	}
}
