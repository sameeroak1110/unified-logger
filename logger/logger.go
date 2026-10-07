package logger

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sameeroak1110/unified-logger/proto"
	/* "unified-logger/data"
	"unified-logger/types"
	"unified-logger/logger" */
)


// NewConnector connects to the gRPC server and opens the persistent stream
func NewLogger(loglevel uint8, target string) (*Logger, error) {
	// 1. Establish the gRPC connection
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	// 2. Instantiate the generated gRPC client
	pbClient := pb.NewMessageServiceClient(conn)

	// 3. Open the long-lived client stream using context.Background()
	ctxParent := context.Background()
	stream, err := pbClient.SendMessageStream(ctxParent)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}

	return &Logger {
		conn:   conn,
		stream: stream,
		log_level: loglevel,
	}, nil
}


// Send transmits a message over the persistent stream
//func (c *Logger) Send(level uint8, msg string) error {
func (c *Logger) Send(msg string) error {
	req := &pb.MessageRequest{
		Content: msg, // Ensures the 'Content' field in your .proto is populated
	}

	return c.stream.Send(req)
}


// Close gracefully terminates the stream and underlying network connection
func (c *Logger) Close() error {
	if c.stream != nil {
		// Signal to server that sending is finished and await server summary response
		summary, err := c.stream.CloseAndRecv()
		if err != nil {
			log.Printf("Error closing stream: %v", err)
		} else {
			log.Printf("Stream closed cleanly. Server queued %d messages.", summary.GetTotalQueued())
		}
	}

	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}
