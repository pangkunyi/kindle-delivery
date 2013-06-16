package stackoverflow

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	questions_api_url=`https://api.stackexchange.com/2.1/questions/%s?order=desc&sort=activity&site=stackoverflow&filter=!T5TVQmNC(GmExNmJa6`
)

type Questions struct {
	Items []Question
}

type Question struct {
	Question_id int64
	Body string
	Title string
	Accepted_answer_id int64
}

func question(q *Questions, id string) (error){
	content, err := httpRequest(fmt.Sprintf(questions_api_url, id))
	if err!=nil{
		return err
	}
	err =json.Unmarshal(content, q)
	if err!=nil{
		return err
	}
	return nil
}

func httpRequest(url string) ([]byte, error){
	resp, err := http.Get(url)
	if err!=nil{
		return nil, err
	}
	body, err :=ioutil.ReadAll(resp.Body)
	if err!=nil{
		return nil, err
	}
	return body, nil
}
