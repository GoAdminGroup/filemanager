package previewer

import (
	"html/template"

	"github.com/GoAdminGroup/html"
)

type Code struct {
	Ext string
}

func NewCode(ext string) *Code {
	return &Code{
		Ext: ext,
	}
}

var ExtJSMap = map[string]template.HTML{
	".go":     "golang",
	".php":    "php",
	".py":     "python",
	".json":   "json",
	".md":     "markdown",
	".sql":    "sql",
	".js":     "javascript",
	".css":    "css",
	".less":   "less",
	".sass":   "sass",
	".cpp":    "c_cpp",
	".rb":     "ruby",
	".ini":    "ini",
	".yaml":   "yaml",
	".yml":    "yaml",
	".xml":    "xml",
	".coffee": "coffee",
	".sh":     "sh",
}

func (i *Code) Preview(content []byte) template.HTML {

	var (
		extJS = template.HTML(`<script src="https://cdn.bootcss.com/ace/1.4.9/ace.js"></script>`)
		ext   = template.HTML("html")
	)

	if e, ok := ExtJSMap[i.Ext]; ok {
		ext = e
	}

	return extJS + html.DivEl().SetClass("preview-content").
		SetStyle("margin", "10px auto 10px auto").
		SetStyle("width", "90%").
		SetContent(`
	<pre id="preview-code" class="ace_editor" style="min-height:580px;">
        <textarea class="ace_text-input"></textarea>
    </pre>
	<div id="ace-code-content" style="display:none;">`+template.HTML(content)+`</div>
    <script>
        editor = ace.edit("preview-code");
		editor.setValue($("#ace-code-content").html());
        editor.setTheme("ace/theme/monokai");
        editor.session.setMode("ace/mode/`+ext+`");
        editor.setFontSize(14);
        editor.setReadOnly(true);
        editor.setOptions({useWorker: false});
		editor.session.selection.clearSelection();
    </script>
`).Get()
}
