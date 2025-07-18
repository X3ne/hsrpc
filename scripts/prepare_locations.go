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

func formatString(str string) string {
	str = strings.ReplaceAll(str, "\"", "")
	str = strings.TrimSpace(str)
	return str
}

func formatAssetId(str string) string {
	regex := regexp.MustCompile("[^a-zA-Z0-9_]+")
	str = strings.ReplaceAll(strings.ToLower(str), " ", "_")
	str = regex.ReplaceAllString(str, "")
	str = "loc_" + str
	return str
}

func PrepareLocations(config *ScriptConfig) {
	fmt.Print("Output path: ", config.OutputPath)

	file, err := os.Create(config.OutputPath + "locations.csv")
	if err != nil {
		log.Fatal(err)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	res, err := http.Get("https://honkai-star-rail.fandom.com/wiki/Astral_Express")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	header := []string{"value", "region", "sub_region", "asset_id"}
	if err := writer.Write(header); err != nil {
		log.Fatal(err)
	}

	regions := []string{"Herta Space Station", "Jarilo-VI", "The Xianzhou Luofu", "Penacony", "Amphoreus", "Other Locations"}

	root := doc.Find(".navbox-border").Eq(0)
	for _, region := range regions {
		root.Find(fmt.Sprintf("a[title*=\"%s\"]", region)).Eq(0).ParentsUntil("table").Each(func(i int, table *goquery.Selection) {
			table.Find(".navbox-group").Each(func(j int, row *goquery.Selection) {
				subRegion := row.Text()
				var subData []string
				subData = append(subData, formatString(subRegion), formatString(subRegion), region, formatAssetId(region))
				if err := writer.Write(subData); err != nil {
					log.Fatal(err)
				}
				row.ParentsUntil("tbody").Find("a").Each(func(k int, location *goquery.Selection) {
					var rowData []string

					assetId := formatAssetId(region)

					if location.Text() != "" {
						rowData = append(rowData, formatString(location.Text()), formatString(subRegion), region, assetId)

						if err := writer.Write(rowData); err != nil {
							log.Fatal(err)
						}
					}
				})
			})
		})
	}

	// Additional locations from ./static_locations.csv
	staticFile, err := os.Open("static_locations.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer staticFile.Close()

	staticReader := csv.NewReader(staticFile)
	staticReader.FieldsPerRecord = -1
	staticRecords, err := staticReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, record := range staticRecords {
		if len(record) < 3 {
			continue
		}
		rowData := []string{
			formatString(record[0]),
			formatString(record[1]),
			formatString(record[2]),
			formatAssetId(record[0]),
		}
		if err := writer.Write(rowData); err != nil {
			log.Fatal(err)
		}
	}
}
