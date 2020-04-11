package guard

import (
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"path/filepath"
	"strings"
)

type DeleteParam struct {
	Path  string
	Error error
	Paths []string
}

func (g *Guardian) Delete(ctx *context.Context) {

	if !g.permissions.AllowDelete {
		ctx.SetUserValue(deleteParamKey, &DeleteParam{Error: errors.NoPermission})
		ctx.Next()
		return
	}

	relativePaths := ctx.FormValue("id")
	relativePathArr := strings.Split(relativePaths, ",")

	paths := make([]string, 0)

	for _, relativePath := range relativePathArr {
		path := filepath.Join(g.root, relativePath)

		if relativePath == "" || !strings.Contains(path, g.root) || !util.FileExist(path) || strings.Contains(path, "..") {
			ctx.SetUserValue(deleteParamKey, &DeleteParam{Error: errors.DirIsNotExist})
			ctx.Next()
			return
		}

		paths = append(paths, path)

	}
	ctx.SetUserValue(deleteParamKey, &DeleteParam{
		Path:  relativePaths,
		Paths: paths,
	})
	ctx.Next()
}

func GetDeleteParam(ctx *context.Context) *DeleteParam {
	return ctx.UserValue[deleteParamKey].(*DeleteParam)
}
