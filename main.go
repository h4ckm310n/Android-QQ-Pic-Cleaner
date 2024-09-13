package main

import (
	"fmt"
	"log"
	"os"
)

func usage() {
	fmt.Println("Options:\n" +
		"-h, --help: Show this usage\n" +
		"-q: Your QQ number\n" +
		"-f, --friends: Friend numbers, seperated by ':'\n" +
		"-g, --groups: Group numbers, seperated by ':'\n" +
		"-l, --list: List all friends and groups\n" +
		"--dry-run: Do not delete files\n" +
		"Example: " + os.Args[0] + " -q 1234 -f 2345:3456 -g 4567:5678:6789")
	os.Exit(0)
}

func parseArgs(args []string) map[string]string {
	res := make(map[string]string)
	currKey := ""
	for _, arg := range args {
		switch arg {
		case "-h", "--help":
			usage()
		case "-q":
			currKey = "qq"
			break
		case "-f", "--friends":
			currKey = "friends"
			break
		case "-g", "--groups":
			currKey = "groups"
			break
		case "--key":
			currKey = "key"
			break
		case "-l", "--list":
			isList = true
			break
		case "--dry-run":
			dryRun = true
			break
		default:
			if currKey != "" {
				res[currKey] = arg
			}
			currKey = ""
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(0)
	}
	args := parseArgs(os.Args[1:])
	qq, ok := args["qq"]
	if !ok || qq == "" {
		log.Fatalln("QQ is not specified!")
	}

	connectDB(qq)
	defer closeDB()

	if isList {
		listContacts()
		os.Exit(0)
	}

	groups, ok_g := args["groups"]
	friends, ok_f := args["friends"]
	if !(ok_g || ok_f) {
		log.Fatalln("Neither friends nor groups are specified!")
	}

	_key, ok := args["key"]
	if ok && _key != "" {
		key = _key
	}

	cleaner := newCleaner(friends, groups)
	cleaner.fetchPictures()
	cleaner.clean()
}
