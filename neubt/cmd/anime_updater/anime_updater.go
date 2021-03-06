package main

import (
	"crawler/neubt"
	neubtCrawler "crawler/neubt/crawler"
	"crawler/ptgen"
	"log"
)

type Updater struct {
	*neubt.NeuBT
	ptgen ptgen.PTGen
}

func NewUpdater() *Updater {
	nb := neubt.NewNeuBT()
	pg, err := ptgen.NewBufferedPTGen(nb.KVS)
	if err != nil {
		log.Println("init ptgen client failed")
		return nil
	}
	return &Updater{
		NeuBT: nb,
		ptgen: pg,
	}
}

func main() {
	u := NewUpdater()
	neubtURL := `/forum.php?mod=post&action=edit&fid=44&tid=1698779&pid=27981153&page=1`
	bgmURL := `https://movie.douban.com/subject/35651863/`
	poster, err := neubtCrawler.NewTorrentModifier(neubtURL, u.Client)
	if err != nil {
		log.Fatal("can not init modifier")
	}
	info, err := u.ptgen.GetBangumiInfoByLink(bgmURL)
	if err != nil {
		log.Fatal("can not get bgm detail")
	}
	detail, err := ptgen.GetDetailFromInfo(info)
	if err != nil {
		log.Fatal("failed to get text from ptgen detail: ", err)
	}
	err = poster.SetPTGENContent(detail.Detail)
	if err != nil {
		log.Fatal("failed to SetPTGENContent: ", err)
	}
	err = poster.UpdateTorrentMultiPart()
	if err != nil {
		log.Fatal("failed to post torrent to neu bt: ", err)
	}
}
