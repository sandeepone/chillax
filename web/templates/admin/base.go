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
</body>
</html>
`
}
