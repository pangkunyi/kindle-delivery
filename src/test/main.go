package main

import(
	sof "stackoverflow"
)
func main(){
	err:= sof.Update()
	if err!=nil{
		panic(err)
	}
}

