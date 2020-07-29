package previewer

import (
	"encoding/base64"
	"html/template"

	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/html"
	"github.com/h2non/filetype"
)

type Image struct{}

func (i *Image) Preview(content []byte) template.HTML {
	t, _ := filetype.Image(content)
	b64 := base64.StdEncoding.EncodeToString(content)

	return html.DivEl().SetClass("preview-content").
		SetStyle("margin", "20px auto 20px auto").
		SetStyle("width", "500px").
		SetContent(template2.Default().
			Image().
			SetWidth("500").
			SetHeight("auto").
			SetSrc(template.HTML("data:" + t.MIME.Value + ";base64," + b64)).GetContent()).
		Get()
}
