package engines

import (
	"io"

	"github.com/XORbit01/retro/shared"
)

type Engine interface {
	Search(query string, maxResults int) ([]shared.SearchResult, error)
	Download(url string) (io.ReadCloser, string, error)
	Exists(url string) (bool, error)
	Name() shared.DResults
	MaxResults() int
}
