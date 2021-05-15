// config
package main

import (
	"strconv"

	"github.com/vaughan0/go-ini"
)

var myversion = "0.3.4"

//var useragent = "ARL backfrf/" + myversion

var RunCfg struct {
	feedpath    string
	port        string
	myname      string
	feedname    string
	maxlastlist int
	ftsenabled  bool
}

var Config struct {
	step      int
	debugmode int
	maxlast   int
	pidfile   string
}

func getIniVar(file ini.File, section string, name string) string {
	rt, _ := file.Get(section, name)
	if len(rt) < 1 {
		outerror(2, "FATAL: Variable '%s' not defined\n", name)
	}
	return rt
}

func getIniNum(file ini.File, section string, name string) int {
	rt, err := strconv.Atoi(getIniVar(file, section, name))
	if err != nil {
		return 0
	}
	return rt
}

func ReadConf() {
	if !isexists("pane.ini") {
		Config.debugmode = 0
		Config.step = 30
		Config.pidfile = ""
		RunCfg.port = ":3000"
		return
	}
	file, err := ini.LoadFile("pane.ini")
	if err != nil {
		outerror(1, "\n! Configuration not available\n")
	}
	Config.step = getIniNum(file, "default", "step")
	Config.debugmode = getIniNum(file, "default", "debug")
	Config.pidfile, _ = file.Get("default", "pidfile")
	RunCfg.port, _ = file.Get("default", "port")
}

func ReadBFConf() string {
	file, err := ini.LoadFile("backfrf.ini")
	if err != nil {
		outerror(1, "\n! Configuration (backfrf.ini) not available, @myname is unknown\n")
	}
	RunCfg.myname, _ = file.Get("credentials", "myname")
	if len(RunCfg.myname) < 2 {
		outerror(1, "\n! @myname not defined in backfrf.ini\n")
	}
	return RunCfg.myname
}
