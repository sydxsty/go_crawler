package main

import (
	"crawler/bangumi"
	"crawler/bangumi/crawler"
	"crawler/bangumi/dao"
	"crawler/neubt"
	neubtCrawler "crawler/neubt/crawler"
	"crawler/ptgen"
	"crawler/qbt"
	"crawler/util"
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
	ptgen         ptgen.PTGen
	ani           *dao.AnimeDB
}

func NewPoster() *Poster {
	bgmClient, err := bangumi.NewClient()
	if err != nil {
		log.Println("init bangumi client failed")
		return nil
	}
	nb := neubt.NewNeuBT()
	pg, err := ptgen.NewBufferedPTGen(nb.KVS)
	if err != nil {
		log.Println("init ptgen client failed")
		return nil
	}
	ani, err := dao.NewAnimeDB()
	if err != nil {
		log.Println("init ani db failed")
		return nil
	}
	return &Poster{
		NeuBT:         nb,
		bgmClient:     bgmClient,
		bgm:           crawler.NewBangumi(bgmClient, nb.KVS),
		bgmDownloader: crawler.NewDownloader(bgmClient),
		bgmTrMgr:      dao.NewTorrentManager(nb.KVS),
		ptgen:         pg,
		ani:           ani,
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
				log.Println("torrent is downloading: ", ti.InfoHash)
				return
			}
			// 1. torrent not exist
			// 2. torrent completed
			detail, err := p.GetTorrentPTGenDetail(ti)
			if err != nil {
				log.Println("no matching result in ptgen: ", err)
				return
			}
			log.Printf("CHSName: %s, ENGName: %s, JPNName: %s", ti.MustGetCHSName(), ti.MustGetENGName(), ti.MustGetJPNName())
			// for torrents not exist
			if !p.Webui.Completed(ti.InfoHash) {
				log.Println("start download torrent from bangumi")
				if err := p.bgmTrMgr.CanDownloadFromBangumi(ti); err != nil {
					log.Println("filter failure reason", err)
					return
				}
				torrentURL, err := ti.GetTorrentDownloadURL()
				if err != nil {
					log.Println("filter failure reason", err)
					return
				}
				_, err = crawler.DownloadBangumiTorrentToFile(
					torrentURL,
					p.Config.TorrentPath,
					ti.InfoHash,
					p.bgmDownloader,
					p.Webui)
				if err != nil {
					log.Println("can not download bangumi torrent", err)
				}
				log.Println("downloaded torrent: ", ti.MustGetCHSName())
				return
			}
			// for completed torrents
			log.Println("prepare to post torrent: ", ti.MustGetCHSName())
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
			err = poster.SetPTGENContent(detail.Detail)
			if err != nil {
				log.Println("failed to SetPTGENContent: ", err)
			}
			mediaInfo, thumb, err := GetMediaInfoFromWEBUI(ti.InfoHash, p.Webui)
			if err != nil {
				log.Println("failed to get media info: ", err)
				return
			}
			poster.SetMediaInfoContent(mediaInfo)
			err = poster.SetTorrentThumb(thumb, "jpg")
			if err != nil {
				log.Println("failed to create neubt media thumb, proceed: ", err)
			}
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
			time.Sleep(time.Second * 60)
		}
	}
}

// GetTorrentPTGenDetail the info will be updated
func (p *Poster) GetTorrentPTGenDetail(info *dao.BangumiTorrentInfo) (*ptgen.BangumiInfoDetail, error) {
	alias := p.ani.GetAliasCHSName(info.Title)
	if alias != "" {
		info.SetCHSName(alias)
	}
	links, err := p.ptgen.GetBangumiLinkByNames(
		info.MustGetCHSName(),
		info.MustGetJPNName(),
		info.MustGetENGName())
	if err != nil {
		return nil, err
	}
	for _, v := range links { // we actually use links[0]
		r, err := p.ptgen.GetBangumiInfoByLink(v.Link)
		if err != nil {
			return nil, err
		}
		err = p.ani.InsertNewCHSName(v.ChnName, "")
		if err != nil {
			log.Println("can not write ani", err)
		}
		// update names
		info.SetReleaseCHSName(v.ChnName)
		if v.JpnName != "" && (info.MustGetJPNName() == "" || alias != "") {
			info.SetJPNName(v.JpnName)
		}
		d, err := ptgen.GetDetailFromInfo(r)
		if err != nil {
			return nil, err
		}
		if d.JpnName != "" && (info.MustGetJPNName() == "" || alias != "") {
			info.SetJPNName(d.JpnName)
		}
		if d.EngName != "" && (info.MustGetENGName() == "" || alias != "") {
			info.SetENGName(d.EngName)
		}
		return d, nil
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
	// set poster type first
	func() {
		if info.ContainsFinishedSeasons() {
			poster.SetTidByName("完结动画")
			return
		}
		if info.ContainsMovie() {
			poster.SetTidByName("剧场OVA")
			return
		}
		poster.SetTidByName("连载动画")
		return
	}()
	// set poster title name and comment name
	poster.SetPostFileName(info.GetTorrentName())
	if info.MustGetCHSName() == "" {
		return errors.New("no Chinese name or English name found in info")
	}
	if info.MustGetTeam() == nil {
		return errors.New("no team name found in info")
	}
	poster.SetTitle(
		info.MustGetCHSName(),
		info.MustGetENGName(),
		info.MustGetJPNName(),
		info.MustGetEpisode(),
		info.MustGetFormat(),
		info.MustGetTeamStr(),
		info.MustGetLanguage(),
		info.MustGetResolution(),
	)
	poster.SetCommentContent(
		"[code]",
		"Debug info:",
		"原种标题："+info.Title,
		"种子信息："+info.GetDetail(),
		"种子内容："+info.GetTorrentContent(),
		"[/code]",
	)
	return nil
}

// GetMediaInfoFromWEBUI return the mediaInfo, thumb, and error
func GetMediaInfoFromWEBUI(infoHash string, webui qbt.WEBUIHelper) (string, []byte, error) {
	// generate media info
	log.Println("generating media info")
	torrent, files, err := webui.GetTorrentDetail(infoHash)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to get torrent files")
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
		return "", nil, errors.New("can not find valid video file in torrent")
	}
	info, err := util.GetMediaInfo("./lib", savePath, fileName)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to generate media info of file"+fileName)
	}
	log.Println("generating media thumb")
	data, err := util.GetMediaImage("./lib", savePath, fileName)
	if err != nil {
		log.Println("WARNING: can not get thumb, ", err)
		return info, nil, nil
	}
	return info, data, nil
}

func IsVideoFile(name string) bool {
	result := regexp.MustCompile(`(.mp4|.ts|.mkv)`).FindAllString(name, -1)
	if result == nil {
		return false
	}
	return true
}
