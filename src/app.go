package app

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"github.com/lxn/win"

	"github.com/X3ne/hsrpc/src/config"
	"github.com/X3ne/hsrpc/src/logger"
	"github.com/X3ne/hsrpc/src/utils"
	"github.com/hugolgst/rich-go/client"
)

type App struct {
	Config   config.AppConfig
	GUIApp   fyne.App
	HWND     win.HWND
	AppState AppState
	AppSize  config.Resolution
}

type Combat struct {
	Started time.Time
	IsBoss  bool
	Boss    utils.Data
}

type AppState struct {
	PlayerInfo         *utils.PlayerInfo
	CharacterPos       int32
	Character          utils.Data
	Location           utils.Data
	Menu               utils.Data
	LoopTime           time.Duration
	AppStarted         time.Time
	Combat             Combat
	IsInMenus          bool
	IsGatewayConnected bool
	IsOCRInitialized   bool
}

func CreateApp(config config.AppConfig, guiApp fyne.App) *App {
	utils.LoadGameData()

	return &App{
		Config: config,
		GUIApp: guiApp,
		HWND:   0,
		AppState: AppState{
			Location: utils.Data{
				AssetID: "menu_lost",
				Value:   "Lost in the space-time continuum",
			},
			LoopTime: config.LoopTime,
		},
	}
}

func (app *App) ResetAppState() {
	app.AppState.Character = utils.Data{}
	app.AppState.Location = utils.Data{
		AssetID: "menu_lost",
		Value:   "Lost in the space-time continuum",
	}
	app.AppState.Menu = utils.Data{}
	app.AppState.IsInMenus = false
	app.AppState.IsOCRInitialized = false
}

// ConnectToDiscordGateway
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

// DisconnectFromDiscordGateway
// Disconnects from the Discord gateway
func (app *App) DisconnectFromDiscordGateway() {
	client.Logout()
	app.AppState.IsGatewayConnected = false
}

// GetWindow
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
		app.AppState.LoopTime = 20000
		logger.Logger.Info("No Honkai window found")
		return
	}

	app.AppState.AppStarted = time.Now()
}

// Get Honkai Star Rail window size
func getWindowSize(hwnd win.HWND) (int32, int32) {
	var rect win.RECT
	win.GetClientRect(hwnd, &rect)
	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top
	return width, height
}

// SetSize
// Set App size
func (app *App) SetSize() {
	w, h := getWindowSize(app.HWND)

	if w == int32(app.AppSize.Width) && h == int32(app.AppSize.Height) {
		return
	}

	app.AppSize = config.Resolution{Width: uint32(w), Height: uint32(h)}

	logger.Logger.Info("Window size:", w, h)

	monitor := win.MonitorFromWindow(app.HWND, win.MONITOR_DEFAULTTONEAREST)

	var monitorInfo win.MONITORINFO
	monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
	win.GetMonitorInfo(monitor, &monitorInfo)

	isFullscreen := monitorInfo.RcMonitor.Right-monitorInfo.RcMonitor.Left == w && monitorInfo.RcMonitor.Bottom-monitorInfo.RcMonitor.Top == h
	logger.Logger.Info("Fullscreen:", isFullscreen)

	xAdjustment := 0
	if w <= 1920 {
		xAdjustment = -100
	}

	yAdjustment := 0
	if !isFullscreen {
		yAdjustment = 33
	}

	app.Config.GUICoordsConfig = config.GetGUICoords(config.Resolution{Width: uint32(w), Height: uint32(h)}, xAdjustment, yAdjustment)

	config.SaveConfig(app.Config)
}

// CaptureCharacter
// Detect the current character currently selected
func (app *App) CaptureCharacter() {
	pos := utils.FindCurrentCharacter(app.HWND, app.Config.GUICoordsConfig.CharactersBoxCoords)
	if pos == -1 {
		return
	}
	app.AppState.CharacterPos = pos

	characterText, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.CharactersCoords[pos], "char", true)
	if characterText == "" {
		return
	}

	playerName := app.Config.PlayerName

	if app.AppState.PlayerInfo != nil {
		playerName = app.AppState.PlayerInfo.Player.Nickname
	}

	charactersData := append(utils.GameData.Characters, utils.Data{
		AssetID: "char_trailblazer",
		Value:   playerName,
	})

	characterPred := utils.FindClosestCorrespondence(characterText, charactersData)

	if characterPred.Value != "" {
		app.AppState.Character = characterPred
	}
}

// CaptureLocation
// Detect the current location of the player
func (app *App) CaptureLocation() {
	locationText, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.LocationCoord, "location", true)
	if locationText == "" {
		app.AppState.IsInMenus = true
		return
	}

	// TODO: clean this and move to the csv (edit the script to auto gen new data with this)
	locationData := append(
		utils.GameData.Locations,
		utils.Data{
			AssetID: "loc_parlor_car",
			Value:   "Parlor Car",
		},
		utils.Data{
			AssetID: "loc_parlor_car",
			Value:   "Passenger Cabin",
		},
	)

	locationPred := utils.FindClosestCorrespondence(locationText, locationData)

	// needed to trigger presence updating
	charName := app.Config.PlayerName
	if charName == "" {
		charName = "Trailblazer"
	}
	if locationPred.Value == "Parlor Car" || locationPred.Value == "Passenger Cabin" {
		app.AppState.Character = utils.Data{
			AssetID: "char_trailblazer",
			Value:   charName,
		}
	}

	if locationPred.Value != "" {
		app.AppState.IsInMenus = false
		app.AppState.Menu = utils.Data{}
		app.AppState.Combat = Combat{}
		app.AppState.Location = locationPred
		return
	}

	app.setMenu("menu_lost", "Lost in the space-time continuum", true)
}

// Set the current menu and if the player is in menus
func (app *App) setMenu(assetID, value string, isInMenus bool, subMenu ...string) {
	subMenu = append(subMenu, "")

	app.AppState.Menu = utils.Data{
		AssetID: assetID,
		Value:   value,
		Region:  subMenu[0],
	}

	if len(subMenu) > 1 {
		app.AppState.Character = utils.Data{
			AssetID: subMenu[1],
			Value:   subMenu[0],
		}
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

// CaptureGameMenu
// Capture the current game menu
func (app *App) CaptureGameMenu() {
	escText, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.EscCoord, "esc_menu", true)
	menuText, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.MenusCoord, "menus", true)
	combatText, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.CombatCoord, "combat", true)
	bossText, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.BossCoord, "boss", true)

	escTextPrediction := utils.FindClosestCorrespondence(escText, utils.GameData.Menus)
	menuTextPrediction := utils.FindClosestCorrespondence(menuText, utils.GameData.Menus)
	bossTextPrediction := utils.FindClosestCorrespondence(bossText, utils.GameData.Bosses)

	if escTextPrediction.Value == "Trailblaze Level" {
		menu := getMenu(utils.GameData.Menus, "Trailblaze Level")
		app.setMenu(menu.AssetID, menu.Message, true)
		return
	}

	if menuTextPrediction.Value != "" {
		menu := getMenu(utils.GameData.Menus, menuTextPrediction.Value)
		subMenu, _ := utils.OcrManager.WindowOcr(app.Config.GUICoordsConfig.SubMenuCoord, "sub_menu", true)
		subMenuPrediction := utils.FindClosestCorrespondence(subMenu, utils.GameData.SubMenus)
		app.setMenu(menu.AssetID, menu.Message, true, subMenuPrediction.Message)
		return
	}

	statusText := "In combat"

	if bossTextPrediction.Value != "" && !app.AppState.Combat.IsBoss {
		app.AppState.Combat = Combat{
			IsBoss: true,
			Boss:   bossTextPrediction,
		}
	}

	if app.AppState.Combat.IsBoss {
		statusText = "Defeating a boss"
	}

	if !app.AppState.Combat.Started.IsZero() {
		app.setMenu("menu_combat", statusText, true, app.AppState.Combat.Boss.Value, app.AppState.Combat.Boss.AssetID)
		return
	}

	if (combatText != "" && len(combatText) > 2) || bossTextPrediction.Value != "" {
		app.setMenu("menu_combat", statusText, true, bossTextPrediction.Value, bossTextPrediction.AssetID)
		app.AppState.Combat = Combat{
			IsBoss:  bossTextPrediction.Value != "",
			Boss:    bossTextPrediction,
			Started: time.Now(),
		}
	}
}

// StartLoop
// Main loop of the app
func (app *App) StartLoop() {
	defer func() {
		if r := recover(); r != nil {
			utils.PanicRecover(r, app.GUIApp)
		}
	}()

	// TODO: Refresh data when uid changes
	go func() {
		waitTime := 5 * time.Second
		for {
			<-time.After(waitTime)
			if app.Config.PlayerUID != "" {
				playerInfos, err := utils.GetPlayerInfos(app.Config.PlayerUID)
				if err != nil {
					logger.Logger.Error(err)
				} else {
					app.AppState.PlayerInfo = playerInfos
					waitTime = 5 * time.Minute
				}
			}
		}
	}()

	for {
		<-time.After(app.AppState.LoopTime * time.Millisecond)
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

		// Set the app size
		app.SetSize()

		if app.AppState.LoopTime != app.Config.LoopTime {
			app.AppState.LoopTime = app.Config.LoopTime
		}

		app.CaptureGameData()
		app.UpdateDiscordPresence()
	}
}

// HandleWindowClosed
// Handles the scenario when the Honkai window is closed
func (app *App) HandleWindowClosed() {
	if app.HWND != 0 {
		app.HWND = 0
		app.AppState.LoopTime = 20000
		app.ResetAppState()
		app.DisconnectFromDiscordGateway()
	}
}

// IsWindowFocused
// Checks if the Honkai window is currently focused
func (app *App) IsWindowFocused() bool {
	foregroundWindow := win.GetForegroundWindow()
	if foregroundWindow != app.HWND {
		app.AppState.LoopTime = 5000
		logger.Logger.Debug("Honkai window not focused")
		return false
	}
	logger.Logger.Debug("Honkai window focused")
	return true
}

// InitializeOCR
// Initializes OCR for capturing game data
func (app *App) InitializeOCR() {
	utils.InitOcr(utils.OCRConfig{
		ExecutablePath:      &app.Config.TesseractPath,
		PreprocessThreshold: &app.Config.PreprocessThreshold,
	}, app.HWND)

	app.AppState.IsOCRInitialized = true
}

// CaptureGameData
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

// UpdateDiscordPresence
// Updates the Discord presence based on game data
func (app *App) UpdateDiscordPresence() {
	var character utils.Data
	var location utils.Data

	if app.AppState.IsInMenus {
		location = app.AppState.Menu

		if app.AppState.Menu.Value == "Defeating a boss" {
			character = app.AppState.Character
		} else {
			character = utils.Data{}
		}
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
	largeText := ""
	if app.AppState.PlayerInfo != nil {
		if app.Config.DisplayLevel {
			largeText += fmt.Sprintf("Level %d ", app.AppState.PlayerInfo.Player.Level)
		}
		if app.Config.DisplayNickname {
			largeText += app.AppState.PlayerInfo.Player.Nickname
		}
	}

	// TODO: prevent update the presence if the data is the same

	if character.Value != "" || app.AppState.Menu.Value != "" {
		app.SetPresence(client.Activity{
			State:      location.Value,
			LargeImage: location.AssetID,
			LargeText:  largeText,
			Details:    location.Region,
			SmallImage: character.AssetID,
			SmallText:  character.Value,
		})
	}
}

// SetPresence
// Util to sets the Discord presence
func (app *App) SetPresence(rich client.Activity) {
	var time time.Time
	if app.AppState.Combat.Started.IsZero() {
		time = app.AppState.AppStarted
	} else {
		time = app.AppState.Combat.Started
	}
	err := client.SetActivity(client.Activity{
		State:      rich.State,
		Details:    rich.Details,
		LargeImage: rich.LargeImage,
		LargeText:  rich.LargeText,
		SmallImage: rich.SmallImage,
		SmallText:  rich.SmallText,
		Timestamps: &client.Timestamps{
			Start: &time,
		},
	})

	if err != nil {
		logger.Logger.Error(err)
	}
}
