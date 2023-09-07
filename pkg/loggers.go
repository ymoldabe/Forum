package pkg

import (
	"log"
	"os"
)

type Loggers struct {
	ErrorLog log.Logger
	InfoLog  log.Logger
}

var (
	InfoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
)
