package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/GoAdminGroup/filemanager/guard"
	"github.com/GoAdminGroup/go-admin/context"
)

func (h *Handler) Delete(ctx *context.Context) {
	param := guard.GetDeleteParam(ctx)

	if param.Error != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  param.Error.Error(),
		})
		return
	}

	for _, path := range param.Paths {
		err := os.RemoveAll(filepath.FromSlash(path))

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
