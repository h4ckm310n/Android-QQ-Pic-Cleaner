package main

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

const DB_PATH = "/data/data/com.tencent.mobileqq/databases/"

var qqDB *sql.DB
var slowtableDB *sql.DB

func connectDB(qq string) {
	qqDBPath := filepath.Join(DB_PATH, qq+".db")
	slowtableDBPath := filepath.Join(DB_PATH, "slowtable_"+qq+".db")

	_, err := os.Stat(qqDBPath)
	if !os.IsNotExist(err) {
		qqDB, err = sql.Open("sqlite", qqDBPath)
		if err != nil {
			log.Println("Failed to connect to database: " + qqDBPath)
			log.Println(err)
			qqDB = nil
		}
	}

	_, err = os.Stat(slowtableDBPath)
	if !os.IsNotExist(err) {
		slowtableDB, err = sql.Open("sqlite", slowtableDBPath)
		if err != nil {
			log.Println("Failed to connect to database: " + slowtableDBPath)
			log.Println(err)
			slowtableDB = nil
		}
	}

	if qqDB == nil && slowtableDB == nil {
		log.Fatalln("No database")
	}
}

func closeDB() {
	if qqDB != nil {
		qqDB.Close()
	}
	if slowtableDB != nil {
		slowtableDB.Close()
	}
}
