package ifeng

import(
	"fmt"
	"errors"
	"net/http"
	"io/ioutil"
	"strings"
	"atom"
	"time"
	"os"
	"db"
)

const(
	kai_juan_ba_fen_zhong_url="http://book.ifeng.com/kaijuanbafenzhong/wendang/list_0/0.shtml"
	temp_file="/tmp/ifeng-kaijuanbafenzhong.html"
	temp_article=`<h1 id="artical_topic">%s</h1><div class="artical">%s</div>`
	template=`
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta charset="utf-8" />
<title>开卷八分钟-%s</title>
</head>
<body>
%s
</body>
</html>
`
)

type Article struct{
	Url string
	Title string
	Content string
}

func load2Feed() (*atom.Feed, error){
	resp, err := http.Get(kai_juan_ba_fen_zhong_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	content := string(body)
	feed := &atom.Feed{Id:kai_juan_ba_fen_zhong_url, Title:"开卷八分钟", Entries:make([]atom.Entry,0)}
	for{
		var url, title string
		var ok bool
		if url, content, ok =extract(content, `sys_url" href="`, `"`); !ok {
		break
		}
		if title, content, ok =extract(content, `>`, `</a>`); !ok {
		break
		}
		hasEntry, err := db.HasEntry(feed.Id, url)
		if err != nil{//database query error
			return nil, err
		}
		if hasEntry{//found entry already in database
			continue
		}

		article := &Article{Url:url}
		err = article.parseArticle()
		if err != nil {
			return nil, err
		}
		entry := atom.Entry{Id:url, Title:title, Summary:article.Content}
		feed.Entries = append(feed.Entries, entry)
		err = db.SaveEntry(feed.Title, feed.Id, &entry)
		if err != nil{
			return nil, err
		}
	}
	return feed, nil
}

func UpdateKJBFZ() error{
	feed, err :=load2Feed()
	if err != nil {
		return err
	}
	if len(feed.Entries)< 1 {
		return errors.New("no new Article found")
	}
	var content string
	for _,entry := range feed.Entries {
		content = fmt.Sprintf(temp_article, entry.Title, entry.Summary) + content
	}
	return write(content)
}

func write(content string) error{
	return ioutil.WriteFile(temp_file, []byte(fmt.Sprintf(template,time.Now().Format("2006-01-02"), content)), os.ModePerm)
}

func (this *Article) parseArticle() error{
	resp, err := http.Get(this.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	content := string(body)
	var ok bool
	if this.Title, content, ok =extract(content, `artical_topic">`, `<`); !ok {
		return errors.New("error article")
	}
	if this.Content, _, ok =extract(content, `<!--mainContent begin-->`, `<!--mainContent end-->`); !ok {
		return errors.New("error article")
	}
	return nil
}

func extract(content, prefix, subfix string) (string, string, bool){
	idx := strings.Index(content, prefix)
	if idx <0 {
		return "", "", false
	}
	content = content[idx+len(prefix):]
	idx = strings.Index(content, subfix)
	if idx <0 {
		return "", "", false
	}
	return content[:idx], content[idx:], true
}
