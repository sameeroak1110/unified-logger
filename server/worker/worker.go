package worker

import (
	"fmt"
	"sync"
	"os"
	"context"
)


// Accepts a Context to handle force-cancellation during shutdown
func Worker(ctx context.Context, msgChan <-chan string, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	//fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND|os.O_SYNC, 0644)
	if err != nil {
		fmt.Printf("Worker failed to open file: %s\n", err.Error())
		return
	}
	defer fp.Close()

	for {
		select {
			// Priority 1: Check if shutdown context timed out or was cancelled
			case <-ctx.Done():
				fmt.Printf("Warning: Worker context cancelled (%s). Flushed remaining logs and stopping.\n", ctx.Err().Error())
				// Optional: Drain any last-second items remaining without blocking
				drainChannel(msgChan, fp)
				return

			// Priority 2: Process items from the channel
			case msg, ok := <-msgChan:
				if !ok {
					// Channel was closed by main() and is completely empty
					fmt.Println("Worker processed all queued messages successfully.")
					return
				}
				//fmt.Printf("message received: %s\n", msg)
				if _, err := fp.WriteString(msg + "\n"); err != nil {
					fmt.Printf("ERROR: Failed to write message to file: %s\n", err.Error())
				}

				// Flushes to disk immediately.
				// Prevents OS buffer from holding onto short messages
				if err := fp.Sync(); err != nil {
					fmt.Printf("Failed to sync file: %s\n", err.Error())
				}
		}
	}
}


// Helper to quickly dump remaining buffered items if cancelled
func drainChannel(msgChan <-chan string, fp *os.File) {
	for len(msgChan) > 0 {
		msg := <-msgChan
		fp.WriteString(msg + "\n")
		fp.Sync()
	}
}
