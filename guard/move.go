package guard

import (
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"path/filepath"
)

type MoveParam struct {
	Src    string
	Dist   string
	Prefix string
	Error  error
}

func (g *Guardian) Move(ctx *context.Context) {

	distDir := ctx.FormValue("dist")
	src := ctx.FormValue("src")

	if src == "" || distDir == "" {
		ctx.SetUserValue(deleteParamKey, &MoveParam{Error: errors.EmptyName})
		ctx.Next()
		return
	}

	if distDir == "/" {
		distDir = ""
	}

	distDir = g.roots.GetFromPrefix(ctx) + distDir
	src = g.roots.GetFromPrefix(ctx) + src

	if !util.IsDirectory(distDir) {
		ctx.SetUserValue(deleteParamKey, &MoveParam{Error: errors.IsNotDir})
		ctx.Next()
		return
	}

	ctx.SetUserValue(deleteParamKey, &MoveParam{
		Src:    src,
		Dist:   distDir + "/" + filepath.Base(src),
		Prefix: g.GetPrefix(ctx),
	})
	ctx.Next()
}

func GetMoveParam(ctx *context.Context) *MoveParam {
	return ctx.UserValue[deleteParamKey].(*MoveParam)
}
