package guard

import (
	"path/filepath"

	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
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

	distDir = g.roots.GetPathFromPrefix(ctx) + distDir
	src = g.roots.GetPathFromPrefix(ctx) + src

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
