// server
package main

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/facette/natsort"
	"github.com/labstack/echo"
	"github.com/xartreal/frfpanehtml"
)

var e *echo.Echo

func startServer() {
	e = echo.New()
	e.Static("/media", RunCfg.feedpath+"media")
	e.Static("/x", "template")
	e.GET("/", mainpage)
	e.GET("/t/:id", thandler)
	e.GET("/m/:id", mhandler)
	e.GET("/h/:id", hhandler)
	e.GET("/s", stathandler)
	e.GET("/f", findFrontHandler)
	e.POST("/f", findBackHandler)
	e.GET("/c", changeFeedFront)
	e.GET("/cx/:feed", changeFeedHandler)
	e.GET("/p/:id", phandler)
	e.Start(RunCfg.port)

}

func mainpage(c echo.Context) error {
	return c.Redirect(307, "/t/0")
}

func thandler(c echo.Context) error {
	lnum := c.Param("id")
	if lnum == "all" {
		return tlisthandler(c)
	}
	fbin, _ := ListDB.MyCollection.Get([]byte("list_" + lnum))
	//	fmt.Printf("lnum=%s file=%s\n", lnum, RunCfg.feedpath+"pane/list_"+lnum)
	list := strings.Split(string(fbin), "\n")
	title := "Timeline (from " + lnum + ")"
	return c.HTML(200, genhtml(list, lnum, true, title, ""))
}

func tlisthandler(c echo.Context) error {
	out := "<h2>Timeline</h2>"
	list := recsDB(&ListDB)
	natsort.Sort(list)
	for i := 0; i < len(list); i++ {
		if len(list[i]) > 1 {
			item := strings.TrimPrefix(list[i], "list_")
			out += "<a href=/t/" + item + "> offset " + item + "</a><br>"
		}
	}
	return c.HTML(200, mkhtml(out, "Timeline"))
}

func mhandler(c echo.Context) error {
	lnum := c.Param("id")
	if lnum == "all" {
		return mlisthandler(c)
	}
	fbin, _ := ByMonthDB.MyCollection.Get([]byte(lnum))
	list := strings.Split(string(fbin), "\n")
	title := "Month: " + lnum
	return c.HTML(200, genhtml(list, lnum, false, title, ""))
}

func mlisthandler(c echo.Context) error {
	out := "<h2>By month</h2>"
	out += calendar(recsDB(&ByMonthDB))
	out = "<tr><td>" + out + "</td></tr>"
	return c.HTML(200, mkhtml(out, "By month"))
}

func hhandler(c echo.Context) error {
	htag, _ := url.QueryUnescape(c.Param("id"))
	if htag == "all" {
		return hlisthandler(c)
	}
	fbin, _ := HashtagDB.MyCollection.Get([]byte(htag))
	list := strings.Split(string(fbin), "\n")
	title := "Hashtag #" + htag
	return c.HTML(200, genhtml(list, htag, false, title, ""))
}

func hlisthandler(c echo.Context) error {
	out := "<h2>Hashtags</h2>"
	list := recsDB(&HashtagDB)
	for i := 0; i < len(list); i++ {
		if len(list[i]) > 2 {
			out += fmt.Sprintf("<a href=/h/%s>%s</a> (%d)<br>", list[i], list[i],
				inlistcount(list[i], &HashtagDB))
		}
	}
	return c.HTML(200, mkhtml(out, "Hashtags"))
}

func phandler(c echo.Context) error {
	id := c.Param("id")
	out := ""
	if !isexists(RunCfg.feedpath + "json/posts_" + id) {
		out = mkhtml("Not found error", "Error")
	} else {
		out = genhtml([]string{id}, "", false, id, "")
	}
	return c.HTML(200, out)
}

func statdb(indb *KVBase, title string) string {
	list := recsDB(indb)
	out := fmt.Sprintf("<p><b>%s</b>: %d items</p><p>", title, len(list)-1)
	for i := 0; i < len(list); i++ {
		if len(list[i]) > 3 {
			items, _ := indb.MyCollection.Get([]byte(list[i]))
			xi := strings.Split(string(items), "\n")
			out += fmt.Sprintf("%s (%d), ", list[i], len(xi)-1)
		}
	}
	return out
}

func stathandler(c echo.Context) error {
	out := "<h2>Statistics</h2>"
	list := recsDB(&ListDB)
	out += fmt.Sprintf("<p><b>Pages</b>: %d (~%d records)</p>", len(list)-1, (len(list)-1)*30) +
		statdb(&ByMonthDB, "By month") +
		statdb(&HashtagDB, "By hashtag")
	return c.HTML(200, mkhtml(out, "Statistics"))
}

func findFrontHandler(c echo.Context) error {
	return c.HTML(200, mkhtml(loadtfile("template/template_kfind.html"), "Find"))
}

func findBackHandler(c echo.Context) error {
	qword := c.FormValue("qword")
	if len([]rune(qword)) < 3 {
		return c.HTML(200, mkhtml("Request too small", "Error"))
	}
	if RunCfg.ftsenabled {
		return findFTSHandler(c)
	}
	xlist := recsDB(&ListDB)
	natsort.Sort(xlist)
	founded := []string{}
	for i := 0; i < len(xlist); i++ {
		fbin, _ := ListDB.MyCollection.Get([]byte(xlist[i]))
		list := strings.Split(string(fbin), "\n")
		for j := 0; j < len(list); j++ {
			if len(list[j]) > 4 {
				ftxt := frfpanehtml.LoadJson(jpath + list[j]).TextOnly()
				if strings.Contains(ftxt, qword) {
					founded = append(founded, list[j])
				}
			}
		}
	}
	return c.HTML(200, genhtml(founded, "0", false,
		"Find: "+qword+" ("+strconv.Itoa(len(founded))+")", qword))
}

func findFTSHandler(c echo.Context) error {
	qword := c.FormValue("qword")
	keys := recsDB(&IdxDB)
	var founded = []string{}
	for i := 0; i < len(keys); i++ {
		if strings.Contains(keys[i], qword) {
			vbin, _ := IdxDB.MyCollection.Get([]byte(keys[i]))
			vx := strings.Split(string(vbin), "\n")
			founded = append(founded, vx...)
		}
	}
	closeDB(&IdxDB)
	founded = uniqueNonEmptyElementsOf(founded)
	founded = timesort(founded)
	return c.HTML(200, genhtml(founded, "0", false,
		"Find: "+qword+" ("+strconv.Itoa(len(founded))+")", qword))
}

func changeFeedFront(c echo.Context) error {
	feedlist := getFeedList()
	out := "<h3>Select feed</h3>"
	out += "<table>"
	for k, v := range feedlist {
		if strings.Contains(v, "#") || !isexists("feeds/"+k+"/pane/list.db") {
			continue //no active jsons or no pane index
		}
		out += fmt.Sprintf(`<tr><td><a href=/cx/%s>%s</a></td><td>%s</td></tr>`, k, k, v)
	}
	out += "</table>"
	return c.HTML(200, mkhtml(out, "Change feed"))
}

func changeFeedHandler(c echo.Context) error {
	newfeed := c.Param("feed")
	if !isexists("feeds/" + newfeed + "/pane/list.db") {
		return c.HTML(200, mkhtml("Incorrect or noindexed feed", "Change feed"))
	}
	closeDB(&ListDB)
	closeDB(&HashtagDB)
	closeDB(&ByMonthDB)
	if RunCfg.ftsenabled {
		closeDB(&IdxDB)
		closeDB(&TlxDB)
	}
	MkFeedPath(newfeed)
	dbpath := RunCfg.feedpath + "pane/"
	openDB(dbpath+"list.db", "pane", &ListDB)
	openDB(dbpath+"hashtag.db", "pane", &HashtagDB)
	openDB(dbpath+"tym.db", "pane", &ByMonthDB)
	RunCfg.ftsenabled = isexists(dbpath + "index.db")
	if RunCfg.ftsenabled {
		openDB(dbpath+"index.db", "pane", &IdxDB)
		openDB(dbpath+"timelx.db", "pane", &TlxDB)
	}
	loadtemplates()
	RunCfg.maxlastlist = (len(recsDB(&ListDB)) - 1) * Config.step
	e.Static("/media", RunCfg.feedpath+"media") //rewrite server rules
	out := "Changed to feed '" + newfeed + "'"
	return c.HTML(200, mkhtml(out, "Change feed"))
}
