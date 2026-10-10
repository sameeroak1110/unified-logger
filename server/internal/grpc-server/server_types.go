//package types
package main

import (
	pb "github.com/sameeroak1110/unified-logger/proto"
)

type server struct {
	pb.UnimplementedMessageServiceServer
	msg_chan chan string
}
