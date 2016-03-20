package api

import (
	"github.com/qor/qor"
	"github.com/isairz/hive/app/models"
	"github.com/isairz/hive/db"
	"github.com/qor/admin"
)

var API *admin.Admin

func init() {
	API = admin.New(&qor.Config{DB: db.DB})

	Manga := API.AddResource(&models.Manga{})

	ChapterMeta := Manga.Meta(&admin.Meta{Name: "Chapters"})
	Chapter := ChapterMeta.Resource
	Chapter.IndexAttrs("ID", "Chapter", "Images")
	Chapter.ShowAttrs("Chapter", "Images")

	API.AddResource(&models.Order{})
	API.AddResource(&models.User{})
}
