package language

import (
	"html/template"

	"github.com/GoAdminGroup/go-admin/modules/language"
)

func Get(key string) string {
	return language.GetWithScope(key, "filemanager")
}

func GetHTML(key string) template.HTML {
	return template.HTML(language.GetWithScope(key, "filemanager"))
}
