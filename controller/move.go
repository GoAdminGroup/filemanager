package controller

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/GoAdminGroup/filemanager/guard"
	"github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
)

func (h *Handler) Move(ctx *context.Context) {

	param := guard.GetMoveParam(ctx)

	if param.Error != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  param.Error.Error(),
		})
		return
	}

	err := os.Rename(filepath.FromSlash(param.Src), filepath.FromSlash(param.Dist))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code": http.StatusInternalServerError,
			"msg":  err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 0,
		"msg":  language.Get("move success"),
	})
}

func (h *Handler) MovePopup(ctx *context.Context) {

	var (
		popupID  = ctx.FormValue("popup_id")
		fileName = ctx.FormValue("id")

		relativePath, _ = url.QueryUnescape(ctx.Query("path"))

		path    = filepath.Join(h.roots.GetPathFromPrefix(ctx), relativePath)
		options = ""
		prefix  = h.Prefix(ctx)

		fileInfos, err = ioutil.ReadDir(filepath.FromSlash(path))
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"code": http.StatusBadRequest,
			"msg":  err.Error(),
		})
		return
	}

	if relativePath == "" {
		relativePath = "/"
	}

	for _, fileInfo := range fileInfos {

		if util.IsHiddenFile(fileInfo.Name()) {
			continue
		}

		if fileInfo.IsDir() {
			filePath := filepath.Join(relativePath, fileInfo.Name())
			if filePath != fileName {
				options += `<option value='` + filePath + `'>` + fileInfo.Name() + `</option>`
			}
		}
	}

	if relativePath != "." && relativePath != "/" {
		parentDir := filepath.Dir(relativePath)
		options += `<option value='` + parentDir + `'>.</option>`
	}

	popupForm := `<form>
          <div class="form-group">
            <select class="form-control select2-hidden-accessible" style="width: 100%;"
            data-multiple="false" data-placeholder="` + language.Get("input name") + `" tabindex="-1" aria-hidden="true" id="fm_move_select">
				<option></option>
				` + options + `
			</select>
          </div>
        </form>
<script>
	$("#fm_move_select").select2();
	$('#` + popupID + ` button.btn.btn-primary').on('click', function (event) {
		$.ajax({
                            method: 'post',
                            url: "` + GetUrl(prefix, "/move") + `",
                            data: {
								dist: $('#fm_move_select').val(),
								src: "` + fileName + `"
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
