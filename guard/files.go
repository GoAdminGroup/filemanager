package guard

import (
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
)

type FilesParam struct {
	*Base
}

func (g *Guardian) Files(ctx *context.Context) {

	relativePath, path, err := g.getPaths(ctx)

	if !util.IsDirectory(path) {
		err = errors.IsNotDir
	}

	ctx.SetUserValue(filesParamKey, &FilesParam{
		Base: &Base{
			Path:     relativePath,
			FullPath: path,
			Error:    err,
			Prefix:   g.GetPrefix(ctx),
		},
	})
	ctx.Next()
}

func GetFilesParam(ctx *context.Context) *FilesParam {
	return ctx.UserValue[filesParamKey].(*FilesParam)
}
