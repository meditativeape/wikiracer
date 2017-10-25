package util

import (
	"bufio"
	"log"
	"os"
	"path"
)

const LogDir string = "log"
const LogFilename string = "service.log"

var Logger *log.Logger = initLogger()

func initLogger() *log.Logger {
	err := os.MkdirAll(LogDir, os.ModePerm|os.ModeDir)
	panicIfError(err)
	f, err := os.Create(path.Join(LogDir, LogFilename))
	panicIfError(err)
	writer := bufio.NewWriter(f)
	return log.New(writer, "", log.Ltime|log.Lmicroseconds)
}
