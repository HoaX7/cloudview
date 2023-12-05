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

func (l LoggerStruct) Log(v ...any) {
	log.Println(l.Name, prepareLog(v))
}

func (l LoggerStruct) Error(v ...any) {
	l.Log(l.Name, "Error", v)
}

var Logger = LoggerStruct{Name: "[log]"}
