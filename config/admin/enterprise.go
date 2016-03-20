// +build enterprise

package admin

import (
	"fmt"

	"github.com/qor/qor"
	"github.com/isairz/hive/app/models"
	"github.com/isairz/hive/db"
	"github.com/qor/admin"
	"github.com/theplant/qor-enterprise/promotion"
)

func init() {
	// Benefits Definations
	discountRateArgument := Admin.NewResource(&struct {
		Percentage uint
	}{})
	discountRateArgument.Meta(&admin.Meta{Name: "Percentage", Label: "Percentage (e.g enter 10 for a 10% discount)"})
	promotion.RegisterBenefitHandler(promotion.BenefitHandler{
		Name:     "Discount Rate",
		Resource: discountRateArgument,
	})

	discountAmountArgument := Admin.NewResource(&struct {
		Amount float32
	}{})
	discountAmountArgument.Meta(&admin.Meta{Name: "Amount", Label: "Amount (e.g enter 10 for a $10 discount)"})
	promotion.RegisterBenefitHandler(promotion.BenefitHandler{
		Name:     "Discount Amount",
		Resource: discountAmountArgument,
	})

	promotion.RegisterBenefitHandler(promotion.BenefitHandler{
		Name: "Shipping Fee",
		Resource: Admin.NewResource(&struct {
			Price float32
		}{}),
	})

	promotion.RegisterBenefitHandler(promotion.BenefitHandler{
		Name: "2nd Day Shipping Fee",
		Resource: Admin.NewResource(&struct {
			Price float32
		}{}),
	})

	mangaCodeCollection := func(value interface{}, context *qor.Context) [][]string {
		var mangas []models.Manga
		var results [][]string
		context.GetDB().Find(&mangas)
		for _, manga := range mangas {
			results = append(results, []string{fmt.Sprint(manga.ID), manga.Code})
		}
		return results
	}
	combinedDiscountArgument := Admin.NewResource(&struct {
		MangaCodes []string
		Category     string
		Quantity     uint
		Price        float32
		Percentage   uint
		Discount     uint
	}{})
	combinedDiscountArgument.Meta(&admin.Meta{Name: "MangaCodes", Type: "select_many", Collection: mangaCodeCollection})
	combinedDiscountArgument.Meta(&admin.Meta{Name: "Category", Type: "select_one", Collection: []string{"All Mangas", "Bags", "Summer Shirts", "Pants"}})
	combinedDiscountArgument.Meta(&admin.Meta{Name: "Percentage", Label: "Discount Percentage (e.g enter 10 for a 10% discount)"})
	combinedDiscountArgument.Meta(&admin.Meta{Name: "Discount", Label: "Discount Amount (e.g enter 10 for a $10 discount)"})
	promotion.RegisterBenefitHandler(promotion.BenefitHandler{
		Name:     "Combined Discounts",
		Resource: combinedDiscountArgument,
	})

	// Rules Definations
	amountGreaterThanArgument := Admin.NewResource(&struct {
		Amount   int
		Category string
	}{})
	amountGreaterThanArgument.Meta(&admin.Meta{Name: "Category", Type: "select_one", Collection: []string{"All Mangas", "Bags", "Summer Shirts", "Pants"}})
	promotion.RegisterRuleHandler(promotion.RuleHandler{
		Name:     "Amount Greater Than",
		Resource: amountGreaterThanArgument,
	})

	quantityGreaterThanArgument := Admin.NewResource(&struct {
		MangaCodes []string
		Category     string
		Quantity     int
	}{})
	quantityGreaterThanArgument.Meta(&admin.Meta{Name: "MangaCodes", Type: "select_many", Collection: mangaCodeCollection})
	quantityGreaterThanArgument.Meta(&admin.Meta{Name: "Category", Type: "select_one", Collection: []string{"All Mangas", "Bags", "Summer Shirts", "Pants"}})
	promotion.RegisterRuleHandler(promotion.RuleHandler{
		Name:     "Quantity Greater Than",
		Resource: quantityGreaterThanArgument,
	})

	userGroupArgument := Admin.NewResource(&struct {
		Group string
	}{})
	userGroupArgument.Meta(&admin.Meta{Name: "Group", Type: "select_one", Collection: []string{"VIP", "Employee", "Normal"}})
	promotion.RegisterRuleHandler(promotion.RuleHandler{
		Name:     "User Group",
		Resource: userGroupArgument,
	})

	promotion.RegisterRuleHandler(promotion.RuleHandler{
		Name: "From Link",
		Resource: Admin.NewResource(&struct {
			VariableName string
			Value        string
		}{}),
	})

	hasMangargument := Admin.NewResource(&struct {
		MangaCodes []string
		Category     string
	}{})
	hasMangargument.Meta(&admin.Meta{Name: "MangaCodes", Type: "select_many", Collection: mangaCodeCollection})
	hasMangargument.Meta(&admin.Meta{Name: "Category", Type: "select_one", Collection: []string{"All Mangas", "Bags", "Summer Shirts", "Pants"}})
	promotion.RegisterRuleHandler(promotion.RuleHandler{
		Name:     "Has Manga",
		Resource: hasMangargument,
	})

	// Auto migrations
	promotion.AutoMigrate(db.DB)

	// Add Promotions to Admin
	Admin.AddResource(&promotion.PromotionDiscount{}, &admin.Config{Name: "Promotions", Menu: []string{"Site Management"}})
}
