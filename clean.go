package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"strings"
)

var dryRun = false

type Cleaner struct {
	friends []string
	groups  []string
}

func newCleaner(friends, groups string) *Cleaner {
	if key == "" {
		log.Fatalln("Key is not provided")
	}
	c := &Cleaner{}
	if friends != "" {
		c.friends = strings.Split(friends, ":")
	}
	if groups != "" {
		c.groups = strings.Split(groups, ":")
	}
	return c
}

func (c *Cleaner) fetchPictures() {
	for _, f := range c.friends {
		c.fetchPicturesByFriendOrGroup(f, 0)
	}

	for _, g := range c.groups {
		c.fetchPicturesByFriendOrGroup(g, 1)
	}
}

func (c *Cleaner) fetchPicturesByFriendOrGroup(id string, mode int) {
	//mode: 0=friend, 1=group
	var prefix string
	if mode == 0 {
		prefix = "mr_friend_"
	} else {
		prefix = "mr_troop_"
	}
	table := prefix + getMD5(id) + "_New"
	query := fmt.Sprintf("SELECT msgData, msgType FROM %s WHERE msgType=-2000 OR msgType=-1035 ORDER BY time;", table)
	rows1, err1 := slowtableDB.Query(query)
	rows2, err2 := qqDB.Query(query)

	if err1 == nil {
		defer rows1.Close()
		for rows1.Next() {
			c.handleMsgRow(rows1)
		}
	}
	if err2 == nil {
		defer rows2.Close()
		for rows2.Next() {
			c.handleMsgRow(rows2)
		}
	}

	if err1 != nil && err2 != nil {
		log.Println("Failed to fetch pictures from " + id)
	}
}

func (c *Cleaner) handleMsgRow(row *sql.Rows) {
	var msgData []byte
	var msgType int
	err := row.Scan(&msgData, &msgType)
	if err != nil {
		log.Println(err)
		return
	}

	data := decryptBytes(msgData)
	if msgType == -2000 {
		decodePic(data)
	} else if msgType == -1035 {
		decodeMix(data)
	}
}

func (c *Cleaner) clean() {
	found := 0
	removed := 0
	size := int64(0)
	for _, pic := range pics {
		if stat, err := os.Stat(pic); !os.IsNotExist(err) {
			found += 1
			if !dryRun {
				err1 := os.Remove(pic)
				if err1 != nil {
					log.Printf("Failed to delete file %s: %x", pic, err1)
				} else {
					removed += 1
				}
			}
			size += stat.Size()
		}
	}

	sizeH := float64(size)
	units := []string{"bytes", "kb", "mb", "gb"}
	var unit string
	for i, u := range units {
		unit = u
		if sizeH < 1024 {
			break
		}
		if i < 4 {
			sizeH /= 1024
		}
	}
	log.Printf("Totally %d pictures in databases, %d found (%f %s), %d removed\n",
		len(pics), found, sizeH, unit, removed)
}

func getMD5(num string) string {
	m := md5.Sum([]byte(num))
	mStr := hex.EncodeToString(m[:])
	return strings.ToUpper(mStr)
}
