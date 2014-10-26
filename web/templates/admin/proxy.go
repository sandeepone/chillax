package admin

import (
	"strings"

	chillax_web_templates "github.com/chillaxio/chillax/web/templates"
)

func NewAdminProxy() *AdminProxy {
	ap := &AdminProxy{}
	ap.Name = "/proxies/{name}"
	ap.Src = ap.StringWithInheritance()
	ap.BaseTemplate = NewAdminBase()
	return ap
}

type AdminProxy struct {
	chillax_web_templates.GoTemplate
	BaseTemplate *AdminBase
}

func (p *AdminProxy) BaseString() string {
	return p.BaseTemplate.String()
}

func (p *AdminProxy) String() string {
	return `
<div class="row">
    <div class="large-12 columns">
        <h2>Proxy {{ .ProxyBackend.Domain }}{{ .ProxyBackend.Path }}</h2>
    </div>
</div>

<div class="row">
    <div class="large-12 columns">
        <table class="full-width">
            <thead>
                <tr>
                    <th>Domain</th>
                    <th>Path</th>
                    <th>Command</th>
                    <th width="50">Numprocs</th>
                </tr>
            </thead>

            <tbody>
                <tr>
                    <td>{{ .ProxyBackend.Domain }}</td>
                    <td>{{ .ProxyBackend.Path }}</td>
                    <td>{{ .ProxyBackend.Command }}</td>
                    <td>{{ .ProxyBackend.Numprocs }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</div>
`
}

func (p *AdminProxy) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
