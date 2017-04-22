package util

import "github.com/jcelliott/lumber"
import "os"

var log *lumber.FileLogger

func init() {
	var err error
	log, err = lumber.NewFileLogger("log.txt", lumber.INFO, lumber.APPEND, 10000, 10000, 100)
	if err != nil {
		os.Exit(1) //TODO: Make the program not just crash here.
	}
}

//Logs an error to log.txt. TODO: Functions for more levels. Lumber supports 6 levels.
func LogError(m string) {
	log.Error(m)
}
