package main

import (
	"parser"
)

const (
	file = "/tmp/zhihu.html"
)

func main() {
	println("testing...")
	err := html.EncodeImg(file)
	if err != nil {
		panic(err)
	}
	println("done.")
}
