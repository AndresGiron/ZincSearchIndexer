package zinc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ExecuteAll() {
	CheckIndex()
	CreateIndex()

}

const user = "admin"
const password = "Complexpass#123"
const indexPath = "http://localhost:4080/api/index/"

func CheckIndex() (bool, error) {
	req, err := http.NewRequest("HEAD", indexPath+"emailsT", nil)
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(user, password)

	resp, err := http.DefaultClient.Do(req)
	fmt.Println(resp)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return true, nil
}

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
					"aggregatable": false,
					"highlightable": false
				},
				"date": {
					"type": "date",
					"format": "2006-01-02T15:04:05Z07:00",
					"index": true,
					"store": false,
					"sortable": true,
					"aggregatable": true,
					"highlightable": false
				},
				"from": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": true,
					"aggregatable": true,
					"highlightable": false
				},
				"to": {
					"type": "keyword",
					"index": true,
					"store": true,
					"sortable": true,
					"aggregatable": true,
					"highlightable": false
				},
				"subject": {
					"type": "text",
					"index": true,
					"store": false,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"body": {
					"type": "text",
					"index": true,
					"store": false,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
				},
				"isRead": {
					"type": "boolean",
					"index": true,
					"store": false,
					"sortable": false,
					"aggregatable": false,
					"highlightable": false
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
