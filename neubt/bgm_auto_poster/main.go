package main

import (
	"bytes"
	"crawler/bangumi"
	"crawler/bangumi/crawler"
	"crawler/bangumi/dao"
	"crawler/neubt"
	neubtCrawler "crawler/neubt/crawler"
	"crawler/ptgen"
	"crawler/qbt"
	"crawler/util"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"os"
	"regexp"
	"time"
)

type Poster struct {
	*neubt.NeuBT
	bgmClient     bangumi.Client
	bgm           crawler.Bangumi
	bgmDownloader crawler.Downloader
	bgmTrMgr      *dao.TorrentManager
	ptgenClient   ptgen.Client
	ptgen         ptgen.PTGen
}

func NewPoster() *Poster {
	bgmClient, err := bangumi.NewClient()
	if err != nil {
		log.Println("init bangumi client failed")
		return nil
	}
	ptgenClient, err := ptgen.NewClient()
	if err != nil {
		log.Println("init ptgen client failed")
		return nil
	}
	nb := neubt.NewNeuBT()
	return &Poster{
		NeuBT:         nb,
		bgmClient:     bgmClient,
		bgm:           crawler.NewBangumi(bgmClient, nb.KVS),
		bgmDownloader: crawler.NewDownloader(bgmClient),
		bgmTrMgr:      dao.NewTorrentManager(nb.KVS),
		ptgenClient:   ptgenClient,
		ptgen:         ptgen.NewBufferedPTGen(ptgenClient, nb.KVS),
	}
}

func main() {
	p := NewPoster()
	for {
		err := crawler.ScanBangumiTorrent(p.bgm, func(ti *dao.BangumiTorrentInfo) {
			log.Println("--------Analysing torrent: ", ti.Title, "--------")
			if p.bgmTrMgr.TorrentIsPosted(ti.InfoHash) {
				// if posted, continue
				log.Println("torrent is already posted, skip")
				return
			}
			if p.Webui.Contains(ti.InfoHash) && !p.Webui.Completed(ti.InfoHash) {
				log.Println("torrent is downloading: ", ti.Title)
				return
			}
			// 1. torrent not exist
			// 2. torrent completed
			detail, err := p.GetTorrentPTGenDetail(ti)
			if err != nil {
				log.Println("no matching result in ptgen: ", err)
				return
			}
			// for torrents not exist
			if !p.Webui.Completed(ti.InfoHash) {
				log.Println("prepare to download torrent from bangumi: ", ti.Title)
				if err := p.bgmTrMgr.CanDownloadFromBangumi(ti); err != nil {
					log.Println("filter failure reason", err)
					return
				}
				_, err = crawler.DownloadBangumiTorrentToFile(
					ti.Detail.TorrentDownloadURL,
					p.Config.TorrentPath,
					ti.InfoHash,
					p.bgmDownloader,
					p.Webui)
				if err != nil {
					log.Println("can not download bangumi torrent", err)
				}
				log.Println("downloaded bangumi torrent: ", ti.Title)
				return
			}
			// for completed torrents
			log.Println("prepare to post torrent to neubt: ", ti.Title)
			// pause torrent to reduce network overhead
			err = p.Webui.Pause(ti.InfoHash)
			if err != nil {
				log.Println("can not pause torrent, webui may have fault: ", err)
				return
			}
			poster, err := neubtCrawler.NewTorrentPoster("44", p.Client)
			if err != nil {
				log.Println("failed to create neubt poster: ", err)
				return
			}
			err = UpdateWithTorrentInfo(poster, ti)
			if err != nil {
				log.Println("failed to update bangumi torrent info: ", err)
				return
			}
			text, err := ptgen.GetTextFromDetail(detail)
			if err != nil {
				log.Println("failed to get text from ptgen detail: ", err)
				return
			}
			poster.SetPTGENContent(text)
			mediaInfo, err := GetMediaInfoFromWEBUI(ti.InfoHash, p.Webui)
			if err != nil {
				log.Println("failed to get media info: ", err)
				return
			}
			poster.SetMediaInfoContent(mediaInfo)
			data, err := crawler.LoadTorrentFromFile(p.Config.TorrentPath, ti.InfoHash)
			if err != nil {
				log.Println("failed to load torrent from disk: ", err)
				return
			}
			// wait for 5 second
			time.Sleep(time.Second * 5)
			// mark the torrent posted
			p.bgmTrMgr.SetTorrentPostedState(ti.InfoHash)
			url, err := poster.PostTorrentMultiPart(data)
			if err != nil {
				log.Println("failed to post torrent to neu bt: ", err)
				return
			}
			// wait for 5 second
			time.Sleep(time.Second * 5)
			err = p.downloadTorrentByLink(url)
			if err != nil {
				log.Println("failed to download post torrent to neu bt: ", err)
				return
			}
		})
		if err != nil {
			log.Println("can not load bangumi latest torrents")
		}
		time.Sleep(time.Second * 600)
	}
}

func (p *Poster) GetTorrentPTGenDetail(info *dao.BangumiTorrentInfo) (map[string]interface{}, error) {
	links, err := p.ptgen.GetBangumiLinkByNames(
		info.Detail.TorrentJpnName,
		info.Detail.TorrentChsName,
		info.Detail.TorrentEngName)
	if err != nil {
		return nil, err
	}
	for _, v := range links {
		result, err := p.ptgen.GetBangumiDetailByLink(v)
		if err == nil {
			return result, nil
		}
	}
	return nil, errors.New("no matching result")
}

func (p *Poster) downloadTorrentByLink(link string) error {
	detail := neubtCrawler.NewThreadDetail(p.Client)
	torrentURLs, err := detail.GetFloorDetailFromThread(link)
	if err != nil {
		log.Println(err, "wrong torrent info")
	}
	downloader := neubtCrawler.NewDownloader(p.Client)
	for _, torrentURL := range torrentURLs {
		data, err := downloader.DownloadFromNestedURL(torrentURL.Comment.TorrentLink)
		if err != nil {
			log.Println(err, "can not download torrent from url")
			continue
		}
		err = os.WriteFile(p.Config.TorrentPath+torrentURL.Comment.TorrentName, data, 0666)
		if err != nil {
			log.Println(err, "write file failure")
			return err
		}
		err = p.Webui.AddTorrentFromFile(p.Config.TorrentPath, torrentURL.Comment.TorrentName)
		if err != nil {
			log.Println(err, "failed to upload torrent to webui")
			continue
		}
		log.Println("the posted torrent is added to webui: ", link)
		return nil
	}
	return errors.New("no torrent found in neu bt")
}

func UpdateWithTorrentInfo(poster neubtCrawler.TorrentPoster, info *dao.BangumiTorrentInfo) error {
	poster.SetTidByName("连载动画")
	poster.SetPostFileName(info.Title)
	if _, err := info.GetCHNName(); err != nil {
		return errors.New("no Chinese name or English name found in info")
	}
	if info.Detail.TeamName == "" {
		return errors.New("no team name found in info")
	}
	poster.SetTitle(info.Detail.TorrentChsName,
		info.Detail.TorrentEngName,
		info.Detail.TorrentJpnName,
		info.Detail.Episode,
		info.Detail.Format,
		info.Detail.TeamName,
		info.Detail.Language,
		info.Detail.Resolution,
	)
	form := func(v interface{}) string {
		detail, _ := json.Marshal(v)
		var out bytes.Buffer
		if err := json.Indent(&out, detail, "", "\t"); err != nil {
			return ""
		}
		return out.String()
	}
	poster.SetCommentContent(
		"[code]",
		"Debug info:",
		"原种标题："+info.Title,
		"种子信息："+form(info.Detail),
		"种子内容："+form(info.Content),
		"[/code]",
	)
	return nil
}

func GetMediaInfoFromWEBUI(infoHash string, webui qbt.WEBUIHelper) (string, error) {
	// generate media info
	log.Println("generating media info")
	torrent, files, err := webui.GetTorrentDetail(infoHash)
	if err != nil {
		return "", errors.Wrap(err, "failed to get torrent files")
	}
	savePath := torrent.SavePath
	var fileName string
	for _, file := range files {
		if IsVideoFile(file.Name) {
			fileName = file.Name
			break
		}
		log.Println(file.Name, "is not a video file")
	}
	if len(fileName) == 0 {
		return "", errors.New("can not find valid video file in torrent")
	}
	info, err := util.GetMediaInfo("./lib", savePath, fileName)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate media info of file"+fileName)
	}
	return info, nil
}

func IsVideoFile(name string) bool {
	result := regexp.MustCompile(`(.mp4|.ts|.mkv)`).FindAllString(name, -1)
	if result == nil {
		return false
	}
	return true
}
