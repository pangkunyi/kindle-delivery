/* vim: set ts=4 sw=2 enc=utf-8: */
package stackoverflow

import(
	"strings"
	"io/ioutil"
	"time"
	"errors"
	"os"
	"fmt"
	"db"
	"atom"
)

var (
	url ="http://stackexchange.com/feeds/tagsets/88069/golang?sort=active"
	idPrefixLen=len("http://stackoverflow.com/questions/")
	questionTemplate=`<H1>%s</H1>%s<br/><H2>Answer</H2>%s<br/>`
	tempFile="/tmp/stackoverflow.html"
	template=`
<html>
<head>
	<title>stackoverflow-%s</title>
	<style>
		pre{
			 background-color: #eeeeee;
		}
	</style>
</head>
<body>
%s
</body>
</html>
`
)

func Update() error{
	var feed atom.Feed
	err :=feed.Load(url)
	if err!= nil {
		return err
	}

	if len(feed.Entries)<1 {
		return errors.New("no entry found.")
	}

	var content string
	count := 0
	for _, entry := range feed.Entries {
		id := getItemQuestionId(entry.Id)
		var qs Questions
		err = question(&qs, id)
		if err != nil || len(qs.Items) < 1 {
			return errors.New("failure to request stackoverflow questions api")
		}
		if qs.Items[0].Accepted_answer_id > 0 {
			hasEntry, err := db.HasEntry(feed.Id, entry.Id)
			if err != nil{//database query error
				return err
			}
			if hasEntry{//found entry already in database
				continue
			}
			//found no entry
			var ans Answers
			err = answer(&ans, qs.Items[0].Accepted_answer_id)
			if err != nil || len(ans.Items) < 1 {
				return errors.New("failure to request stackoverflow answers api")
			}
			content=content+fmt.Sprintf(questionTemplate, qs.Items[0].Title, qs.Items[0].Body, ans.Items[0].Body)
			count = count + 1
			err = db.SaveEntry(feed.Title, feed.Id, &entry)
			if err != nil{
				return err
			}
		}
	}
	if count < 1 {
		return errors.New("no new feed found")
	}
	return write(content)
}

func write(content string) error{
	return ioutil.WriteFile(tempFile, []byte(fmt.Sprintf(template,time.Now().Format("2006-01-02"), content)), os.ModePerm)
}

func getItemQuestionId(entryId string) string{
	tmp := entryId[idPrefixLen:]
	return tmp[:strings.Index(tmp,"/")] 
}
