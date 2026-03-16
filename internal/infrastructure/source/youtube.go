package source

import (
	"context"

	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/domain/gateway"
	"github.com/lucas-de-lima/domain-agnostic-llm-multi-agent-pipeline/internal/infrastructure/downloader"
)

// YouTubeSource implements gateway.ContentSource by fetching subtitles via yt-dlp.
type YouTubeSource struct {
	dl *downloader.YtDlpClient
}

// NewYouTubeSource creates a ContentSource that uses yt-dlp for YouTube (and compatible) URLs.
// SubLangs on the client define which subtitle languages to try (e.g. []string{"pt", "en"}).
func NewYouTubeSource(ytdlp *downloader.YtDlpClient) *YouTubeSource {
	return &YouTubeSource{dl: ytdlp}
}

// Fetch returns raw subtitle text for the given URL. Implements gateway.ContentSource.
func (y *YouTubeSource) Fetch(ctx context.Context, source string) (string, error) {
	// Context not yet passed to yt-dlp; could add exec.CommandContext(ctx,...) in downloader later
	_ = ctx
	return y.dl.Download(source)
}

// Ensure YouTubeSource implements gateway.ContentSource at compile time.
var _ gateway.ContentSource = (*YouTubeSource)(nil)
