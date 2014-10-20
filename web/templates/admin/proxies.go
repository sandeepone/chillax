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
<div class="row">
	<div class="large-12 columns">
		<h2>Proxies</h2>
	</div>
</div>

<div class="row">
	<div class="large-12 columns">
		<table class="full-width">
			<thead>
				<tr>
					<th>Path</th>
					<th width="50">Procs</th>
					<th width="100">Actions</th>
				</tr>
			</thead>

			<tbody>
				{{ range $element := .ProxyHandlers }}
				<tr>
					<td>{{ $element.Backend.Domain }}{{ $element.Backend.Path }}</td>
					<td>{{ $element.Backend.UpNumprocs }}/{{ $element.Backend.Numprocs }}</td>
					<td><a href="#" class="button tiny radius round expand no-margin">Details</a></td>
				</tr>
				{{ end }}
			</tbody>
		</table>
	</div>
</div>
`
}

func (p *AdminProxies) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
