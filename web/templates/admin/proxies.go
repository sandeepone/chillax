package admin

import (
	"strings"

	chillax_web_templates "github.com/chillaxio/chillax/web/templates"
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
					<th>Pings</th>
					<th width="50">Procs</th>
					<th width="100">Actions</th>
				</tr>
			</thead>

			<tbody>
				{{ range $element := .ProxyHandlers }}
				<tr>
					<td>{{ $element.Backend.Domain }}{{ $element.Backend.Path }}</td>
					<td>
					{{ range $host, $isUp := $element.PingBool }}
						<span class="{{ if $isUp }}success{{ else }}alert{{ end }} label" title="{{ $element.PingLastCheck $host }}">{{ $host }}</span>
						<script>
							var unixNanoLastCheck = $('span:last').attr('title');
							var date = new Date(unixNanoLastCheck/1000/1000);
							$('span:last').attr('title', 'Checked at: ' + date.toString());
						</script>
					{{ end }}
					</td>
					<td>{{ $element.Backend.UpNumprocs }}/{{ $element.Backend.Numprocs }}</td>
					<td><a href="/chillax/admin/proxies/{{ $element.Backend.ProxyName }}" class="button tiny radius round expand no-margin">Details</a></td>
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
