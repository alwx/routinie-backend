// Code generated by qtc from "base_page.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

// This is a base page template. All the other template pages implement this interface.
//

//line ../../pkg/templates/base_page.qtpl:3
package templates

//line ../../pkg/templates/base_page.qtpl:3
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line ../../pkg/templates/base_page.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line ../../pkg/templates/base_page.qtpl:3
type Page interface {
//line ../../pkg/templates/base_page.qtpl:3
	Head() string
//line ../../pkg/templates/base_page.qtpl:3
	StreamHead(qw422016 *qt422016.Writer)
//line ../../pkg/templates/base_page.qtpl:3
	WriteHead(qq422016 qtio422016.Writer)
//line ../../pkg/templates/base_page.qtpl:3
	Body() string
//line ../../pkg/templates/base_page.qtpl:3
	StreamBody(qw422016 *qt422016.Writer)
//line ../../pkg/templates/base_page.qtpl:3
	WriteBody(qq422016 qtio422016.Writer)
//line ../../pkg/templates/base_page.qtpl:3
	HTMLClass() string
//line ../../pkg/templates/base_page.qtpl:3
	StreamHTMLClass(qw422016 *qt422016.Writer)
//line ../../pkg/templates/base_page.qtpl:3
	WriteHTMLClass(qq422016 qtio422016.Writer)
//line ../../pkg/templates/base_page.qtpl:3
}

//line ../../pkg/templates/base_page.qtpl:10
func StreamPageTemplate(qw422016 *qt422016.Writer, p Page) {
//line ../../pkg/templates/base_page.qtpl:10
	qw422016.N().S(`
`)
//line ../../pkg/templates/base_page.qtpl:11
	qw422016.N().S(`<!DOCTYPE html><html class="`)
//line ../../pkg/templates/base_page.qtpl:13
	p.StreamHTMLClass(qw422016)
//line ../../pkg/templates/base_page.qtpl:13
	qw422016.N().S(`"><head>`)
//line ../../pkg/templates/base_page.qtpl:15
	p.StreamHead(qw422016)
//line ../../pkg/templates/base_page.qtpl:15
	qw422016.N().S(`</head>`)
//line ../../pkg/templates/base_page.qtpl:17
	p.StreamBody(qw422016)
//line ../../pkg/templates/base_page.qtpl:17
	qw422016.N().S(`</html>`)
//line ../../pkg/templates/base_page.qtpl:19
	qw422016.N().S(`
`)
//line ../../pkg/templates/base_page.qtpl:20
}

//line ../../pkg/templates/base_page.qtpl:20
func WritePageTemplate(qq422016 qtio422016.Writer, p Page) {
//line ../../pkg/templates/base_page.qtpl:20
	qw422016 := qt422016.AcquireWriter(qq422016)
//line ../../pkg/templates/base_page.qtpl:20
	StreamPageTemplate(qw422016, p)
//line ../../pkg/templates/base_page.qtpl:20
	qt422016.ReleaseWriter(qw422016)
//line ../../pkg/templates/base_page.qtpl:20
}

//line ../../pkg/templates/base_page.qtpl:20
func PageTemplate(p Page) string {
//line ../../pkg/templates/base_page.qtpl:20
	qb422016 := qt422016.AcquireByteBuffer()
//line ../../pkg/templates/base_page.qtpl:20
	WritePageTemplate(qb422016, p)
//line ../../pkg/templates/base_page.qtpl:20
	qs422016 := string(qb422016.B)
//line ../../pkg/templates/base_page.qtpl:20
	qt422016.ReleaseByteBuffer(qb422016)
//line ../../pkg/templates/base_page.qtpl:20
	return qs422016
//line ../../pkg/templates/base_page.qtpl:20
}

//line ../../pkg/templates/base_page.qtpl:22
type BasePage struct {
	Title         string
	Description   string
	Keywords      *string
	AssetsVersion string
}

//line ../../pkg/templates/base_page.qtpl:29
func (p *BasePage) StreamHead(qw422016 *qt422016.Writer) {
//line ../../pkg/templates/base_page.qtpl:29
	qw422016.N().S(`
`)
//line ../../pkg/templates/base_page.qtpl:30
	qw422016.N().S(`<title>`)
//line ../../pkg/templates/base_page.qtpl:31
	qw422016.E().S(pageTitle(withDefault(p.Title, "Untitled")))
//line ../../pkg/templates/base_page.qtpl:31
	qw422016.N().S(`</title><meta name="description" content="`)
//line ../../pkg/templates/base_page.qtpl:32
	qw422016.E().S(p.Description)
//line ../../pkg/templates/base_page.qtpl:32
	qw422016.N().S(`" />`)
//line ../../pkg/templates/base_page.qtpl:33
	if p.Keywords != nil {
//line ../../pkg/templates/base_page.qtpl:33
		qw422016.N().S(`<meta name="keywords" content="`)
//line ../../pkg/templates/base_page.qtpl:34
		qw422016.E().S(*p.Keywords)
//line ../../pkg/templates/base_page.qtpl:34
		qw422016.N().S(`" />`)
//line ../../pkg/templates/base_page.qtpl:35
	}
//line ../../pkg/templates/base_page.qtpl:35
	qw422016.N().S(`<meta charset="UTF-8" /><meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1" /><meta name="viewport" content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=no" /><link rel="apple-touch-icon" sizes="180x180" href="/resources/apple-touch-icon.png" /><link rel="icon" type="image/png" sizes="32x32" href="/resources/favicon-32x32.png" /><link rel="icon" type="image/png" sizes="64x64" href="/resources/favicon-64x64.png" /><link rel="icon" type="image/png" sizes="194x194" href="/resources/favicon-194x194.png" /><link rel="icon" type="image/png" sizes="16x16" href="/resources/favicon-16x16.png" />`)
//line ../../pkg/templates/base_page.qtpl:44
	qw422016.N().S(`
`)
//line ../../pkg/templates/base_page.qtpl:45
}

//line ../../pkg/templates/base_page.qtpl:45
func (p *BasePage) WriteHead(qq422016 qtio422016.Writer) {
//line ../../pkg/templates/base_page.qtpl:45
	qw422016 := qt422016.AcquireWriter(qq422016)
//line ../../pkg/templates/base_page.qtpl:45
	p.StreamHead(qw422016)
//line ../../pkg/templates/base_page.qtpl:45
	qt422016.ReleaseWriter(qw422016)
//line ../../pkg/templates/base_page.qtpl:45
}

//line ../../pkg/templates/base_page.qtpl:45
func (p *BasePage) Head() string {
//line ../../pkg/templates/base_page.qtpl:45
	qb422016 := qt422016.AcquireByteBuffer()
//line ../../pkg/templates/base_page.qtpl:45
	p.WriteHead(qb422016)
//line ../../pkg/templates/base_page.qtpl:45
	qs422016 := string(qb422016.B)
//line ../../pkg/templates/base_page.qtpl:45
	qt422016.ReleaseByteBuffer(qb422016)
//line ../../pkg/templates/base_page.qtpl:45
	return qs422016
//line ../../pkg/templates/base_page.qtpl:45
}

//line ../../pkg/templates/base_page.qtpl:47
func (p *BasePage) StreamBody(qw422016 *qt422016.Writer) {
//line ../../pkg/templates/base_page.qtpl:47
	qw422016.N().S(`
<body>
    Needs to be overridden.
</body>
`)
//line ../../pkg/templates/base_page.qtpl:51
}

//line ../../pkg/templates/base_page.qtpl:51
func (p *BasePage) WriteBody(qq422016 qtio422016.Writer) {
//line ../../pkg/templates/base_page.qtpl:51
	qw422016 := qt422016.AcquireWriter(qq422016)
//line ../../pkg/templates/base_page.qtpl:51
	p.StreamBody(qw422016)
//line ../../pkg/templates/base_page.qtpl:51
	qt422016.ReleaseWriter(qw422016)
//line ../../pkg/templates/base_page.qtpl:51
}

//line ../../pkg/templates/base_page.qtpl:51
func (p *BasePage) Body() string {
//line ../../pkg/templates/base_page.qtpl:51
	qb422016 := qt422016.AcquireByteBuffer()
//line ../../pkg/templates/base_page.qtpl:51
	p.WriteBody(qb422016)
//line ../../pkg/templates/base_page.qtpl:51
	qs422016 := string(qb422016.B)
//line ../../pkg/templates/base_page.qtpl:51
	qt422016.ReleaseByteBuffer(qb422016)
//line ../../pkg/templates/base_page.qtpl:51
	return qs422016
//line ../../pkg/templates/base_page.qtpl:51
}

//line ../../pkg/templates/base_page.qtpl:53
func (p *BasePage) StreamHTMLClass(qw422016 *qt422016.Writer) {
//line ../../pkg/templates/base_page.qtpl:53
	qw422016.N().S(`ts`)
//line ../../pkg/templates/base_page.qtpl:53
}

//line ../../pkg/templates/base_page.qtpl:53
func (p *BasePage) WriteHTMLClass(qq422016 qtio422016.Writer) {
//line ../../pkg/templates/base_page.qtpl:53
	qw422016 := qt422016.AcquireWriter(qq422016)
//line ../../pkg/templates/base_page.qtpl:53
	p.StreamHTMLClass(qw422016)
//line ../../pkg/templates/base_page.qtpl:53
	qt422016.ReleaseWriter(qw422016)
//line ../../pkg/templates/base_page.qtpl:53
}

//line ../../pkg/templates/base_page.qtpl:53
func (p *BasePage) HTMLClass() string {
//line ../../pkg/templates/base_page.qtpl:53
	qb422016 := qt422016.AcquireByteBuffer()
//line ../../pkg/templates/base_page.qtpl:53
	p.WriteHTMLClass(qb422016)
//line ../../pkg/templates/base_page.qtpl:53
	qs422016 := string(qb422016.B)
//line ../../pkg/templates/base_page.qtpl:53
	qt422016.ReleaseByteBuffer(qb422016)
//line ../../pkg/templates/base_page.qtpl:53
	return qs422016
//line ../../pkg/templates/base_page.qtpl:53
}
