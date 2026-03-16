package downloader

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// YtDlpClient implements gateway.ContentDownloader
type YtDlpClient struct {
	OutputDir string
	// SubLangs: subtitle languages to try in order (e.g. []string{"pt", "en"}). If empty, defaults to ["en"].
	SubLangs []string
}

func NewYtDlpClient(outputDir string) *YtDlpClient {
	return &YtDlpClient{OutputDir: outputDir, SubLangs: []string{"en"}}
}

// NewYtDlpClientWithLangs creates a client with explicit subtitle language preference.
func NewYtDlpClientWithLangs(outputDir string, langs []string) *YtDlpClient {
	if len(langs) == 0 {
		langs = []string{"en"}
	}
	return &YtDlpClient{OutputDir: outputDir, SubLangs: langs}
}

func (y *YtDlpClient) Download(url string) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(y.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Pre-clean old files (optional)
	y.cleanOldFiles()

	outputTemplate := filepath.Join(y.OutputDir, "raw_legenda.%(id)s")
	subLangs := y.SubLangs
	if len(subLangs) == 0 {
		subLangs = []string{"en"}
	}
	langArg := strings.Join(subLangs, ",")

	cmd := exec.Command("yt-dlp",
		"--write-auto-subs",
		"--write-subs",
		"--skip-download",
		"--sub-lang", langArg,
		"-o", outputTemplate,
		url,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("yt-dlp failed: %v | Log: %s", err, stderr.String())
	}

	return y.findAndReadVTT()
}

func (y *YtDlpClient) cleanOldFiles() {
	files, _ := filepath.Glob(filepath.Join(y.OutputDir, "raw_legenda*"))
	for _, f := range files {
		os.Remove(f)
	}
}

func (y *YtDlpClient) findAndReadVTT() (string, error) {
	matches, _ := filepath.Glob(filepath.Join(y.OutputDir, "raw_legenda*.vtt"))
	if len(matches) > 0 {
		content, err := os.ReadFile(matches[0])
		return string(content), err
	}
	return "", fmt.Errorf("no .vtt subtitle found after execution")
}
