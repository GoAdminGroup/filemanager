package root

import (
	"github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/go-admin/context"
)

type Roots map[string]string

func (r Roots) Add(key, value string) {
	r[key] = value
}

func (r Roots) GetFromPrefix(ctx *context.Context) string {
	prefix := ctx.Query("__prefix")
	if root, ok := r[prefix]; ok {
		return root
	}
	panic(errors.WrongPrefix)
}
