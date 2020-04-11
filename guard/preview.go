package guard

import (
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
)

type PreviewParam struct {
	Base
}

func (g *Guardian) Preview(ctx *context.Context) {

	var (
		relativePath, path, err = g.getPaths(ctx)
	)

	if !util.IsFile(path) {
		err = errors.IsNotFile
	}

	ctx.SetUserValue(previewParamKey, &PreviewParam{
		Base: Base{
			Path:     relativePath,
			FullPath: path,
			Error:    err,
		},
	})
	ctx.Next()
}

func GetPreviewParam(ctx *context.Context) *PreviewParam {
	return ctx.UserValue[previewParamKey].(*PreviewParam)
}
