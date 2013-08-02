package main

import(
	"mail"
	"os"
	"fmt"
	"os/exec"
	"io/ioutil"
	"bytes"
	"net/http"
//	sof "stackoverflow"
	"ifeng"
)
func main(){
	http.HandleFunc("/send2kindle", send2kindleHandler)
	err := http.ListenAndServe(":19999", nil)
	if err != nil{
		panic(err)
	}
}
func send2kindleHandler(w http.ResponseWriter, r *http.Request){
//	err := send2kindle("/tmp/", "stackoverflow.html", "stackoverflow.mobi", func()error{return sof.Update()})
//	if err != nil {
//		fmt.Printf("error:%s\n", err.Error())
//		fmt.Fprintf(w,"error:%s\n", err.Error())
//	}else{
//		fmt.Fprintf(w,"success\n")
//	}
	err := send2kindle("/tmp/", "ifeng-kaijuanbafenzhong.html", "ifeng-kaijuanbafenzhong.mobi", func()error{return ifeng.UpdateKJBFZ()})
	if err != nil {
		fmt.Printf("error:%s\n", err.Error())
		fmt.Fprintf(w,"error:%s\n", err.Error())
	}else{
		fmt.Fprintf(w,"success\n")
	}
}
func send2kindle(dir, html, mobi string, update func() error) error{
	filename :=dir+html
	os.Remove(filename)
	err:= update()
	if err!=nil{
		return err
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no stackoverflow html file create: %s", filename)
		return err
	}
	filename =dir+mobi
	os.Remove(filename)
	kindlegen(dir, html, mobi)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no stackoverflow mobi file create: %s", filename)
		return err
	}

	body, err :=ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = mail.Send(body, mobi)
	if err != nil {
		return err
	}
	return nil
}

func kindlegen(dir, html, mobi string){
	cmd := exec.Command("kindlegen", dir+html, "-o", mobi)
	var in bytes.Buffer
	cmd.Stdin = &in
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	fmt.Println(string(out.Bytes()))
}

