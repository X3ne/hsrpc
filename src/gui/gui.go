package gui

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	rpcApp "github.com/X3ne/hsrpc/src"
	"github.com/X3ne/hsrpc/src/consts"
	"github.com/X3ne/hsrpc/src/logger"
	"github.com/creativeprojects/go-selfupdate"
)

type GUI struct {
	App			fyne.App
	Window	fyne.Window
	RPCApp	*rpcApp.App
}

func (g *GUI) UpdateApplicationGUI( updateCompleted chan bool) {
	updateLabel := widget.NewLabel("Downloading update...")
	progressBar := widget.NewProgressBarInfinite()

	g.Window.SetContent(container.NewVBox(
		updateLabel,
		progressBar,
	))

	g.Window.SetFixedSize(true)
	g.Window.Resize(fyne.NewSize(300, 100))
	g.Window.CenterOnScreen()

	exe, err := os.Executable()
	if err != nil {
		HandleUpdateError("Could not locate executable path: %s", err, updateLabel, updateCompleted, g.Window)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	update, err := GetLatestUpdate()
	if err != nil {
		HandleUpdateError("Error occurred while checking for updates: %s", err, updateLabel, updateCompleted, g.Window)
		return
	}

	updateLabel.SetText("Updating to version " + update.Version())
	if err := selfupdate.UpdateTo(ctx, update.AssetURL, update.AssetName, exe); err != nil {
		HandleUpdateError("Error occurred while updating binary: %s", err, updateLabel, updateCompleted, g.Window)
		return
	}

	logger.Logger.Infof("Successfully updated to version %s", update.Version)
	updateLabel.SetText("Update completed. Restarting...")
	g.Window.SetFixedSize(false)

	updateCompleted <- true
}

func CreateGUI(rpcApp *rpcApp.App) {
	a := app.New()
	w := a.NewWindow("Honkai RPC " + consts.Version)

	g := &GUI{
		App:		a,
		Window:	w,
		RPCApp:	rpcApp,
	}

	g.ConfigApp()

	update, err := GetLatestUpdate()
	if err != nil {
		logger.Logger.Errorf("Error occurred while checking for updates: %s", err)
	}

	g.MainScreen(update)

	w.ShowAndRun()
}

func (g *GUI) ConfigApp() {
	icon := ImportIcon()
	// g.Window.SetFixedSize(true)
	// g.Window.Resize(fyne.NewSize(800, 600))

	g.App.SetIcon(fyne.NewStaticResource("icon.png", icon))
	g.Window.SetIcon(g.App.Icon())
	g.Window.CenterOnScreen()

	if desk, ok := g.App.(desktop.App); ok {
		m := fyne.NewMenu("Honkai RPC",
			fyne.NewMenuItem("Config", func() {
				g.Window.Show()
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(g.App.Icon())
	}

	g.Window.SetCloseIntercept(func() {
		g.Window.Hide()
	})

	g.App.Settings().SetTheme(&appTheme{})
}

func (g *GUI) updateApp(confirm bool) {
	if !confirm {
		return
	}

	updateCompleted := make(chan bool)

	go func() {
		success := <-updateCompleted
		log.Println("Update completed")

		time.Sleep(1 * time.Second)

		if success {
			Restart()
		} else {
			g.Window.SetFixedSize(false)

			g.MainScreen()

			g.Window.CenterOnScreen()
		}
	}()

	g.UpdateApplicationGUI(updateCompleted)

	g.Window.SetFixedSize(true)
	g.Window.Resize(fyne.NewSize(300, 100))
	g.Window.CenterOnScreen()
}

func (g *GUI) MainScreen(update ...*selfupdate.Release) {
	tabs := container.NewAppTabs(
		container.NewTabItem("Presence", container.NewPadded(
			g.createPresenceTab(),
		)),
		container.NewTabItem("Config", container.NewPadded(
			g.createGlobalConfigTab(),
		)),
		// container.NewTabItem("Characters", container.NewVBox(
		// 	g.createCharactersTab(),
		// )),
		// container.NewTabItem("Locations", container.NewVBox(
		// 	g.createGuiCoordsConfigTab(),
		// )),
	)

	if len(update) > 0 && update[0] != nil {
		url, _ := url.Parse("https://github.com/X3ne/hsrpc/releases/latest")

		dialog.ShowCustomConfirm("Update available", "Update", "Update Later", container.NewVBox(
			widget.NewLabel(fmt.Sprintf("An update is available: %s", update[0].Version())),
			widget.NewHyperlink("Release notes", url),
		), g.updateApp, g.Window)
	}

	g.Window.SetContent(tabs)
}
