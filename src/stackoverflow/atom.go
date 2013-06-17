/* vim: set ts=4 sw=2 enc=utf-8: */
package stackoverflow

import(
	feeder "github.com/jteeuwen/go-pkg-rss"
	"strings"
	"io/ioutil"
	"os"
	"fmt"
)

var items []*feeder.Item
var url ="http://stackexchange.com/feeds/tagsets/88069/golang?sort=active"
var feed =feeder.New(0, false, chanHandler, itemHandler)
var idPrefixLen=len("http://stackoverflow.com/questions/")
var template=`
<html>
<head>
	<title>stackoverflow</title>
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
var questionTemplate=`<H1>%s</H1>%s<br/><H2>Answer</H2>%s<br/>`
func chanHandler(feed *feeder.Feed, newchannels []*feeder.Channel) {
	//println(len(newchannels), "new channel(s) in", feed.Url)
}

func itemHandler(feed *feeder.Feed, ch *feeder.Channel, newitems []*feeder.Item) {
	items = newitems
	println(len(newitems), "new item(s) in", ch.Title, "of", feed.Url)
	var content string
	for _,item := range items{
		id := getItemQuestionId(item)
		var qs Questions
		question(&qs, id)
		if qs.Items[0].Accepted_answer_id > 0 {
		  var ans Answers
		  answer(&ans, qs.Items[0].Accepted_answer_id)
		  content=content+fmt.Sprintf(questionTemplate, qs.Items[0].Title, qs.Items[0].Body, ans.Items[0].Body)
		  //fmt.Printf("title:%s\nbody:%s\nanswer:%s\n----------------------------", qs.Items[0].Title, qs.Items[0].Body, ans.Items[0].Body)
		}
	}
	write(content)
}

func write(content string){
	err := ioutil.WriteFile(os.Getenv("HOME")+"/stackoverflow.html", []byte(fmt.Sprintf(template,content)), os.ModePerm)
	if err != nil {
		fmt.Printf("fail to write file, cause by: %v\n", err)
		return
	}
}

func getItemQuestionId(item *feeder.Item) string{
	tmp := item.Id[idPrefixLen:]
	return tmp[:strings.Index(tmp,"/")] 
}

func Update() error{
	err := feed.Fetch(url, nil)
	return err
}

