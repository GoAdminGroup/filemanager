package filemanager

import (
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/auth"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/service"
)

func (f *FileManager) initRouter(srv service.List) *context.App {

	app := context.NewApp()
	route := app.Group(config.GetUrlPrefix())
	authRoute := route.Group("/", auth.Middleware(f.Conn))

	authRoute.GET("/fm/:__prefix/list", f.guard.Files, f.handler.ListFiles)
	authRoute.GET("/fm/:__prefix/download", f.handler.Download)
	authRoute.POST("/fm/:__prefix/upload", f.guard.Upload, f.handler.Upload)
	authRoute.POST("/fm/:__prefix/create/dir/popup", f.handler.CreateDirPopUp)
	authRoute.POST("/fm/:__prefix/create/dir", f.guard.CreateDir, f.handler.CreateDir)
	authRoute.POST("/fm/:__prefix/delete", f.guard.Delete, f.handler.Delete)
	authRoute.POST("/fm/:__prefix/move/popup", f.handler.MovePopup)
	authRoute.POST("/fm/:__prefix/move", f.guard.Move, f.handler.Move)
	authRoute.GET("/fm/:__prefix/preview", f.guard.Preview, f.handler.Preview)
	authRoute.POST("/fm/:__prefix/rename/popup", f.handler.RenamePopUp)
	authRoute.POST("/fm/:__prefix/rename", f.guard.Rename, f.handler.Rename)

	return app
}
