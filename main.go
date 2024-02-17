package main

import (
	app "github.com/X3ne/hsrpc/src"
	"github.com/X3ne/hsrpc/src/config"
	"github.com/X3ne/hsrpc/src/gui"
	"github.com/X3ne/hsrpc/src/logger"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	app := app.CreateApp(config)

	go app.StartLoop()

	gui.CreateGUI(app)
}
