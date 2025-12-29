package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// GDriveDirectDownloadURL returns a direct Google Drive download URL
// that works for both small and large files (handles confirm + uuid params).
func GDriveDirectDownloadURL(fileID string) (string, error) {
	baseURL := fmt.Sprintf(
		"https://drive.google.com/uc?export=download&id=%s",
		fileID,
	)

	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.Get(baseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return "", err
		}

		uuid, exists := doc.Find(`input[type="hidden"][name="uuid"]`).Attr("value")
		if !exists {
			return "", errors.New("could not extract Google Drive uuid token")
		}

		return fmt.Sprintf(
			"https://drive.usercontent.google.com/download?id=%s&export=download&authuser=0&confirm=t&uuid=%s",
			fileID,
			uuid,
		), nil
	}

	return baseURL, nil
}

// DownloadFile downloads a file from a direct URL to a local filepath.
// It supports very large files (multi-GB) and displays a progress bar.
func DownloadFile(url string, filepath string) error {
	// Create HTTP client with no timeout (important for large files)
	client := &http.Client{
		Timeout: 0,
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to start download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create destination file
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Content length may be -1 if unknown
	contentLength := resp.ContentLength

	// Initialize progress container
	p := mpb.New(
		mpb.WithWidth(60),
		mpb.WithRefreshRate(180*time.Millisecond),
	)

	bar := p.New(
		contentLength,
		mpb.BarStyle(),
		mpb.PrependDecorators(
			decor.Name("Downloading: ", decor.WC{W: 13}),
			decor.CountersKibiByte("% .2f / % .2f"),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WC{W: 5}),
			decor.EwmaSpeed("kb", "% .2f", 60),
			decor.EwmaETA(decor.ET_STYLE_GO, 60),
		),
	)

	// Wrap response body with progress bar reader
	reader := bar.ProxyReader(resp.Body)
	defer reader.Close()

	// Stream copy (constant memory usage)
	_, err = io.Copy(out, reader)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	// Wait for progress bar to finish rendering
	p.Wait()

	return nil
}
