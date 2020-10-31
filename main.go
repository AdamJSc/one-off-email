package main

import (
	"log"
	"one-off-email/app"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	deps := dependencies{
		config: app.MustParseConfig(".env"),
	}

	log.Println(deps.Config())
}

// dependencies implements app.Dependencies
type dependencies struct {
	config *app.Config
}

// Config implements app.Dependencies.Config()
func (d *dependencies) Config() *app.Config { return d.config }
