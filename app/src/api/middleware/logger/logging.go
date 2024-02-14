package logger

import (
	"bytes"
	"fmt"
	"log"
)

type LoggerStruct struct {
	Name string
}

func prepareLog(v ...any) string {
	var logOutput bytes.Buffer
	for i, value := range v {
		if i > 0 {
			logOutput.WriteString(" ")
		}
		logOutput.WriteString(fmt.Sprintf("%v", value))
	}
	return logOutput.String()
}
func (l *LoggerStruct) SetName(name string) *LoggerStruct {
	l.Name = name
	return l
}

func (l *LoggerStruct) Log(v ...any) {
	logStr := append([]any{l.Name + ":"}, v...)
	log.Println(logStr...)
}

func (l *LoggerStruct) Error(v ...any) {
	logStr := append([]any{"Error:"}, v...)
	l.Log(logStr...)
}

var Logger = LoggerStruct{Name: "[log]"}

func NewLogger() *LoggerStruct {
	return new(LoggerStruct)
}
