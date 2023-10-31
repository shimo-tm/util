// 文件系统监听器

package explorer

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

const (
	OpCreate = "Create"
	OpWrite  = "Write"
	OpRemove = "Remove"
	OpRename = "Rename"
)

type Event struct {
	FileName string // 文件名
	OPType   string // 变更类型
	err      error  // 发生错误
}

var watcher *fsnotify.Watcher
var eventChan = make(chan Event)

func Add(name string) error {
	return watcher.Add(name)
}

func Remove(name string) error {
	return watcher.Remove(name)
}

func ChangeChan(name string) chan Event {
	return eventChan
}

func init() {
	if watcher != nil {
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Init Explorer Watcher failed: ", err)
	}
	defer watcher.Close()
	go func() {
		defer close(eventChan)
		for {
			select {
			case event := <-watcher.Events:
				if event.Op == fsnotify.Chmod {
					continue
				}
				eventChan <- Event{
					FileName: event.Name,
					OPType:   event.Op.String(),
				}
			case err := <-watcher.Errors:
				eventChan <- Event{
					err: err,
				}
			}
		}
	}()
}
