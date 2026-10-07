//package types
package main

import (
	pb "unified-logger/proto"
)

type server struct {
	pb.UnimplementedMessageServiceServer
	msg_chan chan string
}
