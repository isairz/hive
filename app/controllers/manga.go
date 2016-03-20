package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isairz/hive/app/models"
	"github.com/isairz/hive/db"
	"github.com/qor/seo"
)

func MangaIndex(ctx *gin.Context) {
	var (
		mangas   []models.Manga
		seoSetting models.SEOSetting
	)

	db.DB.Limit(10).Find(&mangas)
	db.DB.First(&seoSetting)

	ctx.HTML(
		http.StatusOK,
		"manga_index.tmpl",
		gin.H{
			"Mangas": mangas,
			"SeoTag":   seoSetting.DefaultPage.Render(seoSetting),
			"MicroSearch": seo.MicroSearch{
				URL:    "http://demo.getqor.com",
				Target: "http://demo.getqor.com/search?q=",
			}.Render(),
			"MicroContact": seo.MicroContact{
				URL:         "http://demo.getqor.com",
				Telephone:   "080-0012-3232",
				ContactType: "Customer Service",
			}.Render(),
		},
	)
}

func MangaShow(ctx *gin.Context) {
	var (
		manga          models.Manga
		chapter        models.Chapter
		seoSetting     models.SEOSetting
		mangaID        = ctx.Param("manga_id")
		chapterID      = ctx.Param("chapter_id")
	)

	db.DB.First(&manga, mangaID)
	db.DB.Preload("Images").Preload("Manga").First(&chapter, chapterID)
	db.DB.First(&seoSetting)

	ctx.HTML(
		http.StatusOK,
		"manga_show.tmpl",
		gin.H{
			"Manga":        manga,
			"Chapter":      chapter,
			"SeoTag":       seoSetting.MangaPage.Render(seoSetting, manga),
		},
	)
}
