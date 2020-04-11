package controller

import (
	"github.com/GoAdminGroup/filemanager/guard"
	"github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/file"
	"net/http"
	"net/url"
	"os"
)

func (h *Handler) Upload(ctx *context.Context) {
	param := guard.GetUploadParam(ctx)
	for k := range param.Files {
		for _, fileObj := range param.Files[k] {

			err := file.SaveMultipartFile(fileObj, param.FullPath+"/"+fileObj.Filename)
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

	param := guard.GetCreateDirParam(ctx)

	if param.Error != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  param.Error.Error(),
		})
		return
	}

	err := os.MkdirAll(param.Dir, os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  language.Get("create success"),
	})
}
