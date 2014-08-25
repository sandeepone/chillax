package settings

import (
    "os"
    "testing"
)

func TestSetDefaults(t *testing.T) {
    settings := NewServerSettings()

    if settings.HttpPort != "80" {
        t.Errorf("By default, settings.HttpPort should == 80. settings.HttpPort: %v", settings.HttpPort)
    }
    if settings.DefaultAssetsPath == "" {
        t.Errorf("settings.DefaultAssetsPath should not be empty")
    }
}

func TestSetEnvOverrides(t *testing.T) {
    originalSettings := NewServerSettings()

    os.Setenv("HTTP_PORT", "8080")
    os.Setenv("PROXY_HANDLERS_PATH", "/aaa")
    os.Setenv("DEFAULT_ASSETS_PATH", "/aab")

    settings := NewServerSettings()

    if settings.HttpPort != "8080" {
        t.Errorf("settings.HttpPort should == 8080. settings.HttpPort: %v", settings.HttpPort)
    }
    if settings.ProxyHandlersPath != "/aaa" {
        t.Errorf("settings.ProxyHandlersPath should == /aaa. settings.ProxyHandlersPath: %v", settings.ProxyHandlersPath)
    }
    if settings.DefaultAssetsPath != "/aab" {
        t.Errorf("settings.DefaultAssetsPath should == /aab. settings.DefaultAssetsPath: %v", settings.DefaultAssetsPath)
    }

    os.Setenv("HTTP_PORT", originalSettings.HttpPort)
    os.Setenv("PROXY_HANDLERS_PATH", originalSettings.ProxyHandlersPath)
    os.Setenv("DEFAULT_ASSETS_PATH", originalSettings.DefaultAssetsPath)
}

func TestHttpAddress(t *testing.T) {
    settings := NewServerSettings()

    if settings.HttpAddress() != ":80" {
        t.Errorf("settings.HttpAddress should == :80. settings.HttpAddress(): %v", settings.HttpAddress())
    }
}