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
		<div id="chart-primary"><svg></svg></div>
	</div>
	<script>
	d3.json("/chillax/api/stats/requests/latency?duration={{ .DurationString }}",function(error,data) {
		nv.addGraph(function() {
			var chart = nv.models.lineChart()
				.margin({top: 30, right: 60, bottom: 50, left: 70})
				//We can set x data accessor to use index. Reason? So the bars all appear evenly spaced.
				.x(function(d,i) { return i })
				.y(function(d,i) {return d[1] })
				;

			chart.xAxis.tickFormat(function(d) {
				var dx = data[0].values[d] && data[0].values[d][0] || 0;
				return d3.time.format('%x')(new Date(dx))
			});

			chart.yAxis.tickFormat(function(d) { return '$' + d3.format(',f')(d) });

			console.log(data);

			d3.select('#chart-primary svg')
				.datum([data])
				.call(chart);

			nv.utils.windowResize(chart.update);

			return chart;
		});
	});
	</script>
</div>
`
}

func (p *AdminStats) StringWithInheritance() string {
	return strings.Replace(p.BaseString(), "{{ .Body }}", p.String(), 1)
}
