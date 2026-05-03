package mosh

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// ScrapePage fetches the given URL and extracts the __NEXT_DATA__ JSON.
func ScrapePage(url string, target any) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("error closing response body: %v\n", err)
		}
	}(res.Body)

	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s for URL: %s", res.StatusCode, res.Status, url)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	selection := doc.Find("#__NEXT_DATA__").First()
	if selection.Length() == 0 {
		return fmt.Errorf("__NEXT_DATA__ not found")
	}

	content := selection.Text()
	jsonData := []byte(content)
	if err := json.Unmarshal(jsonData, target); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return nil
}
