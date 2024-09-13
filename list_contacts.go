package main

import (
	"database/sql"
	"fmt"
	"log"
)

var isList = false

type FriendListItem struct {
	uin    string
	name   string
	remark string
}

type GroupListItem struct {
	uin  string
	name string
}

func listContacts() {
	listFriends()
	listGroups()
}

func listFriends() {
	friends := []*FriendListItem{}
	query := "SELECT uin, name, IFNULL(remark, '') FROM Friends;"
	rows, err := qqDB.Query(query)
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			friends = append(friends, handleFriendRow(rows))
		}
	}
	fmt.Println("Friends:")
	for _, f := range friends {
		fmt.Printf("%s\t%s\t(%s)\n", f.uin, f.name, f.remark)
	}
}

func listGroups() {
	groups := []*GroupListItem{}
	query := "SELECT troopuin, troopname FROM TroopInfoV2;"
	rows, err := qqDB.Query(query)
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			groups = append(groups, handleGroupRow(rows))
		}
	}
	fmt.Println("Groups:")
	for _, g := range groups {
		fmt.Printf("%s\t%s\n", g.uin, g.name)
	}
}

func handleFriendRow(row *sql.Rows) *FriendListItem {
	friend := &FriendListItem{}
	err := row.Scan(&friend.uin, &friend.name, &friend.remark)
	if err != nil {
		log.Println(err)
		return nil
	}

	friend.uin = decryptString(friend.uin)
	friend.name = decryptString(friend.name)
	friend.remark = decryptString(friend.remark)
	return friend
}

func handleGroupRow(row *sql.Rows) *GroupListItem {
	group := &GroupListItem{}
	err := row.Scan(&group.uin, &group.name)
	if err != nil {
		log.Println(err)
		return nil
	}
	group.uin = decryptString(group.uin)
	group.name = decryptString(group.name)
	return group
}
