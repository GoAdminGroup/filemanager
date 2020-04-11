package previewer

import (
	template2 "github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/html"
	"html/template"
)

type Code struct{}

func (i *Code) Preview(content []byte) template.HTML {
	return html.DivEl().SetClass("preview-content").
		SetStyle("margin", "10px auto 10px auto").
		SetStyle("width", "90%").
		SetContent(`
	<pre id="preview-code" class="ace_editor" style="min-height:580px;">
        <textarea id="preview-code-textarea" class="ace_text-input">` + template2.HTML(string(content)) + `</textarea>
    </pre>
    <script>
        editor = ace.edit("preview-code");
        editor.setTheme("ace/theme/monokai");
        editor.session.setMode("ace/mode/html");
        editor.setFontSize(14);
        editor.setReadOnly(true);
        editor.setOptions({useWorker: false});
    </script>
`).Get()
}
