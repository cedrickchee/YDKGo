// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package logger

import (
	"fmt"
	"io"
	"sync"
)

// Logger ...
type Logger struct {
	ch chan string
	wg sync.WaitGroup
}

// New is a factory function that will return a logger.
// We should be using pointer semantics when you can't make copies of loggers
// anymore because we can't make a copy of the WaitGroup, that would create a
// different WaitGroup. Plus, there's only one logger. We only want one logger.
// w is what device we want to write to.
// cap is capacity of our buffer and only the caller can tell us.
func New(w io.Writer, cap int) *Logger {
	// Using value semantic on construction because we are assigning it to
	// variable
	l := Logger{
		ch: make(chan string, cap),
	}

	l.wg.Add(1)

	// Goroutine job is to perform these writes.
	go func() {
		defer l.wg.Done()
		for v := range l.ch {
			// Write whatever data we receive off the channel to the device
			// that we passed that will do it through the fmt.Fprint function.
			fmt.Fprintln(w, v)
		}
	}()

	// Escape analysis: potential allocation
	return &l
}

// Close give our API the ability to clean shutdown the Goroutine in order
// of fashion in our factory function.
func (l *Logger) Close() {
	// Terminate the `for` loop in the Goroutine which will the call
	// WaitGroup.Done()
	close(l.ch)
	l.wg.Wait()
}

// Println prints to the logger
func (l *Logger) Println(v string) {
	select {
	case l.ch <- v:
	default:
		// If line 56 is going to block there's no room in the ch buffer, then
		// we want to drop, dropping these logs.
		fmt.Println("DROP") // Print so we can see we no longer blocking
	}
}
