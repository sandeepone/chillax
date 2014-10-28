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
					<th>Environment Variables</th>
					<th width="50">Processes</th>
					<th width="50">Process Ping Interval</th>
					<th width="50">Process Ping Delay</th>
				</tr>
			</thead>

			<tbody>
				<tr>
					<td>{{ .ProxyBackend.Domain }}</td>
					<td>{{ .ProxyBackend.Path }}</td>
					<td><code>{{ .ProxyBackend.Command }}</code></td>
					<td><code>{{ .ProxyBackend.Env }}</code></td>
					<td>{{ .ProxyBackend.UpNumprocs }}/{{ .ProxyBackend.Numprocs }}</td>
					<td>{{ .ProxyBackend.Ping }}</td>
					<td>{{ .ProxyBackend.Delay }}</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>

{{ if .ProxyBackend.IsDocker }}
{{ else }}
<div class="row">
	<div class="large-12 columns">
		<h3>Processes</h3>
	</div>
</div>

<div class="row">
	<div class="large-12 columns">
		<table class="full-width">
			<thead>
				<tr>
					<th>PID</th>
					<th>Command</th>
					<th>Environment Variables</th>
				</tr>
			</thead>

			<tbody>
				{{ range $element := .ProxyBackend.Process.Instances }}
				<tr>
					<td>{{ $element.ProcessWrapper.Pid }}</td>
					<td><code>{{ $element.Command }}</code></td>
					<td><code>{{ $element.Env }}</code></td>
				</tr>
				{{ end }}
			</tbody>
		</table>
	</div>
</div>
{{ end }}
`
}

func (p *AdminProxy) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
