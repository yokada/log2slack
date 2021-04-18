package main

import (
	"bufio"
	"errors"
	"log"
	"log2slack/tools"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
)

var (
	fp      *os.File = nil
	err     error    = nil
	scanner *bufio.Scanner
)

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func createScanner(filepath string) (*bufio.Scanner, error) {
	closeFile()

	if pathExists(filepath) {
		fp, err = os.Open(filepath)
		if err != nil {
			return nil, err
		}

		_, err := fp.Seek(0, os.SEEK_END)
		if err != nil {
			return nil, err
		}

		s := bufio.NewScanner(fp)
		return s, nil
	}

	return nil, errors.New("file or dir should not be exists")
}

func closeFile() {
	if fp != nil {
		fp.Close()
	}
}

func main() {
	if len(os.Args) != 2 {
		err = errors.New("ログファイルを指定してください。")
		log.Fatal(err)
	}

	filepath := os.Args[1]
	dir := path.Dir(filepath)
	if !pathExists(dir) {
		log.Fatalf("directory not found dir=%v", dir)
	}

	scanner, _ := createScanner(filepath)
	defer closeFile()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	done := make(chan bool)
	defer close(done)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				log.Println("event=", event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("event.Name=%v\n", event.Name)
					if event.Name == filepath {
						scanner, _ = createScanner(event.Name)
					}
				} else if event.Op&fsnotify.Write == fsnotify.Write {
					scanner.Scan()
					log.Printf("scanner error=%#v", scanner.Err())
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
				} else if event.Op&fsnotify.Remove == fsnotify.Remove {
					if event.Name == filepath {
						closeFile()
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

	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Hello world.")
	<-done
}
