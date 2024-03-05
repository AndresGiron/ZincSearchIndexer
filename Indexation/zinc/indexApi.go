package zinc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ExecuteAll() {
	// CheckIndex()
	CreateIndex()

}

const user = "admin"
const password = "Complexpass#123"
const indexPath = "http://localhost:4080/api/index/"

func CreateIndex() error {
	const emailsIndexMapping = `
	{
		"name": "emails",
		"storage_type": "disk",
		"mappings": {
			"properties": {
				"messageId": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": true,
				},
				"date": {
					"type": "date",
					"format": "2006-01-02T15:04:05Z07:00",
					"index": true,
					"store": true,
					"sortable": true,
				},
				"from": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": true,
				},
				"to": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": true,
				},
				"subject": {
					"type": "text",
					"index": true,
					"store": true,
					"sortable": false,
				},
				"body": {
					"type": "text",
					"index": true,
					"store": true,
					"sortable": false,
				}
			}
		}
	}`

	req, err := http.NewRequest("POST", indexPath, bytes.NewReader([]byte(emailsIndexMapping)))
	if err != nil {
		return err
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return nil
}
