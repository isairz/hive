package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/l10n"
	"github.com/qor/media_library"
	"github.com/qor/publish"
	"github.com/qor/slug"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

type Manga struct {
	gorm.Model
	l10n.Locale
	publish.Status
	sorting.SortingDESC

	Name            string
	NameWithSlug    slug.Slug        `l10n:"sync"`
	CategoryID      uint             `l10n:"sync"`
	Category        Category         `l10n:"sync"`
	Tags            []Tag            `l10n:"sync" gorm:"many2many:manga_tags"`
	MadeCountry     string           `l10n:"sync"`
	Description     string           `sql:"size:2000"`
	Chapters        []Chapter        `l10n:"sync"`
	Enabled         bool
}

func (manga Manga) DefaultPath() string {
	defaultPath := "/"
	if len(manga.Chapters) > 0 {
		defaultPath = fmt.Sprintf("/mangas/%s_%s", manga.ID, manga.Chapters[0].ID)
	}
	return defaultPath
}

func (manga Manga) MainImageUrl() string {
	return manga.Chapters[0].MainImageUrl()
}

func (manga Manga) Validate(db *gorm.DB) {
	if strings.TrimSpace(manga.Name) == "" {
		db.AddError(validations.NewError(manga, "Name", "Name can not be empty"))
	}
}

type Chapter struct {
	gorm.Model
    Title        string
	MangaID      uint
    Price        int
	Manga        Manga
	Images       []ChapterImage
}

type ChapterImage struct {
	gorm.Model
	ChapterID uint
	Image            ChapterImageStorage `sql:"type:varchar(409600)"`
}

type ChapterImageStorage struct{ media_library.FileSystem }

func (chapter Chapter) MainImageUrl() string {
	imageURL := "/images/default_manga.png"
	if len(chapter.Images) > 0 {
		imageURL = chapter.Images[0].Image.URL()
	}
	return imageURL
}

func (ChapterImageStorage) GetSizes() map[string]media_library.Size {
	return map[string]media_library.Size{
		"small":  {Width: 480, Height: 480},
		"middle": {Width: 720, Height: 720},
		"big":    {Width: 1080, Height: 1080},
	}
}