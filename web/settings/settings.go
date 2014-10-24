package settings

import (
	"bufio"
	"github.com/BurntSushi/toml"
	"github.com/chillaxio/chillax/libenv"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type ServerSettings struct {
	// HTTP port to listen to
	HttpPort string
	KeyFile  string
	CertFile string

	// Timeout is the duration to allow outstanding requests to survive
	// before forcefully terminating them.
	RequestTimeoutOnRestart string

	ProxyHandlersPath string
	ProxyHandlerTomls [][]byte
}

func NewServerSettings() (*ServerSettings, error) {
	var err error

	settings := &ServerSettings{}

	configPath, err := settings.ConfigPath()
	if err != nil {
		return nil, err
	}

	if configPath != "" {
		fileHandle, _ := os.Open(configPath)
		bufReader := bufio.NewReader(fileHandle)
		definition, _ := ioutil.ReadAll(bufReader)

		_, err = toml.Decode(string(definition), settings)
		if err != nil {
			return nil, err
		}
	}

	settings.SetDefaults()
	settings.SetEnvOverrides()

	err = settings.LoadProxyHandlerTomls()

	return settings, err
}

func (ss *ServerSettings) ConfigPath() (string, error) {
	var err error

	configPath := libenv.EnvWithDefault("CONFIG_PATH", "")
	if configPath != "" {
		configPath, err = filepath.Abs(configPath)
	}

	return configPath, err
}

func (ss *ServerSettings) ConfigDir() (string, error) {
	configPath, err := ss.ConfigPath()
	if err != nil {
		return "", err
	}
	return filepath.Dir(configPath), nil
}

func (ss *ServerSettings) SetDefaults() {
	if ss.HttpPort == "" {
		ss.HttpPort = "80"
	}
	if ss.RequestTimeoutOnRestart == "" {
		ss.RequestTimeoutOnRestart = "3s"
	}
}

func (ss *ServerSettings) SetEnvOverrides() {
	ss.HttpPort = libenv.EnvWithDefault("HTTP_PORT", ss.HttpPort)
	ss.RequestTimeoutOnRestart = libenv.EnvWithDefault("REQUEST_TIMEOUT_ON_RESTART", ss.RequestTimeoutOnRestart)
	ss.ProxyHandlersPath = libenv.EnvWithDefault("PROXY_HANDLERS_PATH", ss.ProxyHandlersPath)
}

func (ss *ServerSettings) LoadProxyHandlerTomls() error {
	if ss.ProxyHandlersPath != "" {
		var globPath string

		if filepath.IsAbs(ss.ProxyHandlersPath) {
			globPath = path.Join(ss.ProxyHandlersPath, "*.toml")
		} else {
			configDir, err := ss.ConfigDir()
			if err != nil {
				return err
			}

			resolvedProxyHandlersPath := filepath.Join(configDir, ss.ProxyHandlersPath)

			resolvedProxyHandlersPath, err = filepath.Abs(resolvedProxyHandlersPath)
			if err != nil {
				return err
			}

			globPath = path.Join(resolvedProxyHandlersPath, "*.toml")
		}

		files, err := filepath.Glob(globPath)
		if err != nil {
			return err
		}

		ss.ProxyHandlerTomls = make([][]byte, len(files))

		for i, fullFilename := range files {
			fileHandle, err := os.Open(fullFilename)

			if err != nil {
				return err
			}

			bufReader := bufio.NewReader(fileHandle)
			definition, err := ioutil.ReadAll(bufReader)

			if err != nil {
				return err
			}

			ss.ProxyHandlerTomls[i] = definition
		}
	}
	return nil
}

func (ss *ServerSettings) HttpAddress() string {
	return ":" + ss.HttpPort
}
