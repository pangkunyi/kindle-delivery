package main

import(
	"os/exec"
	"bytes"
	"fmt"
)
func main(){
	err :=kindlegen()
	if err!=nil{
		panic(err)
	}
}

func kindlegen() error{
	cmd := exec.Command("kindlegen", "/tmp/stackoverflow.html", "-o", "stackoverflow.mobi")
	var in bytes.Buffer
	cmd.Stdin = &in
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	fmt.Println(string(out.Bytes()))
	if err != nil {
		return err
	}
	return nil
}
