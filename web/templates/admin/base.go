package admin

import (
	chillax_web_templates "github.com/didip/chillax/web/templates"
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

func (p *AdminBase) String() string {
	return `
<html>
<head>
	<link href="http://cdn.foundation5.zurb.com/foundation.css" rel="stylesheet" type="text/css" media="all">
	<script src="https://code.jquery.com/jquery-2.1.1.min.js"></script>
	<script src="http://cdn.foundation5.zurb.com/foundation.js"></script>

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
						<a href="/chillax/admin/handlers">Handlers</a>
					</li>
					<li class="chillax-tab">
						<a href="/chillax/admin/proxies">Proxies</a>
					</li>
					<li class="chillax-tab">
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
				$('.chillax-tab a[href="' + location.pathname + '"]').parent().addClass('active');
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
						<li><a href="https://github.com/didip/chillax">GitHub</a></li>
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
