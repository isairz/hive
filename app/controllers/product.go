package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/isairz/hive/app/models"
	"github.com/isairz/hive/db"
	"github.com/qor/seo"
)

func ProductIndex(ctx *gin.Context) {
	var (
		products   []models.Product
		seoSetting models.SEOSetting
	)

	db.DB.Limit(10).Find(&products)
	db.DB.First(&seoSetting)

	ctx.HTML(
		http.StatusOK,
		"product_index.tmpl",
		gin.H{
			"Products": products,
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

func ProductShow(ctx *gin.Context) {
	var (
		product        models.Product
		chapter        models.Chapter
		seoSetting     models.SEOSetting
		codes          = strings.Split(ctx.Param("code"), "_")
		productCode    = codes[0]
		chapterCode    string
	)

	if len(codes) > 1 {
		chapterCode = codes[1]
	}

	db.DB.Where(&models.Product{Code: productCode}).First(&product)
	db.DB.Preload("Images").Preload("Product").Where(&models.Chapter{ProductID: product.ID, ChapterCode: chapterCode}).First(&chapter)
	db.DB.First(&seoSetting)

	ctx.HTML(
		http.StatusOK,
		"product_show.tmpl",
		gin.H{
			"Product":        product,
			"ChapterVariation": chapter,
			"SeoTag":         seoSetting.ProductPage.Render(seoSetting, product),
			"MicroProduct": seo.MicroProduct{
				Name:        product.Name,
				Description: product.Description,
				BrandName:   product.Category.Name,
				SKU:         product.Code,
				Price:       float64(product.Price),
				Image:       chapter.MainImageUrl(),
			}.Render(),
		},
	)
}
