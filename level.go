package goLogger

// 日志级别
const (
	All       = iota
	DEBUG
	INFO
	NOTICE
	WARNING
	ERROR
	CRITICAL
	ALERT
	EMERGENCY
)

// 日志级别标识
var (
	levelTitle = map[int]string{
		DEBUG:     "[DEBUG]",
		INFO:      "[INFO]",
		NOTICE:    "[NOTICE]",
		WARNING:   "[WARNING]",
		ERROR:     "[ERROR]",
		CRITICAL:  "[CRITICAL]",
		ALERT:     "[ALERT]",
		EMERGENCY: "[EMERGENCY]",
	}
)
