package utils

import (
	"encoding/csv"
	"os"

	"github.com/X3ne/hsrpc/src/logger"
)

type GameDataStruct struct {
	Characters	[]Data
	Locations		[]Data
	Menus				[]Data
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
	GameData = &GameDataStruct{
		Characters:	characters,
		Locations:	locations,
		Menus:			menus,
	}
}

func loadCharacters() ([]Data, error) {
	file, err := os.Open("data/characters.csv")
	if err != nil {
		return []Data{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return []Data{}, err
	}

	var characters []Data

	for _, record := range records {
		characters = append(characters, Data{
			AssetID:	record[1],
			Value:		record[0],
		})
	}

	return characters, nil
}

func loadLocations() ([]Data, error) {
	var locations []Data

	file, err := os.Open("data/locations.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		location := Data{
			AssetID:	record[3],
			Value:		record[0],
			Region:		record[2],
			SubRegion:	record[1],
		}

		locations = append(locations, location)
	}

	return locations, nil
}

func loadGameMenus() ([]Data, error) {
	var menus []Data

	file, err := os.Open("data/menus.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		menu := Data{
			AssetID:	record[2],
			Value:		record[0],
			Message:	record[1],
		}

		menus = append(menus, menu)
	}

	return menus, nil
}
