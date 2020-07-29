package controller

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/GoAdminGroup/filemanager/models"
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
)

func (h *Handler) Download(ctx *context.Context) {

	var (
		relativePath, _ = url.QueryUnescape(ctx.Query("path"))
		raw             = ctx.Query("raw") == "true"
		path            = filepath.Join(h.roots.GetPathFromPrefix(ctx), relativePath)
	)

	var filesOfDir = make(models.Files, 0)

	if !strings.Contains(path, h.roots.GetPathFromPrefix(ctx)) {
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

	agent := ctx.Request.Header.Get("User-Agent")
	if strings.Contains(agent, "MSIE") {
		filename = url.QueryEscape(filename)
		filename = strings.Replace(filename, "+", "%20", -1)
	}
	if strings.Contains(agent, "Edge") && strings.Contains(agent, "Gecko") {
		filename = url.QueryEscape(filename)
		filename = strings.Replace(filename, "+", "%20", -1)
	}

	contentType := util.ParseFileContentType(filename)
	ctx.SetContentType(contentType)

	if !raw {
		ctx.AddHeader("content-disposition", `attachment; filename=`+filename)
	}

	_ = ctx.ServeFile(filepath.FromSlash(path), false)
}
