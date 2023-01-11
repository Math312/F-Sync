package main

import (
	"github.com/fsync/config"
	"github.com/fsync/watcher"
	"github.com/fsync/web"
	"net/http"
	"time"
)

var closeMain chan int

//var closeServer chan int

func main() {

	router := web.FsyncManageRouter()
	closeMain = make(chan int)
	config.GetBaiduYunConfigContainer(config.DEFAULT_BAIDUYUN_CONFIG_FILE)
	//closeServer = make(chan int)
	go func() {
		s := &http.Server{
			Addr:         ":8080",
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		s.ListenAndServe()
	}()
	watcher.CreateNewWatcher("D:/File")
	defer watcher.CloseAll()
	// Block main goroutine forever.
	<-closeMain
	//close(closeMain)
}
