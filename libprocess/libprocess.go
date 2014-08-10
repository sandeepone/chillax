package libprocess

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "syscall"
    "time"
)

const DEFAULT_PING = "1m"

type ProcessWrapper struct {
    Name     string
    Command  string
    Args     []string
    Path     string
    Respawn  int
    Delay    string
    Ping     string
    Pid      int
    Status   string
    Handler  *os.Process
    Respawns int
    Children ProcessWrapperChildren
}

func (p *ProcessWrapper) RunAndWatch(name string) chan *ProcessWrapper {
    ch := make(chan *ProcessWrapper)
    go func() {
        if name == "" { name = p.Name }

        p.Start(name)

        p.DoPing(DEFAULT_PING, func(time time.Duration, p *ProcessWrapper) {
            if p.Pid > 0 {
                p.Respawns = 0
                p.Status   = "running"

                fmt.Printf("%s refreshed after %s.\n", p.Name, time)
            }
        })

        go p.Watch()
        ch <- p
    }()
    return ch
}

func (p *ProcessWrapper) String() string {
    js, err := json.Marshal(p)
    if err != nil {
        log.Print(err)
        return ""
    }
    return string(js)
}

// Start process by name
func (p *ProcessWrapper) Start(name string) error {
    wd, err := os.Getwd()
    if err != nil { return err }

    procAttr := &os.ProcAttr{
        Dir: wd,
        Env: os.Environ(),
        Files: []*os.File{
            os.Stdin,
            os.Stdout,
            os.Stderr,
        },
    }

    args := append([]string{name}, p.Args...)
    process, err := os.StartProcess(p.Command, args, procAttr)

    p.Name    = name
    p.Handler = process
    p.Pid     = process.Pid
    p.Status  = "started"

    return err
}

// Stop process and all its children
func (p *ProcessWrapper) Stop() error {
    if p.Handler != nil {
        p.Children.Stop("all")

        err := p.Handler.Signal(syscall.SIGINT)
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
    p.Pid = 0
    p.Status = status
}

//Restart the process
func (p *ProcessWrapper) Restart() chan *ProcessWrapper {
    p.Stop()
    procWrapperChan := p.RunAndWatch("")
    p.Status = "restarted"

    return procWrapperChan
}

//Run callback on the process after given duration.
func (p *ProcessWrapper) DoPing(duration string, f func(t time.Duration, p *ProcessWrapper)) {
    if p.Ping != "" {
        duration = p.Ping
    }
    t, err := time.ParseDuration(duration)
    if err != nil {
        t, err = time.ParseDuration(DEFAULT_PING)
    }
    go func() {
        select {
        case <-time.After(t):
            f(t, p)
        }
    }()
}

//Watch the process
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

        p.Respawns++

        if p.Respawns > p.Respawn {
            p.Release("exited")
            return
        }

        if p.Delay != "" {
            t, _ := time.ParseDuration(p.Delay)
            time.Sleep(t)
        }
        p.Restart()
    case <-diedChan:
        p.Release("killed")
    }
}

//Run child processes
func (p *ProcessWrapper) Run() {
    for name, p := range p.Children {
        p.RunAndWatch(name)
    }
}

//Child processes.
type ProcessWrapperChildren map[string]*ProcessWrapper

//Stringify
func (c ProcessWrapperChildren) String() string {
    js, err := json.Marshal(c)
    if err != nil { return "" }
    return string(js)
}

//Get child processes names.
func (c ProcessWrapperChildren) Keys() []string {
    keys := []string{}
    for k, _ := range c {
        keys = append(keys, k)
    }
    return keys
}

//Get child process.
func (c ProcessWrapperChildren) Get(key string) *ProcessWrapper {
    if v, ok := c[key]; ok {
        return v
    }
    return nil
}

func (c ProcessWrapperChildren) Stop(name string) {
    if name == "all" {
        for name, p := range c {
            p.Stop()
            delete(c, name)
        }
        return
    }
    p := c.Get(name)
    p.Stop()
    delete(c, name)
}
