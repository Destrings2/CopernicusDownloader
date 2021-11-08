package utils

import (
	"encoding/csv"
	"log"
	"os"
)

func ReadCsvFile(fileName string) [][]string {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Unable to read file", err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse CSV file", err)
	}

	return records
}
