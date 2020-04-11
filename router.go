package filemanager

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/service"
)

func (f *FileManager) initRouter(srv service.List) *context.App {

	app := context.NewApp()
	route := app.Group(config.GetUrlPrefix())
	authRoute := route.Group("/", auth.Middleware(db.GetConnection(srv)))

	authRoute.GET("/fm/files", f.guard.Files, f.handler.ListFiles)
	authRoute.GET("/fm/download", f.handler.Download)
	authRoute.POST("/fm/upload", f.guard.Upload, f.handler.Upload)
	authRoute.POST("/fm/create/dir/popup", f.handler.CreateDirPopUp)
	authRoute.POST("/fm/create/dir", f.guard.CreateDir, f.handler.CreateDir)
	authRoute.POST("/fm/delete", f.guard.Delete, f.handler.Delete)
	authRoute.POST("/fm/move/popup", f.handler.MovePopup)
	authRoute.POST("/fm/move", f.guard.Move, f.handler.Move)
	authRoute.GET("/fm/preview", f.guard.Preview, f.handler.Preview)
	authRoute.POST("/fm/rename/popup", f.handler.RenamePopUp)
	authRoute.POST("/fm/rename", f.guard.Rename, f.handler.Rename)

	return app
}
