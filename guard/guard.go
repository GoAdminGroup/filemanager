package guard

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/GoAdminGroup/filemanager/modules/constant"
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	"github.com/GoAdminGroup/filemanager/modules/permission"
	"github.com/GoAdminGroup/filemanager/modules/root"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
)

type Guardian struct {
	conn        db.Connection
	roots       root.Roots
	permissions permission.Permission
}

func New(r root.Roots, c db.Connection, p permission.Permission) *Guardian {
	return &Guardian{
		roots:       r,
		conn:        c,
		permissions: p,
	}
}

func (g *Guardian) Update(r root.Roots, p permission.Permission) {
	g.roots = r
	g.permissions = p
}

const (
	filesParamKey     = "files_param"
	uploadParamKey    = "upload_param"
	createDirParamKey = "create_dir_param"
	deleteParamKey    = "delete_param"
	renameParamKey    = "rename_param"
	previewParamKey   = "preview_param"
)

type Base struct {
	Path     string
	Prefix   string
	FullPath string
	Error    error
}

func (g *Guardian) GetPrefix(ctx *context.Context) string {
	prefix := ctx.Query(constant.PrefixKey)
	if prefix == "" {
		if _, ok := g.roots["def"]; ok {
			return "def"
		} else {
			for name := range g.roots {
				return name
			}
		}
	}
	return prefix
}

func (g *Guardian) getPaths(ctx *context.Context) (string, string, error) {
	var (
		err error

		relativePath, _ = url.QueryUnescape(ctx.Query("path"))
		path            = filepath.Join(g.roots.GetPathFromPrefix(ctx), relativePath)
	)
	if !strings.Contains(path, g.roots.GetPathFromPrefix(ctx)) {
		err = errors.DirIsNotExist
	}

	if !util.FileExist(path) {
		err = errors.DirIsNotExist
	}

	return relativePath, path, err
}
