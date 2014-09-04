package pipelines

import (
    "time"
    "testing"
)

func TestStageNewRequest(t *testing.T) {
    stage := &Stage{}

    req := stage.NewRequest("http://localhost:3000")

    if req.Uri != "http://localhost:3000" {
        t.Error("URI was not set correctly")
    }
    if req.Method != "POST" {
        t.Error("Default method should be POST.")
    }
    if req.Timeout != 1 * time.Second {
        t.Error("Default timeout should be 1 second.")
    }
}

func TestStageRunBadRequest(t *testing.T) {
    stage := &Stage{}

    req := stage.NewRequest("http://localhost:3000")
    req.Timeout = 1 * time.Millisecond

    err := stage.Run(req)

    if err == nil {
        t.Error("Request should fail.")
    }
    if stage.LastFailUnixNano < 0 {
        t.Error("Request should fail and set req.LastFailUnixNano.")
    }
}
