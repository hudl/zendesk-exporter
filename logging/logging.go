package logging

import (
	stdlog "log"
	"os"

	golog "github.com/op/go-logging"
)

// Due to go-logging not having a "Verbose" level, we're
// bumping the lower levels down by 1
// Debug -> Verbose
// Info -> Debug
// Notice -> Info
// Warn, Error, and Critical stay the same
var defaultLevel golog.Level = golog.DEBUG

func Configure() {
	golog.SetFormatter(golog.MustStringFormatter("[0x%{id:x}] [%{level}] [%{module}] %{message}"))
	stdoutLogBackend := golog.NewLogBackend(os.Stdout, "", stdlog.LstdFlags|stdlog.Lshortfile)
	stdoutLogBackend.Color = true
	golog.SetLevel(defaultLevel, "")

	// NOTE these file permissions are restricted by umask, so they probably won't work right.
	err := os.MkdirAll("./log", 0775)
	if err != nil {
		panic(err)
	}
	logFile, err := os.OpenFile("./log/zd-exporter.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	if err != nil {
		panic(err)
	}

	fileLogBackend := golog.NewLogBackend(logFile, "", stdlog.LstdFlags|stdlog.Lshortfile)
	fileLogBackend.Color = false

	golog.SetBackend(stdoutLogBackend, fileLogBackend)
}
