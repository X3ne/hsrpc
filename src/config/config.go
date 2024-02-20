package config

import (
	"time"

	"github.com/X3ne/hsrpc/src/utils"
)

type GUICoordsConfig struct {
	EscCoord						utils.Rect
	MenusCoord					utils.Rect
	CombatCoord					utils.Rect
	LocationCoord				utils.Rect
	CharactersCoords		[]utils.Rect
	CharactersBoxCoords	[]utils.Rect
}

type Resolution struct {
	Width  uint32
	Height uint32
}

type AppConfig struct {
	WindowName					string
	WindowClass					string
	LoopTime						time.Duration
	TesseractPath				string
	DiscordAppId				string
	Resolution					Resolution
	Debug								bool
	StartWithWindows		bool
	GUICoordsConfig			*GUICoordsConfig
}

func GetGUICoords(gameResolution Resolution, xAdjustment, yAjustment int) *GUICoordsConfig {
	referenceResolution := Resolution{Width: 2560, Height: 1080}
	scaleX := float64(gameResolution.Width) / float64(referenceResolution.Width)
	scaleY := float64(gameResolution.Height) / float64(referenceResolution.Height)

	config := &GUICoordsConfig{
		EscCoord: utils.Rect{
			X:      int(float64(1925+(xAdjustment*2)) * scaleX),
			Y:      int(float64(250+yAjustment) * scaleY),
			Width:  180,
			Height: 30,
		},
		MenusCoord: utils.Rect{
			X:      int(float64(100) * scaleX),
			Y:      int(float64(35+yAjustment) * scaleY),
			Width:  300,
			Height: 25,
		},
		CombatCoord: utils.Rect{
			X:      int(float64(2100+xAdjustment) * scaleX),
			Y:      int(float64(25+yAjustment) * scaleY),
			Width:  85,
			Height: 40,
		},
		LocationCoord: utils.Rect{
			X:      int(float64(55) * scaleX),
			Y:      int(float64(15+yAjustment) * scaleY),
			Width:  320,
			Height: 25,
		},
		CharactersCoords: []utils.Rect{
			{X: int(float64(2250+xAdjustment) * scaleX), Y: int(float64(305+yAjustment) * scaleY), Width: 170, Height: 30},
			{X: int(float64(2250+xAdjustment) * scaleX), Y: int(float64(400+yAjustment) * scaleY), Width: 170, Height: 30},
			{X: int(float64(2250+xAdjustment) * scaleX), Y: int(float64(495+yAjustment) * scaleY), Width: 170, Height: 30},
			{X: int(float64(2250+xAdjustment) * scaleX), Y: int(float64(585+yAjustment) * scaleY), Width: 170, Height: 30},
		},
		CharactersBoxCoords: []utils.Rect{
			{X: int(float64(2400+xAdjustment) * scaleX), Y: int(float64(351+yAjustment) * scaleY)},
			{X: int(float64(2400+xAdjustment) * scaleX), Y: int(float64(445+yAjustment) * scaleY)},
			{X: int(float64(2400+xAdjustment) * scaleX), Y: int(float64(538+yAjustment) * scaleY)},
			{X: int(float64(2400+xAdjustment) * scaleX), Y: int(float64(632+yAjustment) * scaleY)},
		},
	}

	return config
}

func NewConfig() AppConfig {
	config :=  AppConfig{
		WindowName:					"Honkai: Star Rail",
		WindowClass:				"UnityWndClass",
		LoopTime:						2000,
		TesseractPath:			`C:\Program Files\Tesseract-OCR\tesseract.exe`,
		Resolution:					Resolution{Width: 1920, Height: 1080},
		DiscordAppId:				"1208212792574869544",
		Debug:							false,
		StartWithWindows:		false,
	}

	return config
}
