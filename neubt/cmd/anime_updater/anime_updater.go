package main

import (
	"crawler/neubt"
	neubtCrawler "crawler/neubt/crawler"
	"crawler/ptgen"
	"log"
)

type Updater struct {
	*neubt.NeuBT
	ptgenClient ptgen.Client
	ptgen       ptgen.PTGen
}

func NewUpdater() *Updater {
	ptgenClient, err := ptgen.NewClient()
	if err != nil {
		log.Println("init ptgen client failed")
		return nil
	}
	nb := neubt.NewNeuBT()
	return &Updater{
		NeuBT:       nb,
		ptgenClient: ptgenClient,
		ptgen:       ptgen.NewBufferedPTGen(ptgenClient, nb.KVS),
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
	detail, err := u.ptgen.GetBangumiDetailByLink(bgmURL)
	if err != nil {
		log.Fatal("can not get bgm detail")
	}
	text, err := ptgen.GetTextFromDetail(detail)
	if err != nil {
		log.Fatal("failed to get text from ptgen detail: ", err)
	}
	err = poster.SetPTGENContent(text)
	if err != nil {
		log.Fatal("failed to SetPTGENContent: ", err)
	}
	err = poster.UpdateTorrentMultiPart()
	if err != nil {
		log.Fatal("failed to post torrent to neu bt: ", err)
	}
}
