package guard

import (
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"path/filepath"
)

type RenameParam struct {
	Src   string
	Dist  string
	Error error
}

func (g *Guardian) Rename(ctx *context.Context) {

	distName := ctx.FormValue("name")
	src := ctx.FormValue("path")

	if src == "" || src == "/" || distName == "" || distName == "/" {
		ctx.SetUserValue(renameParamKey, &RenameParam{Error: errors.EmptyName})
		ctx.Next()
		return
	}

	if filepath.Ext(distName) == "" && util.IsFile(src) {
		distName += filepath.Ext(src)
	}

	ctx.SetUserValue(renameParamKey, &RenameParam{
		Src:  g.root + src,
		Dist: g.root + filepath.Dir(src) + "/" + distName,
	})
	ctx.Next()
}

func GetRenameParam(ctx *context.Context) *RenameParam {
	return ctx.UserValue[renameParamKey].(*RenameParam)
}
