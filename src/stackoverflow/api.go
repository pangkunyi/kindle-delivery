/* vim: set ts=4 sw=2 enc=utf-8: */
package stackoverflow

import(
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

const (
	questions_api_url=`https://api.stackexchange.com/2.1/questions/%s?order=desc&sort=activity&site=stackoverflow&filter=!T5TVQmNC(GmExNmJa6`
	answer_api_url=`https://api.stackexchange.com/2.1/answers/%v?order=desc&sort=activity&site=stackoverflow&filter=!1zKuht46WZuyLgV0nAfmB`
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

type Answers struct {
	Items []Answer
}

type Answer struct {
	Question_id int64
	Body string
	Answer_id int64
}

func question(q *Questions, id string) (error){
	return jsonApi(q, fmt.Sprintf(questions_api_url, id))
}

func answer(a *Answers, id int64) (error){
	return jsonApi(a, fmt.Sprintf(answer_api_url, id))
}

func jsonApi(v interface{}, apiUrl string) (error){
	content, err := httpRequest(apiUrl)
	if err!=nil{
		return err
	}

	err =json.Unmarshal(content, v)
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
