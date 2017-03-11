package goLogger

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type asyncLogType struct {
	files map[string]*LogFile

	sync.RWMutex
}

type LogFile struct {
	filename string
	flag     int    //log.LstdFlags

	sync struct {
		frequency  time.Duration
		beginTime time.Time
		status    syncStatus
	}

	level int

	cache struct {
		enable  bool
		data []string

		mutex sync.Mutex
	}

	logRotate struct {
		rotate LogRotate
		file   *os.File
		suffix string
		mutex  sync.Mutex
	}
}

type syncStatus int

const (
	statusInit  syncStatus = iota
	statusDoing
	statusDone
)

type LogRotate int

const (
	RotateHour LogRotate = iota
	RotateDate
)

const (
	logTimeFormat string = "2006-01-02 15:04:05"

	fileOpenMode = 0666

	fileFlag = os.O_WRONLY | os.O_CREATE | os.O_APPEND

	newlineStr  = "\n"
	newlineChar = '\n'

	cacheInitCap = 128
)

const (
	NoFlag  = 0
	StdFlag = log.LstdFlags
)

var asyncLog *asyncLogType

var nowFunc = time.Now

func init() {
	asyncLog = &asyncLogType{
		files: make(map[string]*LogFile),
	}

	//timer := time.NewTicker(lf.sync.frequency)
	timer := time.NewTicker(time.Second * 1)
	go func() {
		for {
			select {
			case <-timer.C:
				//now := nowFunc()
				asyncLog.RLock()
				for _, file := range asyncLog.files {
					if file.sync.status != statusDoing {
						go file.Flush()
					}
				}
				asyncLog.RUnlock()
			}
		}

	}()
}

func NewFileLogger(filename string) *LogFile {
	asyncLog.Lock()
	defer asyncLog.Unlock()

	if lf, ok := asyncLog.files[filename]; ok {
		return lf
	}

	lf := &LogFile{
		filename: filename,
		flag:     StdFlag,
	}

	asyncLog.files[filename] = lf

	lf.logRotate.rotate = RotateDate

	lf.cache.enable = true

	lf.sync.frequency = time.Second

	return lf
}


// save to file
func (lf *LogFile) Flush() error {
	lf.sync.status = statusDoing
	defer func() {
		lf.sync.status = statusDone
	}()

	file, err := lf.openFileNoCache()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lf.cache.mutex.Lock()
	cache := lf.cache.data
	lf.cache.data = make([]string, 0, cacheInitCap)
	lf.cache.mutex.Unlock()

	if len(cache) == 0 {
		return nil
	}

	for i := 1; i <= 3; i++ {
		_, err = file.WriteString(strings.Join(cache, ""))
		if err != nil {
			if i == 3 {
				panic(err)
			}

		}
		break
	}

	return nil
}

//================== private functions =======================

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

func (lf *LogFile) getFilenameSuffix() string {
	if lf.logRotate.rotate == RotateDate {
		return nowFunc().Format("20060102")
	}
	return nowFunc().Format("2006010215")
}

func (lf *LogFile) directWrite(msg []byte) error {
	file, err := lf.openFile()
	//file, err := lf.openFileNoCache()
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lf.logRotate.mutex.Lock()
	_, err = file.Write(msg)
	lf.logRotate.mutex.Unlock()

	return err
}

func (lf *LogFile) openFile() (*os.File, error) {
	suffix := lf.getFilenameSuffix()

	lf.logRotate.mutex.Lock()
	defer lf.logRotate.mutex.Unlock()

	if suffix == lf.logRotate.suffix {
		return lf.logRotate.file, nil
	}

	logFilename := lf.filename + "." + suffix
	file, err := os.OpenFile(logFilename, fileFlag, fileOpenMode)
	if err != nil {
		panic(err)
	}

	if lf.logRotate.file != nil {
		lf.logRotate.file.Close()
	}

	lf.logRotate.file = file
	lf.logRotate.suffix = suffix
	return file, nil
}

func (lf *LogFile) openFileNoCache() (*os.File, error) {
	logFilename := lf.filename + "." + lf.getFilenameSuffix()

	lf.logRotate.mutex.Lock()
	defer lf.logRotate.mutex.Unlock()

	file, err := os.OpenFile(logFilename, fileFlag, fileOpenMode)
	if err != nil {
		return file, err
	}

	return file, nil
}
