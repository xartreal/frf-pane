// html
package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

    "github.com/xartreal/frfpanehtml"
)

var jpath string //json path

func loadtfile(name string) string {
	fbin, err := ioutil.ReadFile(name)
	if err != nil {
		outerror(2, "FATAL: Template load error: '%s'\n", name)
	}
	return string(fbin)
}

func loadtemplates() {
	frfpanehtml.Templates = &frfpanehtml.THtmlTemplate{
		Comment: loadtfile("template/template_comment.html"),
		Item:    loadtfile("template/template_item.html"),
		File:    loadtfile("template/template_kfile.html"),
		Cal:     loadtfile("template/template_kcal.html"),
	}
	//+params
	frfpanehtml.Params = frfpanehtml.TParams{Feedpath: RunCfg.feedpath, Step: Config.step, Singlemode: false,
		IndexPrefix: "/t/", IndexPostfix: ""}
	jpath = RunCfg.feedpath + "json/posts_"
}

func genhtml(list []string, id string, isindex bool, title string, pen string) string {
	maxx := len(list)
	outtext := "<h2>" + title + "</h2>"
	for i := 0; i < maxx; i++ {
		if len(list[i]) < 2 {
			continue
		}
		outtext += frfpanehtml.LoadJson(jpath+list[i]).ToHtml(list[i], pen) + "<hr>"
	}
	ptitle := RunCfg.feedname + " - " + title + " - frf-pane"
	return frfpanehtml.MkHtmlPage(id, outtext, isindex, RunCfg.maxlastlist, RunCfg.feedname, ptitle)
}

func mkhtml(htmlText string, title string) string {
	ptitle := RunCfg.feedname + " - " + title + " - frf-pane"
	return frfpanehtml.MkHtmlPage("0", htmlText, false, 0, RunCfg.feedname, ptitle)
}

func calendar(list []string) string {
	// scan list
	min := 2030
	max := 2000
	for i := 0; i < len(list); i++ {
		z := strings.Split(list[i], "-")
		year, _ := strconv.Atoi(z[0])
		if year < min {
			min = year
		}
		if year > max {
			max = year
		}
	}
	//	fmt.Printf("min=%d max=%d\n", min, max)
	// builder
	tidx := [13]string{}
	tout := [13]string{}
	ccout := ""
	for i := min; i <= max; i++ { //years
		yout := fmt.Sprintf(`<span class="label is-error is-secondary"><b>%d</b></span>`, i)
		for j := 1; j < 13; j++ { //months
			tidx[j] = fmt.Sprintf("%%%02d", j)
			listname := fmt.Sprintf("%d-%02d", i, j)
			lcnt := inlistcount(listname, &ByMonthDB)
			//	ltxt := strconv.Itoa(lcnt)
			if lcnt > 0 {
				tout[j] = `<br><b><a href=/m/` + listname + `>` + strconv.Itoa(lcnt) + `</a></b>`
			} else {
				tout[j] = ""
			}
		}
		//		fmt.Printf("year=%d, tidx=%v, tout=%v\n", i, tidx, tout)
		xout := frfpanehtml.Templates.Cal
		for j := 1; j < 13; j++ {
			xout = strings.Replace(xout, tidx[j], tout[j], -1)
		}
		ccout += yout + xout
	}
	return ccout
}
