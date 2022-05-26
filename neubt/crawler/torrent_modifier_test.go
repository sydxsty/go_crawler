package crawler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModifyTorrent(t *testing.T) {
	tp, err := NewTorrentModifier("/forum.php?mod=post&action=edit&fid=44&tid=1698779&pid=27981153&page=1", client)
	assert.NoError(t, err, "error create NewTorrentModifier")
	tp.SetTitle("间谍过家家", "SPY×FAMILY", "01", "MKV", "Kamigami", "简繁日内封字幕", "1080p x265 Ma10p AAC")
	tp.SetPostFileName("[间谍过家家][SPY×FAMILY][01][MKV][Kamigami][简繁日内封字幕][1080p x265 Ma10p AAC]")
	err = tp.SetPTGENContent("AAA[img]https://lain.bgm.tv/pic/cover/l/bd/15/343656_j6eWd.jpg[/img]AAA")
	assert.NoError(t, err, "error replace image")
	tp.SetCommentContent("comment content is here")
	err = tp.UpdateTorrentMultiPart()
	assert.NoError(t, err, "error post torrent")
}
