// funcs
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func isexists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func outerror(code int, format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(code)
}

func MkFeedPath(feedname string) {
	feedpath := ""
	if strings.Contains(feedname, "filter:") {
		feedpath = strings.Replace(feedname, ":", "/", -1)
	} else {
		feedpath = "feeds/" + feedname
	}
	RunCfg.feedpath = feedpath + "/"
	RunCfg.feedname = feedname
}

func addtolist(list []byte, item string) []byte {
	qlist := strings.Split(string(list), "\n")
	for _, qitem := range qlist {
		if qitem == item {
			return []byte(strings.Join(qlist, "\n"))
		}
	}
	qlist = append(qlist, item)
	return []byte(strings.Join(qlist, "\n"))
}

func inlistcount(listname string, indb *KVBase) int {
	rec, _ := indb.MyCollection.Get([]byte(listname))
	if len(rec) < 5 { //rec not found
		return 0
	}
	return len(strings.Split(string(rec), "\n")) - 1
}

func writepid(pid int) {
	ioutil.WriteFile(Config.pidfile, []byte(fmt.Sprintf("%d", pid)), 0644)
}

func rmpid() {
	os.Remove(Config.pidfile)
}
