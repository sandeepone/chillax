package admin

import (
	"strings"

	chillax_web_templates "github.com/chillaxio/chillax/web/templates"
)

func NewAdminPipeline() *AdminPipeline {
	ap := &AdminPipeline{}
	ap.Name = "/proxies/{name}"
	ap.Src = ap.StringWithInheritance()
	ap.BaseTemplate = NewAdminBase()
	return ap
}

type AdminPipeline struct {
	chillax_web_templates.GoTemplate
	BaseTemplate *AdminBase
}

func (p *AdminPipeline) BaseString() string {
	return p.BaseTemplate.String()
}

func (p *AdminPipeline) String() string {
	return `
<div class="row">
	<div class="large-12 columns">
		<h2>Pipeline {{ .Pipeline.Id }}</h2>
	</div>
</div>

<div class="row">
	<div class="large-12 columns">
		<table class="full-width">
			<thead>
				<tr>
					<th>Description</th>
					<th width="50">Timeout</th>
					<th width="50">Failures</th>
				</tr>
			</thead>

			<tbody>
				<tr>
					<td>{{ .Pipeline.Description }}</td>
					<td>{{ .Pipeline.TimeoutString }}</td>
					<td>{{ .Pipeline.FailCount }}/{{ .Pipeline.FailMax }}</td>
				</tr>
			</tbody>
		</table>
	</div>
</div>
`
}

func (p *AdminPipeline) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
