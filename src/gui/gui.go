package gui

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	rpcApp "github.com/X3ne/hsrpc/src"
	"github.com/X3ne/hsrpc/src/logger"
)

type GUI struct {
	App			fyne.App
	Window	fyne.Window
	RPCApp	*rpcApp.App
}

func CreateGUI(rpcApp *rpcApp.App) {
	a := app.New()
	w := a.NewWindow("Honkai RPC")

	g := &GUI{
		App:		a,
		Window:	w,
		RPCApp:	rpcApp,
	}

	g.ConfigApp()
	g.MainScreen()

	w.ShowAndRun()
}

func (g *GUI) importIcon() []byte {
	icon, err := os.Open("assets/icon.png")
	if err != nil {
		logger.Logger.Fatal(err)
	}
	defer icon.Close()

	info, err := icon.Stat()
	if err != nil {
		logger.Logger.Fatal(err)
	}

	data := make([]byte, info.Size())
	_, err = icon.Read(data)
	if err != nil {
		logger.Logger.Fatal(err)
	}

	return data
}

func (g *GUI) ConfigApp() {
	icon := g.importIcon()
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
