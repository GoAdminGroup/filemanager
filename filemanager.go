package filemanager

import (
	"github.com/GoAdminGroup/filemanager/controller"
	"github.com/GoAdminGroup/filemanager/guard"
	"github.com/GoAdminGroup/filemanager/modules/error"
	language2 "github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/filemanager/modules/permission"
	"github.com/GoAdminGroup/filemanager/modules/root"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/plugins"
)

type FileManager struct {
	*plugins.Base

	roots root.Roots

	handler *controller.Handler
	guard   *guard.Guardian

	allowUpload    bool
	allowCreateDir bool
	allowDelete    bool
	allowMove      bool
	allowDownload  bool
	allowRename    bool
}

const Name = "filemanager"

func NewFileManager(rootPath string, titles ...string) *FileManager {

	if rootPath == "" {
		panic("filemanager: create fail, wrong path")
	}

	title := Name
	if len(titles) > 0 {
		title = titles[0]
	}
	return &FileManager{
		Base:           &plugins.Base{PlugName: Name},
		roots:          root.Roots{"def": root.Root{Path: rootPath, Title: title}},
		allowUpload:    true,
		allowCreateDir: true,
		allowDelete:    true,
		allowMove:      true,
		allowDownload:  true,
		allowRename:    true,
	}
}

type Config struct {
	AllowUpload    bool   `json:"allow_upload",yaml:"allow_upload",ini:"allow_upload"`
	AllowCreateDir bool   `json:"allow_create_dir",yaml:"allow_create_dir",ini:"allow_create_dir"`
	AllowDelete    bool   `json:"allow_delete",yaml:"allow_delete",ini:"allow_delete"`
	AllowMove      bool   `json:"allow_move",yaml:"allow_move",ini:"allow_move"`
	AllowDownload  bool   `json:"allow_download",yaml:"allow_download",ini:"allow_download"`
	AllowRename    bool   `json:"allow_rename",yaml:"allow_rename",ini:"allow_rename"`
	Path           string `json:"path",yaml:"path",ini:"path"`
	Title          string `json:"title",yaml:"title",ini:"title"`
}

func NewFileManagerWithConfig(cfg Config) *FileManager {

	if cfg.Path == "" {
		panic("filemanager: create fail, wrong path")
	}

	if !util.FileExist(cfg.Path) {
		panic("filemanager: wrong directory path")
	}

	if cfg.Title == "" {
		cfg.Title = Name
	}

	return &FileManager{
		Base:           &plugins.Base{PlugName: Name},
		roots:          root.Roots{"def": root.Root{Path: cfg.Path, Title: cfg.Title}},
		allowUpload:    cfg.AllowUpload,
		allowCreateDir: cfg.AllowCreateDir,
		allowDelete:    cfg.AllowDelete,
		allowMove:      cfg.AllowMove,
		allowRename:    cfg.AllowRename,
		allowDownload:  cfg.AllowDownload,
	}
}

func (f *FileManager) InitPlugin(srv service.List) {

	// DO NOT DELETE
	f.InitBase(srv)

	f.Conn = db.GetConnection(srv)
	p := permission.Permission{
		AllowUpload:    f.allowUpload,
		AllowCreateDir: f.allowCreateDir,
		AllowDelete:    f.allowDelete,
		AllowMove:      f.allowMove,
		AllowRename:    f.allowRename,
		AllowDownload:  f.allowDownload,
	}
	f.handler = controller.NewHandler(f.roots, p)
	f.guard = guard.New(f.roots, f.Conn, p)
	f.App = f.initRouter(srv)
	f.handler.HTML = f.HTML

	language.Lang[language.CN].Combine(language2.CN)
	language.Lang[language.EN].Combine(language2.EN)

	errors.Init()
}

func (f *FileManager) AddRoot(key string, value root.Root) *FileManager {
	f.roots.Add(key, value)
	return f
}
