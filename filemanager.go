package filemanager

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/GoAdminGroup/filemanager/controller"
	"github.com/GoAdminGroup/filemanager/guard"
	errors "github.com/GoAdminGroup/filemanager/modules/error"
	language2 "github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/filemanager/modules/permission"
	"github.com/GoAdminGroup/filemanager/modules/root"
	"github.com/GoAdminGroup/filemanager/modules/util"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/db/dialect"
	"github.com/GoAdminGroup/go-admin/modules/language"
	"github.com/GoAdminGroup/go-admin/modules/service"
	"github.com/GoAdminGroup/go-admin/modules/utils"
	"github.com/GoAdminGroup/go-admin/plugins"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
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
	Website:     "https://www.go-admin.cn",
	Title:       "FileManager",
	Description: "A plugin help you manage files in your server",
	Version:     "v0.0.4",
	Author:      "Official",
	Url:         "https://github.com/GoAdminGroup/filemanager/archive/v0.0.4.zip",
	Cover:       "",
	Agreement:   "",
	Uuid:        "DDN7VxZDTHTeaF8HUU",
	Name:        "filemanager",
	ModulePath:  "github.com/GoAdminGroup/filemanager",
	CreateDate:  utils.ParseTime("2020-04-05"),
	UpdateDate:  utils.ParseTime("2020-07-29"),
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

func (f *FileManager) GetInstallationPage() (bool, table.Generator) {
	return false, func(ctx *context.Context) (fileManagerConfiguration table.Table) {
		fileManagerConfiguration = table.NewDefaultTable(table.DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver).
			SetOnlyNewForm())

		formList := fileManagerConfiguration.GetForm().AddXssJsFilter().
			HideBackButton().
			HideContinueNewCheckBox().
			HideResetButton()

		connNames := config.GetDatabases().Connections()
		ops := make(types.FieldOptions, len(connNames))
		for i, name := range connNames {
			ops[i] = types.FieldOption{Text: name, Value: name}
		}

		formList.AddField(language2.Get("Connection"), "conn", db.Varchar, form.SelectSingle).
			FieldOptions(ops)

		formList.AddRow(func(panel *types.FormPanel) {
			panel.AddField(language2.Get("allow upload"), "allowUpload", db.Int, form.Switch).FieldOptions(types.FieldOptions{
				{Value: "1", Text: language2.Get("yes")},
				{Value: "0", Text: language2.Get("no")},
			}).FieldDefault("1").FieldRowWidth(3)
			panel.AddField(language2.Get("allow createdir"), "allowCreateDir", db.Int, form.Switch).FieldOptions(types.FieldOptions{
				{Value: "1", Text: language2.Get("yes")},
				{Value: "0", Text: language2.Get("no")},
			}).FieldDefault("1").FieldRowWidth(4).FieldHeadWidth(4)
		})

		formList.AddRow(func(panel *types.FormPanel) {
			panel.AddField(language2.Get("allow delete"), "allowDelete", db.Int, form.Switch).FieldOptions(types.FieldOptions{
				{Value: "1", Text: language2.Get("yes")},
				{Value: "0", Text: language2.Get("no")},
			}).FieldDefault("1").FieldRowWidth(3)
			panel.AddField(language2.Get("allow move"), "allowMove", db.Int, form.Switch).FieldOptions(types.FieldOptions{
				{Value: "1", Text: language2.Get("yes")},
				{Value: "0", Text: language2.Get("no")},
			}).FieldDefault("1").FieldRowWidth(4).FieldHeadWidth(4)
		})

		formList.AddRow(func(panel *types.FormPanel) {
			panel.AddField(language2.Get("allow rename"), "allowRename", db.Int, form.Switch).FieldOptions(types.FieldOptions{
				{Value: "1", Text: language2.Get("yes")},
				{Value: "0", Text: language2.Get("no")},
			}).FieldDefault("1").FieldRowWidth(3)
			panel.AddField(language2.Get("allow download"), "allowDownload", db.Int, form.Switch).FieldOptions(types.FieldOptions{
				{Value: "1", Text: language2.Get("yes")},
				{Value: "0", Text: language2.Get("no")},
			}).FieldDefault("1").FieldRowWidth(4).FieldHeadWidth(4)
		})

		formList.AddTable(language2.Get("roots"), "roots", func(panel *types.FormPanel) {
			panel.AddField(language2.Get("name"), "name", db.Varchar, form.Text).FieldHideLabel().
				FieldDisplay(func(value types.FieldModel) interface{} {
					return []string{""}
				})
			panel.AddField(language2.Get("title"), "title", db.Varchar, form.Text).FieldHideLabel().
				FieldDisplay(func(value types.FieldModel) interface{} {
					return []string{""}
				})
			panel.AddField(language2.Get("path"), "path", db.Varchar, form.Text).FieldHideLabel().
				FieldDisplay(func(value types.FieldModel) interface{} {
					return []string{""}
				})
		})

		formList.SetInsertFn(func(values form2.Values) error {
			connName := values.Get("conn")
			if connName == "" {
				return errors.EmptyConnectionName
			}
			tables, err := db.WithDriver(f.Conn).Table(f.Conn.GetConfig(connName).Name).ShowTables()
			if err != nil {
				return err
			}
			var rootsMap = make(root.Roots, len(values["name"]))
			for k, name := range values["name"] {
				rootsMap[name] = root.Root{
					Path:  values["path"][k],
					Title: values["title"][k],
				}
			}
			roots, _ := json.Marshal(rootsMap)
			if !utils.InArray(tables, TableName) {
				err = f.Conn.CreateDB(connName, new(Table))
				if err != nil {
					return err
				}
				_, err = db.WithDriverAndConnection(connName, f.Conn).
					Table(TableName).
					Insert(dialect.H{
						"key":   "roots",
						"value": roots,
					})
				if db.CheckError(err, db.INSERT) {
					return err
				}
				for key, value := range values {
					if key != "name" && key != "path" && key != "roots" {
						_, _ = db.WithDriverAndConnection(connName, f.Conn).
							Table(TableName).
							Insert(dialect.H{
								"key":   key,
								"value": value[0],
							})
					}
				}
			} else {
				_, err = db.WithDriverAndConnection(connName, f.Conn).
					Table(TableName).
					Where("key", "=", "roots").
					Update(dialect.H{
						"value": roots,
					})
				if db.CheckError(err, db.UPDATE) {
					return err
				}

				values = values.RemoveSysRemark()
				for key, value := range values {
					if key != "name" && key != "path" && key != "roots" && !strings.Contains(key, "__checkbox__") {
						_, _ = db.WithDriverAndConnection(connName, f.Conn).
							Table(TableName).
							Where("key", "=", key).
							Update(dialect.H{
								"value": value[0],
							})
					}
				}
			}

			checkExist, _ := db.WithDriver(f.Conn).
				Table("goadmin_site").
				Where("key", "=", ConnectionKey).
				First()

			if checkExist != nil {
				_, _ = db.WithDriver(f.Conn).
					Table("goadmin_site").
					Where("key", "=", ConnectionKey).
					Update(dialect.H{
						"value": connName,
					})
			} else {
				_, _ = db.WithDriver(f.Conn).
					Table("goadmin_site").
					Insert(dialect.H{
						"key":   ConnectionKey,
						"value": connName,
					})
			}

			p := permission.Permission{
				AllowUpload:    values.Get("allowUpload") == "1",
				AllowCreateDir: values.Get("allowCreateDir") == "1",
				AllowDelete:    values.Get("allowDelete") == "1",
				AllowMove:      values.Get("allowMove") == "1",
				AllowRename:    values.Get("allowRename") == "1",
				AllowDownload:  values.Get("allowDownload") == "1",
			}

			f.roots = rootsMap
			f.handler.Update(f.roots, p)
			f.guard.Update(f.roots, p)

			return nil
		})

		formList.EnableAjaxData(types.AjaxData{
			SuccessTitle:   language2.Get("install success"),
			ErrorTitle:     language2.Get("install fail"),
			SuccessJumpURL: config.Prefix() + "/fm",
		}).SetFormNewTitle(language2.GetHTML("filemanager installation")).
			SetTitle(language2.Get("filemanager installation")).
			SetFormNewBtnWord(language2.GetHTML("install"))

		return
	}
}

func (f *FileManager) AddRoot(key string, value root.Root) *FileManager {
	f.roots.Add(key, value)
	return f
}
