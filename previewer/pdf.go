package previewer

import (
	"html/template"
)

type PDF struct{}

// TODO
func (p *PDF) Preview(content []byte) template.HTML {
	//sourceJS := template.HTML(`<script src="https://cdn.bootcss.com/pdf.js/2.4.456/pdf.js"></script>`)
	//
	//return sourceJS + html.DivEl().SetClass("preview-content").
	//	SetStyle("margin", "20px auto 20px auto").
	//	SetStyle("width", "500px").
	//	SetContent(``).
	//	Get()
	panic("implement me")
}
