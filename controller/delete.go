package controller

import (
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handler) Delete(ctx *context.Context) {
	relativePaths := ctx.FormValue("id")

	relativePathArr := strings.Split(relativePaths, ",")

	for _, relativePath := range relativePathArr {
		path := filepath.Join(h.root, relativePath)

		if relativePath == "" || !strings.Contains(path, h.root) || !util.FileExist(path) || strings.Contains(path, "..") {
			ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"code": http.StatusBadRequest,
				"msg":  errors.DirIsNotExist.Error(),
			})
			return
		}

		err := os.RemoveAll(path)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
				"code": http.StatusInternalServerError,
				"msg":  err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 200,
		"msg":  "ok",
	})
}
