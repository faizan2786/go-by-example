package main

import (
	"errors"
	"fmt"
	"os"
)

type logger interface {
	Log(msg string) error
}

type multiLogger []logger

func (ml multiLogger) Log(msg string) error {
	var errs []error
	for _, l := range ml {
		if err := l.Log(msg); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...) // join multiple errors into one
}

type fileLogger struct {
	filePath string // full file path
}

// needs pointer receiver due to efficiency
func (fl *fileLogger) Log(msg string) error {

	// open file with write/create/append mode
	// 0644 give r/w permission to user and only read to others
	f, err := os.OpenFile(fl.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// add a newline if the message doesn't end with it
	if !(len(msg) > 0 && msg[len(msg)-1] == '\n') {
		msg += "\n"
	}

	_, err = f.WriteString(msg) // write to the file
	return err
}

type consoleLogger struct{}

func (cl consoleLogger) Log(msg string) error {
	_, err := fmt.Println(msg)
	return err
}

func main() {

	fl := fileLogger{"logs.txt"}
	cl := consoleLogger{}

	ml := multiLogger{
		&fl, // file logger has pointer receiver for the interface method (Log)
		cl,  // console logger has value receiver so doesn't need pointer (though pointer would work too)
	}

	err := ml.Log("This is a log")
	if err != nil {
		fmt.Println(err)
	}

	err = ml.Log("This is another log")
	if err != nil {
		fmt.Println(err)
	}
}
