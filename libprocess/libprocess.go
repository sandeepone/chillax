package libprocess

import (
    "encoding/json"
    "os"
    "syscall"
    "time"
)

type ProcessWrapper struct {
    Name           string
    Path           string
    Command        string
    Args           []string
    StopDelay      string
    StartDelay     string
    Ping           string
    Pid            int
    Status         string
    Handler        *os.Process
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

func (p *ProcessWrapper) StartAndWatch() {
    go func() {
        p.Start()

        p.DoPing(func(time time.Duration, p *ProcessWrapper) {
            if p.Pid > 0 {
                p.RespawnCounter = 0
                p.Status = "running"
            }
        })

        go p.Watch()
    }()
}

// Start process
func (p *ProcessWrapper) Start() error {
    wd, err := os.Getwd()
    if err != nil { return err }

    delayTime, err := time.ParseDuration(p.StartDelay)
    if err != nil { return err }

    time.Sleep(delayTime)

    procAttr := &os.ProcAttr{
        Dir: wd,
        Env: os.Environ(),
        Files: []*os.File{
            os.Stdin,
            os.Stdout,
            os.Stderr,
        },
    }

    args := append([]string{p.Name}, p.Args...)
    process, err := os.StartProcess(p.Command, args, procAttr)

    p.Handler = process
    p.Pid     = process.Pid
    p.Status  = "started"

    return err
}

// Stop process and all its children
func (p *ProcessWrapper) Stop() error {
    if p.Handler != nil {
        delayTime, err := time.ParseDuration(p.StopDelay)
        if err != nil { return err }

        time.Sleep(delayTime)

        err = p.Handler.Signal(syscall.SIGINT)
        if err != nil { return err }
    }
    p.Release("stopped")
    return nil
}

// Release and remove process pidfile
func (p *ProcessWrapper) Release(status string) {
    if p.Handler != nil {
        p.Handler.Release()
    }
    p.Pid = -1
    p.Status = status
}

func (p *ProcessWrapper) RestartAndWatch() error {
    err := p.Stop()
    if err != nil { return err }

    p.StartAndWatch()
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
func (p *ProcessWrapper) DoPing(f func(t time.Duration, p *ProcessWrapper)) {
    t, err := time.ParseDuration(p.Ping)
    if err == nil {
        go func() {
            select {
            case <-time.After(t):
                f(t, p)
            }
        }()
    }
}

// Watch the process changes and restart if necessary
func (p *ProcessWrapper) Watch() {
    if p.Handler == nil {
        p.Release("stopped")
        return
    }

    procStateChan := make(chan *os.ProcessState)
    diedChan      := make(chan error)

    go func() {
        state, err := p.Handler.Wait()
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
