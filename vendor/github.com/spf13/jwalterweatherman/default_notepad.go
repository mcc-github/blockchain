




package jwalterweatherman

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	TRACE    *log.Logger
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger
	FATAL    *log.Logger

	LOG      *log.Logger
	FEEDBACK *Feedback

	defaultNotepad *Notepad
)

func reloadDefaultNotepad() {
	TRACE = defaultNotepad.TRACE
	DEBUG = defaultNotepad.DEBUG
	INFO = defaultNotepad.INFO
	WARN = defaultNotepad.WARN
	ERROR = defaultNotepad.ERROR
	CRITICAL = defaultNotepad.CRITICAL
	FATAL = defaultNotepad.FATAL

	LOG = defaultNotepad.LOG
	FEEDBACK = defaultNotepad.FEEDBACK
}

func init() {
	defaultNotepad = NewNotepad(LevelError, LevelWarn, os.Stdout, ioutil.Discard, "", log.Ldate|log.Ltime)
	reloadDefaultNotepad()
}


func SetLogThreshold(threshold Threshold) {
	defaultNotepad.SetLogThreshold(threshold)
	reloadDefaultNotepad()
}


func SetLogOutput(handle io.Writer) {
	defaultNotepad.SetLogOutput(handle)
	reloadDefaultNotepad()
}



func SetStdoutThreshold(threshold Threshold) {
	defaultNotepad.SetStdoutThreshold(threshold)
	reloadDefaultNotepad()
}


func SetPrefix(prefix string) {
	defaultNotepad.SetPrefix(prefix)
	reloadDefaultNotepad()
}


func SetFlags(flags int) {
	defaultNotepad.SetFlags(flags)
	reloadDefaultNotepad()
}


func LogThreshold() Threshold {
	return defaultNotepad.logThreshold
}


func StdoutThreshold() Threshold {
	return defaultNotepad.stdoutThreshold
}


func GetLogThreshold() Threshold {
	return defaultNotepad.GetLogThreshold()
}


func GetStdoutThreshold() Threshold {
	return defaultNotepad.GetStdoutThreshold()
}


func LogCountForLevel(l Threshold) uint64 {
	return defaultNotepad.LogCountForLevel(l)
}



func LogCountForLevelsGreaterThanorEqualTo(threshold Threshold) uint64 {
	return defaultNotepad.LogCountForLevelsGreaterThanorEqualTo(threshold)
}


func ResetLogCounters() {
	defaultNotepad.ResetLogCounters()
}
