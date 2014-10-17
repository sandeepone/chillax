package admin

import (
	chillax_web_templates "github.com/didip/chillax/web/templates"
)

func NewAdminProxies() *AdminProxies {
	gt := &AdminProxies{}
	gt.Name = "/proxies"
	gt.Src = gt.String()
	return gt
}

type AdminProxies struct {
	chillax_web_templates.GoTemplate
}

func (p *AdminProxies) String() string {
	return `
{{ range $element := .ProxyHandlers }}
    {{ $element.Backend.Path }}
{{ end }}`
}
