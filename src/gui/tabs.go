package gui

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/X3ne/hsrpc/src/logger"
	"github.com/lxn/walk"
)

func (g *GUI) createCharactersTab() *fyne.Container {
	char1 := createRectForm(&g.RPCApp.Config.GUICoordsConfig.CharactersCoords[0], g.saveConfig)
	char2 := createRectForm(&g.RPCApp.Config.GUICoordsConfig.CharactersCoords[1], g.saveConfig)
	char3 := createRectForm(&g.RPCApp.Config.GUICoordsConfig.CharactersCoords[2], g.saveConfig)
	char4 := createRectForm(&g.RPCApp.Config.GUICoordsConfig.CharactersCoords[3], g.saveConfig)

	box1 := createRectForm(&g.RPCApp.Config.GUICoordsConfig.CharactersBoxCoords[0], g.saveConfig)
	box2 := createRectForm(&g.RPCApp.Config.GUICoordsConfig.CharactersBoxCoords[1], g.saveConfig)
	box3 := createRectForm(&g.RPCApp.Config.GUICoordsConfig.CharactersBoxCoords[2], g.saveConfig)
	box4 := createRectForm(&g.RPCApp.Config.GUICoordsConfig.CharactersBoxCoords[3], g.saveConfig)

	container := container.NewVBox(
		container.NewAdaptiveGrid(
			4,
			container.NewVBox(widget.NewLabel("Character 1 coords"), char1),
			container.NewVBox(widget.NewLabel("Character 2 coords"), char2),
			container.NewVBox(widget.NewLabel("Character 3 coords"), char3),
			container.NewVBox(widget.NewLabel("Character 4 coords"), char4),
		),
		container.NewAdaptiveGrid(
			4,
			container.NewVBox(widget.NewLabel("Character 1 box coords"), box1),
			container.NewVBox(widget.NewLabel("Character 2 box coords"), box2),
			container.NewVBox(widget.NewLabel("Character 3 box coords"), box3),
			container.NewVBox(widget.NewLabel("Character 4 box coords"), box4),
		),
	)

	return container
}

func (g *GUI) createPresenceTab() *fyne.Container {
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
			loopTime.Set((g.RPCApp.AppState.LoopTime * time.Millisecond).String())
			gatewayStatus.Set(getPresenceString(g.RPCApp.AppState.IsGatewayConnected))
			isInMenus.Set(g.RPCApp.AppState.IsInMenus)
			time.Sleep(g.RPCApp.AppState.LoopTime * time.Millisecond)
		}
	}()

	container := container.NewVBox(
		widget.NewButton("Reconnect to Discord", func() {
			gatewayStatus.Set("Connecting to Discord Gateway...")
			go func() {
				time.Sleep(1 * time.Second)
				if err := g.RPCApp.ConnectToDiscordGateway(); err != nil {
					gatewayStatus.Set("Error when connecting to Discord Gateway")
				} else {
					gatewayStatus.Set("Connected to Discord Gateway")
				}
			}()
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

	if !g.RPCApp.AppState.IsGatewayConnected {
		gatewayStatus.Set("Error when connecting to Discord Gateway")
	} else {
		gatewayStatus.Set("Connected to Discord Gateway")
	}

	return container
}

func (g *GUI) createGuiCoordsConfigTab() *fyne.Container {
	locationForm := createRectForm(&g.RPCApp.Config.GUICoordsConfig.LocationCoord, g.saveConfig)
	escLocationForm := createRectForm(&g.RPCApp.Config.GUICoordsConfig.EscCoord, g.saveConfig)
	menuForm := createRectForm(&g.RPCApp.Config.GUICoordsConfig.MenusCoord, g.saveConfig)
	combatForm := createRectForm(&g.RPCApp.Config.GUICoordsConfig.CombatCoord, g.saveConfig)

	container := container.NewAdaptiveGrid(
		2,
		container.NewVBox(widget.NewLabel("Location coords"), locationForm),
		container.NewVBox(widget.NewLabel("ESC coords"), escLocationForm),
		container.NewVBox(widget.NewLabel("Menus coords"), menuForm),
		container.NewVBox(widget.NewLabel("Combat coords"), combatForm),
	)

	return container
}

func (g *GUI) createGlobalConfigTab() *fyne.Container {
	// startupCheckbox := widget.NewCheck("Start with Windows", func(b bool) {
	// 	g.RPCApp.Config.StartWithWindows = b
	// 	g.saveConfig()
	// })
	uuidText := widget.NewEntry()
	playerName := widget.NewEntry()
	timeEntry := widget.NewEntry()
	tesseractPath := widget.NewEntry()
	preprocessThreshold := widget.NewEntry()
	tesseractPathButton := widget.NewButton("Browse", func() {
		dlg := new(walk.FileDialog)

		dlg.Title = "Select Tesseract executable"
		dlg.Filter = "Executable files (*.exe)|*.exe"

		if ok, _ := dlg.ShowOpen(nil); ok {
			tesseractPath.SetText(dlg.FilePath)
			g.RPCApp.Config.TesseractPath = dlg.FilePath
			g.saveConfig()
		}
	})
	windowClassName := widget.NewEntry()
	windowName := widget.NewEntry()
	displayLevelCheckbox := widget.NewCheck("Display your level", func(b bool) {
		g.RPCApp.Config.DisplayLevel = b
		g.saveConfig()
	})
	displayNicknameCheckbox := widget.NewCheck("Display your nickname", func(b bool) {
		g.RPCApp.Config.DisplayNickname = b
		g.saveConfig()
	})

	displayLevelCheckbox.SetChecked(g.RPCApp.Config.DisplayLevel)
	displayNicknameCheckbox.SetChecked(g.RPCApp.Config.DisplayNickname)

	uuidText.SetText(g.RPCApp.Config.PlayerUID)

	uuidText.OnChanged = debounce(func(s string) {
		g.RPCApp.Config.PlayerUID = s
		g.saveConfig()
	}, 500*time.Millisecond)

	playerName.SetText(g.RPCApp.Config.PlayerName)

	playerName.OnChanged = debounce(func(s string) {
		g.RPCApp.Config.PlayerName = s
		g.saveConfig()
	}, 500*time.Millisecond)

	timeEntry.SetText(strconv.Itoa(int((g.RPCApp.AppState.LoopTime * time.Millisecond).Milliseconds())))

	timeEntry.OnChanged = debounce(func(s string) {
		loopTime, err := strconv.Atoi(s)
		if err != nil {
			logger.Logger.Error(err)
		}
		g.RPCApp.Config.LoopTime = time.Duration(loopTime)
		g.saveConfig()
	}, 500*time.Millisecond)

	tesseractPath.SetText(g.RPCApp.Config.TesseractPath)

	tesseractPath.OnChanged = debounce(func(s string) {
		g.RPCApp.Config.TesseractPath = s
		g.saveConfig()
	}, 500*time.Millisecond)

	preprocessThreshold.SetText(strconv.Itoa(g.RPCApp.Config.PreprocessThreshold))

	preprocessThreshold.OnChanged = debounce(func(s string) {
		preprocessThreshold, err := strconv.Atoi(s)
		if err != nil {
			logger.Logger.Error(err)
		}
		g.RPCApp.Config.PreprocessThreshold = preprocessThreshold
		g.saveConfig()
	}, 500*time.Millisecond)

	windowClassName.SetText(g.RPCApp.Config.WindowClass)

	windowClassName.OnChanged = debounce(func(s string) {
		g.RPCApp.Config.WindowClass = s
		g.saveConfig()
	}, 500*time.Millisecond)

	windowName.SetText(g.RPCApp.Config.WindowName)

	windowName.OnChanged = debounce(func(s string) {
		g.RPCApp.Config.WindowName = s
		g.saveConfig()
	}, 500*time.Millisecond)

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Player UID", Widget: uuidText},
			{Text: "Player name", Widget: playerName},
			{Text: "Loop time (ms)", Widget: timeEntry},
			{Text: "Preprocess threshold", Widget: preprocessThreshold},
			{Text: "Tesseract path", Widget: container.NewAdaptiveGrid(2, tesseractPath, tesseractPathButton)},
			{Text: "Window class", Widget: windowClassName},
			{Text: "Window name", Widget: windowName},
		},
	}

	container := container.NewVBox(
		// startupCheckbox,
		container.NewVBox(form),
		displayLevelCheckbox,
		displayNicknameCheckbox,
	)

	return container
}
