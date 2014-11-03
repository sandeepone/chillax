package admin

import (
	chillax_web_templates "github.com/chillaxio/chillax/web/templates"
)

func NewAdminBase() *AdminBase {
	gt := &AdminBase{}
	gt.Name = "/proxies"
	gt.Src = gt.String()
	return gt
}

type AdminBase struct {
	chillax_web_templates.GoTemplate
}

func (p *AdminBase) Bytes() []byte {
	return []byte(p.String())
}

func (p *AdminBase) String() string {
	return `
<html>
<head>
	<link href="//cdnjs.cloudflare.com/ajax/libs/foundation/5.4.6/css/foundation.min.css" rel="stylesheet" type="text/css" media="all">
	<link href="//cdnjs.cloudflare.com/ajax/libs/nvd3/1.1.15-beta/nv.d3.min.css" rel="stylesheet" type="text/css" media="all">

	<script src="//code.jquery.com/jquery-2.1.1.min.js"></script>
	<script src="//cdnjs.cloudflare.com/ajax/libs/foundation/5.4.6/js/foundation.min.js"></script>

	<!-- DynaTable -->
	<!-- Note: Put this on CDNJS -->
	<script type='text/javascript' src='//s3.amazonaws.com/dynatable-docs-assets/js/jquery.dynatable.js'></script>

	<style>
	.full-width {
		width: 100%;
		margin-left: auto;
		margin-right: auto;
		max-width: initial;
	}
	.no-margin {
		margin: 0;
	}
	</style>

	<!-- Refresh admin data every 90 seconds -->
	<meta http-equiv="refresh" content="90">
</head>

<body>

	<div class="sticky">
		<nav class="top-bar" data-topbar role="navigation">
			<ul class="title-area">
				<li class="name">
					<h1><a href="#">Chillax</a></h1>
				</li>
			</ul>

			<section class="top-bar-section">
				<!-- Right Nav Section -->
				<ul class="right">
					<li class="chillax-tab">
						<a href="/chillax/admin/stats">Stats</a>
					</li>
					<li class="chillax-tab proxies">
						<a href="/chillax/admin/proxies">Proxies</a>
					</li>
					<li class="chillax-tab pipelines">
						<a href="/chillax/admin/pipelines">Pipelines</a>
					</li>
					<li class="has-dropdown">
						<a href="#">Didip Kerabat</a>
						<ul class="dropdown">
							<li><a href="/chillax/admin/users/1">Settings</a></li>
							<li class="active"><a href="/chillax/admin/logout">Logout</a></li>
						</ul>
					</li>
				</ul>

				<!-- chillax-tab highlighter -->
				<script>
				$('.chillax-tab').removeClass('active');
				if(location.pathname.indexOf("proxies") > -1) {
					$('.chillax-tab.proxies').addClass('active');
				} else if(location.pathname.indexOf("pipelines") > -1) {
					$('.chillax-tab.pipelines').addClass('active');
				}
				</script>
			</section>
		</nav>
	</div>

	<div class="row full-width">
		<div class="large-12 columns">
			{{ .Body }}
		</div>
	</div>

	<footer class="row full-width">
		<div class="large-12 columns">
			<hr/>
			<div class="row">
				<div class="large-6 columns">
				</div>

				<div class="large-6 columns">
					<ul class="inline-list right">
						<li><a href="https://github.com/chillaxio/chillax">GitHub</a></li>
						<li><a href="#">GoDoc</a></li>
					</ul>
				</div>
			</div>
		</div>
	</footer>



</body>
</html>
`
}
