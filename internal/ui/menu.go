package ui

import (
	"fmt"
	"time"
)

// Spinner runs a braille spinner until `done` is closed.
func Spinner(msg string, done <-chan struct{}) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	for {
		select {
		case <-done:
			fmt.Print("\r\033[K")
			return
		default:
			fmt.Printf("\r  %s %s", Green(frames[i%len(frames)]), DimCyan(msg))
			i++
			time.Sleep(70 * time.Millisecond)
		}
	}
}

// Transition shows a brief spinner for the given duration.
func Transition(msg string, d time.Duration) {
	done := make(chan struct{})
	go Spinner(msg, done)
	time.Sleep(d)
	close(done)
}