package admin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/qor/activity"
	"github.com/qor/admin"
	"github.com/qor/i18n/exchange_actions"
	"github.com/qor/media_library"
	"github.com/qor/qor"
	"github.com/isairz/hive/app/models"
	"github.com/isairz/hive/config"
	"github.com/isairz/hive/db"
	"github.com/qor/qor/resource"
	"github.com/qor/qor/utils"
	"github.com/qor/transition"
	"github.com/qor/validations"
)

var Admin *admin.Admin
var Languages = []string{"Korean", "Japanese"}

func init() {
	Admin = admin.New(&qor.Config{DB: db.Publish.DraftDB()})
	Admin.SetSiteName("Spica Terrible!!")
	Admin.SetAuth(Auth{})

	// Add Dashboard
	Admin.AddMenu(&admin.Menu{Name: "Dashboard", Link: "/admin"})

	// Add Asset Manager, for rich editor
	assetManager := Admin.AddResource(&media_library.AssetManager{}, &admin.Config{Invisible: false})

	// Add Manga
	manga := Admin.AddResource(&models.Manga{}, &admin.Config{Menu: []string{"Manga Management"}})
	manga.Meta(&admin.Meta{Name: "MadeCountry", Type: "select_one", Collection: Languages})
	manga.Meta(&admin.Meta{Name: "Description", Type: "rich_editor", Resource: assetManager})

	chapterMeta := manga.Meta(&admin.Meta{Name: "Chapters"})
	chapter := chapterMeta.Resource
	chapter.NewAttrs("-Manga")
	chapter.EditAttrs("-Manga")

	manga.SearchAttrs("Name", "Category.Name")
	manga.EditAttrs(
		&admin.Section{
			Title: "Basic Information",
			Rows: [][]string{
				{"Name"},
				{"Enabled"},
			}},
		&admin.Section{
			Title: "Organization",
			Rows: [][]string{
				{"Category", "Tags", "MadeCountry"},
			}},
		"Description",
		"Chapters",
	)

	for _, country := range Languages {
		var country = country
		manga.Scope(&admin.Scope{Name: country, Group: "Made Country", Handle: func(db *gorm.DB, ctx *qor.Context) *gorm.DB {
			return db.Where("made_country = ?", country)
		}})
	}

	manga.IndexAttrs("-Chapters")

	manga.Action(&admin.Action{
		Name: "View On Site",
		URL: func(record interface{}, context *admin.Context) string {
			if manga, ok := record.(*models.Manga); ok {
				return fmt.Sprintf("/mangas/%v", manga.ID)
			}
			return "#"
		},
		Modes: []string{"menu_item", "edit"},
	})

	manga.Action(&admin.Action{
		Name: "Disable",
		Handle: func(arg *admin.ActionArgument) error {
			for _, record := range arg.FindSelectedRecords() {
				arg.Context.DB.Model(record.(*models.Manga)).Update("enabled", false)
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if manga, ok := record.(*models.Manga); ok {
				return manga.Enabled == true
			}
			return true
		},
		Modes: []string{"index", "edit", "menu_item"},
	})

	manga.Action(&admin.Action{
		Name: "Enable",
		Handle: func(arg *admin.ActionArgument) error {
			for _, record := range arg.FindSelectedRecords() {
				arg.Context.DB.Model(record.(*models.Manga)).Update("enabled", true)
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if manga, ok := record.(*models.Manga); ok {
				return manga.Enabled == false
			}
			return true
		},
		Modes: []string{"index", "edit", "menu_item"},
	})

	Admin.AddResource(&models.Category{}, &admin.Config{Menu: []string{"Manga Management"}})
	Admin.AddResource(&models.Tag{}, &admin.Config{Menu: []string{"Manga Management"}})

	// Add Order
	order := Admin.AddResource(&models.Order{}, &admin.Config{Menu: []string{"Order Management"}, Invisible: true})
	order.Meta(&admin.Meta{Name: "ShippingAddress", Type: "single_edit"})
	order.Meta(&admin.Meta{Name: "BillingAddress", Type: "single_edit"})
	order.Meta(&admin.Meta{Name: "ShippedAt", Type: "date"})

	orderItemMeta := order.Meta(&admin.Meta{Name: "OrderItems"})
	orderItemMeta.Resource.NewAttrs("-State")
	orderItemMeta.Resource.EditAttrs("-State")

	// define scopes for Order
	for _, state := range []string{"checkout", "cancelled", "paid", "paid_cancelled", "processing", "shipped", "returned"} {
		var state = state
		order.Scope(&admin.Scope{
			Name:  state,
			Label: strings.Title(strings.Replace(state, "_", " ", -1)),
			Group: "Order Status",
			Handle: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where(models.Order{Transition: transition.Transition{State: state}})
			},
		})
	}

	// define actions for Order
	type trackingNumberArgument struct {
		TrackingNumber string
	}

	order.Action(&admin.Action{
		Name: "Processing",
		Handle: func(argument *admin.ActionArgument) error {
			for _, order := range argument.FindSelectedRecords() {
				db := argument.Context.GetDB()
				if err := models.OrderState.Trigger("process", order.(*models.Order), db); err != nil {
					return err
				}
				db.Select("state").Save(order)
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*models.Order); ok {
				return order.State == "paid"
			}
			return false
		},
		Modes: []string{"show", "menu_item"},
	})
	order.Action(&admin.Action{
		Name: "Ship",
		Handle: func(argument *admin.ActionArgument) error {
			var (
				tx                     = argument.Context.GetDB().Begin()
				trackingNumberArgument = argument.Argument.(*trackingNumberArgument)
			)

			if trackingNumberArgument.TrackingNumber != "" {
				for _, record := range argument.FindSelectedRecords() {
					order := record.(*models.Order)
					order.TrackingNumber = &trackingNumberArgument.TrackingNumber
					models.OrderState.Trigger("ship", order, tx, "tracking number "+trackingNumberArgument.TrackingNumber)
					if err := tx.Save(order).Error; err != nil {
						tx.Rollback()
						return err
					}
				}
			} else {
				return errors.New("invalid shipment number")
			}

			tx.Commit()
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*models.Order); ok {
				return order.State == "processing"
			}
			return false
		},
		Resource: Admin.NewResource(&trackingNumberArgument{}),
		Modes:    []string{"show", "menu_item"},
	})

	order.Action(&admin.Action{
		Name: "Cancel",
		Handle: func(argument *admin.ActionArgument) error {
			for _, order := range argument.FindSelectedRecords() {
				db := argument.Context.GetDB()
				if err := models.OrderState.Trigger("cancel", order.(*models.Order), db); err != nil {
					return err
				}
				db.Select("state").Save(order)
			}
			return nil
		},
		Visible: func(record interface{}, context *admin.Context) bool {
			if order, ok := record.(*models.Order); ok {
				for _, state := range []string{"draft", "checkout", "paid", "processing"} {
					if order.State == state {
						return true
					}
				}
			}
			return false
		},
		Modes: []string{"index", "show", "menu_item"},
	})

	order.IndexAttrs("User", "PaymentAmount", "ShippedAt", "CancelledAt", "State", "ShippingAddress")
	order.NewAttrs("-DiscountValue", "-AbandonedReason", "-CancelledAt")
	order.EditAttrs("-DiscountValue", "-AbandonedReason", "-CancelledAt", "-State")
	order.ShowAttrs("-DiscountValue", "-State")
	order.SearchAttrs("User.Name", "User.Email", "ShippingAddress.ContactName", "ShippingAddress.Address1", "ShippingAddress.Address2")

	// Add activity for order
	activity.Register(order)

	// Define another resource for same model
	abandonedOrder := Admin.AddResource(&models.Order{}, &admin.Config{Name: "Abandoned Order", Menu: []string{"Order Management"}, Invisible: true})
	abandonedOrder.Meta(&admin.Meta{Name: "ShippingAddress", Type: "single_edit"})
	abandonedOrder.Meta(&admin.Meta{Name: "BillingAddress", Type: "single_edit"})

	// Define default scope for abandoned orders
	abandonedOrder.Scope(&admin.Scope{
		Default: true,
		Handle: func(db *gorm.DB, context *qor.Context) *gorm.DB {
			return db.Where("abandoned_reason IS NOT NULL AND abandoned_reason <> ?", "")
		},
	})

	// Define scopes for abandoned orders
	for _, amount := range []int{5000, 10000, 20000} {
		var amount = amount
		abandonedOrder.Scope(&admin.Scope{
			Name:  fmt.Sprint(amount),
			Group: "Amount Greater Than",
			Handle: func(db *gorm.DB, context *qor.Context) *gorm.DB {
				return db.Where("payment_amount > ?", amount)
			},
		})
	}

	abandonedOrder.IndexAttrs("-ShippingAddress", "-BillingAddress", "-DiscountValue", "-OrderItems")
	abandonedOrder.NewAttrs("-DiscountValue")
	abandonedOrder.EditAttrs("-DiscountValue")
	abandonedOrder.ShowAttrs("-DiscountValue")

	// Add Store
	store := Admin.AddResource(&models.Store{}, &admin.Config{Menu: []string{"Store Management"}})
	store.AddValidator(func(record interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
		if meta := metaValues.Get("Name"); meta != nil {
			if name := utils.ToString(meta.Value); strings.TrimSpace(name) == "" {
				return validations.NewError(record, "Name", "Name can't be blank")
			}
		}
		return nil
	})

	// Add Translations
	Admin.AddResource(config.Config.I18n, &admin.Config{Menu: []string{"Site Management"}})

	// Add SEOSetting
	Admin.AddResource(&models.SEOSetting{}, &admin.Config{Menu: []string{"Site Management"}, Singleton: true})

	// Add Setting
	Admin.AddResource(&models.Setting{}, &admin.Config{Singleton: true})

	// Add User
	user := Admin.AddResource(&models.User{})
	user.Meta(&admin.Meta{Name: "Gender", Type: "select_one", Collection: []string{"Male", "Female", "Unknown"}})

	user.IndexAttrs("ID", "Email", "Name", "Gender", "Role")
	user.ShowAttrs(
		&admin.Section{
			Title: "Basic Information",
			Rows: [][]string{
				{"Name"},
				{"Email", "Password"},
				{"Gender", "Role"},
			}},
		"Addresses",
	)
	user.EditAttrs(user.ShowAttrs())

	// Add Publish
	Admin.AddResource(db.Publish, &admin.Config{Singleton: true})

	// Add Worker
	Worker := getWorker()
	Admin.AddResource(Worker)
	exchange_actions.RegisterExchangeJobs(config.Config.I18n, Worker)

	// Add Search Center Resources
	Admin.AddSearchResource(manga, user, order)

	initFuncMap()
	initRouter()
}
