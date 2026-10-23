package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/NathanFirmo/gorest/internal/db"
	"github.com/NathanFirmo/gorest/internal/gorest"
)

var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "gorest:", err)
		os.Exit(1)
	}
}

func run() error {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return nil
	}

	path, err := db.DefaultPath()
	if err != nil {
		return err
	}
	store, err := db.Open(path)
	if err != nil {
		return err
	}
	defer store.Close()

	app, err := gorest.CreateApp(store)
	if err != nil {
		return err
	}
	return app.Start()
}
