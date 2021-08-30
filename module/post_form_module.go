package module

import (
	"bytes"
	"github.com/gocolly/colly/v2"
	"goCrawler/dao"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"mime/multipart"
)

type ForumModule interface {
	ScraperModule
	PostMultiPart() error
	UpdateWithTorrentInfo(info *dao.BangumiTorrentInfo) error
	SetText(text string) error
}

type forumModuleImpl struct {
	scraperModuleImpl
	tidList map[string]string

	formHash     string
	postTime     string
	wysiwyg      string
	special      string
	specialExtra string
	tid          string // type id
	subject      string

	torrentFileName  string
	torrentFileBytes []byte

	text     string // message
	fieldID  string //field id
	fileName string
}

func NewForumModule(fid string, fileName string) ForumModule {
	c := &forumModuleImpl{}
	c.init()
	collector := c.getClonedCollector()
	collector.OnResponse(func(r *colly.Response) {
		node, err := NewNodeFromBytes(r.Body)
		if err != nil {
			log.Fatal(err)
			return
		}
		c.formHash = node.GetInnerNode(`.//input[@id="formhash"]/@value`).GetString()
		c.postTime = node.GetInnerNode(`.//input[@id="posttime"]/@value`).GetString()
		c.wysiwyg = node.GetInnerNode(`.//input[@name="wysiwyg"]/@value`).GetString()
		c.special = node.GetInnerNode(`.//input[@name="special"]/@value`).GetString()
		c.specialExtra = node.GetInnerNode(`.//input[@name="specialextra"]/@value`).GetString()
		tidNodeList := node.GetInnerNodeList(`.//select[@name="typeid"]/option`)
		c.tidList = make(map[string]string)
		for _, tidNode := range tidNodeList {
			c.tidList[tidNode.GetString()] = tidNode.GetInnerNode(`./@value`).GetString()
		}

		if _, ok := c.tidList["选择主题分类"]; !ok {
			log.Fatal("no matching default tid")
		} else {
			for _, val := range c.tidList {
				if val != "0" {
					c.tid = val
				}
			}
		}
		// default subject name
		c.subject = node.GetInnerNode(`.//input[@name="subject"]/@value`).GetString()

	})
	url := `forum.php?mod=post&action=newthread&fid=` + fid + `&specialextra=torrent`
	if err := collector.Visit(c.getAbsoluteURL(url)); err != nil {
		log.Fatal(err)
	}
	c.fieldID = fid
	c.fileName = fileName
	return c
}

func (f *forumModuleImpl) PostMultiPart() error {
	collector := f.getClonedCollector()
	// we do not clone controller here
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	w.WriteField("formhash", UTF82GB2312(f.formHash))
	w.WriteField("posttime", UTF82GB2312(f.postTime))
	w.WriteField("wysiwyg", UTF82GB2312(f.wysiwyg))
	w.WriteField("special", UTF82GB2312(f.special))
	w.WriteField("specialextra", UTF82GB2312(f.specialExtra))
	w.WriteField("typeid", UTF82GB2312(f.tid))
	w.WriteField("subject", UTF82GB2312(f.subject))
	w.WriteField("message", UTF82GB2312(f.text))
	w.WriteField("readperm", UTF82GB2312(""))
	w.WriteField("tags", UTF82GB2312(""))
	w.WriteField("allownoticeauthor", UTF82GB2312("1"))
	w.WriteField("usesig", UTF82GB2312("1"))
	w.WriteField("save", UTF82GB2312(""))

	fileData, _ := ioutil.ReadFile(dao.YAMLConfig.TorrentPath + f.fileName)
	pa, _ := w.CreateFormFile("torrent", f.fileName)
	if _, err := pa.Write(fileData); err != nil {
		log.Fatal(err)
	}
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", w.FormDataContentType())

	})
	collector.OnResponse(func(r *colly.Response) {
		log.Println(r.Body)
	})
	if err := collector.PostRaw(
		f.getAbsoluteURL(`forum.php?mod=post&action=reply&fid=44&tid=1684795&extra=&replysubmit=yes`),
		//f.getAbsoluteURL(`forum.php?mod=post&action=newthread&fid=` + f.fieldID + `&extra=&topicsubmit=yes`),
		body.Bytes()); err != nil {
		return err
	}
	return nil
}

func (f *forumModuleImpl) UpdateWithTorrentInfo(info *dao.BangumiTorrentInfo) error {
	f.tid = f.tidList["连载动画"]
	if info.Detail.TorrentChsName == "" {
		log.Fatal("chinese name empty")
	}
	title := "[" + info.Detail.TorrentChsName + "]"
	if info.Detail.TorrentEngName != "" {
		title += "[" + info.Detail.TorrentEngName + "]"
	}
	if info.Detail.TorrentJpnName != "" {
		title += "[" + info.Detail.TorrentJpnName + "]"
	}
	if info.Detail.Episode != "" {
		title += "[" + info.Detail.Episode + "]"
	}
	if info.Detail.Format != "" {
		title += "[" + info.Detail.Format + "]"
	}
	if info.Detail.TeamName == "" {
		log.Fatal("team name empty")
	}
	title += "[" + info.Detail.TeamName + "]"

	if info.Detail.Language != "" {
		title += "[" + info.Detail.Language + "]"
	}
	if info.Detail.Resolution != "" {
		title += "[" + info.Detail.Resolution + "]"
	}

	f.subject = title
	return nil
}

func (f *forumModuleImpl) SetText(text string) error {
	f.text = text
	return nil
}

func UTF82GB2312(s string) string {
	reader := transform.NewReader(bytes.NewReader([]byte(s)), simplifiedchinese.GB18030.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return ""
	}
	return string(d)
}
