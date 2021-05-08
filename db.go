package main

import (
	//	"fmt"
	//	"io/ioutil"
	"log"
	"os"

	"github.com/steveyen/gkvlite"
)

type KVBase struct {
	StoreFile    *os.File
	MyStore      *gkvlite.Store
	MyCollection *gkvlite.Collection
}

func openDB(dbname string, colname string, outdb *KVBase) {
	var err error
	outdb.StoreFile, err = os.OpenFile(dbname, os.O_RDWR, 0755)
	if err != nil {
		println("FATAL: Can't open database")
		log.Printf("FATAL: Can't open database %s\n", dbname)
		os.Exit(1)
	}

	outdb.MyStore, _ = gkvlite.NewStore(outdb.StoreFile)
	outdb.MyCollection = outdb.MyStore.GetCollection(colname)
}

func closeDB(outdb *KVBase) {
	outdb.MyStore.Flush()
	outdb.StoreFile.Sync()
	outdb.StoreFile.Close()
}

func syncDB(outdb *KVBase) {
	outdb.MyStore.Flush()
	outdb.StoreFile.Sync()
}

func createDB(dbname string, colname string, indb *KVBase) {
	f, err := os.Create(dbname)
	if err != nil {
		println("FATAL: Can't create database")
		log.Printf("FATAL: Can't create database %s\n", dbname)
		os.Exit(1)
	}
	indb.StoreFile = f
	indb.MyStore, _ = gkvlite.NewStore(indb.StoreFile)
	indb.MyStore.SetCollection(colname, nil)
	indb.MyStore.Flush()
	indb.StoreFile.Sync()
	indb.StoreFile.Close()
}

func recsDB(indb *KVBase) []string {
	var outlist = []string{}
	indb.MyCollection.VisitItemsAscend([]byte(""), false, func(i *gkvlite.Item) bool {
		// This visitor callback will be invoked with every item
		// with key "ford" and onwards, in key-sorted order.
		// So: "mercedes", "tesla" are visited, in that ascending order,
		// but not "bmw".
		// If we want to stop visiting, return false;
		// otherwise return true to keep visiting.
		//fmt.Printf("up: %s\n", string(i.Key))
		outlist = append(outlist, string(i.Key))
		//ulisttxt += string(i.Key)
		return true
	})
	return outlist
}
