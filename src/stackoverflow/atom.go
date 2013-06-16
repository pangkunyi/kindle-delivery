package stackoverflow

import(
	feeder "github.com/jteeuwen/go-pkg-rss"
	"strings"
	"fmt"
)

var items []*feeder.Item
var url ="http://stackexchange.com/feeds/tagsets/88069/golang?sort=active"
var feed =feeder.New(0, false, chanHandler, itemHandler)
var idPrefixLen=len("http://stackoverflow.com/questions/")
func chanHandler(feed *feeder.Feed, newchannels []*feeder.Channel) {
	//println(len(newchannels), "new channel(s) in", feed.Url)
}

func itemHandler(feed *feeder.Feed, ch *feeder.Channel, newitems []*feeder.Item) {
	items = newitems

	println(len(newitems), "new item(s) in", ch.Title, "of", feed.Url)
	for _,item := range items{
		id := getItemQuestionId(item)
		var qs Questions
		question(&qs, id)
		fmt.Printf("title:%s\nbody:%s\n", qs.Items[0].Title, qs.Items[0].Body)
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

