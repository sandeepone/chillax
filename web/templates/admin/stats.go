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
		<h2>Requests Log</h2>
	</div>
</div>

<div class="row">
	<div class="large-12 columns">
		<table id="primary-table" class="full-width">
			<thead>
				<th data-dynatable-column="CurrentUnixNano">Timestamp</th>
				<th data-dynatable-column="Latency">Latency (ns)</th>
				<th data-dynatable-column="Method">Method</th>
				<th data-dynatable-column="RemoteAddr">Remote Address</th>
				<th data-dynatable-column="URI">URI</th>
				<th data-dynatable-column="UserAgent">User Agent</th>
			</thead>
			<tbody></tbody>
		</table>

		<script>
		$.dynatableSetup({
			features: {
				paginate: false
			},
		});

		$.ajax({
			url: '/chillax/api/stats/requests.json?duration=-336h&end=2014-11-02T17:38:36.718Z',
			success: function(data) {
				console.log(data)
				$('#primary-table').dynatable({
					dataset: {
						records: data
					}
				});
			}
		});
		</script>
	</div>
</div>
`
}

func (p *AdminStats) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
