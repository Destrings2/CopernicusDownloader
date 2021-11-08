package database

import (
	"FilesQuerier/api"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
)

type Connection struct {
	Database *sql.DB
}

func NewConnection(connectionString string) *Connection {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}
	return &Connection{Database: db}
}

func (c *Connection) Close() {
	c.Database.Close()
}

func strToBytes(str string) int64 {
	var size float64
	var unit string
	var err error

	number := str[:len(str)-3]

	if size, err = strconv.ParseFloat(number, 64); err != nil {
		panic(err)
	}
	unit = str[len(str)-2:]

	switch unit {
	case "KB":
		return int64(size * 1000)
	case "MB":
		return int64(size * 1000 * 1000)
	case "GB":
		return int64(size * 1000 * 1000 * 1000)
	default:
		panic("Unknown unit")
	}
}

func toMysqlFormat(str string) string {
	return strings.Replace(strings.Replace(str, "T", " ", 1), "Z", "", 1)
}

func (c *Connection) AddEntries(feed *api.Feed) {
	for _, entry := range feed.Entry {
		var beginPosition, endPosition, ingestionDate, productLevel, footprint, platformName, productType, fileFormat, filename string
		var orbitNumber, relativeOrbitNumber int
		var size int64
		var online, downloaded bool

		online = false
		downloaded = false

		for _, date := range entry.Date {
			switch date.Name {
			case "beginposition":
				beginPosition = toMysqlFormat(date.Text)
			case "endposition":
				endPosition = toMysqlFormat(date.Text)
			case "ingestiondate":
				ingestionDate = toMysqlFormat(date.Text)
			}
		}

		for _, node := range entry.Int {
			switch node.Name {
			case "orbitnumber":
				orbitNumber, _ = strconv.Atoi(node.Text)
			case "relativeorbitnumber":
				relativeOrbitNumber, _ = strconv.Atoi(node.Text)
			}
		}

		for _, node := range entry.Str {
			switch node.Name {
			case "processinglevel":
				productLevel = node.Text
			case "gmlfootprint":
				footprint = node.Text[182 : len(node.Text)-83]
			case "platformname":
				platformName = node.Text
			case "producttype":
				productType = node.Text
			case "format":
				fileFormat = node.Text
			case "filename":
				filename = node.Text
			case "size":
				size = strToBytes(node.Text)
			}
		}

		query := "INSERT INTO file (uuid, title, beginPosition, endPosition, ingestionData, orbitNumber, relativeOrbitNumber, productLevel, footprint, platformName, productType, fileFormat, filename, size, online, downloaded) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
		_, err := c.Database.Exec(query, entry.ID, entry.Title, beginPosition, endPosition, ingestionDate, orbitNumber, relativeOrbitNumber, productLevel, footprint, platformName, productType, fileFormat, filename, size, online, downloaded)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
