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
	args, flags := parseCL()
	if len(args) != 2 {
		helpprog()
	}
	feedname := args[1]
	if feedname == "@myname" {
		feedname = ReadBFConf()
	}
	if (feedname != "all") && !isexists("feeds/"+feedname) {
		outerror(2, "Feed '%s' not found\n", feedname)
	}
	ReadConf()
	//flags
	for _, v := range flags {
		if v == "i" {
			RunCfg.ftsenabled = true
		}
	}

	MkFeedPath(feedname)
	dbpath := RunCfg.feedpath + "pane/"
	switch args[0] {
	case "build":
		os.RemoveAll(RunCfg.feedpath + "pane")
		os.Mkdir(RunCfg.feedpath+"pane", 0755)
		indexer(dbpath)
	case "list":
		listFeeds(feedname)
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
		RunCfg.ftsenabled = isexists(dbpath + "index.db")
		if RunCfg.ftsenabled {
			openDB(dbpath+"index.db", "pane", &IdxDB)
			openDB(dbpath+"timelx.db", pane, &TlxDB)
			defer closeDB(&IdxDB)
			defer closeDB(&TlxDB)
		}
		loadtemplates()
		RunCfg.maxlastlist = (len(recsDB(&ListDB)) - 1) * Config.step
		if len(Config.pidfile) > 2 {
			writepid()
		}
		startServer()
	default:
		outerror(2, "Unknown command '%s'\n", os.Args[1])
	}
}
