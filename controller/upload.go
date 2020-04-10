package controller

import (
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/file"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handler) Upload(ctx *context.Context) {
	relativePath, _ := url.QueryUnescape(ctx.Query("path"))
	path := filepath.Join(h.root, relativePath)

	form := ctx.Request.MultipartForm
	for k := range form.File {
		for _, fileObj := range form.File[k] {

			err := file.SaveMultipartFile(fileObj, path+"/"+fileObj.Filename)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
					"code": http.StatusInternalServerError,
					"msg":  err.Error(),
				})
			}
		}
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  language.GetHTML("upload success"),
	})
}

func (h *Handler) CreateDirPopUp(ctx *context.Context) {

	popupID := ctx.FormValue("popup_id")
	path, _ := url.QueryUnescape(ctx.Query("path"))

	popupForm := `<form>
          <div class="form-group">
            <input type="text" class="form-control" placeholder="` + language.Get("input name") + `" id="dir_name_input">
          </div>
        </form>
<script>
	$('#` + popupID + ` button.btn.btn-primary').on('click', function (event) {
		$.ajax({
                            method: 'post',
                            url: "` + config.Url("/fm/create/dir") + `",
                            data: {
								name: $('#dir_name_input').val(),
								path: "` + path + `"
							},
                            success: function (data) {
                                if (typeof (data) === "string") {
                                    data = JSON.parse(data);
                                }
								$('#` + popupID + `').hide();
								$('.modal-backdrop.fade.in').hide();
                                if (data.code === 0) {
                                    swal(data.msg, '', 'success');
									$.pjax.reload('#pjax-container');
                                } else {
                                    swal(data.msg, '', 'error');
                                }
                            },
							error: function (data) {
								if (data.responseText !== "") {
									swal(data.responseJSON.msg, '', 'error');								
								} else {
									swal('error', '', 'error');
								}
								setTimeout(function() {
									$('#` + popupID + `').hide();
									$('.modal-backdrop.fade.in').hide();
								}, 500)
							},
                        });
	})
</script>`

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  "ok",
		"data": popupForm,
	})
}

func (h *Handler) CreateDir(ctx *context.Context) {
	relativePath := ctx.FormValue("path")
	name := ctx.FormValue("name")

	path := filepath.Join(h.root, relativePath)

	if name == "" {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  errors.EmptyName.Error(),
		})
		return
	}

	if !strings.Contains(path, h.root) || !util.FileExist(path) {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  errors.DirIsNotExist.Error(),
		})
		return
	}

	if !util.IsDirectory(path) {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  errors.IsNotDir.Error(),
		})
		return
	}

	err := os.MkdirAll(path+"/"+name, os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  language.Get("create success"),
	})
}
