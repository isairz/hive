package admin

import (
	"github.com/qor/exchange"
	"github.com/qor/qor"
	"github.com/isairz/hive/app/models"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/validations"
)

var MangaExchange *exchange.Resource

func init() {
	MangaExchange = exchange.NewResource(&models.Manga{}, exchange.Config{PrimaryField: "Code"})
	MangaExchange.Meta(&exchange.Meta{Name: "Code"})
	MangaExchange.Meta(&exchange.Meta{Name: "Name"})
	MangaExchange.Meta(&exchange.Meta{Name: "Price"})

	MangaExchange.AddValidator(func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
		if utils.ToInt(metaValues.Get("Price").Value) < 100 {
			return validations.NewError(record, "Price", "price can't less than 100")
		}
		return nil
	})
}
