package zhihu

import (
	"atom"
	"db"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	username     = "username"
	password     = "password"
	loginPage    = "http://www.zhihu.com/login?email=" + username + "&password=" + password
	homePage     = "http://www.zhihu.com/"
	temp_file    = "/tmp/zhihu.html"
	temp_article = `<h1 id="artical_topic">%s</h1><div class="artical">%s</div>`
	template     = `
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta charset="utf-8" />
<title>知乎-%s</title>
</head>
<body>
%s
</body>
</html>
`
)

type Article struct {
	Url     string
	Title   string
	Content string
}

func login() (cookie string, err error) {
	client := http.DefaultTransport
	req, err := http.NewRequest("POST", loginPage, nil)
	if err != nil {
		return
	}
	setHeader(req)
	resp, err := client.RoundTrip(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	cookies := resp.Header["Set-Cookie"]
	if cookies == nil {
		err = errors.New("login failed")
		return
	}
	for _, c := range cookies {
		cookie += c[:strings.Index(c, ";")+2]
	}
	return
}

func setHeader(req *http.Request) {
	req.Header.Add("User-Agent", "curl/7.19.7 (i386-redhat-linux-gnu) libcurl/7.19.7 NSS/3.14.0.0 zlib/1.2.3 libidn/1.18 libssh2/1.4.2")
	req.Header.Add("Host", "www.zhihu.com")
	req.Header.Add("Accept", "*/*")
}

func getUrlContent(cookie, url string) (content string, err error) {
	client := http.DefaultTransport
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	setHeader(req)
	req.Header.Add("Cookie", cookie)
	resp, err := client.RoundTrip(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	content = string(body)
	return
}

func load2Feed(cookie string) (feed *atom.Feed, err error) {
	content, err := getUrlContent(cookie, homePage)
	if err != nil {
		return
	}
	feed = &atom.Feed{Id: homePage, Title: "知乎", Entries: make([]atom.Entry, 0)}
	for {
		var url, title string
		var ok bool
		if url, content, ok = extract(content, `question_link" target="_blank" href="`, `"`); !ok {
			break
		}
		url = "http://www.zhihu.com" + url
		fmt.Println("got url:", url)
		if title, content, ok = extract(content, `>`, `</a>`); !ok {
			break
		}
		fmt.Println("got title:", title)
		hasEntry, err := db.HasEntry(feed.Id, url)
		if err != nil { //database query error
			return nil, err
		}
		if hasEntry { //found entry already in database
			continue
		}

		article := &Article{Url: url}
		err = article.parseArticle(cookie)
		if err != nil {
			return nil, err
		}
		entry := atom.Entry{Id: url, Title: title, Summary: article.Content}
		feed.Entries = append(feed.Entries, entry)
		err = db.SaveEntry(feed.Title, feed.Id, &entry)
		if err != nil {
			return nil, err
		}
	}
	return feed, nil
}

func UpdateZhihu() error {
	cookie, err := login()
	if err != nil {
		return err
	}
	feed, err := load2Feed(cookie)
	if err != nil {
		return err
	}
	if len(feed.Entries) < 1 {
		return errors.New("no new Article found")
	}
	var content string
	for _, entry := range feed.Entries {
		content = fmt.Sprintf(temp_article, entry.Title, entry.Summary) + content
	}
	return write(content)
}

func write(content string) error {
	return ioutil.WriteFile(temp_file, []byte(fmt.Sprintf(template, time.Now().Format("2006-01-02"), content)), os.ModePerm)
}

func (this *Article) parseArticle(cookie string) error {
	content, err := getUrlContent(cookie, this.Url)
	if err != nil {
		return err
	}
	fmt.Println("got content:", content)
	var ok bool
	var detail string
	if detail, content, ok = extract(content, `/question/detail">`, `</div>
<div class="zm-item-meta`); !ok {
		return errors.New("error article when parse question detail")
	}
	if strings.Contains(detail, `zm-editable-tip`) {
		detail = ""
	} else {
		detail = fmt.Sprintf("<h3>补充:%s</h3>", detail)
	}
	var people string
	if people, content, ok = extract(content, `" href="/people/`, `>`); !ok {
		return errors.New("error article when parse people 1")
	}
	if people, content, ok = extract(content, ``, `</a>`); !ok {
		return errors.New("error article when parse people 2")
	}
	this.Content = fmt.Sprintf("%s<h4>%s</h4>", detail, people)
	var answer string
	if answer, content, ok = extract(content, `/answer/content">`, `</div>
<a class="zg-anchor-hidden`); !ok {
		return errors.New("error article when parse content")
	}
	this.Content += answer
	// one more time
	if people, content, ok = extract(content, `" href="/people/`, `>`); !ok {
		return nil
	}
	if people, content, ok = extract(content, ``, `</a>`); !ok {
		return nil
	}
	this.Content += fmt.Sprintf("<h4>%s</h4>", people)
	if answer, content, ok = extract(content, `/answer/content">`, `</div>
<a class="zg-anchor-hidden`); !ok {
		return nil
	}
	this.Content += answer
	return nil
}

func extract(content, prefix, subfix string) (string, string, bool) {
	idx := strings.Index(content, prefix)
	if idx < 0 {
		return "", "", false
	}
	content = content[idx+len(prefix):]
	idx = strings.Index(content, subfix)
	if idx < 0 {
		return "", "", false
	}
	return content[:idx], content[idx:], true
}
