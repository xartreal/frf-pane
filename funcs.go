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

func writepid() {
	ioutil.WriteFile(Config.pidfile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
}

func rmpid() {
	os.Remove(Config.pidfile)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func parseCL() ([]string, []string) {
	var args = []string{}
	var flags = []string{}
	for i := 1; i < len(os.Args); i++ {
		if strings.HasPrefix(os.Args[i], "-") { //flags
			flags = append(flags, strings.TrimPrefix(os.Args[i], "-"))
		} else {
			args = append(args, os.Args[i])
		}
	}
	return args, flags
}
