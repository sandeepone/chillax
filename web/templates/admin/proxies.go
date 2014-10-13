package admin

import (
	chillax_web_templates "github.com/didip/chillax/web/templates"
)

func NewProxies() *Proxies {
	gt := &Proxies{}
	gt.Name = "/proxies"
	gt.Src = gt.String()
	return gt
}

type Proxies struct {
	chillax_web_templates.GoTemplate
}

func (p *Proxies) String() string {
	return `
{{ range $element := .ProxyHandlers }}
    {{ $element.Backend.Path }}
{{ end }}`
}
