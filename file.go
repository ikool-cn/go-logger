package goLogger

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type LogFile struct {
	filename string
	flag     int //default log.LstdFlags

	sync struct {
		frequency time.Duration
		status    int
	}

	level int

	cache struct {
		enable bool
		data   []string
		mutex  sync.Mutex
	}

	rotate struct {
		enable     bool
		rotateType int
		mutex      sync.Mutex
	}
}

type AsyncFilesMap struct {
	files map[string]*LogFile
	sync.RWMutex
}

const (
	SyncInit  = iota
	SyncDoing
	SyncDone
)

const (
	RotateHour = iota
	RotateDate
)

const (
	RotateHourFormat = "2006010215"
	RotateDateFormat = "20060102"
)

const (
	logTimeFormat string = "2006-01-02 15:04:05"
	fileOpenMode         = 0666
	fileFlag             = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	newlineStr           = "\n"
	cacheInitCap         = 1024
)

const (
	NoFlag  = 0
	StdFlag = log.LstdFlags
)

var asyncFiles *AsyncFilesMap

var nowFunc = time.Now

func init() {
	asyncFiles = &AsyncFilesMap{
		files: make(map[string]*LogFile),
	}

	//timer := time.NewTicker(lf.sync.frequency)
	timer := time.NewTicker(time.Second * 1)
	go func() {
		for {
			select {
			case <-timer.C:
				asyncFiles.RLock()
				for _, file := range asyncFiles.files {
					if file.sync.status != SyncDoing {
						go file.Flush()
					}
				}
				asyncFiles.RUnlock()
			}
		}

	}()
}

func NewFileLogger(filename string) *LogFile {
	asyncFiles.Lock()
	defer asyncFiles.Unlock()

	if lf, ok := asyncFiles.files[filename]; ok {
		return lf
	}

	lf := &LogFile{
		filename: filename,
		flag:     StdFlag,
	}

	asyncFiles.files[filename] = lf
	lf.rotate.enable = true
	lf.rotate.rotateType = RotateDate
	lf.cache.enable = true

	return lf
}

// save to file
func (lf *LogFile) Flush() error {
	lf.sync.status = SyncDoing
	defer func() {
		lf.sync.status = SyncDone
	}()

	lf.cache.mutex.Lock()
	cache := lf.cache.data
	lf.cache.data = make([]string, 0, cacheInitCap)
	lf.cache.mutex.Unlock()

	if len(cache) == 0 {
		return nil
	}

	lf.rotate.mutex.Lock()
	defer lf.rotate.mutex.Unlock()
	file, err := lf.openFile()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(cache, ""))
	if err != nil {
		panic(err)
	}

	return nil
}

func (lf *LogFile) write(msg string) error {
	if lf.flag == StdFlag {
		msg = nowFunc().Format(logTimeFormat) + " " + msg + newlineStr
	} else {
		msg = msg + newlineStr
	}

	if lf.cache.enable {
		lf.appendCache(msg)
		return nil
	}

	return lf.directWrite([]byte(msg))
}

func (lf *LogFile) appendCache(msg string) {
	lf.cache.mutex.Lock()
	lf.cache.data = append(lf.cache.data, msg)
	lf.cache.mutex.Unlock()
}

func (lf *LogFile) directWrite(msg []byte) error {
	lf.rotate.mutex.Lock()
	defer lf.rotate.mutex.Unlock()
	file, err := lf.openFile()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(msg)
	return err
}

func (lf *LogFile) getFilename(filename string) string {
	if lf.rotate.enable == false {
		return filename
	}

	var suffix string
	if lf.rotate.rotateType == RotateDate {
		suffix = nowFunc().Format(RotateDateFormat)
	} else {
		suffix = nowFunc().Format(RotateHourFormat)
	}
	arr := strings.Split(filename, ".")
	l := len(arr)
	if l <= 1 {
		return filename + "." + suffix
	} else {
		return strings.Join(arr[:l-1], "") + "." + suffix + "." + strings.Join(arr[l-1:], "");
	}
}

func (lf *LogFile) openFile() (*os.File, error) {
	logFilename := lf.getFilename(lf.filename)
	file, err := os.OpenFile(logFilename, fileFlag, fileOpenMode)
	if err != nil {
		panic(err)
	}

	return file, nil
}
