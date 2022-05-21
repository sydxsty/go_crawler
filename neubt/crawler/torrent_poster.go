package crawler

import (
	"bytes"
	"crawler/neubt"
	"crawler/neubt/html"
	"errors"
	"github.com/gocolly/colly/v2"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"log"
	"mime/multipart"
	"strings"
)

type TorrentPoster interface {
	// PostTorrentMultiPart if success, return the url of the new torrent
	PostTorrentMultiPart(data []byte) (string, error)
	// SetTitle forum title
	SetTitle(pieces ...string)
	// SetTid torrent kind
	SetTid(tid string) bool
	SetTidByName(name string) bool
	// SetPostFileName the torrent attach name
	SetPostFileName(name string)
	// SetPTGENContent set content
	SetPTGENContent(text string)
	SetCommentContent(texts ...string)
}

type TorrentPosterImpl struct {
	client neubt.Client

	tidList map[string]string

	formHash     string
	postTime     string
	wysiwyg      string
	special      string
	specialExtra string
	subject      string // title of the torrent
	tid          string // type id

	fieldID      string // field id, set in NewTorrentPoster
	genTxt       string // main message, from pt-gen
	comment      string // header message, torrent detail
	postFileName string // the torrent name showed in forum
}

func (t *TorrentPosterImpl) SetTid(tid string) bool {
	// set tid
	if _, ok := t.tidList["选择主题分类"]; !ok {
		log.Println("no matching default tid for uploading torrent")
		return false
	} else {
		for _, val := range t.tidList {
			if val == tid {
				t.tid = tid
				return true
			}
		}
	}
	return false
}

func (t *TorrentPosterImpl) SetTidByName(name string) bool {
	if value, ok := t.tidList[name]; !ok {
		log.Println("no matching tid for SetTidByName")
		return false
	} else {
		t.tid = value
		return true
	}
}

func (t *TorrentPosterImpl) SetPostFileName(name string) {
	t.postFileName = name
}

func (t *TorrentPosterImpl) SetPTGENContent(text string) {
	t.genTxt = text
}

func (t *TorrentPosterImpl) SetCommentContent(texts ...string) {
	t.comment = ""
	for _, text := range texts {
		t.comment += text + "\n"
	}
}

func NewTorrentPoster(fid string, client neubt.Client) (TorrentPoster, error) {
	t := &TorrentPosterImpl{
		client: client.Clone(),
	}
	t.fieldID = fid
	url := `forum.php?mod=post&action=newthread&fid=` + t.fieldID + `&specialextra=torrent`
	resp, err := t.client.SyncVisit(url)
	if err != nil {
		return nil, err
	}
	node, err := html.NewNodeFromBytes(resp.Body)
	if err != nil {
		return nil, err
	}

	t.formHash, _ = node.GetInnerString(`.//input[@id="formhash"]/@value`)
	t.postTime, _ = node.GetInnerString(`.//input[@id="posttime"]/@value`)
	t.wysiwyg, _ = node.GetInnerString(`.//input[@name="wysiwyg"]/@value`)
	t.special, _ = node.GetInnerString(`.//input[@name="special"]/@value`)
	t.specialExtra, _ = node.GetInnerString(`.//input[@name="specialextra"]/@value`)
	t.subject, _ = node.GetInnerString(`.//input[@name="subject"]/@value`)
	tidNodeList, err := node.GetInnerNodeList(`.//select[@name="typeid"]/option`)
	if err != nil {
		return nil, err
	}
	t.tidList = make(map[string]string)
	for _, tidNode := range tidNodeList {
		t.tidList[tidNode.GetString()], _ = tidNode.GetInnerString(`./@value`)
	}
	return t, nil
}

func (t *TorrentPosterImpl) SetTitle(pieces ...string) {
	var title string
	for _, piece := range pieces {
		if len(piece) == 0 {
			continue
		}
		title += "[" + piece + "]"
	}
	t.subject = title
}

func (t *TorrentPosterImpl) getMessageBody() string {
	note := `[quote]自动发种试运行，有问题请在github上提Issue` + "\n" + `[url=https://github.com/sydxsty/go_crawler/releases]https://github.com/sydxsty/go_crawler/releases[/url][/quote]`
	return note + t.comment + "\n" + t.genTxt + "\n"
}

func (t *TorrentPosterImpl) PostTorrentMultiPart(data []byte) (string, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	w.WriteField(UTF82GB2312("formhash"), UTF82GB2312(t.formHash))
	w.WriteField(UTF82GB2312("posttime"), UTF82GB2312(t.postTime))
	w.WriteField(UTF82GB2312("wysiwyg"), UTF82GB2312(t.wysiwyg))
	w.WriteField(UTF82GB2312("special"), UTF82GB2312(t.special))
	w.WriteField(UTF82GB2312("specialextra"), UTF82GB2312(t.specialExtra))
	w.WriteField(UTF82GB2312("typeid"), UTF82GB2312(t.tid))
	w.WriteField(UTF82GB2312("subject"), UTF82GB2312(t.subject))
	w.WriteField(UTF82GB2312("message"), UTF82GB2312(t.getMessageBody()))
	w.WriteField(UTF82GB2312("readperm"), UTF82GB2312(""))
	w.WriteField(UTF82GB2312("tags"), UTF82GB2312(""))
	w.WriteField(UTF82GB2312("allownoticeauthor"), UTF82GB2312("1"))
	w.WriteField(UTF82GB2312("usesig"), UTF82GB2312("1"))
	w.WriteField(UTF82GB2312("save"), UTF82GB2312(""))

	// data, _ := ioutil.ReadFile(path + file)
	pa, _ := w.CreateFormFile(UTF82GB2312("torrent"), UTF82GB2312(t.postFileName+".torrent"))
	if _, err := pa.Write(data); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	t.client.SetRequestCallback(func(r *colly.Request) {
		r.Headers.Set("Content-Type", w.FormDataContentType())
	})
	resp, err := t.client.SyncPostRaw(
		`forum.php?mod=post&action=newthread&fid=`+t.fieldID+`&extra=&topicsubmit=yes`,
		body.Bytes())
	if err != nil {
		return "", err
	}
	if strings.Contains(resp.Request.URL.Path, "forum.php") {
		return "", errors.New("error publish torrent")
	}
	return resp.Request.URL.Path, nil
}

func UTF82GB2312(s string) string {
	covert, err := gcharset.UTF8To("GBK", s)
	if err != nil {
		log.Fatalln(err)
	}
	return covert
}
