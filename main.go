//go:generate go-winres simply --icon assets/icon.png --manifest gui

package main

import (
	rpcApp "github.com/X3ne/hsrpc/src"
	"github.com/X3ne/hsrpc/src/config"
	"github.com/X3ne/hsrpc/src/gui"
	"github.com/X3ne/hsrpc/src/logger"
	"github.com/X3ne/hsrpc/src/utils"
)

func main() {
	guiApp := gui.CreateApp()

	defer func() {
		if r := recover(); r != nil {
			utils.PanicRecover(r, guiApp)
		}
	}()

	config, err := config.LoadConfig()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	app := rpcApp.CreateApp(config, guiApp)

	go app.StartLoop()

	gui.CreateGUI(app, guiApp)

	guiApp.Run()
}
