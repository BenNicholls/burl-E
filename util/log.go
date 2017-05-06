package util

import "github.com/jcelliott/lumber"
import "os"

var log *lumber.FileLogger

//TODO: Logfile checking and cleanup? Could archive old logs or logs that get too big and
//make a new one??
func init() {
	var err error
	log, err = lumber.NewFileLogger("log.txt", lumber.INFO, lumber.APPEND, 10000, 10000, 100)
	if err != nil {
		//Who logs the errors from the error logger??? -_-
		os.Exit(1) //TODO: Make the program not just crash here.
	}
}

//Logs an error to log.txt. TODO: Functions for more levels. Lumber supports 6 levels.
func LogError(m string) {
	log.Error(m)
}
