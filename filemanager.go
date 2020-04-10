package filemanager

import (
	"github.com/GoAdminGroup/filemanager/controller"
	"github.com/GoAdminGroup/filemanager/modules/error"
	language2 "github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/plugins"
)

type FileManager struct {
	app     *context.App
	name    string
	root    string
	handler *controller.Handler
	conn    db.Connection
}

func NewFileManager(rootPath string) *FileManager {
	return &FileManager{
		name: "filemanager",
		root: rootPath,
	}
}

func (f *FileManager) InitPlugin(srv service.List) {
	f.conn = db.GetConnection(srv)
	f.handler = controller.NewHandler(f.root, f.conn)
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
