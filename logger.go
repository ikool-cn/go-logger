package goLogger

// Logger ...
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Notice(msg string)
	Warning(msg string)
	Error(msg string)
	Critical(msg string)
	Alert(msg string)
	Emergency(msg string)
	SetLevel(level int)
	Flush()
}