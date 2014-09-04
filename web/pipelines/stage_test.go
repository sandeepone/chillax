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
