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

type AppConfig struct {
	WindowName					string
	WindowClass					string
	LoopTime						time.Duration
	TesseractPath				string
	DiscordAppId				string
	Resolution					uint32
	Ultrawide						bool
	Debug								bool
	StartWithWindows		bool
	GUICoordsConfig			*GUICoordsConfig
}

func (c *AppConfig) setGUICoords() {
	if c.Ultrawide {
		c.GUICoordsConfig = &GUICoordsConfig{
			EscCoord: utils.Rect{X: 1925, Y: 250, Width: 180, Height: 30},
			MenusCoord: utils.Rect{X: 100, Y: 35, Width: 300, Height: 25},
			CombatCoord: utils.Rect{X: 2100, Y: 25, Width: 85, Height: 40},
			LocationCoord: utils.Rect{X: 55, Y: 15, Width: 320, Height: 25},
			CharactersCoords: []utils.Rect{
				{X: 2250, Y: 305, Width: 170, Height: 30},
				{X: 2250, Y: 400, Width: 170, Height: 30},
				{X: 2250, Y: 495, Width: 170, Height: 30},
				{X: 2250, Y: 585, Width: 170, Height: 30},
			},
			CharactersBoxCoords: []utils.Rect{
				{X: 2400, Y: 351},
				{X: 2400, Y: 445},
				{X: 2400, Y: 538},
				{X: 2400, Y: 632},
			},
		}
		return
	}
	c.GUICoordsConfig = &GUICoordsConfig{
		EscCoord: utils.Rect{X: 0, Y: 0, Width: 0, Height: 0},
		MenusCoord: utils.Rect{X: 0, Y: 0, Width: 0, Height: 0},
		CombatCoord: utils.Rect{X: 0, Y: 0, Width: 0, Height: 0},
		LocationCoord: utils.Rect{X: 55, Y: 15, Width: 320, Height: 25},
		CharactersCoords: []utils.Rect{
			{X: 1620, Y: 240, Width: 140, Height: 60},
			{X: 1620, Y: 330, Width: 140, Height: 60},
			{X: 1620, Y: 430, Width: 140, Height: 60},
			{X: 1620, Y: 530, Width: 140, Height: 60},
		},
		CharactersBoxCoords: []utils.Rect{
			{X: 1860, Y: 260},
			{X: 1860, Y: 360},
			{X: 1860, Y: 460},
			{X: 1860, Y: 560},
		},
	}
}

func NewConfig() AppConfig {
	config :=  AppConfig{
		WindowName:					"Honkai: Star Rail",
		WindowClass:				"UnityWndClass",
		LoopTime:						2000,
		TesseractPath:			`C:\Program Files\Tesseract-OCR\tesseract.exe`,
		Resolution:					1080,
		Ultrawide:					true,
		DiscordAppId:				"1208212792574869544",
		Debug:							false,
		StartWithWindows:		false,
	}

	config.setGUICoords()

	return config
}
