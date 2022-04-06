package controller

import (
	"goCrawler/dao"
	"goCrawler/module"
	"log"
	"time"
)

func GetListByForumName(nameList ...string) []*dao.TorrentInfo {
	var infoList []*dao.TorrentInfo
	c := module.NewIndexModule()
	forumMap := c.GetForumList()
	log.Println("all available forums:")
	log.Println(forumMap)
	for _, name := range nameList {
		if url, ok := forumMap[name]; ok {
			infoList = append(infoList, c.GetForum(url)...)
		}
	}
	return infoList
}

func CrawlBangumiInfo() {
	animeList := GetLatestAnimeList()
	d := module.NewDownloader()
	for _, anime := range animeList {
		// check if anime has already posted
		// if it is not been posted yet, check if it has been downloaded yet
		if GetPostState(anime) {
			// if posted, continue
			continue
		}
		if d.Contains(anime.InfoHash) && !d.Completed(anime.InfoHash) {
			// if downloading, continue
			continue
		}
		// if the torrent can be posted
		if ptGenDetail := GetTorrentPTGenDetail(anime); ptGenDetail != nil {
			// if it has been downloaded, post it
			if d.Completed(anime.InfoHash) {
				m := module.NewForumModule("44", anime.InfoHash+".torrent")
				if err := m.UpdateWithTorrentInfo(anime); err != nil {
					log.Println(err)
					continue
				}
				if err := m.SetText(ptGenDetail["format"].(string)); err != nil {
					log.Println(err)
					continue
				}
				SetPostedState(anime)
				if response, err := m.PostMultiPart(); err != nil {
					log.Println(err)
					continue
				} else {
					info := &dao.TorrentInfo{
						Link: response.Request.URL.Path,
					}
					downloadTorrentByInfo(info)
					log.Println("post torrent success")
				}
			} else {
				// else, download it from bgm, and add to bittorrent
				t := module.NewBangumiModule()
				t.DownloadTorrentFromUrl(dao.YAMLConfig.TorrentPath, anime.InfoHash+".torrent", anime.Detail.TorrentDownloadURL)
				if err := d.AddTorrentToBitTorrent(dao.YAMLConfig.TorrentPath, anime.InfoHash+".torrent"); err != nil {
					log.Fatal(err)
				}
				log.Println("download torrent from bangumi")
			}
			log.Println("success, sleep 5sec to continue")
			time.Sleep(time.Second * 5)
		} else {
			// we can not get torrent detail
			log.Println("error get torrent detail: " + anime.Title)
			time.Sleep(time.Second * 2)
		}
	}
	log.Println("all torrent finished scanning, return")
	return
}

func Login() error {
	indexModule := module.NewIndexModule()
	if dao.YAMLConfig.UseCookie {
		if err := indexModule.LoadCookie(); err != nil {
			return err
		}
	} else {
		indexModule.Login(dao.YAMLConfig.Username, dao.YAMLConfig.Password)
		if err := indexModule.SaveCookie(); err != nil {
			return err
		}
	}
	return nil
}

func GetPostState(anime *dao.BangumiTorrentInfo) bool {
	var state map[string]interface{}
	if err := dao.LoadFromDB("anime_post_state_"+anime.InfoHash, &state); err != nil {
		return false
	}
	return true
}

func SetPostedState(anime *dao.BangumiTorrentInfo) {
	state := make(map[string]interface{})
	state["posted"] = "true"
	dao.SaveToDB("anime_post_state_"+anime.InfoHash, state)
}

// we actually just need url here
func downloadTorrentByInfo(info *dao.TorrentInfo) {
	detailModule := module.NewDetailModule()
	d := module.NewDownloader()
	form := detailModule.GetDetailFrom(info)
	for _, floor := range form.Floors {
		if floor.Comment.TorrentLink != "" {
			if err := detailModule.DownloadTorrentFromNestedURL(
				floor.Comment.TorrentLink,
				floor.Comment.TorrentName,
				dao.YAMLConfig.TorrentPath); err != nil {
				log.Println(err)
				continue
			}
			if err := d.AddTorrentToBitTorrent(
				dao.YAMLConfig.TorrentPath,
				floor.Comment.TorrentName); err != nil {
				continue
			}
			info.Crawled = true
			log.Println("downloading: ")
			log.Println(info)
			if err := info.SaveToDB(); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

// DownloadTorrentByInfoList apply filter here
func DownloadTorrentByInfoList(infoList []*dao.TorrentInfo) {
	for _, info := range infoList {
		if info.Crawled || info.TorrentID < dao.YAMLConfig.ThreadWaterMark || info.Discount < dao.YAMLConfig.DiscountWaterMark {
			continue
		}
		downloadTorrentByInfo(info)
	}
}
