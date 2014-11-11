package backend

import (
	"errors"
	"fmt"
	"github.com/chillaxio/chillax/libnet"
	"github.com/chillaxio/chillax/libprocess"
	chillax_portkeeper "github.com/chillaxio/chillax/portkeeper"
	"os"
	"strconv"
	"strings"
)

// Process data structure
type ProxyBackendProcessConfig struct {
	HttpPortEnv string `toml:httpportenv`
	Hosts       []string
	Instances   []ProxyBackendProcessInstanceConfig
}

type ProxyBackendProcessInstanceConfig struct {
	Command        string
	Delay          string
	Ping           string
	Env            []string
	Host           string
	MapPorts       map[string]int
	ProcessWrapper *libprocess.ProcessWrapper
}

func (pb *ProxyBackend) NewProcessWrapper(command string) *libprocess.ProcessWrapper {
	pw := &libprocess.ProcessWrapper{
		Name:       pb.ProxyName(),
		Command:    command,
		StopDelay:  pb.Delay,
		StartDelay: pb.Delay,
		Ping:       pb.Ping,
	}
	pw.SetDefaults()
	return pw
}

func (pb *ProxyBackend) NewProxyBackendProcessInstanceConfig(host string, httpPort int) ProxyBackendProcessInstanceConfig {
	pbpi := ProxyBackendProcessInstanceConfig{}

	pbpi.Command = pb.Command
	pbpi.Command = strings.Replace(pbpi.Command, fmt.Sprintf("$%v", pb.Process.HttpPortEnv), strconv.Itoa(httpPort), -1)
	pbpi.Delay = pb.Delay
	pbpi.Ping = pb.Ping
	pbpi.Host = host
	pbpi.Env = pb.Env
	pbpi.Env = append(pbpi.Env, fmt.Sprintf("%v=%v", pb.Process.HttpPortEnv, httpPort))
	pbpi.MapPorts = make(map[string]int)

	pbpi.MapPorts[pb.Process.HttpPortEnv] = httpPort

	pbpi.ProcessWrapper = pb.NewProcessWrapper(pbpi.Command)

	return pbpi
}

func (pb *ProxyBackend) CreateProcesses() []error {
	errorSlice := make([]error, 0)

	if pb.Process == nil {
		missingProcessSectionErr := errors.New("[process] section is missing.")
		errorSlice = append(errorSlice, missingProcessSectionErr)
		return errorSlice
	}

	numHosts := len(pb.Process.Hosts)

	pb.Process.Instances = make([]ProxyBackendProcessInstanceConfig, pb.Numprocs)

	for i := 0; i < pb.Numprocs; i++ {
		hostIndex := i

		if i >= numHosts {
			hostIndex = i % numHosts
		}

		host := pb.Process.Hosts[hostIndex]
		newPort := chillax_portkeeper.ReservePort(host)

		pb.Process.Instances[i] = pb.NewProxyBackendProcessInstanceConfig(host, newPort)

		err := pb.Save()
		if err != nil {
			errorSlice = append(errorSlice, err)
		}
	}
	return errorSlice
}

func (pb *ProxyBackend) StartProcesses() []error {
	errorSlice := make([]error, pb.Numprocs)

	if pb.Process == nil {
		return errorSlice
	}

	for _, env := range pb.Env {
		envParts := strings.Split(env, "=")
		os.Setenv(envParts[0], envParts[1])
	}

	for i, instance := range pb.Process.Instances {
		//
		// Start process only if instance.Host is local
		//
		if libnet.RemoteToLocalHostEquality(instance.Host) {
			err := instance.ProcessWrapper.StartAndWatch()
			if err == nil {
				err = pb.Save()
			}
			errorSlice[i] = err
		}
	}

	return errorSlice
}

func (pb *ProxyBackend) StopProcesses() []error {
	errorSlice := make([]error, pb.Numprocs)

	if pb.Process == nil {
		return errorSlice
	}

	errChan := make(chan error, pb.Numprocs)

	for i, instance := range pb.Process.Instances {
		go func(i int) {
			if instance.ProcessWrapper == nil {
				errChan <- errors.New("Process has not been started.")
			} else {
				err := instance.ProcessWrapper.Stop()

				if err == nil {
					err = pb.Save()
				}
				errChan <- err
			}
		}(i)
	}

	for i := 0; i < len(pb.Process.Instances); i++ {
		err := <-errChan
		errorSlice[i] = err
	}
	close(errChan)

	return errorSlice
}

func (pb *ProxyBackend) RestartProcesses() []error {
	errorSlice := make([]error, pb.Numprocs)

	if pb.Process == nil {
		return errorSlice
	}

	for i, instance := range pb.Process.Instances {
		if instance.ProcessWrapper == nil {
			errorSlice[i] = errors.New("Process has not been started.")
		} else {
			err := instance.ProcessWrapper.RestartAndWatch()

			if err == nil {
				err = pb.Save()
			}
			errorSlice[i] = err
		}
	}

	return errorSlice
}
