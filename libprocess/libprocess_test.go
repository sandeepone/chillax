package libprocess

import (
    "testing"
)

func TestProcessStart(t *testing.T) {
    p := &ProcessWrapper{
        Name:    "bash",
        Command: "/bin/bash",
        Args:    []string{"foo", "bar"},
        Respawn: 3,
    }
    p.Start("bash")
    if p.Pid <= 0 || p.Handler.Pid <= 0 {
        t.Errorf("Process should start with PID > 0")
    }
    if p.Pid != p.Handler.Pid {
        t.Errorf("ProcessWrapper PID should == Process PID")
    }
    p.Stop()
}
