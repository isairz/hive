package models

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/l10n"
	"github.com/qor/media_library"
	"github.com/qor/slug"
	"github.com/qor/sorting"
	"github.com/qor/validations"
)

type Manga struct {
	gorm.Model
	l10n.Locale
	sorting.SortingDESC

	Name            string
	NameWithSlug    slug.Slug        `l10n:"sync"`
	CategoryID      uint             `l10n:"sync"`
	Category        Category         `l10n:"sync"`
	Tags            []Tag            `l10n:"sync" gorm:"many2many:manga_tags"`
	PublishedCountry     string           `l10n:"sync"`
	Description     string           `sql:"size:2000"`
	Chapters        []Chapter
	Enabled         bool
}

func (manga Manga) DefaultPath() string {
	defaultPath := "/"
	if len(manga.Chapters) > 0 {
		defaultPath = fmt.Sprintf("/mangas/%d/%d", manga.ID, manga.Chapters[0].ID)
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
	Name         string
	MangaID      uint
	Manga        Manga
	Storage      ChapterStorage `sql:"type:varchar(40960)"`
}

func (chapter Chapter) DefaultName() string {
    if (len(chapter.Name) > 1) {
        return chapter.Name
    }
	return chapter.Manga.Name
}

func (chapter Chapter) DefaultPath() string {
	return fmt.Sprintf("/mangas/%d/%d", chapter.MangaID, chapter.ID)
}

type ChapterStorage struct {
	media_library.FileSystem
}

func (chapter Chapter) MainImageUrl() string {
	imageURL := "/images/default_manga.png"
	if url := chapter.Storage.GetFirstPage(); len(url) > 0 {
		imageURL = url
	}
	return imageURL
}

func (storage ChapterStorage) GetFirstPage() string {
	if len(storage.Pages) > 0 {
		return storage.Pages[0]
	}
	return ""
}

// func (ChapterStorage) GetSizes() map[string]media_library.Size {
// 	return map[string]media_library.Size{
// 		"small":  {Width: 480, Height: 480},
// 		"middle": {Width: 720, Height: 720},
// 		"big":    {Width: 1080, Height: 1080},
// 	}
// }

func (ChapterStorage) GetURLTemplate(option *media_library.Option) (path string) {
	if path = option.Get("URL"); path == "" {
		path = "/system/{{class}}/{{primary_key}}/{{filename}}"
	}
	return
}

func (b *ChapterStorage) GetPages() []string {
    fmt.Printf("get %v\n", b.Pages)
	return b.Pages
}

func (b *ChapterStorage) SetPages(urls []string) bool {
    fmt.Printf("%v\n",urls)
	b.Pages = urls
	return true
}