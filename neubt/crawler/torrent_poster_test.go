package crawler

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"testing"
)

func TestPostTorrent(t *testing.T) {
	tp, err := NewTorrentPoster("44", client)
	assert.NoError(t, err, "error create NewTorrentPoster")
	tp.SetTitle("间谍过家家", "SPY×FAMILY", "01", "MKV", "Kamigami", "简繁日内封字幕", "1080p x265 Ma10p AAC")
	res := tp.SetTidByName("连载动画")
	assert.True(t, res)
	res = tp.SetTid("222")
	assert.True(t, res)
	tp.SetPostFileName("[间谍过家家][SPY×FAMILY][01][MKV][Kamigami][简繁日内封字幕][1080p x265 Ma10p AAC]")
	tp.SetPTGENContent("pt-gen content is here")
	tp.SetCommentContent("comment content is here")
	data, _ := ioutil.ReadFile("test.torrent")
	url, err := tp.PostTorrentMultiPart(data)
	assert.NoError(t, err, "error post torrent")
	log.Println(url)
}
