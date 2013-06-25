package main

import(
	"mail"
	"os"
	"fmt"
	"os/exec"
	"io/ioutil"
	"bytes"
	"net/http"
	sof "stackoverflow"
)
func main(){
	http.HandleFunc("/send2kindle", send2kindleHandler)
	err := http.ListenAndServe(":19999", nil)
	if err != nil{
		panic(err)
	}
}
func send2kindleHandler(w http.ResponseWriter, r *http.Request){
	err := send2kindle()
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
		fmt.Fprintf(w,"error:%s\n", err.Error())
	}else{
		fmt.Fprintf(w,"success\n")
	}
}
func send2kindle() error{
	filename :="/tmp/stackoverflow.html"
	os.Remove(filename)
	err:= sof.Update()
	if err!=nil{
		return err
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no stackoverflow html file create: %s", filename)
		return err
	}
	filename ="/tmp/stackoverflow.mobi"
	os.Remove(filename)
	kindlegen()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no stackoverflow mobi file create: %s", filename)
		return err
	}

	body, err :=ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = mail.Send(body, "stackoverflow.mobi")
	if err != nil {
		return err
	}
	return nil
}

func kindlegen(){
	cmd := exec.Command("kindlegen", "/tmp/stackoverflow.html", "-o", "stackoverflow.mobi")
	var in bytes.Buffer
	cmd.Stdin = &in
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	fmt.Println(string(out.Bytes()))
}

