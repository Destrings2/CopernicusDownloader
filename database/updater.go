package database

import (
	"FilesQuerier/api"
	"fmt"
	"log"
)

func (c *Connection) UpdateDatabase(client *api.CopernicusClient) {
	results, err, _ := c.Query("SELECT uuid FROM file")
	if err != nil {
		return
	}

	length := len(results)
	for i, row := range results {
		if i%100 == 0 {
			log.Printf(fmt.Sprintf("%d/%d\n", i, length))
		}

		uuid := string(row["uuid"].([]byte))
		c.UpdateByUUID(uuid, client)
	}
}

func (c *Connection) UpdateByUUID(uuid string, client *api.CopernicusClient) *api.ODataEntry {
	entry := client.GetEntryByUUID(uuid)
	online := entry.Properties.Online
	checksum := entry.Properties.Checksum.Value

	_, err := c.Execute("UPDATE file SET online = ?, checkmd5 = ? WHERE uuid = ?", online, checksum, uuid)
	if err != nil {
		return nil
	}

	return entry
}

func (c *Connection) LockByUUID(uuid string, client *api.CopernicusClient) {
	_, err := c.Execute("UPDATE file SET lockedBy = ? WHERE uuid = ?", client.Username, uuid)
	if err != nil {
		return
	}
}

func (c *Connection) UnlockByUUID(uuid string) {
	_, err := c.Execute("UPDATE file SET lockedBy = ? WHERE uuid = ?", nil, uuid)
	if err != nil {
		return
	}
}

func (c *Connection) UnlockNotDownloaded() {
	_, err := c.Execute("UPDATE file SET lockedBy = ? WHERE lockedBy IS NOT NULL AND downloaded = FALSE", nil)
	if err != nil {
		return
	}
}

func (c *Connection) MarkDownloadedByUUID(uuid string, name string) {
	_, err := c.Execute("UPDATE file SET downloaded = ?, localPath = ? WHERE uuid = ?", true, name, uuid)
	if err != nil {
		return
	}
}
