package config

import (
	"time"

	"github.com/X3ne/hsrpc/src/utils"
)

type GUICoordsConfig struct {
	EscCoord						utils.Rect
	MenusCoord					utils.Rect
	SubMenuCoord				utils.Rect
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
	PreprocessThreshold int
	GUICoordsConfig			*GUICoordsConfig
}

func GetGUICoords(gameResolution Resolution, xAdjustment, yAdjustment int) *GUICoordsConfig {
	referenceResolution := Resolution{Width: 2560, Height: 1080}
	scaleX := float64(gameResolution.Width) / float64(referenceResolution.Width)
	scaleY := float64(gameResolution.Height) / float64(referenceResolution.Height)

	adjustSize := func(originalSize int, scale float64) int {
		return int(float64(originalSize) * scale)
	}

	wAdjust := 0
	if gameResolution.Width <= 1920 {
		wAdjust = 50
	}

	config := &GUICoordsConfig{
		EscCoord: utils.Rect{
			X:			adjustSize(1925+(xAdjustment*2), scaleX),
			Y:			adjustSize(250+yAdjustment, scaleY),
			Width:	adjustSize(180, scaleX) + 50,
			Height:	adjustSize(30, scaleY),
		},
		MenusCoord: utils.Rect{
			X:			adjustSize(100, scaleX),
			Y:			adjustSize(35+yAdjustment, scaleY),
			Width:	adjustSize(300, scaleX) + 50,
			Height:	adjustSize(40, scaleY),
		},
		SubMenuCoord: utils.Rect{
			X:			adjustSize(100, scaleX),
			Y:			adjustSize(65+yAdjustment, scaleY),
			Width:	adjustSize(300, scaleX) + 50,
			Height:	adjustSize(25, scaleY),
		},
		CombatCoord: utils.Rect{
			X:			adjustSize(2100+xAdjustment, scaleX),
			Y:			adjustSize(25+yAdjustment, scaleY),
			Width:	adjustSize(85, scaleX) + 50,
			Height:	adjustSize(40, scaleY),
		},
		LocationCoord: utils.Rect{
			X:			adjustSize(55, scaleX),
			Y:			adjustSize(15+yAdjustment, scaleY),
			Width:	adjustSize(320, scaleX) + 50,
			Height:	adjustSize(25, scaleY),
		},
		CharactersCoords: []utils.Rect{
			{X: adjustSize(2250+xAdjustment, scaleX), Y: adjustSize(305+yAdjustment, scaleY), Width: adjustSize(170, scaleX) + wAdjust, Height: adjustSize(30, scaleY)},
			{X: adjustSize(2250+xAdjustment, scaleX), Y: adjustSize(400+yAdjustment, scaleY), Width: adjustSize(170, scaleX) + wAdjust, Height: adjustSize(30, scaleY)},
			{X: adjustSize(2250+xAdjustment, scaleX), Y: adjustSize(495+yAdjustment, scaleY), Width: adjustSize(170, scaleX) + wAdjust, Height: adjustSize(30, scaleY)},
			{X: adjustSize(2250+xAdjustment, scaleX), Y: adjustSize(585+yAdjustment, scaleY), Width: adjustSize(170, scaleX) + wAdjust, Height: adjustSize(30, scaleY)},
		},
		CharactersBoxCoords: []utils.Rect{
			{X: adjustSize(2400+xAdjustment, scaleX), Y: adjustSize(351+yAdjustment, scaleY)},
			{X: adjustSize(2400+xAdjustment, scaleX), Y: adjustSize(445+yAdjustment, scaleY)},
			{X: adjustSize(2400+xAdjustment, scaleX), Y: adjustSize(538+yAdjustment, scaleY)},
			{X: adjustSize(2400+xAdjustment, scaleX), Y: adjustSize(632+yAdjustment, scaleY)},
		},
	}

	return config
}

func NewConfig() AppConfig {
	config :=  AppConfig{
		WindowName:						"Honkai: Star Rail",
		WindowClass:					"UnityWndClass",
		LoopTime:							2000,
		TesseractPath:				`C:\Program Files\Tesseract-OCR\tesseract.exe`,
		Resolution:						Resolution{Width: 1920, Height: 1080},
		DiscordAppId:					"1208212792574869544",
		Debug:								false,
		StartWithWindows:			false,
		PreprocessThreshold:	150,
	}

	return config
}
