package filemanager

import (
	"encoding/json"
	"time"

	"github.com/GoAdminGroup/filemanager/controller"
	"github.com/GoAdminGroup/filemanager/guard"
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	language2 "github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/filemanager/modules/permission"
	"github.com/GoAdminGroup/filemanager/modules/root"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/modules/utils"
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

func init() {
	plugins.Add(&FileManager{Base: &plugins.Base{PlugName: Name}})
}

const (
	Name          = "filemanager"
	TableName     = "filemanager_setting"
	ConnectionKey = "filemanager_connection"
)

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

func (f *FileManager) IsInstalled() bool {
	return len(f.roots) != 0
}

func (f *FileManager) GetIndexURL() string {
	return config.Url("/fm")
}

func (f *FileManager) InitPlugin(srv service.List) {

	// DO NOT DELETE
	f.InitBase(srv, "fm")

	f.Conn = db.GetConnection(srv)

	if len(f.roots) == 0 {
		checkExist, _ := db.WithDriver(f.Conn).
			Table("goadmin_site").
			Where("key", "=", ConnectionKey).
			First()
		if checkExist != nil {
			connName := checkExist["value"].(string)
			records, err := db.WithDriverAndConnection(connName, f.Conn).Table(TableName).All()
			if !db.CheckError(err, db.INSERT) && len(records) > 0 {
				for _, record := range records {
					switch record["key"].(string) {
					case "roots":
						err = json.Unmarshal([]byte(record["value"].(string)), &f.roots)
						if err != nil {
							continue
						}
					case "allowUpload":
						f.allowUpload = record["value"].(string) == "1"
					case "allowCreateDir":
						f.allowCreateDir = record["value"].(string) == "1"
					case "allowDelete":
						f.allowDelete = record["value"].(string) == "1"
					case "allowMove":
						f.allowMove = record["value"].(string) == "1"
					case "allowRename":
						f.allowRename = record["value"].(string) == "1"
					case "allowDownload":
						f.allowDownload = record["value"].(string) == "1"
					}
				}
			}
		}
	}

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
	f.handler.HTML = f.HTMLMenu

	language.Lang[language.CN].Combine(language2.CN)
	language.Lang[language.EN].Combine(language2.EN)

	errors.Init()

	f.SetInfo(info)
}

var info = plugins.Info{
	Website:     "http://www.go-admin.cn/plugins/detail/DDN7VxZDTHTeaF8HUU",
	Title:       "FileManager",
	Description: "A plugin help you manage files in your server",
	Version:     "v0.0.6",
	Author:      "Official",
	Url:         "https://github.com/GoAdminGroup/filemanager/archive/v0.0.6.zip",
	Cover:       "",
	Agreement:   "",
	Uuid:        "DDN7VxZDTHTeaF8HUU",
	Name:        "filemanager",
	ModulePath:  "github.com/GoAdminGroup/filemanager",
	CreateDate:  utils.ParseTime("2020-04-05"),
	UpdateDate:  utils.ParseTime("2020-08-03"),
}

type Table struct {
	Id        int64
	Key       string    `xorm:"VARCHAR(100) 'key'"`
	Value     string    `xorm:"TEXT 'value'"`
	CreatedAt time.Time `xorm:"'created_at' timestamp NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `xorm:"'updated_at' timestamp NULL DEFAULT CURRENT_TIMESTAMP"`
}

func (*Table) TableName() string {
	return TableName
}

func (f *FileManager) AddRoot(key string, value root.Root) *FileManager {
	f.roots.Add(key, value)
	return f
}
