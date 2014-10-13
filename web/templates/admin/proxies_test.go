package admin

import (
	"bufio"
	"bytes"
	"fmt"
	chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
	"io/ioutil"
	"os"
	"testing"
)

func NewProxyHandlerForTest() *chillax_proxy_handler.ProxyHandler {
	fileHandle, _ := os.Open("./tests-data/proxy-backend.toml")
	bufReader := bufio.NewReader(fileHandle)
	definition, _ := ioutil.ReadAll(bufReader)
	handler := chillax_proxy_handler.NewProxyHandler(definition)
	return handler
}

func TestProxies(t *testing.T) {
	p := NewProxies()

	if p.Src() == "" {
		t.Errorf("Template source should not be empty. p.Src(): %v", p.Src())
	}
}

func TestProxiesExecute(t *testing.T) {
	proxyHandlers := make([]*chillax_proxy_handler.ProxyHandler, 1)
	proxyHandlers[0] = NewProxyHandlerForTest()

	data := struct {
		ProxyHandlers []*chillax_proxy_handler.ProxyHandler
	}{
		proxyHandlers,
	}

	template, err := NewProxies().Parse()

	if err != nil {
		t.Errorf("Unable to parse template. Error: %v", err)
	}

	var html bytes.Buffer

	err = template.Execute(&html, data)
	if err != nil {
		t.Errorf("Unable to execute template. Error: %v", err)
	}

	if html.String() == "" {
		t.Errorf("Generated HTML should not be empty. HTML: %v", html.String())
	}
}
