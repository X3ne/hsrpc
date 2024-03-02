package utils

import (
	"bytes"
	"encoding/csv"

	"github.com/X3ne/hsrpc/src/internal/bundle"
	"github.com/X3ne/hsrpc/src/logger"
)

type GameDataStruct struct {
	Characters	[]Data
	Locations		[]Data
	Menus				[]Data
	SubMenus		[]Data
}

type Data struct {
	AssetID		string
	Value			string
	Message		string
	Region		string
	SubRegion	string
}

var GameData *GameDataStruct

func LoadGameData() {
	characters, err := loadCharacters()
	if err != nil {
		logger.Logger.Fatal("Error when loading characters", err)
	}
	locations, err := loadLocations()
	if err != nil {
		logger.Logger.Fatal("Error when loading locations", err)
	}
	menus, err := loadGameMenus()
	if err != nil {
		logger.Logger.Fatal("Error when loading menus", err)
	}
	subMenus, err := loadGameSubMenus()
	if err != nil {
		logger.Logger.Fatal("Error when loading sub menus", err)
	}
	GameData = &GameDataStruct{
		Characters:	characters,
		Locations:	locations,
		Menus:			menus,
		SubMenus:		subMenus,
	}
}

func loadCharacters() ([]Data, error) {
	var characters []Data

	logger.Logger.Info("Loading characters")

	data := bundle.Get("characters.csv")

	reader := csv.NewReader(bytes.NewReader(data))

	records, err := reader.ReadAll()
	if err != nil {
		return []Data{}, err
	}

	var loadedCharacters int
	for _, record := range records {
		characters = append(characters, Data{
			AssetID:	record[1],
			Value:		record[0],
		})
		loadedCharacters++
	}

	logger.Logger.Info("Loaded ", loadedCharacters, " characters")

	return characters, nil
}

func loadLocations() ([]Data, error) {
	var locations []Data

	logger.Logger.Info("Loading locations")

	data := bundle.Get("locations.csv")

	reader := csv.NewReader(bytes.NewReader(data))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var loadedLocations int
	for _, record := range records {
		location := Data{
			AssetID:	record[3],
			Value:		record[0],
			Region:		record[2],
			SubRegion:	record[1],
		}

		locations = append(locations, location)
		loadedLocations++
	}

	logger.Logger.Info("Loaded ", loadedLocations, " locations")

	return locations, nil
}

func loadGameMenus() ([]Data, error) {
	var menus []Data

	logger.Logger.Info("Loading game menus")

	data := bundle.Get("menus.csv")

	reader := csv.NewReader(bytes.NewReader(data))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var loadedMenus int
	for _, record := range records {
		menu := Data{
			AssetID:	record[2],
			Value:		record[0],
			Message:	record[1],
		}

		menus = append(menus, menu)
		loadedMenus++
	}

	logger.Logger.Info("Loaded ", loadedMenus, " game menus")

	return menus, nil
}

func loadGameSubMenus() ([]Data, error) {
	var menus []Data

	logger.Logger.Info("Loading sub game menus")

	data := bundle.Get("sub_menus.csv")

	reader := csv.NewReader(bytes.NewReader(data))

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var loadedMenus int
	for _, record := range records {
		menu := Data{
			AssetID:	record[2],
			Value:		record[0],
			Message:	record[1],
		}

		menus = append(menus, menu)
		loadedMenus++
	}

	logger.Logger.Info("Loaded ", loadedMenus, " game sub menus")

	return menus, nil
}
