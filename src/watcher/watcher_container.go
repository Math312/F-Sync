package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/fsync/common"
	"log"
	"strings"
	"sync"
	"sync/atomic"
)

var containerInitialized uint32

var watchContainer WatchContainer

var watcherContainerMu sync.Mutex

type WatchContainer struct {
	watchers        map[string]*fsnotify.Watcher
	fileSyncLogChan map[string]chan string
}

func GetWatcherContainer() WatchContainer {
	if atomic.LoadUint32(&containerInitialized) == 1 {
		return watchContainer
	}
	watcherContainerMu.Lock()
	defer watcherContainerMu.Unlock()

	if containerInitialized == 0 {
		watchContainer = WatchContainer{
			watchers:        make(map[string]*fsnotify.Watcher, 10),
			fileSyncLogChan: make(map[string]chan string, 10),
		}
		atomic.StoreUint32(&containerInitialized, 1)
	}
	return watchContainer
}

func CreateNewWatcher(fileFolderName string) {
	watchContainer = GetWatcherContainer()
	watcherContainerMu.Lock()
	defer watcherContainerMu.Unlock()
	_, ok := watchContainer.watchers[fileFolderName]
	if ok {
		return
	} else {
		newWatcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		watchContainer.watchers[fileFolderName] = newWatcher
		logChan := make(chan string, 200)
		go listenedFunc(*newWatcher, logChan)
		go consumeLog(logChan, fileFolderName, "/test/")
		err = newWatcher.Add(fileFolderName)
		if err != nil {
			panic("123")
		}
	}

}

func CloseAll() {
	if atomic.LoadUint32(&containerInitialized) == 0 {
		return
	}
	watcherContainerMu.Lock()
	defer watcherContainerMu.Unlock()
	if atomic.LoadUint32(&containerInitialized) == 0 {
		return
	}
	for fileName, watcher := range watchContainer.watchers {
		err := watcher.Close()
		if err != nil {
			log.Fatalf("Watcher of File : %s , close fail.", fileName)
		} else {
			log.Printf("Watcher of File : %s , has closed.", fileName)
		}
		logChan, ok := watchContainer.fileSyncLogChan[fileName]
		if ok {
			close(logChan)
		}
	}
	atomic.StoreUint32(&containerInitialized, 0)

}

func consumeLog(logChan chan string, fileFolderName string, baiduYunBasePath string) {
	for true {
		data, ok := <-logChan
		if data == "" && !ok {
			log.Println(data)
			break
		} else {
			log.Printf("Received LogCommand: %s", data)
			doConsumeLog(data, fileFolderName, baiduYunBasePath)
		}
	}
	log.Printf("LogChan closed, File Folder Name: %s", fileFolderName)
}

func doConsumeLog(fileProcessLog string, basePath string, baiduYunBasePath string) {
	accessToken := "123.7762029b59e60cc85f1d5396ee2fe37f.Yg5RdwKy6ZK6dMct0QAMASHGZicl_1uN-N41F1n.xaSXIw"
	if len(fileProcessLog) > 6 {
		command := fileProcessLog[:6]
		fileName := fileProcessLog[6:]
		switch command {
		case "WRITE ":
			fallthrough
		case "CREATE":
			log.Printf("File upload, File name: %s, Log Command: %s", fileName, fileProcessLog)
			relativeFileName := "/" + fileName[len(basePath):]
			baiduYunFileName := baiduYunBasePath + relativeFileName
			common.UploadFileToBaiduYun(accessToken, fileName, baiduYunFileName)
		}
	}

}

func listenedFunc(watcher fsnotify.Watcher, logChan chan string) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if !strings.HasSuffix(event.Name, ".fsync") {
				log.Println("event:", event)
				logChan <- createLogCommand(event)
				if event.Has(fsnotify.Write) {
					log.Println("Modified file:", event.Name)
				} else if event.Has(fsnotify.Create) {
					log.Println("Create New file:", event.Name)
				} else if event.Has(fsnotify.Rename) {
					// event - rename file can trigger the event - create file
					// so the logic of rename file event is same as the remove file
					log.Printf("Rename file: %s to New File", event.Name)
				} else if event.Has(fsnotify.Remove) {
					log.Printf("Remove file: %s", event.Name)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func createLogCommand(event fsnotify.Event) string {
	if event.Has(fsnotify.Write) {
		return "WRITE " + event.Name
	} else if event.Has(fsnotify.Create) {
		return "CREATE" + event.Name
	} else if event.Has(fsnotify.Rename) {
		// event - rename file can trigger the event - create file
		// so the logic of rename file event is same as the remove file
		return "RENAME" + event.Name
	} else if event.Has(fsnotify.Remove) {
		return "REMOVE" + event.Name
	} else {
		panic("Unsupported Log")
	}
}
