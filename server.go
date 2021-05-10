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
		if len(list[i]) < 2 {
			continue
		}
		item := strings.Replace(list[i], "list_", "", -1)
		out += "<a href=/t/" + item + "> offset " + item + "</a><br>"
	}
	out = mkhtml(out, "Timeline")
	return c.HTML(200, out)
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
	list := recsDB(&ByMonthDB)
	out += calendar(list)
	out = "<tr><td>" + out + "</td></tr>"
	out = mkhtml(out, "By month")
	return c.HTML(200, out)
}

func hhandler(c echo.Context) error {
	lnum, _ := url.QueryUnescape(c.Param("id"))
	if lnum == "all" {
		return hlisthandler(c)
	}
	fbin, _ := HashtagDB.MyCollection.Get([]byte(lnum))
	list := strings.Split(string(fbin), "\n")
	title := "Hashtag #" + lnum
	return c.HTML(200, genhtml(list, lnum, false, title, ""))
}

func hlisthandler(c echo.Context) error {
	out := "<h2>Hashtags</h2>"
	list := recsDB(&HashtagDB)
	for i := 0; i < len(list); i++ {
		var id = list[i]
		if len(id) < 2 {
			continue
		}
		hcnt := inlistcount(id, &HashtagDB)
		out += fmt.Sprintf("<a href=/h/%s>%s</a> (%d)<br>", id, id, hcnt)
	}
	out = mkhtml(out, "Hashtags")
	return c.HTML(200, out)
}

func stathandler(c echo.Context) error {
	out := "<h2>Statistics</h2>"
	list := recsDB(&ListDB)
	out += fmt.Sprintf("<p><b>Pages</b>: %d (~%d records)</p>", len(list)-1, (len(list)-1)*30)
	list = recsDB(&ByMonthDB)
	out += fmt.Sprintf("<p><b>By month</b>: %d items</p><p>", len(list)-1)
	for i := 0; i < len(list); i++ {
		if len(list[i]) < 3 {
			continue
		}
		items, _ := ByMonthDB.MyCollection.Get([]byte(list[i]))
		xi := strings.Split(string(items), "\n")
		out += fmt.Sprintf("%s (%d), ", list[i], len(xi)-1)
	}
	list = recsDB(&HashtagDB)
	out += fmt.Sprintf("</p><p><b>By hashtag</b>: %d items</p><p>", len(list)-1)
	for i := 0; i < len(list); i++ {
		if len(list[i]) < 3 {
			continue
		}
		items, _ := HashtagDB.MyCollection.Get([]byte(list[i]))
		xi := strings.Split(string(items), "\n")
		out += fmt.Sprintf("%s (%d), ", list[i], len(xi)-1)
	}
	//
	out = mkhtml(out, "Statistics")
	return c.HTML(200, out)
}

func findFrontHandler(c echo.Context) error {
	out := `
	<h2>Find</h2>
	<form action=/f method=post>
  	<input name=mode type=hidden value=none>
  	<div class="form-item">
    	<label>Search in feed for:</label>
    	<input class="w50" name=qword type=text>
  	</div>
  	<button type=submit class="button" value=Submit>Find</button>
  	</form>
	`
	out = mkhtml(out, "Find")
	return c.HTML(200, out)
}

func findBackHandler(c echo.Context) error {
	qword := c.FormValue("qword")
	if len([]rune(qword)) < 3 {
		return c.HTML(200, mkhtml("Request too small", "Error"))
	}
	xlist := recsDB(&ListDB)
	natsort.Sort(xlist)
	out := ""
	founded := []string{}
	for i := 0; i < len(xlist); i++ {
		fbin, _ := ListDB.MyCollection.Get([]byte(xlist[i]))
		list := strings.Split(string(fbin), "\n")
		for j := 0; j < len(list); j++ {
			if len(list[j]) < 5 {
				continue
			}
			ftxt := frfpanehtml.LoadJson(jpath + list[j]).TextOnly()
			if strings.Contains(ftxt, qword) {
				founded = append(founded, list[j])
			}
		}
	}
	out = genhtml(founded, "0", false, "Find: "+qword+" ("+strconv.Itoa(len(founded))+")", qword)
	return c.HTML(200, out)
}

func changeFeedFront(c echo.Context) error {
	feedlist := getFeedList()
	out := "<h3>Select feed</h3>"
	out += "<table>"
	for k, v := range feedlist {
		if strings.Contains(v, "#") { //no active jsons
			continue
		}
		if !isexists("feeds/" + k + "/pane/list.db") {
			continue
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
	MkFeedPath(newfeed)
	dbpath := RunCfg.feedpath + "pane/"
	openDB(dbpath+"list.db", "pane", &ListDB)
	openDB(dbpath+"hashtag.db", "pane", &HashtagDB)
	openDB(dbpath+"tym.db", "pane", &ByMonthDB)
	loadtemplates()
	RunCfg.maxlastlist = (len(recsDB(&ListDB)) - 1) * Config.step
	e.Static("/media", RunCfg.feedpath+"media") //rewrite server rules
	out := "Changed to feed '" + newfeed + "'"
	return c.HTML(200, mkhtml(out, "Change feed"))
}
