package guard

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/GoAdminGroup/filemanager/modules/constant"
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
)

type UploadParam struct {
	Base
	Files map[string][]*multipart.FileHeader
}

func (g *Guardian) Upload(ctx *context.Context) {

	if !g.permissions.AllowUpload {
		ctx.SetUserValue(uploadParamKey, &UploadParam{Base: Base{Error: errors.NoPermission}})
		ctx.Next()
		return
	}

	relativePath, path, err := g.getPaths(ctx)

	if !util.IsDirectory(path) {
		ctx.SetUserValue(uploadParamKey, &UploadParam{
			Base: Base{Error: errors.IsNotDir},
		})
		ctx.Next()
		return
	}

	files := ctx.Request.MultipartForm.File

	if len(files) == 0 {
		err = errors.NoFile
	}

	ctx.SetUserValue(uploadParamKey, &UploadParam{
		Base: Base{
			Path:     relativePath,
			FullPath: path,
			Error:    err,
			Prefix:   ctx.Query(constant.PrefixKey),
		},
		Files: files,
	})
	ctx.Next()
}

func GetUploadParam(ctx *context.Context) *UploadParam {
	return ctx.UserValue[uploadParamKey].(*UploadParam)
}

type CreateDirParam struct {
	Base
	Dir string
}

func (g *Guardian) CreateDir(ctx *context.Context) {

	if !g.permissions.AllowCreateDir {
		ctx.SetUserValue(createDirParamKey, &CreateDirParam{Base: Base{Error: errors.NoPermission}})
		ctx.Next()
		return
	}

	var (
		relativePath = ctx.FormValue("path")

		name = ctx.FormValue("name")
		path = filepath.Join(g.roots.GetPathFromPrefix(ctx), relativePath)
	)

	if name == "" || !strings.Contains(path, g.roots.GetPathFromPrefix(ctx)) {
		ctx.SetUserValue(createDirParamKey, &CreateDirParam{
			Base: Base{Error: errors.DirIsNotExist},
		})
		ctx.Next()
		return
	}

	if !util.FileExist(path) {
		ctx.SetUserValue(createDirParamKey, &CreateDirParam{
			Base: Base{Error: errors.DirIsNotExist},
		})
		ctx.Next()
		return
	}

	if !util.IsDirectory(path) {
		ctx.SetUserValue(createDirParamKey, &CreateDirParam{
			Base: Base{Error: errors.IsNotDir},
		})
		ctx.Next()
		return
	}

	ctx.SetUserValue(createDirParamKey, &CreateDirParam{
		Base: Base{
			Path:     relativePath,
			FullPath: path,
			Prefix:   ctx.Query(constant.PrefixKey),
		},
		Dir: path + "/" + name,
	})
	ctx.Next()
}

func GetCreateDirParam(ctx *context.Context) *CreateDirParam {
	return ctx.UserValue[createDirParamKey].(*CreateDirParam)
}
