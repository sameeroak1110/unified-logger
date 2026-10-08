//package main
//package types
package logger

import (
	"google.golang.org/grpc"

	pb "github.com/sameeroak1110/unified-logger/proto"
)


// Logger wraps the gRPC connection and stream instance
type Logger struct {
	conn   *grpc.ClientConn
	stream pb.MessageService_SendMessageStreamClient
	log_level uint8
}

type log_msg struct {
	str   string
	color string
	wt    uint8
}
