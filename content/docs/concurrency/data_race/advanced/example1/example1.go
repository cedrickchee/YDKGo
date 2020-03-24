// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show a more complicated race condition using
// an interface value. This produces a read to an interface value after
// a partial write.
package main

import (
	"fmt"
	"os"
	"sync"
)

// Speaker allows for speaking behavior.
type Speaker interface {
	Speak() bool
}

// Ben is a person who can speak.
type Ben struct {
	name string
}

// Speak allows Ben to say hello. It returns false if the method is
// called through the interface value after a partial write.
func (b *Ben) Speak() bool {
	if b.name != "Ben" {
		fmt.Printf("Ben says, \"Hello my name is %s\"\n", b.name)
		return false
	}

	return true
}

// Jerry is a person who can speak.
type Jerry struct {
	name string
}

// Speak allows Jerry to say hello. It returns false if the method is
// called through the interface value after a partial write.
func (j *Jerry) Speak() bool {
	if j.name != "Jerry" {
		fmt.Printf("Jerry says, \"Hello my name is %s\"\n", j.name)
		return false
	}

	return true
}

func main() {

	// Create values of type Ben and Jerry.
	ben := Ben{"Ben"}
	jerry := Jerry{"Jerry"}

	// Assign the pointer to the Ben value to the interface value.
	person := Speaker(&ben)

	// Have a goroutine constantly assign the pointer of
	// the Ben value to the interface and then Speak.
	go func() {
		for {
			person = &ben
			if !person.Speak() {
				os.Exit(1)
			}
		}
	}()

	// Have a goroutine constantly assign the pointer of
	// the Jerry value to the interface and then Speak.
	go func() {
		for {
			person = &jerry
			if !person.Speak() {
				os.Exit(1)
			}
		}
	}()

	// Just hold main from returning. The data race will
	// cause the program to exit.
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

/*
$ go build -race

$ ./example1
==================
WARNING: DATA RACE
Write at 0x00c000012220 by goroutine 8:
  main.main.func2()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/advanced/example1/example1.go:78 +0x3c

Previous write at 0x00c000012220 by goroutine 7:
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/advanced/example1/example1.go:67 +0x3c

Goroutine 8 (running) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/advanced/example1/example1.go:76 +0x152

Goroutine 7 (running) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/advanced/example1/example1.go:65 +0x126
==================
Jerry says, "Hello my name is Ben"
*/
