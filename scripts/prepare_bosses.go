package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func PrepareBosses(config *ScriptConfig) {
	fmt.Print("Output path: ", config.OutputPath)

	res, err := http.Get("https://honkai-star-rail.fandom.com/wiki/Enemy/Boss")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(config.OutputPath + "bosses.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"value", "asset_id"}
	if err := writer.Write(header); err != nil {
		log.Fatal(err)
	}

	doc.Find(".wikitable").First().Find("tr").Each(func(i int, s *goquery.Selection) {
		var rowData []string

		s.Find(".item-text").Each(func(cellIndex int, cellHtml *goquery.Selection) {
			charName := strings.TrimSpace(cellHtml.Text())

			charNameRegex := regexp.MustCompile("[^a-zA-Z0-9_ ():]+")
			charName = charNameRegex.ReplaceAllString(charName, "")
			charName = strings.Join(strings.Fields(charName), " ")

			assetRegex := regexp.MustCompile(`\([^)]*\)|[^a-zA-Z0-9 ]+`)
			charNameFormatted := assetRegex.ReplaceAllString(charName, "")
			charNameFormatted = strings.TrimSpace(charNameFormatted)
			charNameFormatted = strings.ReplaceAll(charNameFormatted, " ", "_")

			assetName := "boss_" + strings.ToLower(charNameFormatted)

			rowData = append(rowData, charName, assetName)
		})

		if len(rowData) == 0 {
			return
		}

		if err := writer.Write(rowData); err != nil {
			log.Fatal(err)
		}
	})
}
