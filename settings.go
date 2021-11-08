package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type Settings struct {
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	DbConn            string `yaml:"db_conn"`
	BuildDatabase     bool   `yaml:"build_database"`
	ParallelDownloads int    `yaml:"parallel_downloads"`
	QuerySettings     struct {
		PointsFile string    `yaml:"points_file"`
		StartDate  time.Time `yaml:"start_date"`
		EndDate    time.Time `yaml:"end_date"`
	} `yaml:"query_settings"`
	Download    bool   `yaml:"download"`
	Request     bool   `yaml:"request"`
	DownloadDir string `yaml:"download_dir"`
}

func (s Settings) String() string {
	return fmt.Sprintf("Username: %s\nPassword: %s\nDbconn: %s\nBuildDatabase: %t\nPointsFile: %s\nStartDate: %s\nEndDate: %s\nDownload: %t\n", s.Username, s.Password, s.DbConn, s.BuildDatabase, s.QuerySettings.PointsFile, s.QuerySettings.StartDate, s.QuerySettings.EndDate, s.Download)
}

func GetSettings(file string) *Settings {
	var settings Settings
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	err = yaml.Unmarshal(data, &settings)
	if err != nil {
		log.Fatalf("Error unmarshaling to Setting: %v", err)
	}

	return &settings
}
