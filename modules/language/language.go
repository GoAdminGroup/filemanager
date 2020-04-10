package language

import "github.com/GoAdminGroup/go-admin/modules/language"

func Get(key string) string {
	return language.GetWithScope(key, "filemanager")
}
