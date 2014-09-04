package pipelines

import (
    "time"
    "testing"
)

func TestNewStage(t *testing.T) {
    stage := NewStage("http://localhost:3000")

    if stage.Uri != "http://localhost:3000" {
        t.Error("URI was not set correctly")
    }
    if stage.Method != "POST" {
        t.Error("Default method should be POST.")
    }
    if stage.Timeout != 1 * time.Second {
        t.Error("Default timeout should be 1 second.")
    }
}

func TestStageRunBadRequest(t *testing.T) {
    stage := NewStage("http://localhost:3000")
    stage.Timeout = 1 * time.Millisecond

    stagerun, err := stage.Run()

    if err == nil {
        t.Error("Request should fail.")
    }
    if stagerun.TimestampUnixNano < 0 {
        t.Error("Request should fail and set stagerun.TimestampUnixNano.")
    }
}
