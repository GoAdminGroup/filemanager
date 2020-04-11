package filemanager

import (
	"github.com/GoAdminGroup/filemanager/controller"
	"github.com/GoAdminGroup/filemanager/guard"
	"github.com/GoAdminGroup/filemanager/modules/error"
	language2 "github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/filemanager/modules/permission"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/plugins"
)

type FileManager struct {
	app  *context.App
	name string
	root string
	conn db.Connection

	handler *controller.Handler
	guard   *guard.Guardian

	allowUpload    bool
	allowCreateDir bool
	allowDelete    bool
	allowMove      bool
	allowDownload  bool
}

func NewFileManager(rootPath string) *FileManager {
	return &FileManager{
		name:           "filemanager",
		root:           rootPath,
		allowUpload:    true,
		allowCreateDir: true,
		allowDelete:    true,
		allowMove:      true,
		allowDownload:  true,
	}
}

type Config struct {
	AllowUpload    bool
	AllowCreateDir bool
	AllowDelete    bool
	AllowMove      bool
	AllowDownload  bool
	Path           string
}

func NewFileManagerWithConfig(cfg Config) *FileManager {
	return &FileManager{
		name:           "filemanager",
		root:           cfg.Path,
		allowUpload:    cfg.AllowUpload,
		allowCreateDir: cfg.AllowCreateDir,
		allowDelete:    cfg.AllowDelete,
		allowMove:      cfg.AllowMove,
		allowDownload:  cfg.AllowDownload,
	}
}

func (f *FileManager) InitPlugin(srv service.List) {
	f.conn = db.GetConnection(srv)
	p := permission.Permission{
		AllowUpload:    f.allowUpload,
		AllowCreateDir: f.allowCreateDir,
		AllowDelete:    f.allowDelete,
		AllowMove:      f.allowMove,
		AllowDownload:  f.allowDownload,
	}
	f.handler = controller.NewHandler(f.root, f.conn, p)
	f.guard = guard.New(f.root, f.conn, p)
	f.app = f.initRouter(srv)

	language.Lang[language.CN].Combine(language2.CN)
	language.Lang[language.EN].Combine(language2.EN)

	errors.Init()
}

func (f *FileManager) GetRequest() []context.Path {
	return f.app.Requests
}

func (f *FileManager) GetHandler() context.HandlerMap {
	return plugins.GetHandler(f.app)
}

func (f *FileManager) Name() string {
	return f.name
}
