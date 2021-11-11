package database

import (
	"FilesQuerier/api"
	"fmt"
	"log"
	"time"
)

func StartDownloadChannel(ch chan api.DownloadRequest, dbConn string, client *api.CopernicusClient) {
	con := NewConnection(dbConn)
	defer con.Close()
	for {
		select {
		case req := <-ch:
			log.Println(fmt.Sprintf("Starting download request for %s", req.Title))
			name, err := client.Download(req.Uuid, req.Path, req.Title)
			if err != nil {
				log.Println(err.Error())
				con.UnlockByUUID(req.Uuid)
			} else {
				con.MarkDownloadedByUUID(req.Uuid, name)
			}
		default:
			time.Sleep(10 * time.Second)
		}
	}
}

func (c *Connection) RequestDownloadFiles(path string, client *api.CopernicusClient, downloadChan chan<- api.DownloadRequest, allowRequest bool) {
	toDownloadQuery := "SELECT uuid, title, online, requestedDate, checkedDate, priority, downloaded FROM file WHERE lockedBy IS NULL AND downloaded = FALSE AND online = TRUE ORDER BY priority DESC, online DESC, checkedDate DESC, RAND() DESC LIMIT 1;"
	toCheckQuery := "SELECT uuid, title, online, requestedDate, checkedDate, priority, downloaded FROM file WHERE requestedDate IS NOT NULL AND downloaded = FALSE AND checkedDate < DATE_SUB(NOW(), INTERVAL 10 MINUTE) ORDER BY priority DESC, online DESC, checkedDate DESC, RAND() DESC LIMIT 10;"
	toRequestQuery := "SELECT uuid, title, online, requestedDate, checkedDate, priority, downloaded FROM file WHERE requestedDate IS NULL AND downloaded = FALSE ORDER BY priority DESC, online DESC, checkedDate DESC, RAND() DESC LIMIT 10;"

	toDownload, err, _ := c.Query(toDownloadQuery)
	if err != nil {
		log.Fatalf("Error while querying database: %s", err.Error())
	}

	toRequest, err, _ := c.Query(toRequestQuery)
	if err != nil {
		log.Fatalf("Error while querying database: %s", err.Error())
	}

	toCheck, err, _ := c.Query(toCheckQuery)
	if err != nil {
		log.Fatalf("Error while querying database: %s", err.Error())
	}

	noToCheck := len(toCheck) == 0
	noToDownload := len(toDownload) == 0

	if noToCheck && noToDownload && (client.RequestQuotaExceeded || !allowRequest) {
		log.Println("There are no online files to check or download, and the quota is exceeded. Sleeping for 10 minutes.")
		time.Sleep(10 * time.Minute)
		return
	}

	for _, row := range toDownload {
		c.ProcessDownload(row, path, client, downloadChan)
	}

	for _, row := range toCheck {
		c.ProcessCheck(row, client)
	}

	if allowRequest {
		for _, row := range toRequest {
			c.ProcessRequest(row, client)
		}
	}

}

func (c *Connection) ProcessDownload(row map[string]interface{}, path string, client *api.CopernicusClient, downloadChan chan<- api.DownloadRequest) {

	uuid := string(row["uuid"].([]byte))
	title := string(row["title"].([]byte))

	log.Println(fmt.Sprintf("File %s is online. Requesting its download", title))
	c.LockByUUID(uuid, client)
	downloadChan <- api.DownloadRequest{Title: title, Path: path, Uuid: uuid}
}

func (c *Connection) ProcessCheck(row map[string]interface{}, client *api.CopernicusClient) {

	uuid := string(row["uuid"].([]byte))
	title := string(row["title"].([]byte))
	checkedDateStr := row["checkedDate"]

	log.Println(fmt.Sprintf("Checking if requested file %s is online", title))
	entry := c.UpdateByUUID(uuid, client)

	if entry.Properties.Online {
		c.Execute("UPDATE file SET priority = ? WHERE uuid = ?", 3, uuid)
	} else {
		if checkedDateStr != nil {
			checkedDate, _ := time.Parse("2006-01-02 15:04:05", string(checkedDateStr.([]byte)))
			// If the file has not been uploaded in 3 hours, delete the request
			if checkedDate.After(time.Now().Add(-time.Hour * 3)) {
				log.Println(fmt.Sprintf("More than 3 hours since request off %s, setting it up for requesting again", title))
				c.Execute("UPDATE file SET requestedDate = ?, checkedDate = ? WHERE uuid = ?", nil, nil, uuid)
			}
		}
	}
	c.SetChecked(uuid)
}

func (c *Connection) ProcessRequest(row map[string]interface{}, client *api.CopernicusClient) {
	uuid := string(row["uuid"].([]byte))
	title := string(row["title"].([]byte))

	// Do not request new files if we know the quota is exceeded
	if client.RequestQuotaExceeded && !client.RequestQuotaExceededAt.After(time.Now().Add(-time.Minute*10)) {
		client.RequestQuotaExceeded = false
	} else if client.RequestQuotaExceeded {
		return
	}

	log.Println(fmt.Sprintf("File %s is not online. Requesting it", title))

	err := client.Request(uuid)
	if err != nil {
		switch err.(type) {
		case api.QuotaExceededError:
			log.Println(err.Error())
			client.RequestQuotaExceeded = true
			client.RequestQuotaExceededAt = time.Now()
		case api.CantRequestOnlineError:
			//File is actually online, update it
			c.UpdateByUUID(uuid, client)
		case api.TooManyRequestsError:
			log.Println(err.Error())
			time.Sleep(time.Second * 10)
		case api.OtherError:
			log.Println(err.Error())
		}
		return
	}
	c.SetRequested(uuid)
}
