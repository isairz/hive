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

type Product struct {
	gorm.Model
	l10n.Locale
	publish.Status
	sorting.SortingDESC

	Name            string
	NameWithSlug    slug.Slug        `l10n:"sync"`
	Code            string           `l10n:"sync"`
	CategoryID      uint             `l10n:"sync"`
	Category        Category         `l10n:"sync"`
	Tags            []Tag            `l10n:"sync" gorm:"many2many:product_tags"`
	MadeCountry     string           `l10n:"sync"`
	Price           float32          `l10n:"sync"`
	Description     string           `sql:"size:2000"`
	Chapters        []Chapter        `l10n:"sync"`
	Enabled         bool
}

func (product Product) DefaultPath() string {
	defaultPath := "/"
	if len(product.Chapters) > 0 {
		defaultPath = fmt.Sprintf("/products/%s_%s", product.Code, product.Chapters[0].ChapterCode)
	}
	return defaultPath
}

func (product Product) MainImageUrl() string {
	return product.Chapters[0].MainImageUrl()
}

func (product Product) Validate(db *gorm.DB) {
	if strings.TrimSpace(product.Name) == "" {
		db.AddError(validations.NewError(product, "Name", "Name can not be empty"))
	}

	if strings.TrimSpace(product.Code) == "" {
		db.AddError(validations.NewError(product, "Code", "Code can not be empty"))
	}
}

type Chapter struct {
	gorm.Model
	ProductID      uint
	Product        Product
	// ColorID        uint
	Chapter          string
	ChapterCode      string
	Images         []ChapterImage
	// SizeVariations []SizeVariation
}

type ChapterImage struct {
	gorm.Model
	ChapterID uint
	Image            ChapterImageStorage `sql:"type:varchar(4096)"`
}

type ChapterImageStorage struct{ media_library.FileSystem }

func (chapter Chapter) MainImageUrl() string {
	imageURL := "/images/default_product.png"
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