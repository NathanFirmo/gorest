package main

import (
	"github.com/NathanFirmo/gorest/internal/gorest"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	app := gorest.CreateApp()

	if err := app.Start(); err != nil {
		panic(err)
	}
}
