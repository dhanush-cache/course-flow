package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// TODO: Verify the installation of ffmpeg and ffprobe.

var metadataArgs = []string{
	"-map_metadata", "-1",
	"-map_chapters", "-1",
	"-metadata:g", `encoder=" "`,
}

// execFFmpeg runs an ffmpeg command with the given arguments and returns its output or an error.
func execFFmpeg(additionalArgs ...string) (string, error) {
	args := []string{"-y"}
	args = append(args, additionalArgs...)
	args = insertMetadataArgs(args)

	cmd := exec.Command("ffmpeg", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg execution failed: %w: %s", err, string(output))
	}
	return string(output), nil
}

// insertMetadataArgs inserts metadata removal arguments into the ffmpeg command arguments.
func insertMetadataArgs(args []string) []string {
	insertIndex := 0
	for i := 1; i < len(args)-1; i++ {
		if args[i] == "-i" {
			insertIndex = i + 2
		}
	}
	return append(args[:insertIndex], append(metadataArgs, args[insertIndex:]...)...)
}

// GetThumbnail extracts a JPEG thumbnail from a video at the given timestamp (in seconds).
func GetThumbnail(videoPath string, timestamp int, outFile string) (string, error) {
	args := []string{
		"-ss", fmt.Sprintf("%d", timestamp),
		"-i", videoPath,
		"-vframes", "1",
		outFile,
	}

	if _, err := execFFmpeg(args...); err != nil {
		return "", fmt.Errorf("extract thumbnail from %q at %ds: %w", videoPath, timestamp, err)
	}

	return outFile, nil
}

// HasEmbeddedSubs checks if a video file has embedded subtitles.
func HasEmbeddedSubs(videoPath string) (bool, error) {
	cmd := exec.Command("ffprobe", videoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("ffprobe failed for %q: %w: %s", videoPath, err, string(output))
	}
	return strings.Contains(string(output), "Subtitle"), nil
}

// GetTitle extracts the title from a video file name by removing leading numbering.
func GetTitle(videoPath string) string {
	pattern := `^\d+\s*-\s*`
	re := regexp.MustCompile(pattern)

	filename := filepath.Base(videoPath)
	ext := filepath.Ext(filename)
	stem := filename[:len(filename)-len(ext)]

	return re.ReplaceAllString(stem, "")
}

// GetBlankVideo generates a blank video of specified duration (in seconds) and returns its file path.
func GetBlankVideo(duration int) (string, error) {
	tmpFile, err := os.CreateTemp("", "blank_video_*.mp4")
	if err != nil {
		return "", fmt.Errorf("create temp blank video file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return "", fmt.Errorf("close temp blank video file: %w", err)
	}

	tmpFilePath := tmpFile.Name()

	args := []string{
		"-f", "lavfi",
		"-i", fmt.Sprintf("color=c=black:s=1280x720:d=%d", duration),
		"-f", "lavfi",
		"-i", fmt.Sprintf("anullsrc=r=44100:cl=stereo:d=%d", duration),
		"-c:v", "libx264",
		"-c:a", "aac",
		"-pix_fmt", "yuv420p",
		"-shortest",
		tmpFilePath,
	}

	if _, err := execFFmpeg(args...); err != nil {
		return "", fmt.Errorf("generate blank video (%ds): %w", duration, err)
	}

	return tmpFilePath, nil
}

type AdvancedOpts struct {
	ThumbnailTS int
	Subtitles   string
}

// FFprocess processes the videos in `source` and places them in `dest`
func FFprocess(source string, dest string, options AdvancedOpts) error {
	dir := filepath.Dir(dest)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create destination directory %q: %w", dir, err)
	}

	args := []string{"-i", source}

	hasExternalSubs := options.Subtitles != ""
	if hasExternalSubs {
		args = append(args, "-i", options.Subtitles)
	}

	title := GetTitle(dest)

	args = append(args,
		"-metadata", fmt.Sprintf("title=%v", title),
		"-metadata:s:a:0", "language=en",
	)

	hasEmbeddedSubs, err := HasEmbeddedSubs(source)
	if err != nil {
		return fmt.Errorf("check embedded subtitles for %q: %w", source, err)
	}

	hasSubs := hasExternalSubs || hasEmbeddedSubs
	if hasSubs {
		args = append(args, "-metadata:s:s:0", "language=en")
	}
	tmpFile, err := os.CreateTemp("", "thumb-*.jpeg")
	if err != nil {
		return fmt.Errorf("create temp thumbnail file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("close temp thumbnail file: %w", err)
	}
	
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			fmt.Println("temp thumbnail remove error:", err)
		}
	}(tmpFile.Name())

	thumbnail, err := GetThumbnail(source, options.ThumbnailTS, tmpFile.Name())
	if err != nil {
		return fmt.Errorf("generate thumbnail for %q: %w", source, err)
	}

	args = append(args,
		"-attach", thumbnail,
		"-metadata:s:t", fmt.Sprintf("filename=%v", title),
		"-metadata:s:t", "mimetype=image/jpeg",
		"-map", "0:v",
		"-map", "0:a",
		"-c", "copy",
	)

	if hasExternalSubs {
		args = append(args, "-map", "1:s")
	} else if hasEmbeddedSubs {
		args = append(args, "-map", "0:s")
	}

	if hasSubs {
		args = append(args, "-c:s", "srt")
	}

	args = append(args, dest)

	if _, err := execFFmpeg(args...); err != nil {
		return fmt.Errorf("final ffmpeg processing (%q → %q): %w", source, dest, err)
	}

	return nil
}
