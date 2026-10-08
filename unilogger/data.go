//package main
//package data
package logger

import (
	//"google.golang.org/grpc"

	//pb "unified-logger/proto"
	//"unified-logger/types"
)

//var conn *grpc.ClientConn
//var client *Client
//var stream pb.MessageService_SendMessageStreamClient

// log levels
/* const DBGRM string = "DBGRM"      // green
const DEBUG string = "DEBUG"      // normal
const INFO string = "INFO"        // normal
const WARNING string = "WARNING"  // yellow
const ERROR string = "ERROR"      // red */

//var loglevel uint8
const (
	DBGRM uint8 = iota
	DEBUG
	INFO
	WARNING
	ERROR
)

var str_DBGRM string = "DBGRM"
var str_DEBUG string = "DEBUG"
var str_INFO string = "INFO"
var str_WARNING string = "WARNING"
var str_ERROR string = "ERROR"

const color_nornal string = "\033[0m"
const color_dbgrm_green string = "\033[32m"
const color_warn_yellow string = "\033[33m"
const color_error_red string = "\033[31m"

var map_loglevel = map[uint8]log_msg {
	DBGRM: log_msg {
		str: "DBGRM",
		color: color_dbgrm_green,
		wt: 0,
	},
	DEBUG: log_msg {
		str: "DEBUG",
		color: color_nornal,
		wt: 1,
	},
	INFO: log_msg {
		str: "INFO",
		color: color_nornal,
		wt: 2,
	},
	WARNING: log_msg {
		str: "WARNING",
		color: color_warn_yellow,
		wt: 3,
	},
	ERROR: log_msg {
		str: "ERROR",
		color: color_error_red,
		wt: 4,
	},
}
