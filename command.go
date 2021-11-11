package main

import (
	"FilesQuerier/api"
	"FilesQuerier/database"
	"log"
	"os"
	"os/signal"
)

func cleanup(con *database.Connection, errorCode int) {
	log.Println("Terminating program")
	con.UnlockNotDownloaded()
	os.Exit(errorCode)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("An argument is required for the program to run. Input the settings yaml file")
	}

	file := os.Args[1]
	s := GetSettings(file)

	con := database.NewConnection(s.DbConn)
	defer con.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cleanup(con, 0)
	}()
	queryClient := api.NewCopernicusClient(s.Username, s.Password)

	if s.BuildDatabase {
		con.BuildDatabase(s.QuerySettings.PointsFile, s.QuerySettings.StartDate, s.QuerySettings.EndDate, queryClient)
	}

	defer cleanup(con, 1)

	if s.Download {
		downloadClient := api.NewCopernicusClient(s.Username, s.Password)

		ch := downloadClient.OpenDownloadChannel()

		for i := 0; i < s.ParallelDownloads; i++ {
			go database.StartDownloadChannel(ch, s.DbConn, downloadClient)
		}

		for {
			con.RequestDownloadFiles(s.DownloadDir, queryClient, ch, s.Request)
		}
	}
}
