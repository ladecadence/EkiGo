package logging

import (
	"os"
	"time"
)

const (
	LogData = iota
	LogInfo
	LogWarn
	LogError
	LogClean
)

type Logging struct {
	Filename string
}

func New(n string) error {
	l := Logging{}
	l.Filename = n + time.Now().Format(time.RFC3339) + ".log"
	f, err := os.Create(l.Filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

func (l *Logging) Log(logType int, msg string) error {
	f, err := os.OpenFile(l.Filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// write header
	switch logType {
	case LogData:
		_, err = f.WriteString("DATA::")
	case LogInfo:
		_, err = f.WriteString("INFO::")
	case LogWarn:
		_, err = f.WriteString("WARN::")
	case LogError:
		_, err = f.WriteString(" ERR::")
	default:
		{
		}
	}

	// and date
	if logType != LogClean {
		_, err = f.WriteString(time.Now().Format(time.RFC3339))
		_, err = f.WriteString(":: ")
	}

	// now the data
	_, err = f.WriteString(msg)
	_, err = f.WriteString("\n")
	err = f.Sync()

	// if there was any error writing, return it
	return err
}
