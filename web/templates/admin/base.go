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
	<script src="http://cdn.foundation5.zurb.com/foundation.js"></script>
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
					<li>
						<a href="/chillax/admin/handlers">Handlers</a>
					</li>
					<li>
						<a href="/chillax/admin/proxies">Proxies</a>
					</li>
					<li>
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
			</section>
		</nav>
	</div>


  <div class="row">


	<div class="large-9 columns" role="content">

	  <article>

		<h3><a href="#">Blog Post Title</a></h3>
		<h6>Written by <a href="#">John Smith</a> on August 12, 2012.</h6>

		<div class="row">
		  <div class="large-6 columns">
			<p>Bacon ipsum dolor sit amet nulla ham qui sint exercitation eiusmod commodo, chuck duis velit. Aute in reprehenderit, dolore aliqua non est magna in labore pig pork biltong. Eiusmod swine spare ribs reprehenderit culpa.</p>
			<p>Boudin aliqua adipisicing rump corned beef. Nulla corned beef sunt ball tip, qui bresaola enim jowl. Capicola short ribs minim salami nulla nostrud pastrami.</p>
		  </div>
		  <div class="large-6 columns">
			<img src="http://placehold.it/400x240&text=[img]"/>
		  </div>
		</div>

		<p>Pork drumstick turkey fugiat. Tri-tip elit turducken pork chop in. Swine short ribs meatball irure bacon nulla pork belly cupidatat meatloaf cow. Nulla corned beef sunt ball tip, qui bresaola enim jowl. Capicola short ribs minim salami nulla nostrud pastrami. Nulla corned beef sunt ball tip, qui bresaola enim jowl. Capicola short ribs minim salami nulla nostrud pastrami.</p>

		<p>Pork drumstick turkey fugiat. Tri-tip elit turducken pork chop in. Swine short ribs meatball irure bacon nulla pork belly cupidatat meatloaf cow. Nulla corned beef sunt ball tip, qui bresaola enim jowl. Capicola short ribs minim salami nulla nostrud pastrami.</p>

	  </article>

	  <hr/>

	  <article>

		<h3><a href="#">Blog Post Title</a></h3>
		<h6>Written by <a href="#">John Smith</a> on August 12, 2012.</h6>

		<div class="row">
		  <div class="large-6 columns">
			<p>Bacon ipsum dolor sit amet nulla ham qui sint exercitation eiusmod commodo, chuck duis velit. Aute in reprehenderit, dolore aliqua non est magna in labore pig pork biltong. Eiusmod swine spare ribs reprehenderit culpa.</p>
			<p>Boudin aliqua adipisicing rump corned beef. Nulla corned beef sunt ball tip, qui bresaola enim jowl. Capicola short ribs minim salami nulla nostrud pastrami.</p>
		  </div>
		  <div class="large-6 columns">
			<img src="http://placehold.it/400x240&text=[img]"/>
		  </div>
		</div>

		<p>Pork drumstick turkey fugiat. Tri-tip elit turducken pork chop in. Swine short ribs meatball irure bacon nulla pork belly cupidatat meatloaf cow. Nulla corned beef sunt ball tip, qui bresaola enim jowl. Capicola short ribs minim salami nulla nostrud pastrami. Nulla corned beef sunt ball tip, qui bresaola enim jowl. Capicola short ribs minim salami nulla nostrud pastrami.</p>

		<p>Pork drumstick turkey fugiat. Tri-tip elit turducken pork chop in. Swine short ribs meatball irure bacon nulla pork belly cupidatat meatloaf cow. Nulla corned beef sunt ball tip, qui bresaola enim jowl. Capicola short ribs minim salami nulla nostrud pastrami.</p>

	  </article>

	</div>






	<aside class="large-3 columns">

	  <h5>Categories</h5>
	  <ul class="side-nav">
		<li><a href="#">News</a></li>
		<li><a href="#">Code</a></li>
		<li><a href="#">Design</a></li>
		<li><a href="#">Fun</a></li>
		<li><a href="#">Weasels</a></li>
	  </ul>

	  <div class="panel">
		<h5>Featured</h5>
		<p>Pork drumstick turkey fugiat. Tri-tip elit turducken pork chop in. Swine short ribs meatball irure bacon nulla pork belly cupidatat meatloaf cow.</p>
		<a href="#">Read More →</a>
	  </div>

	</aside>


  </div>






  <footer class="row">
	<div class="large-12 columns">
	  <hr/>
	  <div class="row">
		<div class="large-6 columns">
		  <p>© Copyright no one at all. Go to town.</p>
		</div>
		<div class="large-6 columns">
		  <ul class="inline-list right">
			<li><a href="#">Link 1</a></li>
			<li><a href="#">Link 2</a></li>
			<li><a href="#">Link 3</a></li>
			<li><a href="#">Link 4</a></li>
		  </ul>
		</div>
	  </div>
	</div>
  </footer>



</body>
</html>
`
}
