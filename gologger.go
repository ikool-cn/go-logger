package goLogger

func (lf *LogFile) Debug(msg string) error {
	return lf.log(msg, DEBUG)
}

func (lf *LogFile) Info(msg string) error {
	return lf.log(msg, INFO)
}

func (lf *LogFile) Notice(msg string) error {
	return lf.log(msg, NOTICE)
}

func (lf *LogFile) Warning(msg string) error {
	return lf.log(msg, WARNING)
}

func (lf *LogFile) Error(msg string) error {
	return lf.log(msg, ERROR)
}

func (lf *LogFile) Critical(msg string) error {
	return lf.log(msg, CRITICAL)
}

func (lf *LogFile) Alert(msg string) error {
	return lf.log(msg, ALERT)
}

func (lf *LogFile) Emergency(msg string) error {
	return lf.log(msg, EMERGENCY)
}

func (lf *LogFile) SetLevel(logLevel int) {
	lf.level = logLevel
}

func (lf *LogFile) log(msg string, level int) error {
	if level >= lf.level {
		return lf.write(levelTitle[level] + " " + msg)
	}

	return nil
}