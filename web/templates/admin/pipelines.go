package admin

import (
	chillax_web_templates "github.com/chillaxio/chillax/web/templates"
	"strings"
)

func NewAdminPipelines() *AdminPipelines {
	ap := &AdminPipelines{}
	ap.Name = "/pipelines"
	ap.Src = ap.StringWithInheritance()
	ap.BaseTemplate = NewAdminBase()
	return ap
}

type AdminPipelines struct {
	chillax_web_templates.GoTemplate
	BaseTemplate *AdminBase
}

func (p *AdminPipelines) BaseString() string {
	return p.BaseTemplate.String()
}

func (p *AdminPipelines) String() string {
	return `
<div class="row">
	<div class="large-12 columns">
		<h2>Pipelines</h2>
	</div>
</div>

<div class="row">
	<div class="large-12 columns">
		<table class="full-width">
			<thead>
				<tr>
					<th width="50">ID</th>
					<th>Description</th>
					<th width="100">Actions</th>
				</tr>
			</thead>

			<tbody>
				{{ range $element := .Pipelines }}
				<tr>
					<td>{{ $element.Id }}</td>
					<td>{{ $element.Description }}</td>
					<td><a href="#" class="button tiny radius round expand no-margin">Details</a></td>
				</tr>
				{{ end }}
			</tbody>
		</table>
	</div>
</div>
`
}

func (p *AdminPipelines) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
