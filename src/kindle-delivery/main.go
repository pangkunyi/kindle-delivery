package main

import (
	"bytes"
	"fmt"
	"ifeng"
	"io/ioutil"
	"mail"
	"net/http"
	"os"
	"os/exec"
	htmlParser "parser"
	sof "stackoverflow"
	"zhihu"
)

func main() {
	http.HandleFunc("/send2kindle", send2kindleHandler)
	err := http.ListenAndServe(":19999", nil)
	if err != nil {
		panic(err)
	}
}
func send2kindleHandler(w http.ResponseWriter, r *http.Request) {
	err := send2kindle("/tmp/", "stackoverflow.html", "stackoverflow.mobi", "", func() error { return sof.Update() })
	if err != nil {
		fmt.Printf("stackoverflow serror:%s\n", err.Error())
		fmt.Fprintf(w, "stackoverflow serror:%s\n", err.Error())
	} else {
		fmt.Fprintf(w, "stackoverflow success\n")
	}
	err = send2kindle("/tmp/", "ifeng-kaijuanbafenzhong.html", "ifeng-kaijuanbafenzhong.mobi", "zh", func() error { return ifeng.UpdateKJBFZ() })
	if err != nil {
		fmt.Printf("ifeng-kaijuanbafenzhong error:%s\n", err.Error())
		fmt.Fprintf(w, "ifeng-kaijuanbafenzhong error:%s\n", err.Error())
	} else {
		fmt.Fprintf(w, "ifeng-kaijuanbafenzhong success\n")
	}
	err = send2kindle("/tmp/", "zhihu.html", "zhihu.mobi", "zh", func() error { return zhihu.UpdateZhihu() })
	if err != nil {
		fmt.Printf("zhihu error:%s\n", err.Error())
		fmt.Fprintf(w, "zhihu error:%s\n", err.Error())
	} else {
		fmt.Fprintf(w, "zhihu success\n")
	}
}

func send2kindle(dir, html, mobi, locale string, update func() error) error {
	filename := dir + html
	os.Remove(filename)
	err := update()
	if err != nil {
		return err
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no html file create: %s", filename)
		return err
	}
	err = htmlParser.EncodeImg(filename)
	if err != nil {
		return err
	}
	filename = dir + mobi
	os.Remove(filename)
	kindlegen(dir, html, mobi, locale)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no mobi file create: %s", filename)
		return err
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = mail.Send(body, mobi)
	if err != nil {
		return err
	}
	return nil
}

func kindlegen(dir, html, mobi, locale string) {
	var args []string
	if locale == "" {
		args = []string{dir + html, "-o", mobi}
	} else {
		args = []string{dir + html, "-o", mobi, "-locale", locale}
	}
	fmt.Printf("kindlegen %s", args)
	cmd := exec.Command("kindlegen", args...)
	var in bytes.Buffer
	cmd.Stdin = &in
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	fmt.Println(string(out.Bytes()))
}
