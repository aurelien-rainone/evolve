package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jaschaephraim/lrserver"
	"gopkg.in/fsnotify.v1"
)

func checkf(err error, format string, a ...interface{}) {
	if err == nil {
		return
	}

	log.Fatalf("%s: %s", fmt.Sprintf(format, a...), err)
}

func liveReload() {
	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	lr := lrserver.New(lrserver.DefaultName, lrserver.DefaultPort)
	go lr.ListenAndServe()

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					lr.Reload(event.Name)
				}
			case err := <-watcher.Errors:
				lr.Alert(err.Error())
			}
		}
	}()

	err = watcher.Add("./example/tsp/index.html")
	if err != nil {
		log.Fatalln(err)
	}

	select {}
}

func main() {
	// s := []int{1, 2, 3, 4, 5}
	// Perm(s, func(p []int) { fmt.Println(p) })

	host := "localhost:8080"
	flag.StringVar(&host, "host", host, "http server [host]:[port]")
	flag.Parse()

	go liveReload()

	server := newServer()
	server.serve(host)
	err := server.start()
	if err != nil {
		log.Fatal(err)
	}
}
