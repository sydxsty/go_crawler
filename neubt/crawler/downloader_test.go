package crawler

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

// TestDownloadTorrent test download a torrent to memory buffer
func TestDownloadTorrent(t *testing.T) {
	downloader := NewDownloader(client)
	data, err := downloader.DownloadFromNestedURL(`forum.php?mod=attachment&aid=NTQyNjIxN3w3NWJlODg3MHwxNjUzMDYzMTc2fDY4OTk4MHwxNjk4Mjgx`)
	assert.NoError(t, err, "error redirect url")
	assert.True(t, len(data) > 1024*100)
	f, _ := os.Create(`test.torrent`)
	io.Copy(f, bytes.NewReader(data))
	f.Close()
}
