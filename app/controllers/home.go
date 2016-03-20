package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isairz/hive/app/models"
	"github.com/isairz/hive/db"
	"github.com/qor/seo"
)

func HomeIndex(ctx *gin.Context) {
	var products []models.Product
	db.DB.Limit(9).Preload("Chapters").Preload("Chapters.Images").Find(&products)
	seoObj := models.SEOSetting{}
	db.DB.First(&seoObj)

	ctx.HTML(
		http.StatusOK,
		"home_index.tmpl",
		gin.H{
			"SeoTag":   seoObj.HomePage.Render(seoObj, nil),
			"Products": products,
			"MicroSearch": seo.MicroSearch{
				URL:    "http://demo.getqor.com",
				Target: "http://demo.getqor.com/search?q={keyword}",
			}.Render(),
			"MicroContact": seo.MicroContact{
				URL:         "http://demo.getqor.com",
				Telephone:   "080-0012-3232",
				ContactType: "Customer Service",
			}.Render(),
		},
	)
}
