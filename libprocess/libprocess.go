package libprocess

import (
    "fmt"
    "time"
    "os"
    "os/exec"
    "os/signal"
    "syscall"
    "strings"
    "encoding/json"
    "github.com/didip/chillax/libtime"
)

type ProcessWrapper struct {
    Name           string
    Command        string
    Args           []string
    StopDelay      string
    StartDelay     string
    Ping           string
    Pid            int
    Status         string
    CmdStruct      *exec.Cmd
    Respawn        int
    RespawnCounter int
}

func (p *ProcessWrapper) ToJson() ([]byte, error) {
    return json.Marshal(p)
}

func (p *ProcessWrapper) SetDefaults() {
    p.Ping       = "30s"
    p.StopDelay  = "0s"
    p.StartDelay = "0s"
    p.Pid        = -1
    p.Respawn    = -1
}

func (p *ProcessWrapper) NewCmd(command string) *exec.Cmd {
    wd, _ := os.Getwd()

    parts := strings.Fields(command)
    head  := parts[0]
    parts  = parts[1:len(parts)]

    cmd := exec.Command(head,parts...)
    cmd.Dir    = wd
    cmd.Env    = os.Environ()
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin  = os.Stdin

    return cmd
}

func (p *ProcessWrapper) IsProcessStarted() bool {
    return p.CmdStruct.Process != nil
}

func (p *ProcessWrapper) StartAndWatch() error {
    err := p.Start()
    if err != nil { return err }

    p.DoPing(func() {
        if p.Pid > 0 {
            p.RespawnCounter = 0
            p.Status = "running"
        }
    })

    go p.Watch()

    return nil
}

// Start process
func (p *ProcessWrapper) Start() error {
    err := libtime.SleepString(p.StartDelay)
    if err != nil { return err }

    p.CmdStruct = p.NewCmd(p.Command)

    err = p.CmdStruct.Run()
    if err != nil { return err }

    p.Pid    = p.CmdStruct.Process.Pid
    p.Status = "started"

    fmt.Printf("DA FUQ Start: %v", p.CmdStruct)

    p.ListenStopSignals()

    return err
}

// Stop process and all its children
func (p *ProcessWrapper) Stop() error {
    var err error

    if p.CmdStruct != nil && p.CmdStruct.Process != nil {
        err := libtime.SleepString(p.StopDelay)
        if err != nil { return err }

        err = p.CmdStruct.Process.Kill()

        if err == nil {
            p.Release("stopped")
        }
    }

    return err
}

// Release and remove process pidfile
func (p *ProcessWrapper) Release(status string) {
    if p.CmdStruct != nil && p.CmdStruct.Process != nil {
        p.CmdStruct.Process.Release()
    }
    p.Pid = -1
    p.Status = status
}

func (p *ProcessWrapper) RestartAndWatch() error {
    err := p.Stop()
    if err != nil { return err }

    err = p.StartAndWatch()
    if err != nil {return err }

    p.Status = "restarted"

    return nil
}

// Restart process
// Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func (p *ProcessWrapper) Restart() error {
    err := p.Stop()
    if err != nil { return err }

    err = p.Start()
    if err != nil { return err }

    p.Status = "restarted"
    return nil
}

//Run callback on the process after *ProcessWrapper.Ping duration.
func (p *ProcessWrapper) DoPing(callback func()) {
    t, err := time.ParseDuration(p.Ping)
    if err == nil {
        go func() {
            select {
            case <- time.After(t):
                callback()
            }
        }()
    }
}

// Watch the process changes and restart if necessary
func (p *ProcessWrapper) Watch() {
    if p.CmdStruct.Process == nil {
        p.Release("stopped")
        return
    }

    procStateChan := make(chan *os.ProcessState)
    diedChan      := make(chan error)

    go func() {
        state, err := p.CmdStruct.Process.Wait()
        if err != nil {
            diedChan <- err
            return
        }
        procStateChan <- state
    }()

    select {
    case <-procStateChan:
        if p.Status == "stopped" { return }

        p.RespawnCounter++

        if (p.Respawn != -1) && p.RespawnCounter > p.Respawn {
            p.Release("exited")
            return
        }

        p.RestartAndWatch()

    case <-diedChan:
        p.Release("killed")
    }
}

func (p *ProcessWrapper) ListenStopSignals() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

    go func() {
        <-sigChan

        if p.IsProcessStarted() {
            p.Stop()
            close(sigChan)
        }
    }()
}