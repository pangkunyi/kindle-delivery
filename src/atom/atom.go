package atom

import (
	"encoding/xml"
	"net/http"
	"io/ioutil"
)

type Feed struct{
	XMLName xml.Name `xml:"feed"`
	Title string `xml:"title"`
	Id string `xml:"id"`
	Entries []Entry `xml:"entry"`
}

type Entry struct{
	Id string `xml:"id"`
	Title string `xml:"title"`
	Summary string `xml:"summary"`
}

func (this *Feed) Load(url string) error{
	resp, err := http.Get(url)
	if err!=nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err!=nil {
		return err
	}
	err = xml.Unmarshal(body, this)
	return err
}
