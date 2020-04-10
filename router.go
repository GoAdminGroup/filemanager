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

	authRoute.GET("/fm/files", f.handler.ListFiles)
	authRoute.GET("/fm/download", f.handler.Download)
	authRoute.POST("/fm/upload", f.handler.Upload)
	authRoute.POST("/fm/create/dir/popup", f.handler.CreateDirPopUp)
	authRoute.POST("/fm/create/dir", f.handler.CreateDir)
	authRoute.POST("/fm/delete", f.handler.Delete)

	return app
}
