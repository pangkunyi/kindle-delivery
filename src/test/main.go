package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	sof "stackoverflow"
)

func main() {
	var e sof.Entry
	err :=GetEntry(`http://stackexchange.com/feeds/tagsets/88069/golang?sort=active`, `http://stackoverflow.com/questions/17265463/how-do-i-convert-a-database-row-into-a-struct-in-go12`, &e)
	if err!=nil {
		panic(err)
}
	fmt.Println(e)
}
