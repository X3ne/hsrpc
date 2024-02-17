package app

import (
	"syscall"
	"time"

	"github.com/lxn/win"

	"github.com/X3ne/hsrpc/src/config"
	"github.com/X3ne/hsrpc/src/logger"
	"github.com/X3ne/hsrpc/src/utils"
	"github.com/hugolgst/rich-go/client"
)

type App struct {
	Config			config.AppConfig
	HWND				win.HWND
	AppState		AppState
}

type AppState struct {
	CharacterPos							int32
	Character									utils.Data
	Location									utils.Data
	Menu											utils.Data
	LoopTime									time.Duration
	AppStarted								time.Time
	CombatStarted							time.Time
	IsInMenus									bool
	IsGatewayConnected				bool
	IsOCRInitialized					bool
}

func CreateApp(config config.AppConfig) *App {
	utils.LoadGameData()

	return &App{
		Config:				config,
		HWND:					0,
		AppState:			AppState{
			Location: utils.Data{
				AssetID: "menu_lost",
				Value: "Lost in the space-time continuum",
			},
			LoopTime: config.LoopTime,
		},
	}
}

func (app *App) ResetAppState() {
	app.AppState.Character = utils.Data{}
	app.AppState.Location = utils.Data{
		AssetID: "menu_lost",
		Value: "Lost in the space-time continuum",
	}
	app.AppState.Menu = utils.Data{}
	app.AppState.IsInMenus = false
	app.AppState.IsOCRInitialized = false
}

// Connects to the Discord gateway
func (app *App) ConnectToDiscordGateway() error {
	logger.Logger.Info("Connecting to Discord Gateway")
	err := client.Login(app.Config.DiscordAppId)
	if err != nil {
		app.AppState.IsGatewayConnected = false
		return err
	}

	app.AppState.IsGatewayConnected = true
	return nil
}

// Disconnects from the Discord gateway
func (app *App) DisconnectFromDiscordGateway() {
	client.Logout()
	app.AppState.IsGatewayConnected = false
}

// Get the Honkai window handle
func (app *App) GetWindow() {
	winClassPtr, err := syscall.UTF16PtrFromString(app.Config.WindowClass)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	winTitlePtr, err := syscall.UTF16PtrFromString(app.Config.WindowName)
	if err != nil {
		logger.Logger.Error(err)
		return
	}
	app.HWND = win.FindWindow(winClassPtr, winTitlePtr)
	if app.HWND == 0 {
		app.AppState.LoopTime = 20 * time.Second
		logger.Logger.Info("No Honkai window found")
		return
	}

	app.AppState.AppStarted =	time.Now()
}

// Detect the current character currently selected
func (app *App) CaptureCharacter() {
	pos := utils.FindCurrentCharacter(app.Config.CharacterBoxCoords)
	if pos == -1 {
		return
	}
	app.AppState.CharacterPos = pos

	characterText, _ := utils.OcrManager.WindowOcr(app.Config.CharactersCoords[pos], "char")
	if characterText == "" {
		return
	}

	characterPred := utils.FindClosestCorrespondence(characterText, utils.GameData.Characters)

	if characterPred.Value != "" {
		app.AppState.Character = characterPred
	}
}

// Detect the current location of the player
func (app *App) CaptureLocation() {
	locationText, _ := utils.OcrManager.WindowOcr(app.Config.LocationCoord, "location")
	if locationText == "" {
		app.AppState.IsInMenus = true
		return
	}

	locationPred := utils.FindClosestCorrespondence(locationText, utils.GameData.Locations)

	if locationPred.Value != "" {
		app.AppState.IsInMenus = false
		app.AppState.Menu = utils.Data{}
		app.AppState.CombatStarted = time.Time{}
		app.AppState.Location = locationPred
		return
	}

	app.setMenu("menu_lost", "Lost in the space-time continuum", true)
}

// Set the current menu and if the player is in menus
func (app *App) setMenu(assetID, value string, isInMenus bool) {
	app.AppState.Menu = utils.Data{
		AssetID: assetID,
		Value:   value,
	}
	app.AppState.IsInMenus = isInMenus
}

// Get the menu from the list of menus
func getMenu(menus []utils.Data, value string) utils.Data {
	for _, menu := range menus {
		if menu.Value == value {
			return menu
		}
	}
	return utils.Data{}
}

// Capture the current game menu
func (app *App) CaptureGameMenu() {
	escText, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.EscCoord, "esc_menu")
	menuText, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.MenusCoord, "menus")
	combatText, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.CombatCoord, "combat")

	escTextPrediction := utils.FindClosestCorrespondence(escText, utils.GameData.Menus)
	menuTextPrediction := utils.FindClosestCorrespondence(menuText, utils.GameData.Menus)

	if escTextPrediction.Value == "Trailblaze Level" {
		menu := getMenu(utils.GameData.Menus, "Trailblaze Level")
		app.setMenu(menu.AssetID, menu.Message, true)
		return
	}

	if menuTextPrediction.Value != "" {
		menu := getMenu(utils.GameData.Menus, menuTextPrediction.Value)
		app.setMenu(menu.AssetID, menu.Message, true)
		return
	}

	if !app.AppState.CombatStarted.IsZero() {
		app.setMenu("menu_combat", "In combat", true)
		return
	} else if combatText != "" && app.AppState.CombatStarted.IsZero() {
		app.setMenu("menu_combat", "In combat", true)
		app.AppState.CombatStarted = time.Now()
		return
	}

	app.setMenu("menu_lost", "Lost in the space-time continuum", true)
}

// Main loop of the app
func (app *App) StartLoop() {
	for {
		<-time.After(app.AppState.LoopTime)
		// Reconnect to Discord gateway if not connected
		if !app.AppState.IsGatewayConnected {
			app.ConnectToDiscordGateway()
			continue
		}

		// Get the window handle if not obtained
		if app.HWND == 0 {
			app.GetWindow()
			continue
		}

		// Disconnect if window is closed
		if !win.IsWindowEnabled(app.HWND) {
			app.HandleWindowClosed()
			continue
		}

		// Check if Honkai window is focused
		if !app.IsWindowFocused() {
			continue
		}

		// Initialize OCR if not done already
		if !app.AppState.IsOCRInitialized {
			app.InitializeOCR()
		}

		if app.AppState.LoopTime != app.Config.LoopTime {
			app.AppState.LoopTime = app.Config.LoopTime
		}

		app.CaptureGameData()
		app.UpdateDiscordPresence()
	}
}

// Handles the scenario when the Honkai window is closed
func (app *App) HandleWindowClosed() {
	if app.HWND != 0 {
		app.HWND = 0
		app.AppState.LoopTime = 20 * time.Second
		app.ResetAppState()
		app.DisconnectFromDiscordGateway()
	}
}

// Checks if the Honkai window is currently focused
func (app *App) IsWindowFocused() bool {
	foregroundWindow := win.GetForegroundWindow()
	if foregroundWindow != app.HWND {
		app.AppState.LoopTime = 5 * time.Second
		logger.Logger.Info("Honkai window not focused")
		return false
	}
	logger.Logger.Info("Honkai window focused")
	return true
}

// Initializes OCR for capturing game data
func (app *App) InitializeOCR() {
	utils.InitOcr(utils.OCRConfig{
		ExecutablePath: app.Config.TesseractPath,
	}, app.HWND)

	app.AppState.IsOCRInitialized = true
}

// Captures game data such as character and location
func (app *App) CaptureGameData() {
	// TODO: view the possibility of running these in parallel
	app.CaptureLocation()
	if app.AppState.IsInMenus {
		app.CaptureGameMenu()
		return
	}
	app.CaptureCharacter()
}

// Updates the Discord presence based on game data
func (app *App) UpdateDiscordPresence() {
	var character utils.Data
	var location utils.Data

	if app.AppState.IsInMenus {
		location = app.AppState.Menu
		character = utils.Data{}
	} else {
		location = app.AppState.Location
		if location.Value == "" {
			location = utils.Data{
				AssetID: "menu_lost",
				Value:   "Lost in the space-time continuum",
			}
		}
		character = app.AppState.Character
	}

	// TODO: prevent update the presence if the data is the same

	if character.Value != "" || app.AppState.Menu.Value != "" {
		app.SetPresence(client.Activity{
			State:      location.Value,
			LargeImage: location.AssetID,
			LargeText:  location.Region,
			Details:    location.Region,
			SmallImage: character.AssetID,
			SmallText:  character.Value,
		})
	}
}

// Util to sets the Discord presence
func (app *App) SetPresence(rich client.Activity) {
	var time time.Time
	if app.AppState.CombatStarted.IsZero() {
		time = app.AppState.AppStarted
	} else {
		time = app.AppState.CombatStarted
	}
	err := client.SetActivity(client.Activity{
		State:			rich.State,
		Details:		rich.Details,
		LargeImage:	rich.LargeImage,
		LargeText:	rich.LargeText,
		SmallImage:	rich.SmallImage,
		SmallText:	rich.SmallText,
		Timestamps:	&client.Timestamps{
			Start:	&time,
		},
	})

	if err != nil {
		logger.Logger.Error(err)
	}
}
