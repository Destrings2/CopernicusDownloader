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

func (c *Connection) DownloadFiles(path string, client *api.CopernicusClient, downloadChan chan<- api.DownloadRequest, allowRequest bool) {
	query := "SELECT uuid, title, online, requestedDate, priority, downloaded FROM file WHERE lockedBy IS NULL AND downloaded = FALSE AND (requestedDate IS NULL OR checkedDate < DATE_SUB(NOW(), INTERVAL 10 MINUTE) OR online = TRUE) ORDER BY priority DESC, online DESC, checkedDate ASC, RAND() DESC LIMIT 10;"

	results, err, _ := c.Query(query)
	if err != nil {
		log.Fatalf("Error while querying database: %s", err.Error())
	}

	if !allowRequest {
		client.RequestQuotaExceeded = true
	}

	if client.RequestQuotaExceeded {
		onlineFilesQuery, _, _ := c.Query("SELECT COUNT(uuid) count FROM file WHERE online = TRUE AND downloaded = FALSE AND lockedBy IS NULL")
		noOnlineFiles := string(onlineFilesQuery[0]["count"].([]byte)) == "0"

		requestedFilesQuery, _, _ := c.Query("SELECT COUNT(uuid) count FROM file WHERE online = FALSE AND lockedBy IS NULL AND requestedDate IS NOT NULL AND checkedDate < DATE_SUB(NOW(), INTERVAL 10 MINUTE)")
		noRequestedFiles := string(requestedFilesQuery[0]["count"].([]byte)) == "0"

		if noOnlineFiles && noRequestedFiles {
			log.Println("There are no online files to request download and the quota is exceeded, sleeping for 10 minutes.")
			time.Sleep(10 * time.Minute)
		}
	}

	for _, row := range results {
		c.ProcessDownload(row, path, client, downloadChan, allowRequest)
	}
}

func (c *Connection) ProcessDownload(row map[string]interface{}, path string, client *api.CopernicusClient, downloadChan chan<- api.DownloadRequest, allowRequest bool) {

	uuid := string(row["uuid"].([]byte))
	title := string(row["title"].([]byte))
	online := string(row["online"].([]byte)) == "1"
	requestedDateStr := row["requestedDate"]
	checkedDateStr := row["checkedDate"]

	if online {
		log.Println(fmt.Sprintf("File %s is online. Requesting its download", title))
		c.LockByUUID(uuid, client)
		downloadChan <- api.DownloadRequest{Title: title, Path: path, Uuid: uuid}
	} else if allowRequest && requestedDateStr == nil {
		// Do not request new files if we know the quota is exceeded
		if client.RequestQuotaExceeded && !client.RequestQuotaExceededAt.Before(time.Now().Add(-time.Minute*10)) {
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
			case api.OtherError:
				log.Println(err.Error())
			}
			c.SetRequested(uuid)
		}
	} else {
		log.Println(fmt.Sprintf("Checking if requested file %s is online", title))
		entry := c.UpdateByUUID(uuid, client)

		if entry.Properties.Online {
			c.Execute("UPDATE file SET priority = ? WHERE uuid = ?", 3, uuid)
		} else {
			if checkedDateStr != nil {
				checkedDate, _ := time.Parse("2006-01-02 15:04:05", string(checkedDateStr.([]byte)))
				// If the file has not been uploaded in 3 hours, delete the request
				if checkedDate.Before(time.Now().Add(-time.Hour * 3)) {
					c.Execute("UPDATE file SET requestedDate = ?, checkedDate = ? WHERE uuid = ?", nil, nil, uuid)
				}
			}
		}

		c.SetChecked(uuid)
	}
}
