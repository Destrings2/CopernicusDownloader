package database

import (
	"FilesQuerier/api"
	"FilesQuerier/utils"
	"log"
	"strconv"
	"time"
)

func (c *Connection) BuildDatabase(pointsFile string, start time.Time, end time.Time, client *api.CopernicusClient) {
	records := utils.ReadCsvFile(pointsFile)
	geoPoints := api.GeoPoints{}
	for _, v := range records {

		lat, _ := strconv.ParseFloat(v[0], 64)
		lon, _ := strconv.ParseFloat(v[1], 64)

		geoPoints.Add(api.GeoPoint{
			Latitude:  lat,
			Longitude: lon,
		})
	}

	var pointQueries []string
	for _, v := range geoPoints {
		pointQueries = append(pointQueries, api.IntersectGeoPoint(v))
	}

	// If the parameter url is too big, it errors. Split it up
	for i := 0; i < len(pointQueries); i += 41 {
		var endIndex int
		if len(pointQueries)-1 <= i+40 {
			endIndex = len(pointQueries) - 1
		} else {
			endIndex = i + 40
		}

		pointSlice := pointQueries[i:endIndex]

		builder := api.NewQueryBuilder()
		builder.And(api.OrGroup(pointSlice...))
		builder.And(api.SensingRange(start, end))
		//Platform and product type
		builder.And(api.PlatformName("Sentinel-2"))
		builder.And(api.ProductType("S2MSI1C"))

		query := builder.GetQuery()

		index := 0
		for {
			feed := client.Search(query, index, 100)
			log.Println(feed.Subtitle)

			startIndex, _ := strconv.Atoi(feed.StartIndex)
			itemsPerPage, _ := strconv.Atoi(feed.ItemsPerPage)

			index = startIndex + itemsPerPage
			c.AddEntries(feed)

			if !feed.HasNextPage() {
				break
			}
		}
	}
}
