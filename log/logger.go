package log

import (
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"time"
	"github.com/zbindenren/logrus_mail"
	"apollo/config"
)
var Logger,MailLogger *logrus.Logger

func init()  {
	Logger = NewLogger()
	MailLogger = NewMailLogger()
}

func NewLogger() *logrus.Logger {
	if Logger != nil {
		return Logger
	}

	dateStr := time.Now().Format("2006-01-02")
	appName := config.C.GetString("app_name")
	path := config.C.GetString("log_path") + appName + "/" + dateStr + "/"
	pathMap := lfshook.PathMap {
		logrus.InfoLevel:  path + "info.log",
		logrus.ErrorLevel: path + "error.log",
	}

	Logger = logrus.New()
	Logger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))

	return Logger
}

func NewMailLogger() *logrus.Logger {
	if MailLogger != nil {
		return MailLogger
	}

	MailLogger = logrus.New()
	hook, err := logrus_mail.NewMailHook("system_manager", "smtp.gmail.com", 587, "hcj8080@gmail.com", "hcj8080@gmail.com")

	if err != nil {
		MailLogger.Hooks.Add(hook)
	}

	return MailLogger
}

