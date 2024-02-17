package gui

import (
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	rpcApp "github.com/X3ne/hsrpc/src"
	"github.com/X3ne/hsrpc/src/utils"
	"go.uber.org/zap"
)

type GUI struct {
	App			fyne.App
	Window	fyne.Window
	RPCApp	*rpcApp.App
}

var Logger *zap.SugaredLogger

func CreateGUI(rpcApp *rpcApp.App) {
	a := app.New()
	w := a.NewWindow("Honkai RPC")

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	Logger = logger.Sugar()

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
		Logger.Fatal(err)
	}
	defer icon.Close()

	info, err := icon.Stat()
	if err != nil {
		Logger.Fatal(err)
	}

	data := make([]byte, info.Size())
	_, err = icon.Read(data)
	if err != nil {
		Logger.Fatal(err)
	}

	return data
}

func createRectForm(rect *utils.Rect) *widget.Form {
	xEntry := widget.NewEntry()
	yEntry := widget.NewEntry()
	widthEntry := widget.NewEntry()
	heightEntry := widget.NewEntry()

	xEntry.SetText(strconv.Itoa(int(rect.X)))
	yEntry.SetText(strconv.Itoa(int(rect.Y)))
	widthEntry.SetText(strconv.Itoa(int(rect.Width)))
	heightEntry.SetText(strconv.Itoa(int(rect.Height)))

	go func() {
		for {
			if xEntry.Text != "" {
				rect.X = stringToInt(xEntry.Text)
			}

			if yEntry.Text != "" {
				rect.Y = stringToInt(yEntry.Text)
			}

			if widthEntry.Text != "" {
				rect.Width = stringToInt(widthEntry.Text)
			}

			if heightEntry.Text != "" {
				rect.Height = stringToInt(heightEntry.Text)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "X", Widget: xEntry},
			{Text: "Y", Widget: yEntry},
			{Text: "Width", Widget: widthEntry},
			{Text: "Height", Widget: heightEntry},
		},
	}

	return form
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

func presenceString(isGatewayConnected bool) string {
	if isGatewayConnected {
		return "Connected to Discord Gateway"
	}
	return "Not connected to Discord Gateway"
}

func (g *GUI) PresenceTab() *fyne.Container {
	gatewayStatus := binding.NewString()
	location := binding.NewString()
	characterPos := binding.NewString()
	character := binding.NewString()
	detectedMenu := binding.NewString()
	loopTime := binding.NewString()
	isInMenus := binding.NewBool()

	go func() {
		for {
			location.Set(g.RPCApp.AppState.Location.Value)
			posString := strconv.Itoa(int(g.RPCApp.AppState.CharacterPos))
			characterPos.Set(posString)
			character.Set(g.RPCApp.AppState.Character.Value)
			detectedMenu.Set(g.RPCApp.AppState.Menu.Value)
			loopTime.Set(g.RPCApp.AppState.LoopTime.String())
			gatewayStatus.Set(presenceString(g.RPCApp.AppState.IsGatewayConnected))
			isInMenus.Set(g.RPCApp.AppState.IsInMenus)
			time.Sleep(time.Second * 1)
		}
	}()

	container := container.NewVBox(
		widget.NewButton("Reconnect to Discord", func() {
			gatewayStatus.Set("Connecting to Discord Gateway...")
			time.Sleep(1 * time.Second)
			if err := g.RPCApp.ConnectToDiscordGateway(); err != nil {
				gatewayStatus.Set("Error when connecting to Discord Gateway")
			} else {
				gatewayStatus.Set("Connected to Discord Gateway")
			}
		}),
		widget.NewLabelWithData(gatewayStatus),
		container.NewHBox(
			widget.NewLabel("Location:"), widget.NewLabelWithData(location),
		),
		container.NewHBox(
			widget.NewLabel("Character position:"), widget.NewLabelWithData(characterPos),
		),
		container.NewHBox(
			widget.NewLabel("Character:"), widget.NewLabelWithData(character),
		),
		container.NewHBox(
			widget.NewLabel("Detected menu:"), widget.NewLabelWithData(detectedMenu),
		),
		container.NewHBox(
			widget.NewLabel("Loop time:"), widget.NewLabelWithData(loopTime),
		),
		container.NewHBox(
			widget.NewLabel("Is in menus:"), widget.NewLabelWithData(binding.BoolToString(isInMenus)),
		),
	)

	if err := g.RPCApp.ConnectToDiscordGateway(); err != nil {
		gatewayStatus.Set("Error when connecting to Discord Gateway")
	} else {
		gatewayStatus.Set("Connected to Discord Gateway")
	}

	return container
}

func stringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		Logger.Error(err)
	}
	return i
}

func (g *GUI) ConfigCharactersTab() *fyne.Container {
	char1 := createRectForm(&g.RPCApp.Config.CharactersCoords[0])
	char2 := createRectForm(&g.RPCApp.Config.CharactersCoords[1])
	char3 := createRectForm(&g.RPCApp.Config.CharactersCoords[2])
	char4 := createRectForm(&g.RPCApp.Config.CharactersCoords[3])

	box1 := createRectForm(&g.RPCApp.Config.CharacterBoxCoords[0])
	box2 := createRectForm(&g.RPCApp.Config.CharacterBoxCoords[1])
	box3 := createRectForm(&g.RPCApp.Config.CharacterBoxCoords[2])
	box4 := createRectForm(&g.RPCApp.Config.CharacterBoxCoords[3])

	container := container.NewVBox(
		container.NewAdaptiveGrid(
			2,
			container.NewVBox(widget.NewLabel("Character 1 coords"), char1),
			container.NewVBox(widget.NewLabel("Character 2 coords"), char2),
			container.NewVBox(widget.NewLabel("Character 3 coords"), char3),
			container.NewVBox(widget.NewLabel("Character 4 coords"), char4),
		),
		container.NewAdaptiveGrid(
			2,
			container.NewVBox(widget.NewLabel("Character 1 box coords"), box1),
			container.NewVBox(widget.NewLabel("Character 2 box coords"), box2),
			container.NewVBox(widget.NewLabel("Character 3 box coords"), box3),
			container.NewVBox(widget.NewLabel("Character 4 box coords"), box4),
		),
	)

	return container
}

func (g *GUI) ConfigLocationsTab() *fyne.Container {
	form1 := createRectForm(&g.RPCApp.Config.LocationCoord)

	container := container.NewAdaptiveGrid(
		2,
		container.NewVBox(widget.NewLabel("Location coords"), form1),
	)

	return container
}

func (g *GUI) ConfigLoopTab() *fyne.Container {
	timeEntry := widget.NewEntry()

	timeEntry.SetText(strconv.Itoa(int(g.RPCApp.AppState.LoopTime.Milliseconds())))

	go func() {
		for {
			if timeEntry.Text != "" {
				loopTime, err := strconv.Atoi(timeEntry.Text)
				if err != nil {
					Logger.Error(err)
				}
				g.RPCApp.Config.LoopTime = time.Duration(loopTime) * time.Millisecond
			}

			time.Sleep(1 * time.Second)
		}
	}()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Loop time (ms)", Widget: timeEntry},
		},
	}

	container := container.NewVBox(
		container.NewVBox(form),
	)

	return container
}

func (g *GUI) MainScreen() {
	// TODO: Add store to save the config
	tabs := container.NewAppTabs(
		container.NewTabItem("Presence", container.NewVBox(
			g.PresenceTab(),
		)),
		container.NewTabItem("Characters", container.NewVBox(
			g.ConfigCharactersTab(),
		)),
		container.NewTabItem("Locations", container.NewVBox(
			g.ConfigLocationsTab(),
		)),
		container.NewTabItem("Loop", container.NewVBox(
			g.ConfigLoopTab(),
		)),
	)

	g.Window.SetContent(tabs)
}
