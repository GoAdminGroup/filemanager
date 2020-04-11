package controller

import (
	"fmt"
	"github.com/GoAdminGroup/filemanager/models"
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"net/url"
	"path/filepath"
	"strings"
)

func (h *Handler) Download(ctx *context.Context) {

	var (
		relativePath, _ = url.QueryUnescape(ctx.Query("path"))
		raw             = ctx.Query("raw") == "true"
		path            = filepath.Join(h.roots.GetFromPrefix(ctx), relativePath)
	)

	var filesOfDir = make(models.Files, 0)

	if !strings.Contains(path, h.roots.GetFromPrefix(ctx)) {
		h.table(ctx, filesOfDir, errors.DirIsNotExist)
		return
	}

	if !util.FileExist(path) {
		h.table(ctx, filesOfDir, errors.DirIsNotExist)
		return
	}

	if util.IsDirectory(path) {
		h.table(ctx, filesOfDir, errors.IsNotFile)
		return
	}

	filename := filepath.Base(path)
	contentType := util.ParseFileContentType(filename)
	ctx.SetContentType(contentType)

	if !raw {
		ctx.AddHeader("content-disposition", `attachment; filename=`+filename)
	}

	fmt.Println("err", ctx.ServeFile(path, false))
}
