package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type QuerySearch struct {
	Term      string `json:"term"`
	From      string `json:"from"`
	MaxResult string `json:"max_results"`
}

type QueryAll struct {
	From      string `json:"from"`
	MaxResult string `json:"max_results"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	res := map[string]interface{}{"message": "helloWorld"}
	_ = json.NewEncoder(w).Encode(res)
}

func Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	var body QuerySearch
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Fatal(err)
	}

	bodyJson, err := json.Marshal(body.Term)
	if err != nil {
		log.Fatal(err)
	}
	//recordar agregar el asterisco al final de la sentencia en el front
	fmt.Println(string(bodyJson))
	query := `{
		"search_type": "querystring",
        "query":
        {
            "term": ` + string(bodyJson) + `
        },
        "from": ` + string(body.From) + `,
        "max_results": ` + string(body.MaxResult) + `,
		"sort_fields": ["date"],
        "_source": []
    }`
	fmt.Println(strings.NewReader(query))

	req, err := http.NewRequest("POST", searchPath, strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	results, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(results)

}

func ListAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	var body QueryAll
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Fatal(err)
	}

	query := `{
        "search_type": "alldocuments",
        "from": ` + string(body.From) + `,
        "max_results":` + string(body.MaxResult) + ` ,
		"sort_fields": [],
        "_source": []
    }`
	fmt.Println(strings.NewReader(query))

	req, err := http.NewRequest("POST", searchPath, strings.NewReader(query))
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(user, password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	results, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(results)

}
