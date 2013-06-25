package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"atom"
)

var (
	DB_FILE string = os.Getenv("HOME")+"/stackoverflow.db"
)

func InitDatabase() error{
	os.Remove(DB_FILE)
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		return err
	}
	defer db.Close()

	sqls := []string{`create table feed (
			id integer not null primary key autoincrement, 
			feed_title text,
			feed_id text,
			entry_title text,
			entry_id text,
			entry_summary text)`,
			`create index idx_feed_entry_ids on feed (
			feed_id, entry_id)`}
	for _, sql := range sqls {
		_, err = db.Exec(sql)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveEntry(feedTitle, feedId string, entry *atom.Entry) error{
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(
		`insert into feed(feed_title, feed_id, entry_title, entry_id, entry_summary) values(?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(feedTitle, feedId, entry.Title, entry.Id, entry.Summary)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func GetEntry(feedId, entryId string, entry *atom.Entry) error{
	db, err := sql.Open("sqlite3", DB_FILE)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("select entry_Title, entry_summary from feed where feed_id = ? and entry_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	row := stmt.QueryRow(feedId, entryId)
	entry.Id=entryId
	return row.Scan(&entry.Title, &entry.Summary)
}

func HasEntry(feedId, entryId string) (bool,error){
	var entry atom.Entry
	err := GetEntry(feedId, entryId, &entry)
	if err == nil{//found entry already in database
		return true, nil
	}else if err != sql.ErrNoRows{//database query error
		return false, err
	}
	return false, nil
}
