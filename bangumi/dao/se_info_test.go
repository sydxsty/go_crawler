package dao

import (
	bgm "crawler/bangumi/anime_control"
	"log"
	"testing"
)

func TestSEInfo(t *testing.T) {
	f := bgm.NewBangumiFilter()
	for _, in := range testCase {
		log.Println(NewSEInfoFromTitle(in, f).GetEpisodeStringList())
	}
}

var testCase []string

func init() {
	testCase = []string{
		"【喵萌奶茶屋】★04月新番★[Estab Life: Great Escape][01-12END][1080p][简体][招募翻译校对]",
		"[桜都字幕组] 间谍过家家 / Spy x Family [09v2][1080p][简繁内封]",
		"[桜都字幕组] 间谍过家家 / Spy x Family [2022][09][1080p][简繁内封]",
		"[桜都字幕组] 间谍过家家 / Spy x Family - OVA[1080p][简繁内封]",
		"[喵萌奶茶屋&VCB-Studio] Bokutachi wa Benkyō ga Dekinai / 我们无法一起学习 / ぼくたちは勉強ができない 10-bit 1080p HEVC BDRip [S1 + S2 + OVA Fin]",
	}
}
