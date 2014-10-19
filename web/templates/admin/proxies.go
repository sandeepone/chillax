package admin

import (
	chillax_web_templates "github.com/didip/chillax/web/templates"
	"strings"
)

func NewAdminProxies() *AdminProxies {
	ap := &AdminProxies{}
	ap.Name = "/proxies"
	ap.Src = ap.StringWithInheritance()
	ap.BaseTemplate = NewAdminBase()
	return ap
}

type AdminProxies struct {
	chillax_web_templates.GoTemplate
	BaseTemplate *AdminBase
}

func (p *AdminProxies) BaseString() string {
	return p.BaseTemplate.String()
}

func (p *AdminProxies) String() string {
	return `
{{ range $element := .ProxyHandlers }}
    {{ $element.Backend.Path }}
{{ end }}`
}

func (p *AdminProxies) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
