package main

import(
	sof "stackoverflow"
	"mail"
	"os"
	"fmt"
	"os/exec"
	"io/ioutil"
	"bytes"
	"net/http"
)
func main(){
	http.HandleFunc("/send2kindle", send2kindleHandler)
	err := http.ListenAndServe(":19999", nil)
	if err != nil{
		panic(err)
	}
}
func send2kindleHandler(w http.ResponseWriter, r *http.Request){
	send2kindle()
	fmt.Fprintf(w,"success")
}
func send2kindle(){
	filename :="/tmp/stackoverflow.html"
	err:= sof.Update()
	if err!=nil{
		panic(err)
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
	    fmt.Printf("no stackoverflow html file create: %s", filename)
	    panic(err)
	}
	kindlegen()
	filename ="/tmp/stackoverflow.mobi"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
	    fmt.Printf("no stackoverflow mobi file create: %s", filename)
	    panic(err)
	}
	
	body, err :=ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	err = mail.Send(body, "stackoverflow", ".mobi")
	if err != nil {
		panic(err)
	}
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

