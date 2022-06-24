package dao

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestTFL(t *testing.T) {
	tl, err := NewTorrentFileList(fl)
	assert.NoError(t, err, "can not init tl")
	_, err = tl.PrintToStringList()
	assert.NoError(t, err, "can not print tl")
	log.Println(tl.GetTorrentName())
	toString, err := tl.PrintToString(10)
	assert.NoError(t, err, "can not print tl")
	log.Println(toString)
	toString, err = tl.PrintToString(-1)
	assert.NoError(t, err, "can not print tl")
	log.Println(toString)
}

var fl []interface{}

func init() {
	err := json.Unmarshal([]byte(`[
  [
    "[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [Ma10p_1080p]/CDs/[210506] ｢EVERYBODY! EVERYBODY!／YOU YOU YOU｣ [24bit_96kHz] (flac)/01. EVERYBODY! EVERYBODY!.flac",
    "122.28 MB"
  ],
  [
    "[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [Ma10p_1080p]/CDs/[210506] ｢EVERYBODY! EVERYBODY!／YOU YOU YOU｣ [24bit_96kHz] (flac)/天使动漫自购转载声明.txt",
    "1021 B"
  ],
  [
    "[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [Ma10p_1080p]/CDs/[210625] SPCD 01 (flac)/01. Dangerous Escape.flac",
    "12.24 MB"
  ],
  [
    "[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [Ma10p_1080p]/CDs/[210625] SPCD 01 (flac)/02. Chaotic Haze.flac",
    "6.63 MB"
  ],
  [
    "[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [Ma10p_1080p]/SPs/[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [CM][Ma10p_1080p][x265_flac].mkv",
    "7.85 MB"
  ],
  [
    "[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [Ma10p_1080p]/[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [10][Ma10p_1080p][x265_flac].tc.ass",
    "37.22 KB"
  ],
  [
    "[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [Ma10p_1080p]/[Nekomoe kissaten&VCB-Studio] Isekai Maou to Shoukan Shoujo no Dorei Majutsu Omega [Fonts].7z",
    "29.85 MB"
  ]
]`), &fl)
	if err != nil {
		log.Fatalln("failed to init data")
	}
}
