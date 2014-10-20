package admin

import (
	chillax_web_templates "github.com/didip/chillax/web/templates"
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
		<div class="large-12 columns">
			<ul class="small-block-grid-2 medium-block-grid-4 large-block-grid-8">
				{{ range $element := .Pipelines }}
				<li style="border: 1px solid #dddddd; text-align: center">
					<h6>ID: {{ $element.Id }}</h6>
					<a href="#" class="button tiny radius round expand no-margin">Details</a>
				</li>
				{{ end }}
			</ul>
		</div>
	</div>
</div>
`
}

func (p *AdminPipelines) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
