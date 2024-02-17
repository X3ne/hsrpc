package main

import (
	app "github.com/X3ne/hsrpc/src"
	"github.com/X3ne/hsrpc/src/config"
	"github.com/X3ne/hsrpc/src/gui"
)

func main() {
	config := config.NewConfig()

	app := app.CreateApp(config)

	go app.StartLoop()

	gui.CreateGUI(app)
}
