---
title: "Context"
weight: 4
---

# Context

## Context - Part 1

### Store and retrieve values from a context

The Context package is an important package because it's where we will implement
cancellation and deadline.

**Value bag**

The other things that the Context package has is a generic value bag. I'm very
afraid of the generic value bag because you don't want to be using it for things
like local requests storage. You never want to hide data in a context, but
putting data in a context can help for things that functions don't need.

So the way the value bag works is, data is stored inside the context with a
combination of a key, and then the value. But the key isn't just some value,
it's also will be based on type.

Everything here is value semantics with Context.

[Sample program](example1/example1.go).

```go
import (
	"context"
	"fmt"
)

// TraceID is represents the trace id.
type TraceID string
```

```go
// TraceIDKey is the type of value to use for the key. The key is type specific
// and only values of the same type will match.
type TraceIDKey int
```

When we store a value inside a context, what getting stored is not just a value
but also a type associated with the storage. We can only pull a value out of
that context if we know the type of value that we are looking for. The idea of
this `TraceIDKey` type becomes really important when we want to store a value
inside the context.

```go
// Create a traceID for this request.
traceID := TraceID("f47ac10b-58cc-0372-8567-0e02b2c3d479")

// Declare a key with the value of zero of type TraceIDKey.
const traceIDKey TraceIDKey = 0

// Store the `traceID` value inside the context with a value of zero of
// type `TraceIDKey`.
ctx := context.WithValue(context.Background(), traceIDKey, traceID)
```

This program is starting with the empty context, we passed it into a function
called `WithValue`. WithValue's how you add state to a context. We are using
`context.WithValue` because a new context value and we want to initialize that
with data to begin with. Anytime we work with a context, the context has to have
a parent context. This is where the `Background` function comes in. We will
store the key `traceIDKey` to its value (which is 0 in this case), and value
of `traceID`.

```go
// Retrieve that traceID value from the Context value bag.
if uuid, ok := ctx.Value(traceIDKey).(TraceID); ok {
    fmt.Println("TraceID:", uuid)
}
```

`Value` allows us to pass the key of the corrected type (in our case is
`traceIDKey` of `TraceIDKey` type) and returns an empty interface. Because we
are working with an interface, we have to perform a type assertion to pull the
value that we store in there out the interface so we can work with the concrete
again.

```go
// Retrieve that traceID value from the Context value bag not
// using the proper key type.
if _, ok := ctx.Value(0).(TraceID); !ok {
    fmt.Println("TraceID Not Found")
}
```

Attempt to retrieve the value again using the same value but of a different type.
Even though the key value is 0, if we just pass 0 into this function call, we
are not going to get back that value to the `traceID` because 0 is based on
integer type, not our `TraceIDKey` type. It's important that when we store the
value inside the context to not use the built-in type. Declare our own key type.
That way, only us and who understand that type can pull that out. Because what
if multiple partial programs want to use that value of 0, we are all being
tripped up on each other. That type extends an extra level of protection on
being able to store and retrieve value out of context. If we are using this, we
want to raise a flag because we have to ask twice why do we want to do that
instead of passing down the call stack. Because if we can pass it down the call
stack, it would be much better for readability and maintainability for our
legacy code in the future.

### WithTimeout

The idea: We will have some Goroutine go off and do some work. We will be using
that cancellation pattern. Goroutine only has so much time to get the work done.

[Sample program](example4/example4.go) to show how to use the `WithTimeout`
function of the Context package.

```go
type data struct {
	UserID string
}

func main() {
	// Set a duration.
	duration := 150 * time.Millisecond

	// Create a context that is both manually cancellable and will signal
    // a cancel at the specified duration.
    // cancel is our handler function. You have to call cancel. It must be
    // called to clean up resources or you will leak.
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	// Create a channel to received a signal that work is done.
	ch := make(chan data, 1)

	// Ask the goroutine to do some work for us.
	go func() {
		// Simulate work.
		time.Sleep(50 * time.Millisecond)

		// Report the work is done.
		ch <- data{"123"}
	}()

    // Wait for the work to finish.
    // If the Done function returns before the channel receive `d := <-ch`,
    // we've timed out. If it takes too long (more than 150 ms) move on.
	select {
    case d := <-ch:
        // Channel receive
		fmt.Println("work complete", d)

    case <-ctx.Done():
        // We're on the receive of Done function on the context.
		fmt.Println("work cancelled")
    }
}
```

## Context - Part 2

### Request/Response Context Timeout

A practical example using a standard library's http and net packages on how we
can do cancellations.

[Sample program](example5/example5.go).

```go
// Sample program that implements a web request with a
// context that is used to timeout the request if it takes too long.

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Create a new request.
	req, err := http.NewRequest("GET", "https://www.ardanlabs.com/blog/post/index.xml", nil)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	// Create a context with a timeout of 50 milliseconds.
	ctx, cancel := context.WithTimeout(req.Context(), 50*time.Millisecond)
	defer cancel()

    // Declare a new transport and client for the call.
    // Transport is basically a pool of socket connections to some sort of
    // resource.
    // This is also where you will set up your timeouts. These are the default
    // values for transport, but we at some point have to come back and measure it.
	tr := http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

    // The client wraps our transport and it gives us an API to be able to go
    // and use the transport to access all these resources.
	client := http.Client{
		Transport: &tr,
    }

    // If we allow our main thread to do this work, we'll never be able to
    // cancel it.

    // Make the web call in a separate Goroutine so it can be cancelled.
	ch := make(chan error, 1)
	go func() {
        // Second path of execution

		log.Println("Starting Request")

		// Make the web call and return any error.
        // client.Do is going out and trying to hit the request URL.
        // It's probably blocked right now because it will need to wait for the
        // entire document to comeback.
		resp, err := client.Do(req)

        // It the error occurs, we perform a send on the channel to report that
        // we are done. We are going to use this channel at some point to report
        // back what is happening.
		if err != nil {
			ch <- err
			return
		}

        // If it doesn't fail, we close the response body on the return or we're
        // going to leak memory.
		defer resp.Body.Close()

		// Write the response to stdout.
		io.Copy(os.Stdout, resp.Body)

		// Then send back the nil instead of error.
		ch <- nil
    }()

    // Wait the request or timeout.
    // We perform a receive on ctx.Done saying that we want to wait 50 ms for
    // that whole process above to happen. If it doesn't, we signal back to that
    // Goroutine to cancel the sending request. We don't have to just walk away
    // and let that eat up resources and finish because we are not going to need it.
    // We are able to call CancelRequest and underneath, we are able to kill
	// that connection.
	select {
	case <-ctx.Done(): // use the select and start the clock for our 50 ms.
		log.Println("timeout, cancel work...")

        // Call transport and cancel the request and wait for it to complete.
        tr.CancelRequest(req)
        // Basically, we're able to shut down that Goroutine immediately and
        // waiting to validate that the Goroutine did shut down.
		log.Println(<-ch)
	case err := <-ch:
		if err != nil {
			log.Println(err)
		}
	}
}
```

## Failure Detection

Example: we have an application that is streaming data for any particular device, let's say, the World Cup
We got a socket connection in our server and we're streaming and we're logging as that happens.
Since we're letting every Goroutine write the standard out to do their own logs, what happens if standard out suddenly causes blocking?
Every Goroutine is going to be stopped, in a waiting state. Basically, there is no more video streaming.

Let's simulate a situation where the logging blocks and it causes our entire service to deadlock all at one time.

[Sample program](../patterns/advanced/main.go).

```go
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

// device allows us to mock a device we write logs to.
type device struct {
	mu      sync.RWMutex
	problem bool
}

// Write implements the io.Writer interface.
func (d *device) Write(p []byte) (n int, err error) {
	// Simulate disk problems.
	for d.isProblem() {
		time.Sleep(time.Second)
	}

	fmt.Print(string(p))
	return len(p), nil
}

// isProblem checks in a safe way if there is a problem.
func (d *device) isProblem() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.problem
}

// flipProblem reverses the problem flag to the opposite value.
func (d *device) flipProblem() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.problem = !d.problem
}

func main() {
	// Number of goroutines that will be writing logs.
	const grs = 10

	// Create a logger value with a buffer of capacity
	// for each goroutine that will be logging.
	var d device
	l := log.New(&d, "prefix", 0)

	// Generate goroutines, each writing to disk.
	for i := 0; i < grs; i++ {
		go func(id int) {
			for {
				l.Println(fmt.Sprintf("%d: log data", id))
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// We want to control the simulated disk blocking. Capture
	// interrupt signals to toggle device issues. Use <ctrl> z
	// to kill the program.

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		<-sigChan
		d.flipProblem()
	}
}
```

We can't stop streaming because we can't write logs, not in this case. So using this standard library logger and allowing Goroutines to log to the device directly is not a solution we can use anymore. We now have to add some extra complexity to deal with this problem. Basically, we've gotta write some code now that lets Goroutines even up to 10,000, write to the log, but to be able to detect if that log write is going to block. If it is, don't let the Goroutine block. Just bypass or skip the logs.

How can we do this? If we use a drop pattern, then I think we could come up with a very simple solution that solves the problem.

```
//       G
//       |   buf channel  +--G
//  +----|   +---------+  |--G
//  |    |---|   10K   |--|--G
//  x    |   +---------+  |--G
//  |    |                +--G
//  DB   |
```

Why does the buffer pattern work for us? If we use a drop pattern, then life is amazing here.
While there's room in the buffer, these Goroutines can perform their send and there's no latency between the send and the receive.
There will be some latency between multiple sends, but that should be quick.
If there's room in the buffer, that means that we are good. The system is healthy.
But what happens if the DB (log) blocks? Then this buffer will get full very quickly.
The drop pattern will attempt to signal with data on this buffer channel. What it will do then is drop that log line and not log it, but will not stop the video stream.

This buffer channel is giving us the ability to detect when we're not able to write to the device.
But once the device is writing again, this buffer will clear pretty quickly and we'll be able to recover.

So let's use a drop pattern to solve this problem:

[Code](../patterns/advanced/logger/logger.go)

```go
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
			fmt.Fprint(w, v)
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
```

So let's bring this logger that we just wrote into the previous app and see how it works.

```go
package main

import (
	"fmt"
	// "log"
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/cedrickchee/ultimate-go/concurrency/patterns/advanced/logger"
)

// device allows us to mock a device we write logs to.
type device struct {
	mu      sync.RWMutex
	problem bool
}

// Write implements the io.Writer interface.
func (d *device) Write(p []byte) (n int, err error) {
	// Simulate disk problems.
	for d.isProblem() {
		time.Sleep(time.Second)
	}

	fmt.Print(string(p))
	return len(p), nil
}

// isProblem checks in a safe way if there is a problem.
func (d *device) isProblem() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.problem
}

// flipProblem reverses the problem flag to the opposite value.
func (d *device) flipProblem() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.problem = !d.problem
}

func main() {
	// Number of goroutines that will be writing logs.
	const grs = 10

	// Create a logger value with a buffer of capacity
	// for each goroutine that will be logging.
	var d device
	// Replace the construction of the logger by using our new logger function
	// l := log.New(&d, "prefix", 0)
	l := log.New(&d, grs)

	// Generate goroutines, each writing to disk.
	for i := 0; i < grs; i++ {
		go func(id int) {
			for {
				l.Println(fmt.Sprintf("%d: log data", id))
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// We want to control the simulated disk blocking. Capture
	// interrupt signals to toggle device issues. Use <ctrl> z
	// to kill the program.

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	for {
		<-sigChan
		d.flipProblem()
	}
}
```

Start running the code. We got the "DROP". We are still streaming video. Unlike
before when we were completely blocked.

## WithCancel

[Sample program](example2/example2.go) to show different ways we can do
cancellation, timeout using the WithCancel function.

```go
func main() {
	// Create a context that is cancellable only manually.
	// The cancel function must be called regardless of the outcome.
    // WithCancel allows us to create a context and provides us a cancel
    // function that can be called in order to report a signal, a signal without
    // data, that we want whatever that Goroutine is doing to stop right away.
    // Again, we are using Background as our parents context.
    ctx, cancel := context.WithCancel(context.Background())

    // The cancel function must be called regardless of the outcome.
    // The Goroutine that creates the context must always call cancel. These are
    // things that have to be cleaned up. It's the responsibility that the
    // Goroutine creates the context the first time to make sure to call cancel
    // after everything is done.
	// The use of the defer keyword is perfect here for this use case.
	defer cancel()

	// We launch a Goroutine to do some work for us.
    // It will sleep for 50 milliseconds and then call cancel. It is reporting
    // that it want to signal a cancel without data.
	go func() {
		// Simulate work.
		// If we run the program using 50 ms, we expect the work to be complete. But if it is 150
		// ms, then we move on.
		time.Sleep(50 * time.Millisecond)

		// Report the work is done.
		cancel()
	}()

    // The original Goroutine that creates that channel is in its select case.
    // It will receive after time.After. We will wait 100 milliseconds for
    // something to happen. We are also waiting on context.Done. We are just
    // going to sit here, and if we are told to Done, we know that work up
    // there is complete.
	select {
	case <-time.After(100 * time.Millisecond):
		fmt.Println("moving on")
	case <-ctx.Done():
		fmt.Println("work complete")
	}
}
```

## WithDeadline

[Sample program](example3/example3.go) to show how to use the WithDeadline
function.
