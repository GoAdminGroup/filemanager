package controller

import (
	"bytes"
	"github.com/GoAdminGroup/filemanager/models"
	"github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/filemanager/modules/permission"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/plugins"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/paginator"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/parameter"
	"github.com/GoAdminGroup/go-admin/template"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/action"
	template2 "html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

type Handler struct {
	root        string
	conn        db.Connection
	navButtons  types.Buttons
	permissions permission.Permission
}

func NewHandler(root string, conn db.Connection, p permission.Permission) *Handler {
	return &Handler{
		root:        root,
		conn:        conn,
		permissions: p,
	}
}

func (h *Handler) Execute(ctx *context.Context, panel types.Panel, animation ...bool) *bytes.Buffer {
	return plugins.Execute(ctx, h.conn, h.navButtons, auth.Auth(ctx), panel, animation...)
}

func (h *Handler) preview(ctx *context.Context, content template2.HTML, relativePath, path string, err error) {

	comp := template.Default()

	alert := template.HTML("")
	if err != nil {
		alert = comp.Alert().Warning(err.Error())
	}

	isSubDir := relativePath != "" && err == nil
	lastDir := ""
	if isSubDir {
		dir := filepath.Dir(relativePath)
		if dir != "." && dir != "/" {
			lastDir = dir
		}
	}

	btns := make(types.Buttons, 0)

	if isSubDir {
		homeBtn := types.GetDefaultButton(language.GetHTML("home"), icon.Home, action.Jump(config.Url("/fm/files")))
		btns = append(btns, homeBtn)
		if lastDir != "" {
			lastBtn := types.GetDefaultButton(language.GetHTML("last"), icon.Backward, action.Jump(config.Url("/fm/files?path="+url.QueryEscape(lastDir))))
			btns = append(btns, lastBtn)
		}
	}

	btnHTML, _ := btns.Content()

	table := comp.DataTable().
		SetHideRowSelector(true).
		SetButtons(btnHTML + btns.FooterContent())

	buf := h.Execute(ctx, types.Panel{
		Content: alert + comp.Box().
			SetBody(content).
			SetHeader(table.GetDataTableHeader()).
			SetNoPadding().
			WithHeadBorder().
			GetContent(),
		Title:       language.Get("filemanager"),
		Description: fixedDescription(path),
	}, false, true)
	ctx.HTML(http.StatusOK, buf.String())
}

func fixedDescription(des string) string {
	return des
	//return string(html.SpanEl().SetAttr("title", des).SetContent(template2.HTML(des)).Get())
}

func (h *Handler) table(ctx *context.Context, files models.Files, err error) {
	buf := h.Execute(ctx, h.tablePanel(ctx, files, err), false)
	ctx.HTML(http.StatusOK, buf.String())
}

func link(u string, c template2.HTML, pjax bool) template2.HTML {
	if pjax {
		return template.Default().Link().SetURL(u).SetContent(c).GetContent()
	}
	return template.Default().Link().NoPjax().SetURL(u).SetContent(c).GetContent()
}

func linkWithAttr(class, c template2.HTML, pjax bool, attr template2.HTMLAttr) template2.HTML {
	if pjax {
		return template.Default().Link().SetClass(class).SetAttributes(attr).SetContent(c).GetContent()
	}
	return template.Default().Link().NoPjax().SetClass(class).SetAttributes(attr).SetContent(c).GetContent()
}

func (h *Handler) hasOperation() bool {
	return h.permissions.HasOperation()
}

func (h *Handler) tablePanel(ctx *context.Context, files models.Files, err error) types.Panel {
	comp := template.Default()
	path, _ := url.QueryUnescape(ctx.Query("path"))

	defaultPageSize := 10
	param := parameter.GetParam(ctx.Request.URL, defaultPageSize)
	total := len(files)

	if len(files) > param.PageSizeInt {
		if len(files) > param.PageSizeInt*param.PageInt {
			files = files[param.PageSizeInt*(param.PageInt-1) : param.PageSizeInt*param.PageInt]
		} else {
			files = files[param.PageSizeInt*(param.PageInt-1):]
		}
	}

	length := len(files)

	isSubDir := path != "" && err == nil
	lastDir := ""
	if isSubDir {
		length++
		dir := filepath.Dir(path)
		if dir != "." && dir != "/" {
			length++
			lastDir = dir
		}
	}
	if path == "" {
		path = "."
	}
	list := make([]map[string]types.InfoItem, length)

	name := template2.HTML("")
	op := template2.HTML("")

	var movePopUp = new(action.PopUpAction)
	movePopUpJs := template2.JS("")
	moveFooter := template2.HTML("")
	if h.permissions.AllowMove && len(files) > 0 {
		movePopUp = action.PopUp("_", language.Get("move"), nil).
			SetBtnTitle(language.GetHTML("move")).
			SetUrl(config.Url("/fm/move/popup?path=" + ctx.Query("path")))
		movePopUp.SetBtnId("fm-move-btn")
		movePopUpJs = movePopUp.Js()
		moveFooter = movePopUp.FooterContent()
	}

	for k, f := range files {
		if f.Path[0] != '/' {
			f.Path = "/" + f.Path
		}

		if f.IsDirectory {
			name = icon.Icon(icon.FolderO, 2) + link(config.Url("/fm/files?path="+url.QueryEscape(f.Path)), template2.HTML(f.Name), true)
		} else {
			name = icon.Icon(icon.File, 2) + link(config.Url("/fm/preview?path="+url.QueryEscape(f.Path)), template2.HTML(f.Name), true)
		}

		list[k] = map[string]types.InfoItem{
			"name":        {Content: name},
			"size":        {Content: template.HTML(util.ByteCountIEC(f.Size))},
			"modify_time": {Content: template.HTML(time.Unix(f.LastModified, 0).Format("2006-01-02 15:04:05"))},
			"path":        {Content: template.HTML(f.Path)},
		}

		if h.hasOperation() {

			del := template.HTML("")
			if h.permissions.AllowDelete {
				del = linkWithAttr("grid-row-delete", language.GetHTML("delete"), false, template2.HTMLAttr("data-id="+f.Path))
			}
			move := template.HTML("")
			if h.permissions.AllowMove {
				move = linkWithAttr("fm-move-btn", language.GetHTML("move"), false,
					template2.HTMLAttr(`data-toggle="modal" data-target="#`+movePopUp.Id+`" data-id="`+f.Path+`"`))
			}
			download := template.HTML("")
			if h.permissions.AllowDownload {
				download = link(config.Url("/fm/download?path="+url.QueryEscape(f.Path)), template.HTML(language.Get("download")), false)
			}

			sep := template2.HTML(" | ")

			if f.IsDirectory {
				if del != "" && move != "" {
					op = del + sep + move
				} else if del == "" && move != "" {
					op = move
				} else if del != "" && move == "" {
					op = del
				} else {
					op = "-"
				}
			} else {
				if download != "" && del != "" && move != "" {
					op = download + sep + del + sep + move
				} else if download != "" && del == "" && move != "" {
					op = download + sep + move
				} else if download != "" && del != "" && move == "" {
					op = download + sep + del
				} else if download != "" && del == "" && move == "" {
					op = download
				} else if download == "" && del == "" && move != "" {
					op = move
				} else if download == "" && del != "" && move == "" {
					op = del
				} else if download == "" && del != "" && move != "" {
					op = del + sep + move
				} else {
					op = "-"
				}
			}

			list[k]["operation"] = types.InfoItem{Content: op}
		}
	}

	if isSubDir {
		list[length-1] = map[string]types.InfoItem{
			"name":        {Content: link(config.Url("/fm/files"), template2.HTML("."), true)},
			"size":        {Content: "-"},
			"modify_time": {Content: "-"},
		}

		if h.hasOperation() {
			list[length-1]["operation"] = types.InfoItem{Content: "-"}
		}

		if lastDir != "" {
			list[length-2] = map[string]types.InfoItem{
				"name":        {Content: link(config.Url("/fm/files?path="+url.QueryEscape(lastDir)), template2.HTML("..."), true)},
				"size":        {Content: "-"},
				"modify_time": {Content: "-"},
			}

			if h.hasOperation() {
				list[length-2]["operation"] = types.InfoItem{Content: "-"}
			}
		}
	}

	escapeLastDir := url.QueryEscape(lastDir)

	btns := make(types.Buttons, 0)

	if h.permissions.AllowCreateDir {
		btns = append(btns, types.GetDefaultButton(language.GetHTML("new directory"), icon.Plus,
			action.PopUp("_", language.Get("new directory"), nil).
				SetBtnTitle(language.GetHTML("create")).
				SetUrl(config.Url("/fm/create/dir/popup?path="+escapeLastDir))))
	}

	if h.permissions.AllowUpload {
		btns = append(btns, types.GetDefaultButton(language.GetHTML("upload"), icon.Upload,
			action.FileUpload("_", nil).SetUrl(config.Url("/fm/upload?path="+url.QueryEscape(lastDir)))))
	}

	if isSubDir {
		homeBtn := types.GetDefaultButton(language.GetHTML("home"), icon.Home, action.Jump(config.Url("/fm/files")))
		btns = append(btns, homeBtn)
		if lastDir != "" {
			lastBtn := types.GetDefaultButton(language.GetHTML("last"), icon.Backward, action.Jump(config.Url("/fm/files?path="+url.QueryEscape(lastDir))))
			btns = append(btns, lastBtn)
		}
	}

	btnHTML, btnsJs := btns.Content()

	thead := types.Thead{
		{Head: language.Get("filename"), Field: "name"},
		{Head: language.Get("size"), Field: "size"},
		{Head: language.Get("last modify time"), Field: "modify_time"},
	}

	if h.hasOperation() {
		thead = append(thead, types.TheadItem{Head: language.Get("operation"), Field: "operation"})
	}

	delUrl := ""

	if h.permissions.AllowDelete {
		delUrl = config.Url("/fm/delete")
	}

	table := comp.DataTable().
		SetHideFilterArea(true).
		SetButtons(btnHTML + btns.FooterContent() + moveFooter).
		SetDeleteUrl(delUrl).
		SetActionJs(btnsJs + movePopUpJs).
		SetPrimaryKey("path").
		SetNoAction().
		SetHideRowSelector(true).
		SetInfoList(list).
		SetThead(thead)

	return h.panel(ctx, path, err, table, total, defaultPageSize)
}

func (h *Handler) panel(ctx *context.Context, path string, err error, table types.DataTableAttribute, total, defaultPageSize int) types.Panel {
	alert := template.HTML("")
	if err != nil {
		alert = template.Default().Alert().Warning(err.Error())
	}

	return types.Panel{
		Content: alert + template.Default().Box().
			SetBody(table.GetContent()).
			SetNoPadding().
			SetHeader(table.GetDataTableHeader()).
			WithHeadBorder().
			SetFooter(paginator.Get(paginator.Config{
				Size:         total,
				PageSizeList: []string{"10", "20", "30", "50"},
				Param:        parameter.GetParam(ctx.Request.URL, defaultPageSize),
			}).GetContent()).
			GetContent(),
		Title:       language.Get("filemanager"),
		Description: fixedDescription(path),
	}
}
