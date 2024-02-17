package config

import (
	"time"

	"github.com/X3ne/hsrpc/src/utils"
)

type GUICoordsConfig struct {
	EscCoord						utils.Rect
	MenusCoord					utils.Rect
	CombatCoord					utils.Rect
}

type AppConfig struct {
	WindowName					string
	WindowClass					string
	LoopTime						time.Duration
	TesseractPath				string
	DiscordAppId				string
	Resolution					uint32
	Ultrawide						bool
	CharactersCoords		[]utils.Rect
	CharacterBoxCoords	[]utils.Rect
	LocationCoord				utils.Rect
	GUICoordsConfig			GUICoordsConfig
}

func (c *AppConfig) setCharactersCoords() {
	if c.Ultrawide {
		c.CharactersCoords = []utils.Rect{
			{X: 2250, Y: 305, Width: 170, Height: 30},
			{X: 2250, Y: 400, Width: 170, Height: 30},
			{X: 2250, Y: 490, Width: 170, Height: 30},
			{X: 2250, Y: 590, Width: 170, Height: 30},
		}
		return
	}
	c.CharactersCoords = []utils.Rect{
		{X: 1620, Y: 240, Width: 150, Height: 60},
		{X: 1620, Y: 330, Width: 150, Height: 60},
		{X: 1620, Y: 430, Width: 150, Height: 60},
		{X: 1620, Y: 530, Width: 150, Height: 60},
	}
}

func (c *AppConfig) setCharacterBoxCoords() {
	if c.Ultrawide {
		c.CharacterBoxCoords = []utils.Rect{
			{X: 2400, Y: 351, Width: 200, Height: 1},
			{X: 2400, Y: 445, Width: 200, Height: 1},
			{X: 2400, Y: 538, Width: 200, Height: 1},
			{X: 2400, Y: 632, Width: 200, Height: 1},
		}
		return
	}
	c.CharacterBoxCoords = []utils.Rect{
		{X: 1860, Y: 260},
		{X: 1860, Y: 360},
		{X: 1860, Y: 460},
		{X: 1860, Y: 560},
	}
}

func (c *AppConfig) setLocationCoord() {
	c.LocationCoord = utils.Rect{
		X: 55,
		Y: 15,
		Width: 320,
		Height: 25,
	}
}

func (c *AppConfig) setGUICoords() {
	if c.Ultrawide {
		c.GUICoordsConfig = GUICoordsConfig{
			EscCoord: utils.Rect{
				X: 1925,
				Y: 250,
				Width: 180,
				Height: 30,
			},
			MenusCoord: utils.Rect{
				X: 100,
				Y: 35,
				Width: 300,
				Height: 25,
			},
			CombatCoord: utils.Rect{
				X: 2100,
				Y: 25,
				Width: 85,
				Height: 40,
			},
		}
		return
	}
	c.GUICoordsConfig = GUICoordsConfig{
		EscCoord: utils.Rect{
			X: 0,
			Y: 0,
			Width: 0,
			Height: 0,
		},
		MenusCoord: utils.Rect{
			X: 0,
			Y: 0,
			Width: 0,
			Height: 0,
		},
	}
}

func NewConfig() AppConfig {
	config :=  AppConfig{
		WindowName:					"Honkai: Star Rail",
		WindowClass:				"UnityWndClass",
		LoopTime:						2 * time.Second,
		TesseractPath:			`C:\Program Files\Tesseract-OCR\tesseract.exe`,
		Resolution:					1080,
		Ultrawide:					true,
		DiscordAppId:				"1208212792574869544",
	}

	config.setCharacterBoxCoords()
	config.setCharactersCoords()
	config.setLocationCoord()
	config.setGUICoords()

	return config
}
