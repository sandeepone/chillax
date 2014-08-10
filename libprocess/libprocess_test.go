package libprocess

import (
    "testing"
)

func TestProcessStartStop(t *testing.T) {
    p := &ProcessWrapper{
        Name:    "bash",
        Command: "/bin/bash",
        Args:    []string{"foo", "bar"},
        Respawn: 3,
    }
    err := p.Start()
    if err != nil {
        t.Errorf("Unable to start process")
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
}
