package main

import (
	"fmt"
	"net/http"

	"github.com/isairz/hive/config"
	"github.com/isairz/hive/config/admin"
	"github.com/isairz/hive/config/api"
	"github.com/isairz/hive/config/routes"
	_ "github.com/isairz/hive/db/migrations"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", routes.Router())
	admin.Admin.MountTo("/admin", mux)
	api.API.MountTo("/api", mux)

	for _, path := range []string{"system", "downloads", "javascripts", "stylesheets", "images"} {
		mux.Handle(fmt.Sprintf("/%s/", path), http.FileServer(http.Dir("public")))
	}

	fmt.Printf("Listening on: %v\n", config.Config.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), mux); err != nil {
		panic(err)
	}
}
