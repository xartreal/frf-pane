// sort
package main

import (
	"sort"
)

type lxrec struct {
	id    string
	tmark string
}

func timesort(list []string) []string {
	//return list //debug
	slist := []lxrec{}
	for i := 0; i < len(list); i++ { //mk sortlist
		tbin, _ := TlxDB.MyCollection.Get([]byte(list[i]))
		slist = append(slist, lxrec{list[i], string(tbin)})
	}

	sort.SliceStable(slist, func(i, j int) bool {
		return slist[i].tmark > slist[j].tmark
	})

	xlist := []string{}
	for i := 0; i < len(slist); i++ { //mk outlist
		xlist = append(xlist, slist[i].id)
	}
	return xlist
}

// from https://gist.github.com/johnwesonga/6301924
func uniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}
	return us
}
