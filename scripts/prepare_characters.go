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

func PrepareCharacters(config *ScriptConfig) {
	fmt.Print("Output path: ", config.OutputPath)

	res, err := http.Get("https://honkai-star-rail.fandom.com/wiki/Character")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(config.OutputPath + "characters.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	doc.Find(".article-table").First().Find("tr").Each(func(i int, s *goquery.Selection) {
		var rowData []string

		s.Find("td:nth-child(2)").Each(func(cellIndex int, cellHtml *goquery.Selection) {
			charName := strings.TrimSpace(cellHtml.Text())
			regex := regexp.MustCompile("[^a-zA-Z0-9_ ]+")
			charName = regex.ReplaceAllString(charName, "")
			charName = strings.Join(strings.Fields(charName), " ")

			charNameFormatted := strings.ReplaceAll(charName, " ", "_")
			assetName := "char_" + strings.ToLower(charNameFormatted)

			rowData = append(rowData, charName, assetName)
		})

		if err := writer.Write(rowData); err != nil {
			log.Fatal(err)
		}
	})
}
