# go-logger
a simple golang file logger, support async save to file

#Usage
```go
	package main

	import (
		"github.com/ikool-cn/go-logger"
	)

	func main() {
		log := goLogger.NewFileLogger("app.log")
		defer log.Flush()
		log.SetLevel(goLogger.WARNING)
		log.Debug("Debug")
		log.Info("Info")
		log.Notice("Notice")
		log.Warning("Warning")
		log.Error("Error")

		//output
		//2017-03-11 23:31:43 [WARNING] Warning
		//2017-03-11 23:31:43 [ERROR] Error
	}
```