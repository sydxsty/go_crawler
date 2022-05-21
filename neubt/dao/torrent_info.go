package dao

type rawTorrentInfo struct {
	Title        string
	Link         string
	Size         string
	Discount     string
	Signal       string
	Uploader     string
	UploaderPage string
}

type TorrentInfo struct {
	Title     string
	Link      string
	TorrentID int
	Size      float64 // MB
	Discount  int     // 0 to 100
	Crawled   bool
}
