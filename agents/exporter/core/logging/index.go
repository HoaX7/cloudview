package logging

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

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

func setLogFile(fn string) {
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}

func Log(v ...any) {
	setLogFile("usage.log")
	writeToFile(v...)
}

func Error(v ...any) {
	setLogFile("error.log")
	writeToFile(v...)
}

func writeToFile(v ...any) {
	log.Println(v...)
}
