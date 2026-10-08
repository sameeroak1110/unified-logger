//package main
package logger

import (
	"fmt"
	"time"
	"runtime"
	"path"
	"strings"

	//pb "unified-logger/proto" // Replace with your actual module path
	//"unified-logger/data"
)


func get_log_msg(loglevel uint8, msg string, args ...interface{}) string {
	str_loglevel := map_loglevel[loglevel]

	t := time.Now()
	zonename, _ := t.In(time.Local).Zone()
	msg_ts := fmt.Sprintf("%02d-%02d-%d-%02d%02d%02d-%06d-%s", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), zonename)
	//_pc, fn, line, _ := runtime.Caller(0)
	//_pc, fn, line, _ := runtime.Caller(2)
	_pc, fn, line, _ := runtime.Caller(1)

	tmp1 := strings.Split((runtime.FuncForPC(_pc).Name()), ".")
	pkgname := tmp1[0]
	src_file := pkgname + "/" + path.Base(fn)
	func_name := tmp1[1]

	msg_prefix := ""
	if loglevel == DBGRM {
		msg_prefix = "#### "
	}

	log_msg := fmt.Sprintf("[%s] [%s] [%s +%d]@[%s]:\n", msg_ts, str_loglevel.str, src_file, line, func_name)
	log_msg = fmt.Sprintf(log_msg + msg, args...)
	log_msg = msg_prefix + log_msg + "\n"

	return log_msg
}

func (pc *Logger) Dbgrm(msg string, args ...interface{}) error {
	log_msg := get_log_msg(DBGRM, msg, args)

	return pc.Send(log_msg)
}

func (pc *Logger) Debug(msg string, args ...interface{}) error {
	if DEBUG < pc.log_level {
		return nil
	}

	log_msg := get_log_msg(DEBUG, msg, args)
	return pc.Send(log_msg)
}


func (pc *Logger) Info(msg string, args ...interface{}) error {
	if INFO < pc.log_level {
		return nil
	}

	log_msg := get_log_msg(INFO, msg, args)
	return pc.Send(log_msg)
}


func (pc *Logger) Warning(msg string, args ...interface{}) error {
	if WARNING < pc.log_level {
		return nil
	}

	log_msg := get_log_msg(WARNING, msg, args)
	return pc.Send(log_msg)
}


func (pc *Logger) Error(msg string, args ...interface{}) error {
	if ERROR < pc.log_level {
		return nil
	}

	log_msg := get_log_msg(ERROR, msg, args)
	return pc.Send(log_msg)

	return nil
}
