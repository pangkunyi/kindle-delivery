package main

import (
	"zhihu"
)

func main() {
	println("testing...")
	err := zhihu.UpdateZhihu()
	if err != nil {
		panic(err)
	}
	println("done.")
}
