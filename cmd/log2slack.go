package main

import (
	"bufio"
	"errors"
	"log"
	"log2slack/tools"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if len(os.Args) != 2 {
		err := errors.New("ログファイルを指定してください。")
		log.Fatal(err)
	}

	filepath := os.Args[1]
	fp, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	cur, err := fp.Seek(0, os.SEEK_END)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("seek succeeded by whence cur=", cur)

	scanner := bufio.NewScanner(fp)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event=", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					scanner.Scan()
					msg := scanner.Text()
					log.Printf("msg=%v\n", msg)
					b := []byte(msg)
					doneSend, err := tools.SendAsync(b)
					if err != nil {
						log.Printf("err=%#v\n", err)
					}
					if doneSend != nil {
						<-doneSend
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error=", err)
			}
		}
	}()
	err = watcher.Add(filepath)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Hello world.")
	<-done
}
