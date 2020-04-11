package guard

import (
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/permission"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"net/url"
	"path/filepath"
	"strings"
)

type Guardian struct {
	conn        db.Connection
	root        string
	permissions permission.Permission
}

func New( r string, c db.Connection, p permission.Permission) *Guardian {
	return &Guardian{
		root:        r,
		conn:        c,
		permissions: p,
	}
}

const (
	filesParamKey     = "files_param"
	uploadParamKey    = "upload_param"
	createDirParamKey = "create_dir_param"
	deleteParamKey    = "delete_param"
	previewParamKey     = "preview_param"
)

type Base struct {
	Path     string
	FullPath string
	Error    error
}

func (g *Guardian) getPaths(ctx *context.Context) (string, string, error) {
	var (
		err error

		relativePath, _ = url.QueryUnescape(ctx.Query("path"))
		path            = filepath.Join(g.root, relativePath)
	)
	if !strings.Contains(path, g.root) {
		err = errors.DirIsNotExist
	}

	if !util.FileExist(path) {
		err = errors.DirIsNotExist
	}

	return relativePath, path, err
}
