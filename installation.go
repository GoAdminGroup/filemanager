package filemanager

import (
	"encoding/json"
	"github.com/GoAdminGroup/go-admin/modules/db/dialect"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	"github.com/GoAdminGroup/go-admin/modules/menu"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/template/icon"
	"strings"

	errors "github.com/GoAdminGroup/filemanager/modules/error"
	language2 "github.com/GoAdminGroup/filemanager/modules/language"
	"github.com/GoAdminGroup/filemanager/modules/permission"
	"github.com/GoAdminGroup/filemanager/modules/root"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/config"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/utils"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
)

func (f *FileManager) GetSettingPage() table.Generator {
	return func(ctx *context.Context) (fileManagerConfiguration table.Table) {

		cfg := table.DefaultConfigWithDriver(config.GetDatabases().GetDefault().Driver)

		if !f.IsInstalled() {
			cfg = cfg.SetOnlyNewForm()
		} else {
			cfg = cfg.SetOnlyUpdateForm()
		}

		fileManagerConfiguration = table.NewDefaultTable(cfg)

		formList := fileManagerConfiguration.GetForm().
			AddXssJsFilter().
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

				_, err := f.NewMenu(menu.NewMenuData{
					Order:      int64(k),
					Title:      values["title"][k],
					Icon:       icon.FolderO,
					PluginName: f.Name(),
					Uri:        "/" + f.URLPrefix + "/" + name + "/list",
					Uuid:       "fm_" + name,
				})

				if err != nil {
					logger.Error("filemanager insert menu error: ", err)
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
				values = values.RemoveSysRemark()
				for key, value := range values {
					if key != "name" && key != "path" && key != "title" && key != "roots" &&
						!strings.Contains(key, "__checkbox__") {
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
					if key != "name" && key != "path" && key != "title" && key != "roots" &&
						!strings.Contains(key, "__checkbox__") {
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
