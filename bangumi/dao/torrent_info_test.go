package dao

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var torrentInfoSample map[string]interface{}

func init() {
	err := json.Unmarshal([]byte(`{
  "_id": "62b3d0c1f248710007c5eab6",
  "category_tag_id": "549ef207fe682f7549f1ea90",
  "title": "[百冬练习组&LoliHouse] 身为女主角！～被讨厌的女主角和秘密的工作～ / Heroine Tarumono! - 10 [WebRip 1080p HEVC-10bit AAC][简繁内封字幕]",
  "tag_ids": [
    "581be821ee98e9ca20730eae",
    "624774a370e8f89b23a1855a",
    "548ee2ce4ab7379536f56358",
    "549ef207fe682f7549f1ea90",
    "5c6c2248475bb7b273306824"
  ],
  "comments": 0,
  "downloads": 56,
  "finished": 248,
  "leechers": 3,
  "seeders": 51,
  "uploader_id": "5c4fce0060a958730e124848",
  "team_id": "581b44bfee98e9ca20730e9a",
  "publish_time": "2022-06-23T02:32:33.231Z",
  "magnet": "magnet:?xt=urn:btih:0c2d0c785d3fdbff1fcd963197bf5271e06c7061",
  "infoHash": "0c2d0c785d3fdbff1fcd963197bf5271e06c7061",
  "file_id": "62b3d0bdf248710007c5eab1",
  "teamsync": true,
  "content": [
    [
      "[hyakuhuyu&LoliHouse] Heroine Tarumono! Kiraware Heroine to Naisho no Oshigoto - 10 [WebRip 1080p HEVC-10bit AAC ASSx2].mkv",
      "272.94 MB"
    ]
  ],
  "size": "272.94 MB",
  "btskey": "",
  "sync": {
    "dmhy": "http://share.dmhy.org/topics/view/603687_LoliHouse_Heroine_Tarumono%21_-_10_WebRip_1080p_HEVC-10bit_AAC.html",
    "acgrip": "https://acg.rip/t/257657",
    "nyaa": "unexpected body content"
  }
}`), &torrentInfoSample)
	if err != nil {
		log.Fatalln("failed to init data")
	}
}

func TestTorrentInfo(t *testing.T) {
	ti, err := NewTorrentInfoFromMap(torrentInfoSample)
	assert.NoError(t, err, "can not init TorrentInfo")
	ti.InitTorrentDetail(nil)
	teams := ti.MustGetTeamStr()
	assert.True(t, teams == "百冬练习组&LoliHouse")
	res := ti.MustGetResolution()
	assert.True(t, res == "1080p")
	ep := ti.MustGetEpisode()
	assert.True(t, ep == "10")
	co := ti.GetContent()
	log.Println(co)
}
