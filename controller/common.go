package controller

import (
	"bytes"
	"github.com/GoAdminGroup/filemanager/models"
	"github.com/GoAdminGroup/filemanager/modules/language"
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
	template2 "html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"time"
)

type Handler struct {
	root       string
	conn       db.Connection
	navButtons types.Buttons
}

func NewHandler(root string, conn db.Connection) *Handler {
	return &Handler{
		root: root,
		conn: conn,
	}
}

func (h *Handler) Execute(ctx *context.Context, panel types.Panel, animation ...bool) *bytes.Buffer {
	return plugins.Execute(ctx, h.conn, h.navButtons, auth.Auth(ctx), panel, animation...)
}

func (h *Handler) table(ctx *context.Context, files models.Files, err error) {
	buf := h.Execute(ctx, h.tablePanel(ctx, files, err), false)
	ctx.HTML(http.StatusOK, buf.String())
}

func link(u string, c template2.HTML, pjax bool) template2.HTML {
	if pjax {
		return template.Get(config.GetTheme()).Link().SetURL(u).SetContent(c).GetContent()
	}
	return template.Get(config.GetTheme()).Link().NoPjax().SetURL(u).SetContent(c).GetContent()
}

func (h *Handler) tablePanel(ctx *context.Context, files models.Files, err error) types.Panel {
	comp := template.Get(config.GetTheme())

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
	path, _ := url.QueryUnescape(ctx.Query("path"))
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

	for k, f := range files {
		if f.Path[0] != '/' {
			f.Path = "/" + f.Path
		}
		if f.IsDirectory {
			name = icon.Icon(icon.FolderO, 2) + link(config.Url("/fm/files?path="+url.QueryEscape(f.Path)), template2.HTML(f.Name), true)
			op = "-"
		} else {
			name = icon.Icon(icon.FileO, 2) + template2.HTML(f.Name)
			op = link(config.Url("/fm/download?path="+url.QueryEscape(f.Path)), template.HTML(language.Get("download")), false)
		}
		list[k] = map[string]types.InfoItem{
			"name":        {Content: name},
			"size":        {Content: template.HTML(util.ByteCountIEC(f.Size))},
			"modify_time": {Content: template.HTML(time.Unix(f.LastModified, 0).Format("2006-01-02 15:04:05"))},
			"operation":   {Content: op},
		}
	}

	if isSubDir {
		list[length-1] = map[string]types.InfoItem{
			"name":        {Content: link(config.Url("/fm/files"), template2.HTML("."), true)},
			"size":        {Content: "-"},
			"modify_time": {Content: "-"},
			"operation":   {Content: "-"},
		}

		if lastDir != "" {
			list[length-2] = map[string]types.InfoItem{
				"name":        {Content: link(config.Url("/fm/files?path="+url.QueryEscape(lastDir)), template2.HTML("..."), true)},
				"size":        {Content: "-"},
				"modify_time": {Content: "-"},
				"operation":   {Content: "-"},
			}
		}
	}

	table := comp.DataTable().
		SetHideFilterArea(true).
		SetHideRowSelector(true).
		SetInfoList(list).
		SetThead(types.Thead{
			{Head: language.Get("filename"), Field: "name"},
			{Head: language.Get("size"), Field: "size"},
			{Head: language.Get("last modify time"), Field: "modify_time"},
			{Head: language.Get("operation"), Field: "operation"},
		})

	alert := template.HTML("")
	if err != nil {
		alert = template.Get(config.GetTheme()).Alert().Warning(err.Error())
	}

	return types.Panel{
		Content: alert + comp.Box().
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
		Description: path,
	}
}
