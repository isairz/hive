package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/l10n"
)

type Author struct {
	gorm.Model
	l10n.Locale
	Name string
}
