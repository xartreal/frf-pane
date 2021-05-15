// indexer
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/xartreal/frfpanehtml"
)

var (
	ListDB    KVBase
	HashtagDB KVBase
	ByMonthDB KVBase
	IdxDB     KVBase //fts
	TlxDB     KVBase
)

var dbgout string

func addToDBList(list string, id string, indb *KVBase) {
	fbin, _ := indb.MyCollection.Get([]byte(list))
	indb.MyCollection.Set([]byte(list), addtolist(fbin, id))
}

func getHashList(text string) []string {
	hlist := frfpanehtml.RegHashtag.FindAllString(text, -1) //FindAllStringSubmatch ?
	for i := 0; i < len(hlist); i++ {
		e := hlist[i]
		idx := strings.Index(e, "#") + 1
		hlist[i] = e[idx:] //clean excess chars
		if Config.debugmode == 1 {
			dbgout += fmt.Sprintf("%q --> %q\n", e, hlist[i])
		}
	}
	return uniqueNonEmptyElementsOf(hlist)
}

const pane = "pane"

func indexer(dbpath string) {
	createDB(dbpath+"list.db", pane, &ListDB)
	createDB(dbpath+"hashtag.db", pane, &HashtagDB)
	createDB(dbpath+"tym.db", pane, &ByMonthDB)

	openDB(dbpath+"list.db", pane, &ListDB)
	openDB(dbpath+"hashtag.db", pane, &HashtagDB)
	openDB(dbpath+"tym.db", pane, &ByMonthDB)

	if RunCfg.ftsenabled {
		createDB(dbpath+"index.db", pane, &IdxDB)
		createDB(dbpath+"timelx.db", pane, &TlxDB)
		openDB(dbpath+"index.db", pane, &IdxDB)
		openDB(dbpath+"timelx.db", pane, &TlxDB)
	}

	logtxt := ""
	hstart := 0
	idx := make(index)
	start := time.Now()
	ilist := RunCfg.feedpath + "index/list_"
	jposts := RunCfg.feedpath + "json/posts_"
	for isexists(ilist + strconv.Itoa(hstart)) {
		jpost := new(FrFjson)
		logtxt += fmt.Sprintf("offset %d\n", hstart)
		hcnt := strconv.Itoa(hstart)
		fbin, _ := ioutil.ReadFile(ilist + hcnt)
		postList := strings.Split(string(fbin), "\n")
		ListDB.MyCollection.Set([]byte("list_"+hcnt), fbin) // add to list.db
		for i := 0; i < len(postList); i++ {
			if len(postList[i]) < 2 {
				continue
			}
			postBin, _ := ioutil.ReadFile(jposts + postList[i])
			json.Unmarshal(postBin, jpost)
			//timelist y-m
			utime, _ := strconv.ParseInt(jpost.Posts.CreatedAt, 10, 64)
			qCreated := time.Unix(utime/1000, 0).Format("2006-01")
			addToDBList(qCreated, postList[i], &ByMonthDB)
			//hashlist
			postText := jpost.Posts.Body
			for j := 0; j < len(jpost.Comments); j++ {
				postText += "\n" + jpost.Comments[j].Body
			}
			hashList := getHashList(postText)
			for j := 0; j < len(hashList); j++ {
				logtxt += fmt.Sprintf("%q\n", hashList[j])
				addToDBList(hashList[j], postList[i], &HashtagDB)
			}
			if RunCfg.ftsenabled {
				TlxDB.MyCollection.Set([]byte(postList[i]), []byte(jpost.Posts.CreatedAt))
				// index it!
				idx.add(postText, postList[i])
			}
		}
		fmt.Printf("\roffset: %d", hstart)
		hstart += Config.step
	}
	closeDB(&ListDB)
	closeDB(&HashtagDB)
	closeDB(&ByMonthDB)
	if Config.debugmode == 1 {
		ioutil.WriteFile("pane.log", []byte(logtxt), 0755)
		ioutil.WriteFile("pane2.log", []byte(dbgout), 0755)
	}
	if RunCfg.ftsenabled {
		//mem
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		//save
		idx.save(&IdxDB)
		fmt.Printf("\nFTS: indexed %d items in %v (MemAlloc = %d MiB)\n", len(idx), time.Since(start), bToMb(m.TotalAlloc))
		// coda
		closeDB(&IdxDB)
		closeDB(&TlxDB)
	}
	fmt.Printf("\n")
}
