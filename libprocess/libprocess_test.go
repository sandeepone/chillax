package libprocess

import (
    "encoding/json"
    "testing"
)

func NewProcessWrapperForTest() *ProcessWrapper {
    return &ProcessWrapper{
        Name:    "bash",
        Command: "/bin/bash",
        Args:    []string{"foo", "bar"},
        Respawn: 3,
    }
}

func TestToJson(t *testing.T) {
    p := NewProcessWrapperForTest()

    err := p.Start()
    if err != nil {
        t.Errorf("Unable to start process")
    }

    inJson, _ := p.ToJson()

    var deserializedData map[string]interface{}

    err = json.Unmarshal(inJson, &deserializedData)
    if err != nil {
        t.Errorf("Unable to deserialize JSON. Error: %v", err)
    }

    if deserializedData["Name"].(string) != p.Name {
        t.Errorf("Bad deserialization")
    }

    err = p.Stop()
    if err != nil {
        t.Errorf("Unable to stop process")
    }
}

func TestProcessStartStop(t *testing.T) {
    p := NewProcessWrapperForTest()

    err := p.Start()
    if err != nil {
        t.Errorf("Unable to start process")
    }
    if p.Status != "started" {
        t.Errorf("process status is set incorrectly")
    }
    if p.Pid <= 0 || p.Handler.Pid <= 0 {
        t.Errorf("Process should start with PID > 0")
    }
    if p.Pid != p.Handler.Pid {
        t.Errorf("ProcessWrapper PID should == Process PID")
    }

    err = p.Stop()
    if err != nil {
        t.Errorf("Unable to stop process")
    }
    if p.Status != "stopped" {
        t.Errorf("process status is set incorrectly")
    }
}
