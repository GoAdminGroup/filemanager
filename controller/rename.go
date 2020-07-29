package controller

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/GoAdminGroup/filemanager/guard"
	"github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/go-admin/context"
)

func (h *Handler) Rename(ctx *context.Context) {

	param := guard.GetRenameParam(ctx)

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
		"msg":  language.Get("rename success"),
	})
}

func (h *Handler) RenamePopUp(ctx *context.Context) {

	var (
		popupID = ctx.FormValue("popup_id")
		path    = ctx.FormValue("id")
		prefix  = h.Prefix(ctx)
	)

	popupForm := `<form>
          <div class="form-group">
            <input type="text" class="form-control" placeholder="` + language.Get("input name") + `" id="rename_input">
          </div>
        </form>
<script>
	$('#` + popupID + ` button.btn.btn-primary').on('click', function (event) {
		NProgress.start();
		$.ajax({
                            method: 'post',
                            url: "` + GetUrl(prefix, "/rename") + `",
                            data: {
								name: $('#rename_input').val(),
								path: "` + path + `"
							},
                            success: function (data) {
								NProgress.done();
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
								NProgress.done();
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
