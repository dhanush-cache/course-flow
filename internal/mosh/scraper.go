package mosh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// GetToken fetches the buildId token from the codewithmosh website.
func GetToken() (*Token, error) {
	url := "https://codewithmosh.com/"

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	selection := doc.Find("#__NEXT_DATA__").First()
	jsonData := []byte(selection.Text())
	var obj map[string]any
	if err := json.Unmarshal(jsonData, &obj); err != nil {
		return nil, err
	}

	token, ok := obj["buildId"].(string)
	if !ok {
		return nil, fmt.Errorf("buildId not found or not a string")
	}
	return &Token{
		Value:     token,
		ExpiresAt: time.Now().AddDate(0, 0, 1),
	}, nil
}
