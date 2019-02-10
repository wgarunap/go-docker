package main

import (
	"fmt"
	"github.com/tecbot/gorocksdb"
	"github.com/wgarunap/go-docker/config"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"time"
)

/*
*	Hello...
*	This Test is written for beginners who works with Golang and docker
*	to understand how to setup the environment and how to reduce the
*	container size with Multi-Stage docker building. Dockerfile.build
*	file is used to build the base image and used it for every intermediate
*	build images, Dockerfile.rocks file is used to build the final
*	production ready image which includes RocksDB only.
 */

var (
	key   = []byte(`aruna`)
	value = []byte(`{"age":25,"country":"sri-lanka"}`)
)

func main() {
	db, err := newRocksDB(`userdata`)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Put(gorocksdb.NewDefaultWriteOptions(), key, value)
	if err != nil {
		log.Println(`error writing to rocks database`, err)
	}

	err = db.Flush(gorocksdb.NewDefaultFlushOptions())
	if err != nil {
		log.Println(`error flushing to rocks database`, err)
	}

	readSlice, err := db.GetBytes(gorocksdb.NewDefaultReadOptions(), key)
	if err != nil {
		log.Println(`error reading the rocks database`, err)
	}

	//db output print 10 times
	go func(readSlice []byte) {
		for i := 0; i < 10; i++ {
			fmt.Println(`Data stored to RocksDB:`, string(readSlice))
			time.Sleep(1 * time.Second)
		}
	}(readSlice)

	//Go profiling router start
	go func() {
		log.Println(http.ListenAndServe(fmt.Sprintf("localhost:%d", config.Config.DebugPort), nil))
	}()

	//waiting for a keyboard response to end the service
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		select {
		case <-sig:
			wg.Done()
		}
	}()
	wg.Wait()
}

func newRocksDB(name string) (db *gorocksdb.DB, err error) {
	conf := gorocksdb.NewDefaultOptions()
	conf.SetCreateIfMissing(true)
	conf.SetUseFsync(true)
	conf.SetWALTtlSeconds(uint64(30 * 24 * time.Hour))

	storage := `storage`
	err = os.MkdirAll(storage, os.ModePerm)
	if err != nil {
		return nil, err
	}

	newdb, err := gorocksdb.OpenDb(conf, storage+`/`+name)
	if err != nil {
		log.Fatal(`cannot open rocksdb storage`)
	}
	return newdb, nil

}
