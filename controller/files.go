package controller

import (
	"github.com/GoAdminGroup/filemanager/models"
	"github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"io/ioutil"
	"mime"
	"net/url"
	"path/filepath"
	"strings"
)

func (h *Handler) ListFiles(ctx *context.Context) {
	relativePath, _ := url.QueryUnescape(ctx.Query("path"))

	path := filepath.Join(h.root, relativePath)

	var filesOfDir = make(models.Files, 0)

	if !strings.Contains(path, h.root) {
		h.table(ctx, filesOfDir, errors.DirIsNotExist)
		return
	}

	if !util.FileExist(path) {
		h.table(ctx, filesOfDir, errors.DirIsNotExist)
		return
	}

	if !util.IsDirectory(path) {
		h.table(ctx, filesOfDir, errors.IsNotDir)
		return
	}

	fileInfos, err := ioutil.ReadDir(path)

	if err != nil {
		h.table(ctx, filesOfDir, err)
		return
	}

	for _, fileInfo := range fileInfos {

		if util.IsHiddenFile(fileInfo.Name()) {
			continue
		}

		file := models.File{
			IsDirectory:  fileInfo.IsDir(),
			Name:         fileInfo.Name(),
			Size:         int(fileInfo.Size()),
			Extension:    strings.TrimLeft(filepath.Ext(fileInfo.Name()), "."),
			Path:         filepath.Join(relativePath, fileInfo.Name()),
			Mime:         mime.TypeByExtension(filepath.Ext(fileInfo.Name())),
			LastModified: fileInfo.ModTime().Unix(),
		}

		filesOfDir = append(filesOfDir, file)
	}

	h.table(ctx, filesOfDir, nil)
	return
}
