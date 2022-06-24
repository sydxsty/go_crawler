package crawler

import (
	"bytes"
	"crawler/neubt"
	"crawler/util"
	"crawler/util/html"
	"github.com/gocolly/colly/v2"
	"github.com/gogf/gf/v2/encoding/gcharset"
	"github.com/pkg/errors"
	"log"
	"mime/multipart"
	"strings"
	"time"
)

// TorrentPoster post torrent
type TorrentPoster interface {
	// PostTorrentMultiPart if success, return the url of the new torrent
	PostTorrentMultiPart(data []byte) (string, error)
	// SetTitle thread title
	SetTitle(pieces ...string)
	// SetTid torrent kind
	SetTid(tid string) bool
	SetTidByName(name string) bool
	// SetPostFileName the torrent attach name
	SetPostFileName(name string)
	// SetPTGENContent set content
	SetPTGENContent(text string) error
	// SetTorrentThumb add a thumb for torrent
	SetTorrentThumb(image []byte, suffix string) error
	SetMediaInfoContent(text string)
	SetCommentContent(texts ...string)
}

type TorrentPosterImpl struct {
	client neubt.Client

	tidList      map[string]string
	aidList      []string
	thumbAidList []string

	formHash     string
	postTime     string
	wysiwyg      string
	special      string
	specialExtra string
	subject      string // title of the torrent
	tid          string // type id

	fieldID        string // field id, set in NewTorrentPoster
	genTxt         string // main message, from pt-gen
	mediaInfoText  string // media info
	thumbImageText string // aids for movie thumbs
	comment        string // header message, torrent detail
	postFileName   string // the torrent name showed in thread

	imgUploader ImageUploader
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
	if name == "" {
		name = "default_torrent_name"
	}
	t.postFileName = name
}

func (t *TorrentPosterImpl) SetPTGENContent(text string) error {
	t.genTxt = text
	for _, imgURL := range GetAllImgFromPTGen(text) {
		imageDownloader, err := util.NewImageDownloader(imgURL)
		if err != nil {
			return errors.Wrap(err, "can not init image downloader")
		}
		data, fileType, err := imageDownloader.Download()
		if err != nil {
			return errors.Wrap(err, "can not init download image")
		}
		aid, err := t.imgUploader.UploadImage(data, "poster", fileType)
		if err != nil {
			return errors.Wrap(err, "can not upload poster to neubt")
		}
		t.aidList = append(t.aidList, aid)
		log.Println("uploaded poster to neubt")
		time.Sleep(time.Second * 5)
	}
	replaced, err := ReplaceImgWithTagID(text, t.aidList)
	if err != nil {
		for _, aid := range t.aidList {
			_ = t.imgUploader.RemoveImage(aid)
		}
		return errors.Wrap(err, "can not replace the original txt")
	}
	t.genTxt = replaced
	return nil
}

func (t *TorrentPosterImpl) SetTorrentThumb(image []byte, suffix string) error {
	aid, err := t.imgUploader.UploadImage(image, "thumb", suffix)
	if err != nil {
		return errors.Wrap(err, "can not upload TorrentThumb to neubt")
	}
	t.thumbAidList = append(t.thumbAidList, aid)
	log.Println("uploaded TorrentThumb to neubt")
	time.Sleep(time.Second * 5)
	t.thumbImageText += GetAIDText(aid) + "\n"
	return nil
}

func (t *TorrentPosterImpl) SetMediaInfoContent(text string) {
	t.mediaInfoText = "[code] " + text + " [/code]"
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
	// build the image loader
	uid, _ := node.GetInnerString(`.//input[@name="uid"]/@value`)
	hash, _ := node.GetInnerString(`.//input[@name="hash"]/@value`)
	t.imgUploader = NewImageUploader(client, uid, hash)

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
	return note + t.comment + "\n" + t.genTxt + "\n" + t.mediaInfoText + "\n" + t.thumbImageText + "\n"
}

func (t *TorrentPosterImpl) PostTorrentMultiPart(data []byte) (string, error) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)

	_ = w.WriteField(UTF82GB2312("formhash"), UTF82GB2312(t.formHash))
	_ = w.WriteField(UTF82GB2312("posttime"), UTF82GB2312(t.postTime))
	_ = w.WriteField(UTF82GB2312("wysiwyg"), UTF82GB2312(t.wysiwyg))
	_ = w.WriteField(UTF82GB2312("special"), UTF82GB2312(t.special))
	_ = w.WriteField(UTF82GB2312("specialextra"), UTF82GB2312(t.specialExtra))
	_ = w.WriteField(UTF82GB2312("typeid"), UTF82GB2312(t.tid))
	_ = w.WriteField(UTF82GB2312("subject"), UTF82GB2312(t.subject))
	_ = w.WriteField(UTF82GB2312("message"), UTF82GB2312(t.getMessageBody()))
	_ = w.WriteField(UTF82GB2312("readperm"), UTF82GB2312(""))
	_ = w.WriteField(UTF82GB2312("tags"), UTF82GB2312(""))
	_ = w.WriteField(UTF82GB2312("allownoticeauthor"), UTF82GB2312("1"))
	_ = w.WriteField(UTF82GB2312("usesig"), UTF82GB2312("1"))
	_ = w.WriteField(UTF82GB2312("save"), UTF82GB2312(""))
	for _, aid := range append(t.aidList, t.thumbAidList...) {
		_ = w.WriteField(UTF82GB2312("attachupdate["+aid+"]"), UTF82GB2312(""))
		_ = w.WriteField(UTF82GB2312("attachnew["+aid+"][description]"), UTF82GB2312(""))
		_ = w.WriteField(UTF82GB2312("attachnew["+aid+"][readperm]"), UTF82GB2312(""))
		_ = w.WriteField(UTF82GB2312("attachnew["+aid+"][price]"), UTF82GB2312("0"))
	}
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
	var covert string
	for _, v := range s {
		res, err := gcharset.UTF8To("GBK", string(v))
		if err != nil {
			log.Println(`error when transform word`, string(v))
			continue
		}
		covert += res
	}
	return covert
}
