package settings

import (
    "os"
    "path"
    "bufio"
    "runtime"
    "io/ioutil"
    "path/filepath"
    "github.com/BurntSushi/toml"
    "github.com/didip/chillax/libenv"
)

type ServerSettings struct {
    // HTTP port to listen to
    HttpPort string

    // Timeout is the duration to allow outstanding requests to survive
    // before forcefully terminating them.
    RequestTimeoutOnRestart string

    ProxyHandlersPath   string
    DefaultAssetsPath   string
    ProxyHandlerTomls   [][]byte
}

func NewServerSettings() (*ServerSettings, error) {
    settings := &ServerSettings{}

    configPath := libenv.EnvWithDefault("CONFIG_PATH", "")
    if configPath != "" {
        fileHandle, _ := os.Open(configPath)
        bufReader     := bufio.NewReader(fileHandle)
        definition, _ := ioutil.ReadAll(bufReader)

        toml.Decode(string(definition), settings)
    }

    settings.SetDefaults()
    settings.SetEnvOverrides()

    err := settings.LoadProxyHandlerTomls()

    return settings, err
}

func (ss *ServerSettings) SetDefaults() {
    if ss.HttpPort == "" {
        ss.HttpPort = "80"
    }
    if ss.RequestTimeoutOnRestart == "" {
        ss.RequestTimeoutOnRestart = "3s"
    }
    if ss.DefaultAssetsPath == "" {
        _, currentFilePath, _, _ := runtime.Caller(1)
        currentFileFullPath, _   := filepath.Abs(currentFilePath)

        ss.DefaultAssetsPath = path.Join(path.Dir(currentFileFullPath), "..", "default-assets")
    }
}

func (ss *ServerSettings) SetEnvOverrides() {
    ss.HttpPort                = libenv.EnvWithDefault("HTTP_PORT", ss.HttpPort)
    ss.RequestTimeoutOnRestart = libenv.EnvWithDefault("REQUEST_TIMEOUT_ON_RESTART", ss.RequestTimeoutOnRestart)
    ss.ProxyHandlersPath       = libenv.EnvWithDefault("PROXY_HANDLERS_PATH", ss.ProxyHandlersPath)
    ss.DefaultAssetsPath       = libenv.EnvWithDefault("DEFAULT_ASSETS_PATH", ss.DefaultAssetsPath)
}

func (ss *ServerSettings) LoadProxyHandlerTomls() error {
    if ss.ProxyHandlersPath != "" {
        files, err := filepath.Glob(path.Join(ss.ProxyHandlersPath, "*.toml"))
        if err != nil { return err }

        ss.ProxyHandlerTomls = make([][]byte, len(files))

        for i, fullFilename := range files {
            fileHandle, err := os.Open(fullFilename)

            if err != nil { return err }

            bufReader       := bufio.NewReader(fileHandle)
            definition, err := ioutil.ReadAll(bufReader)

            if err != nil { return err }

            ss.ProxyHandlerTomls[i] = definition
        }
    }
    return nil
}

func (ss *ServerSettings) HttpAddress() string {
    return ":" + ss.HttpPort
}
