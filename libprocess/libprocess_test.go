package libprocess

import (
    "encoding/json"
    "testing"
)

func NewProcessWrapperForTest() *ProcessWrapper {
    p := &ProcessWrapper{
        Name:    "bash",
        Command: "/bin/bash",
        Args:    []string{"foo", "bar"},
    }
    p.SetDefaults()
    return p
}

func TestToJson(t *testing.T) {
    p := NewProcessWrapperForTest()

    err := p.Start()
    if err != nil {
        t.Errorf("Unable to start process. Error: %v", err)
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
        t.Errorf("Unable to stop process. Error: %v", err)
    }
}

func TestProcessStartRestartStop(t *testing.T) {
    p := NewProcessWrapperForTest()

    err := p.Start()
    if err != nil {
        t.Errorf("Unable to start process. Error: %v", err)
    }
    if p.Status != "started" {
        t.Errorf("process status is set incorrectly")
    }
    if p.Pid <= 0 || p.CmdStruct.Process.Pid <= 0 {
        t.Errorf("Process should start with PID > 0")
    }
    if p.Pid != p.CmdStruct.Process.Pid {
        t.Errorf("ProcessWrapper PID should == Process PID")
    }

    err = p.Restart()
    if err != nil {
        t.Errorf("Unable to restart process. Error: %v", err)
    }
    if p.Status != "restarted" {
        t.Errorf("process status is set incorrectly")
    }

    err = p.Stop()
    if err != nil {
        t.Errorf("Unable to stop process. Error: %v", err)
    }
    if p.Status != "stopped" {
        t.Errorf("process status is set incorrectly")
    }
}
