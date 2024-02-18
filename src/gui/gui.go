package gui

import (
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	rpcApp "github.com/X3ne/hsrpc/src"
	"github.com/X3ne/hsrpc/src/consts"
)

type GUI struct {
	App			fyne.App
	Window	fyne.Window
	RPCApp	*rpcApp.App
}

// func update() error {
// 	version := consts.Version
// 	latest, found, err := selfupdate.DetectLatest(context.Background(), selfupdate.ParseSlug("X3ne/hsrpc"))
// 	if err != nil {
// 		return fmt.Errorf("error occurred while detecting version: %w", err)
// 	}
// 	if !found {
// 		return fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
// 	}

// 	if latest.LessOrEqual(version) {
// 		log.Printf("Current version (%s) is the latest", version)
// 		return nil
// 	}

// 	exe, err := os.Executable()
// 	if err != nil {
// 		return errors.New("could not locate executable path")
// 	}
// 	if err := selfupdate.UpdateTo(context.Background(), latest.AssetURL, latest.AssetName, exe); err != nil {
// 		return fmt.Errorf("error occurred while updating binary: %w", err)
// 	}
// 	logger.Logger.Infof("Successfully updated to version %s", latest.Version)
// 	return nil
// }

func (g *GUI) UpdateApplicationGUI(updateCompleted chan bool) {
	updateLabel := widget.NewLabel("Looking for updates...")
	progressBar := widget.NewProgressBarInfinite()

	g.Window.SetFixedSize(true)
	g.Window.Resize(fyne.NewSize(300, 100))

	g.Window.SetContent(container.NewVBox(
		updateLabel,
		progressBar,
	))

	go func() {
		// if err := update(); err != nil {
		// 	updateLabel.SetText("Update failed")
		// 	logger.Logger.Error(err)
		// 	g.Window.SetFixedSize(false)
		// 	updateCompleted <- true
		// 	return
		// }

		updateLabel.SetText("Update completed")
		g.Window.SetFixedSize(false)

		updateCompleted <- true
	}()
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

	updateCompleted := make(chan bool)

	go g.UpdateApplicationGUI(updateCompleted)

	go func()  {
		<-updateCompleted
		log.Println("Update completed")

		time.Sleep(1 * time.Second)

		g.MainScreen()
		g.Window.CenterOnScreen()
	}()

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
}

func (g *GUI) MainScreen() {
	tabs := container.NewAppTabs(
		container.NewTabItem("Presence", container.NewVBox(
			g.createPresenceTab(),
		)),
		container.NewTabItem("Config", container.NewVBox(
			g.createGlobalConfigTab(),
		)),
		container.NewTabItem("Characters", container.NewVBox(
			g.createCharactersTab(),
		)),
		container.NewTabItem("Locations", container.NewVBox(
			g.createGuiCoordsConfigTab(),
		)),
	)

	g.Window.SetContent(tabs)
}
