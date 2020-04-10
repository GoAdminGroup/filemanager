package errors

import (
	"errors"
	"github.com/GoAdminGroup/filemanager/modules/language"
)

var (
	DirIsNotExist error
	IsNotDir      error
	IsNotFile     error
)

func Init() {
	DirIsNotExist = errors.New(language.Get("not exist"))
	IsNotDir = errors.New(language.Get("is not a dir"))
	IsNotFile = errors.New(language.Get("is not a file"))
}
