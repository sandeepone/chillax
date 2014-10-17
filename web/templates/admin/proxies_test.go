package admin

import (
	"bufio"
	"bytes"
	chillax_proxy_handler "github.com/didip/chillax/proxy/handler"
	"io/ioutil"
	"os"
	"strings"
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
	p := NewAdminProxies()

	if p.String() == "" {
		t.Errorf("Template string should not be empty. p.String(): %v", p.String())
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

	template, err := NewAdminProxies().Parse()

	if err != nil {
		t.Errorf("Unable to parse template. Error: %v", err)
	}

	var buffer bytes.Buffer

	err = template.Execute(&buffer, data)
	if err != nil {
		t.Errorf("Unable to execute template. Error: %v", err)
	}

	html := buffer.String()

	if html == "" {
		t.Errorf("Generated HTML should not be empty. HTML: %v", html)
	}

	if !strings.Contains(html, "/for-testing-only") {
		t.Errorf("Generated HTML did not contain /for-testing-only. HTML: %v", html)
	}
}
