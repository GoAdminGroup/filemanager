package errors

import (
	"errors"

	"github.com/GoAdminGroup/filemanager/modules/language"
)

var (
	DirIsNotExist       error
	WrongPrefix         error
	IsNotDir            error
	IsNotFile           error
	EmptyName           error
	NoFile              error
	NoPermission        error
	EmptyConnectionName error
)

func Init() {
	DirIsNotExist = errors.New(language.Get("not exist"))
	IsNotDir = errors.New(language.Get("is not a dir"))
	IsNotFile = errors.New(language.Get("is not a file"))
	EmptyName = errors.New(language.Get("empty name"))
	NoFile = errors.New(language.Get("no files"))
	NoPermission = errors.New(language.Get("no permission"))
	WrongPrefix = errors.New(language.Get("wrong prefix"))
	EmptyConnectionName = errors.New(language.Get("empty connection name"))
}
