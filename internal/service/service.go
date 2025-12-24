package service

import (
	"fmt"
	"os"
	"strings"

	config "github.com/dhanush-cache/course-flow/internal"
	"github.com/dhanush-cache/course-flow/internal/adapters"
	"github.com/dhanush-cache/course-flow/internal/utils"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// TODO: Understand the bars better and refactor the code to make it cleaner.

func Process(zipFiles []string, targets []string, cfg *config.Config) error {
	var line string
	lines := utils.BuildTree(targets)
	lineIndex := 0

	tempDir, err := os.MkdirTemp("", "zip-extract-*")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			fmt.Println("failed to remove tempDir:", err)
		}
	}()

	for _, zip := range zipFiles {
		extractBar := getExtractBar()
		err = utils.Unzip(zip, tempDir, func(done, total int, _ string) {
			extractBar.SetTotal(int64(total), total == done)
			extractBar.SetCurrent(int64(done))
		})
		if err != nil {
			return err
		}
	}

	ch := make(chan any)
	processBar := getProcessBar(&line, ch)

	err = utils.ProcessVideos(tempDir, targets, cfg, func(done, total int, current string) {
		processBar.SetTotal(int64(total), total == done)
		processBar.SetCurrent(int64(done))
		var b strings.Builder
		for {
			s := lines[lineIndex]
			b.WriteString(s)
			lineIndex++
			if strings.Contains(s, cfg.VideoExt) || lineIndex >= len(lines) {
				break
			} else {
				b.WriteString("\n")
			}
		}
		line = b.String()
		ch <- nil
	})
	if err != nil {
		return err
	}
	return nil
}

func getExtractBar() *mpb.Bar {
	p := mpb.New()
	return p.AddBar(0,
		mpb.BarRemoveOnComplete(),
		mpb.PrependDecorators(
			decor.Name("Extracting", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Percentage(decor.WC{W: 5}), "done"),
		),
	)
}
func getProcessBar(line *string, ch chan any) *mpb.Bar {
	p := mpb.New(mpb.WithManualRefresh(ch))

	return p.AddBar(0,
		mpb.BarNoPop(),
		mpb.PrependDecorators(
			adapters.Line(func(statistics decor.Statistics) string { return *line }),
			decor.Name("Processing", decor.WC{C: decor.DindentRight | decor.DextraSpace}),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
			decor.OnComplete(decor.AverageETA(decor.ET_STYLE_GO, decor.WCSyncWidth), ""),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Percentage(decor.WC{W: 5}), "done"),
		),
	)
}
