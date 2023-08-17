package pkg

import (
	"fmt"
	filename "github.com/keepeye/logrus-filename"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"io"
	"log/syslog"
	"os"
	"path"
	"time"
)

const (
	logFileName = "source-connection.log"
)

var (
	Log     *logrus.Logger
	logFile *os.File
)

func InitLogger() {
	Log = logrus.New()
	initSyslog()
	filenameHook := filename.NewHook()
	filenameHook.Field = "file"
	Log.AddHook(filenameHook)

	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   true,
		FullTimestamp:   true,
	})

	Log.Debugf("init with args %s", os.Args)
}

func initSyslog() {
	if hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, ""); err != nil {
		Log.Error("Unable to connect to local syslog daemon")
	} else {
		Log.AddHook(hook)
	}
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

func init() {
	InitLogger()
	file := path.Join("../log", logFileName)
	logWriter, err := rotatelogs.New(
		file+".%Y-%m-%d_%H-%M-%S",
		rotatelogs.WithLinkName(file),
		rotatelogs.WithRotationTime(time.Hour*time.Duration(1)),
		rotatelogs.WithMaxAge(time.Hour*time.Duration(24*7)),
	)

	if err != nil {
		fmt.Println("Failed to init log file settings..." + err.Error())
		Log.Infof("Failed to log to file, using default stderr.")
	} else {
		mw := io.MultiWriter(os.Stdout, logWriter)
		Log.SetOutput(mw)
	}
}
