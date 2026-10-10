package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"sync/atomic"

	"google.golang.org/grpc"

	pb "github.com/sameeroak1110/unified-logger/proto"

	"github.com/sameeroak1110/unified-logger/server/worker"
)

/* type server struct {
	pb.UnimplementedMessageServiceServer
	msg_chan chan string
} */


func main() {
	// 1. Initialize buffered channel (capacity of 100 messages)
	msg_chan := make(chan string, 100)
	var wg sync.WaitGroup

	// 2. Create a root cancellation context for the worker pipeline
	workerCtx, cancelWorker := context.WithCancel(context.Background())
	defer cancelWorker()

	// 3. Launch background file worker
	wg.Add(1)
	go worker.Worker(workerCtx, msg_chan, "output.txt", &wg)

	// 4. Set up TCP listener for gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("ERROR: Failed to listen: %s\n", err.Error())
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessageServiceServer(grpcServer, &server{msg_chan: msg_chan})

	// 5. Register OS signal listener for graceful shutdown (SIGINT, SIGTERM)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// 6. Start gRPC server in background goroutine
	go func() {
		fmt.Println("Server running on port :50051...")
		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			fmt.Printf("ERROR: Failed to start server: %s\n", err.Error())
			return
		}
	}()

	// 7. Block main until termination signal is captured
	<-stop
	fmt.Println("Shutdown signal received. Starting graceful cleanup.")

	// 8. Gracefully stop gRPC server (reject new streams, allow active RPCs to wrap up)
	grpcShutdownCtx, cancelGrpc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelGrpc()

	stoppedChan := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stoppedChan)
	}()

	select {
		case <-stoppedChan:
			fmt.Println("gRPC server stopped cleanly.")
		case <-grpcShutdownCtx.Done():
			fmt.Println("Graceful stop timed out. Forcing gRPC server stop.")
			grpcServer.Stop()
	}

	// 9. Close channel to signal worker range loop that no more messages are incoming
	close(msg_chan)

	// 10. Wait for worker to flush queued messages to disk within a deadline
	fmt.Println("Waiting for file worker to flush remaining queue.")
	flushTimeout := 5 * time.Second
	workerDone := make(chan struct{})

	go func() {
		wg.Wait()
		close(workerDone)
	}()

	select {
		case <-workerDone:
			fmt.Println("Worker successfully flushed all remaining messages to disk.")
		case <-time.After(flushTimeout):
			fmt.Println("Flush deadline exceeded! Force-cancelling worker context.")
			cancelWorker()
			<-workerDone
	}

	fmt.Println("Server process exited gracefully.")
}


func (s *server) SendMessageStream_1(stream pb.MessageService_SendMessageStreamServer) error {
	var totalReceived, totalQueued int32

	for {
		// Context cancellation or timeout will automatically interrupt blocked Recv/Sends
		/* select {
			case <-stream.Context().Done():
				fmt.Println("Stream context cancelled during shutdown or client disconnect.")
				return stream.Context().Err()
			default:
		} */

		req, err := stream.Recv()
		/* if req != nil {
			fmt.Printf("DBGRM: message received: %s\n", req.Content)
		} */
		//fmt.Printf("DBGRM: message received: %s\n", req.Content)
		fmt.Printf("DBGRM: message received: %s\n", req.GetContent())
		if err == io.EOF {
			return stream.SendAndClose(&pb.MessageSummary{
				TotalReceived: totalReceived,
				TotalQueued:   totalQueued,
			})
		}
		if err != nil {
			return err
		}

		totalReceived++

		// Explicit Backpressure:
		// Blocks when channel is full, or cancels if stream Context terminates.
		select {
			case s.msg_chan <- req.GetContent():
				totalQueued++
			case <-stream.Context().Done():
				// Respect context cancellation if buffer is full and server is shutting down
				return stream.Context().Err()
		}
	}
}


// SendMessageStream receives incoming messages and immediately enqueues them onto s.msg_chan
func (s *server) SendMessageStream_2(stream pb.MessageService_SendMessageStreamServer) error {
	var totalReceived, totalQueued int32

	for {
		// 1. Read message off the gRPC network stream
		req, err := stream.Recv()
		if err == io.EOF {
			// Client completed the stream; send final summary back
			return stream.SendAndClose(&pb.MessageSummary{
				TotalReceived: totalReceived,
				TotalQueued:   totalQueued,
			})
		}
		if err != nil {
			return err
		}

		totalReceived++

		// 2. IMMEDIATELY add to the buffered channel.
		// Blocking send ensures zero latency delay, while enforcing backpressure if full.
		select {
			case s.msg_chan <- req.GetContent():
				//data.Total_enqueued++
				Total_enqueued++

			case <-stream.Context().Done():
				// Cancels if stream context ends while waiting for buffer space
				return stream.Context().Err()
		}
	}
}


func (s *server) SendMessageStream(stream pb.MessageService_SendMessageStreamServer) error {
	fmt.Println("DBGRM: Client stream connecsted")

	for {
		// 1. Blocks HERE until this specific client transmits a message over TCP.
		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("<<< Client closed stream cleanly.")
			return stream.SendAndClose(&pb.MessageSummary{})
		}

		fmt.Printf("DBGRM: message received: %s\n", req.GetContent())

		// If the client network connection drops or terminates:
		if err != nil {
			fmt.Printf("WARNING: Client disconnected: %s\n", err.Error())
			//return err
		}
		fmt.Printf("[RECV] Got payload from client: %q\n", req.GetContent())

		atomic.AddUint64(&Total_received, 1)
		atomic.AddUint64(&Total_bytes_received, 1)

		// 2. The moment bytes arrive, immediately hand off to the shared channel
		select {
			case s.msg_chan <- req.GetContent():
				// Enqueued instantly!
				fmt.Println("[ENQUEUED] Pushed to msg_chan successfully")
				atomic.AddUint64(&Total_received, 1)
				atomic.AddUint64(&Total_bytes_enqueued, 1)
			case <-stream.Context().Done():
				// Respect server shutdown/context cancellation
			return stream.Context().Err()
		}
	}
}
