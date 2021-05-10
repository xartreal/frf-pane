package main

import (
	"fmt"
	"os"
)

func helpprog() {
	fmt.Printf("Usage: \n")
	fmt.Printf("build indexes: frf-pane build <feed>\n")
	fmt.Printf("list feeds:    frf-pane list all\n")
	fmt.Printf("server:        frf-pane server <feed>\n")
	os.Exit(1)
}

func main() {
	fmt.Printf("Frf-pane %s\n\n", myversion)
	if len(os.Args) != 3 {
		helpprog()
	}
	feedname := os.Args[2]
	if os.Args[1] == "list" {
		listFeeds(feedname)
		os.Exit(0)
	}
	if feedname == "@myname" {
		feedname = ReadBFConf()
	}
	if !isexists("feeds/" + feedname) {
		outerror(2, "Feed '%s' not found\n", feedname)
	}
	ReadConf()

	MkFeedPath(feedname)
	dbpath := RunCfg.feedpath + "pane/"
	switch os.Args[1] {
	case "build":
		os.RemoveAll(RunCfg.feedpath + "pane")
		os.Mkdir(RunCfg.feedpath+"pane", 0755)
		indexer(dbpath)
	case "server":
		if !isexists(dbpath + "list.db") {
			outerror(2, "FATAL: No index DB found\n")
		}
		openDB(dbpath+"list.db", "pane", &ListDB)
		openDB(dbpath+"hashtag.db", "pane", &HashtagDB)
		openDB(dbpath+"tym.db", "pane", &ByMonthDB)
		defer closeDB(&ListDB)
		defer closeDB(&HashtagDB)
		defer closeDB(&ByMonthDB)
		loadtemplates()
		RunCfg.maxlastlist = (len(recsDB(&ListDB)) - 1) * Config.step
		if len(Config.pidfile) > 2 {
			pid := os.Getpid()
			writepid(pid)
		}
		startServer()
	default:
		outerror(2, "Unknown command '%s'\n", os.Args[1])
	}
}
