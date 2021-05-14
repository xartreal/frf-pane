// lister
package main

import (
	"fmt"
	"io/ioutil"
)

func getFeedList() map[string]string {
	feedmap := map[string]string{}
	files, _ := ioutil.ReadDir("feeds")
	for _, item := range files {
		if !item.IsDir() {
			continue
		}
		flag := ""
		if !isexists("feeds/" + item.Name() + "/json/profile") {
			flag = " # no json data"
		} else {
			jfiles, _ := ioutil.ReadDir("feeds/" + item.Name() + "/json")
			flag = fmt.Sprintf(" %d jsons", len(jfiles)-1)
			if !isexists("feeds/" + item.Name() + "/pane/list.db") {
				flag += " (not indexed)"
			} else if isexists("feeds/" + item.Name() + "/pane/index.db") {
				flag += " (fts)"
			}
		}
		feedmap[item.Name()] = flag
	}
	return feedmap
}

func listFeeds(feedname string) {
	if !isexists("feeds") {
		outerror(1, "No feeds\n")
	}
	if feedname != "all" {
		outerror(1, "Incorrect list cmd\n")
	}

	fmt.Printf("Feeds:\n\n")
	feedmap := getFeedList()
	for k, v := range feedmap {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Printf("\n")
}
