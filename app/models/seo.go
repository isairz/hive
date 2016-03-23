package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/seo"
)

type SEOSetting struct {
	gorm.Model
	SiteName    string
	DefaultPage seo.Setting
	HomePage    seo.Setting
	MangaPage seo.Setting `seo:"Name,ID"`
}
