package admin

import (
	"strings"

	chillax_web_templates "github.com/chillaxio/chillax/web/templates"
)

func NewAdminStats() *AdminStats {
	ap := &AdminStats{}
	ap.Name = "/stats"
	ap.Src = ap.StringWithInheritance()
	ap.BaseTemplate = NewAdminBase()
	return ap
}

type AdminStats struct {
	chillax_web_templates.GoTemplate
	BaseTemplate *AdminBase
}

func (p *AdminStats) BaseString() string {
	return p.BaseTemplate.String()
}

func (p *AdminStats) String() string {
	return `
<div class="row">
	<div class="large-12 columns">
		<h2>Requests</h2>
	</div>
</div>

<div class="row">
	<div class="large-12 columns">
	</div>
</div>
`
}

func (p *AdminStats) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
