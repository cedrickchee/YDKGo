---
title: Channels
weight: 3
---

# Channels

Channels are a way of doing orchestration in our multi-threaded software.

Reminder: Synchronization is about getting in line and taking a turn.
Orchestration is about the interaction.

Channels will allow us to move data across these Goroutine boundaries, and do
the orchestration which is not in line but filing an interaction.

Prepare your mind:
**Don't think about channels as a data structure, a queue, a synchronous queue.**

Look at channel behavior. The behavior we will talk about is signaling.

## Signaling Semantics

Signaling is the semantic. Channels are for signaling. The idea is that a
Goroutine is going to send a signal to another Goroutine.

### Signaling guarantees of delivery

The idea is do you need a guarantee that a signal being sent by one Goroutine
has been received by the other.

Go has 2 types of channels: unbuffered and buffered. They both allow us to
signal with data. The key difference is that, when we use unbuffered channel, we
are signaling and getting a guarantee the signal was received. We are not going
to be sure if that Goroutine is done whatever work we assign it to do but we do
have the guarantee. The trade off for the guarantee that the signal was received
is higher latency because we have to wait to make sure that the Goroutine on the
other side of that unbuffered channel receive the data.

This is how the unbuffered channel going to work.
There will be a Goroutine comes to the channel. The channel wants to signal with
some piece of data. It will put the data right there in the channel. However,
the data is locked in and cannot move because channel has to know if there is
another Goroutine is on the other side to receive it. Eventually a Goroutine
come and say that it want to receive the data. Both of Goroutines are not
putting their hands in the channel. The data now can be transferred. Here is the
key to why that unbuffered channel gives us that guarantee: the receive happens
first. When the receive happens, we know that the data transfer has occurred and
we can walk away.

```
//  G                      G
//  |        Channel       |
//  |     +----------+     |
//  |     |   D  D   |     |
//  |-----|--->  <---|-----|
//  |     |          |     |
//  |     +----------+     |
//  |                      |
```

The unbuffered channel is a very powerful channel. We want to leverage that
guarantee as much as possible. But again, the cost of the guarantee is higher
latency because we have to wait for this.

The buffered channel is a bit different: we do not get the guarantee but we get
to reduce the amount of latencies on any given send or receive.

Back to the previous example, we replace the unbuffered channel with a buffered
channel. We are going with a buffered channel of just 1. It means there is a
space in this channel for 1 piece of data that we are using the signal and we
don't have to wait for the other side to get it. So now a Goroutine comes in,
put the data in and then move away immediately. In other word, the send is
happening before the receive. All the sending Goroutine know is that it issues
the signal, put that data but has no clue when the signal is going to be
received. Now hopefully a Goroutine comes in. It see that there is a data there,
receive it and move on.

```
//  G                      G
//  |      Channel (1)     |
//  |     +----------+     |
//  |---->|    D     |<----|
//  |     +----------+     |
//  |
```

We use a buffered of 1 when dealing with these type of latency. We may need
buffers that are larger but there are some design rules that we will learn later
on we use buffers that are greater than 1. But if we are in a situation where we
can have these sends coming in and they could potentially be locked then we have
to think again: if the channel of 1 is fast enough to reduce the latency that we
are dealing with. Because what will happen is the following:
What we are hoping is, the buffered channel is always empty every time we
perform a send.

Buffered channel is not for performance. What the buffered channel need to be
used for is continuity, to keep the wheel moving. One thing we have to
understand is that, everybody can write a piece of software that works when
everything is going well. When things are going bad, it's where the architecture
and engineer really come in. Our software doesn't enclose and it doesn't cost
stress. We need to be responsible.

Diagrams:
- [Guarantee Of Delivery](https://github.com/ardanlabs/gotraining/blob/master/topics/go/concurrency/channels/README.md#guarantee-of-delivery)
- [Signaling With Or Without Data](https://github.com/ardanlabs/gotraining/blob/master/topics/go/concurrency/channels/README.md#signaling-with-or-without-data)
- [State](https://github.com/ardanlabs/gotraining/blob/master/topics/go/concurrency/channels/README.md#state)

### Language Mechanics

Back to the example, it's not important that we know exactly the signaling data
was received but we do have to make sure that it was. The buffered channel of 1
gives us almost guarantee because what happen is: it performs a send, puts the
data in there, turns around and when it comes back, it sees that the buffered is
empty. Now we know that it was received. We don't know immediately at the time
that we sent but by using a buffer of 1, we do know that is empty when we come
back. Then it is okay to put another piece of data in there and hopefully when
we come back again, it is gone. If it's not gone, we have a problem. There is a
problem upstream. We cannot move forward until the channel is empty. This is
something that we want to report immediately because we want to know why the
data is still there. That's how we can build systems that are reliable. We don't
take more work at any give time. We identify upstream when there is a problem so
we don't put more stress on our systems. We don't take more responsibilities for
things that we shouldn't be.

## Basic Patterns - Part 1 (Wait for Task)

**Unbuffered channel: Signaling with data**

[Basic mechanics](example1/example1.go)

`waitForTask` shows the basics of a send and receive.
We are using `make` function to create a channel. We have no other way of
creating a channel that is usable until we use `make`. Channel is also based on
type, a type of data that we will do the signaling. In this case, we use
`string`. That channel is a reference type. `ch` is just a pointer variable to
larger data structure underneath. We will be using value semantics all the time.

```go
// waitForTask: You are a manager and you hire a new employee. Your new
// employee doesn't know immediately what they are expected to do and waits for
// you to tell them what to do. You prepare the work and send it to them. The
// amount of time they wait is unknown because you need a guarantee that the
// work your sending is received by the employee.
func waitForTask() {
    // This is an unbuffered channel.
	ch := make(chan string)

	go func() {
        // This is a receive: also an arrow but it is a unary operation where it is
        // attached to the left hand side of the channel to show that is coming out.
        // We are now have a unbuffered channel where the send and receive have to
        // come together. Both will block until both come together so
        // the exchange can happen.
		p := <-ch
		fmt.Println("employee : recv'd signal :", p)
	}()

    time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    // This is a send: a binary operation with the arrow pointing into the channel.
    // We are signaling with a string "paper".
	ch <- "paper"
	fmt.Println("manager : sent signal")

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------------------")
}
```

If I ran this a few more times, we might get lucky and see the receive output
happen before the send.

This is why so many people get confused with concurrency, because they're
looking at the order of things and the only order is atomically at that actual
send and receive. You **cannot be using print statements to look at ordering.**
I've left the print statements in here to prove that to you. Your print
statements are not going to help you. They will confuse you.

## Basic Patterns - Part 2 (Wait for Result)

We will use this pattern in things like drop patterns and fan-out patterns.

```go
// waitForResult: You are a manager and you hire a new employee. Your new
// employee knows immediately what they are expected to do and starts their
// work. You sit waiting for the result of the employee's work. The amount
// of time you wait on the employee is unknown because you need a
// guarantee that the result sent by the employee is received by you.
func waitForResult() {
    // This is an unbuffered channel.
	ch := make(chan string)

	go func() {
        time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
        // This is a send: a binary operation with the arrow pointing into the channel.
		// We are signaling with a string "paper".
		ch <- "paper"
		fmt.Println("employee : sent signal")
	}()

    // This is a receive: also an arrow but it is a unary operation where it is
    // attached to the left hand side of the channel to show that is coming out.
    // We are now have a unbuffered channel where the send and receive have to
    // come together. We also know that the signal has been received because the
    // receive happens first. Both will block until both come together so
    // the exchange can happen.
	p := <-ch
	fmt.Println("manager : recv'd signal :", p)

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------------------")
}
```

## Basic Patterns - Part 3 (Wait for Finished)

**Unbuffered channel: Signaling without data**

Our next basic and our last basic channel pattern here is wait for finished.
Wait for finished will show us how we can signal without data.
Now I just want to make it clear that what I'm showing you would really be
better served with a wait group. A wait group would make this code cleaner but I
need to show you the mechanics so later on when we talk about cancellation and
deadlines with the context package, it'll give you a little bit more sense.

```go
// waitForFinished: Think about being a manager and hiring a new employee. In
// this scenario, you want your new employee to perform a task immediately when
// they are hired, and you need to wait for the result of their work. You need
// to wait because you can't move on until you know they are but you don't need
// anything from them.
func waitForFinished() {
	// We are making a channel using an empty struct. This is a signal without data.
	ch := make(chan struct{})

    // We will launch a Goroutine to do some work. Then, it wants to signal
    // another Goroutine that it's done. It will close the channel to report
    // that it's done without the need of data.

    // When we create a channel, buffered or unbuffered, that channel can be in
    // 2 different states. All channels start out in open state so we can send
    // and receive data. When we change the state to be closed, it cannot be
    // opened. We also cannot close the channel twice because that is an
    // integrity issue. We cannot signal twice without data twice.
	go func() {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		close(ch)
		fmt.Println("employee : sent signal")
	}()

    // We can get a second flag, which is a boolean that I use wd for "with
    // data". What it's telling us is we receive successfully off this channel
    // and the wd flag is true, then we have data.
	// When the channel is closed, the receive will immediately return.
    // When we receive on a channel that is open, we cannot return until we
    // receive the data signal. But if we receive on a channel that is closed,
    // we are able to receive the signal without data. We know that event is
    // occurred. Every receive on that channel will immediately return.
	_, wd := <-ch
    fmt.Println("manager : recv'd signal :", wd)

    time.Sleep(time.Second)
    fmt.Println("-------------------------------------------------------------")
}
```

## Pooling Pattern

**Unbuffered channel: Close and range**

`pooling` shows how to use `range` to receive value and using `close` to
terminate the loop. We're both signaling with data and without data.

```go
// pooling: You are a manager and you hire a team of employees. None of the new
// employees know what they are expected to do and wait for you to provide work.
// When work is provided to the group, any given employee can take it and you
// don't care who it is. The amount of time you wait for any given employee to
// take your work is unknown because you need a guarantee that the work your
// sending is received by an employee.
func pooling() {
    // We want to have guarantees and we're signaling with string data.
    // Pooling and the whole idea is Goroutines waiting for work to do. So the
    // work is gonna be string based and we want guarantees.
    // You absolutely want guarantees with pooling, because later on we want to
    // be able to apply deadlines or timeouts when a pool is, let's say,
    // underload and not responding fast enough.
	ch := make(chan string)

    // We will create this pool of Goroutines right here.
    // We will end up with two paths of execution in the pool, Goroutine1,
    // Goroutine2 and they're both will be here in this waiting state.
    // What makes them wait is the `for` range. We are ranging over the channel.
    // When you range over a channel, you are basically in a channel receive.
    // We are now blocked in a channel receive.
    // Order doesn't matter, because once data comes into the channel, the
    // scheduler can choose any Goroutine that it wants to do the work.
    // The `for` range terminate through a signaling change, going from open to
    // closed and that will terminate the loop.
	g := runtime.NumCPU()
	for e := 0; e < g; e++ {
		go func(emp int) {
			for p := range ch {
				fmt.Printf("employee %d : recv'd signal : %s\n", emp, p)
			}
			fmt.Printf("employee %d : recv'd shutdown signal\n", emp)
		}(e)
	}

    // Populate channel with data.
    // Work loop. We will pass work into the pool, 100 pieces of work.
    const work = 100
	for w := 0; w < work; w++ {
        // Setting the work into the channel send.
        // In order for the send to complete, we need a corresponding receive
        // and again, the scheduler has to choose a Goroutine.
		ch <- "paper"
		fmt.Println("manager : sent signal :", w)
	}

    // Close the channel
    // We're simulating a program shutting down.
	close(ch)
	fmt.Println("manager : sent shutdown signal")

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------------------")
}
```

## Fan Out Pattern - Part 1

**Buffered channel: send and receive**

Fan out allow you to take a piece of work and to distribute it across 'n' number
of Goroutines that can run in parallel and we will use that wait for result as
our base pattern to build this.

But I need to stress that fan out patterns are dangerous patterns, especially in
web services where you might already be having 10K Goroutines already running in
the service and suddenly you have some Goroutine fan out another 'n' number of
Goroutines. These things can multiply very quickly, so I get afraid of fan out
patterns in long running apps like services, but for background apps that run on
cron jobs in certain intervals, and for maybe CLI tooling it's great.

```go
// fanOut: You are a manager and you hire one new employee for the exact amount
// of work you have to get done. Each new employee knows immediately what they
// are expected to do and starts their work. You sit waiting for all the results
// of the employees work. The amount of time you wait on the employees is
// unknown because you need a guarantee that all the results sent by employees
// are received by you. No given employee needs an immediate guarantee that you
// received their result.
func fanOut() {
    emps := 20
    // This is a buffered channel of 20.
    // A buffer for every Goroutine.
	ch := make(chan string, emps)

    // We launched 20 Goroutines.
    // Here, we're using wait for result pattern.
	for e := 0; e < emps; e++ {
		go func(emp int) {
            time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond) // work. e.g. DB work.
            // Send data to channel.
            // Goroutine signaling with data. Every Goroutine has their own slot
            // in the buffer.
			ch <- "paper" // e.g. DB work send back their result.
			fmt.Println("employee : sent signal :", emp)
		}(e)
	}

    // emps is our own little local counter that acts like a local wait group
    // and we just sit here and we loop 20 times.
	for emps > 0 {
        // Perform the received.
        // Pull data out from the buffered channel.
        // There's no latency between the send and receive because the send
        // happened before the receive.
		p := <-ch
		emps-- // if we get a piece of data, we decrement the work count.
		fmt.Println(p)
		fmt.Println("manager : recv'd signal :", emps)
    }
    // How long are we going to be waiting here? That is really still unknown
    // latency. We still are kind of getting guarantees here, even though we're
    // using a buffered channel, because we're kinda guaranteed that we're not
    // moving on until all of the data comes back. We will not return from main
    // until the wait group goes down to zero.

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------------------")
}
```

## Fan Out Pattern - Part 2

**Buffered channel: fan out semaphore pattern**

There's another version of the fan out pattern that we call the fan out
semaphore.

This does let's fan out in a number of Goroutines, but we don't necessarily want
all of the Goroutines to be able to run at the same time. This kind of dynamic
pool. This is another way of being able to reduce latency in terms of Goroutine
creation, but at the same time limit the impact that these Goroutines are having
on another resource.

```go
// fanOutSem: You are a manager and you hire one new employee for the exact amount
// of work you have to get done. Each new employee knows immediately what they
// are expected to do and starts their work. However, you don't want all the
// employees working at once. You want to limit how many of them are working at
// any given time. You sit waiting for all the results of the employees work.
// The amount of time you wait on the employees is unknown because you need a
// guarantee that all the results sent by employees are received by you. No
// given employee needs an immediate guarantee that you received their result.
func fanOutSem() {
    emps := 20
    // Buffer channel of 20.
    // We will reduce the latency on the sends between the sends and receives.
	ch := make(chan string, emps)

    g := runtime.NumCPU()
    // Our semaphore channel and we've set it to NumCPU, let's say 5. I only
    // want 5 out of the 20 Goroutines to be actually executing at the same time.
    // We will save some latency by creating all the Goroutines we need up front.
	sem := make(chan bool, g)

    // Create 20 Goroutines, same thing that we did before.
	for e := 0; e < emps; e++ {
		go func(emp int) {
            // Except, this call is a channel send on a buffered channel. It will
            // cause all these Goroutines to block.
            // Blocked here because this channel send will only complete if
            // there's enough room in these semaphore.
			sem <- true
			{
                // Goroutine do their work and send result back.
				time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
				ch <- "paper"
				fmt.Println("employee : sent signal :", emp)
            }
            // Semaphore channel receive. Pull data out from the semaphore
            // buffer to give room for another Goroutine to come in and then do
            // their work.
			<-sem
		}(e)
	}

	for emps > 0 {
		p := <-ch
		emps--
		fmt.Println(p)
		fmt.Println("manager : recv'd signal :", emps)
	}

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------------------")
}
```

## Drop Pattern

**Buffered channel: Select and drop**

Drop pattern helps us reduce back pressure when things are going bad.

What the drop pattern tries to identify is what our capacity is. What is the
maximum number of pending anything, requests, tasks, that we can take in before
we have to say no.

It's going to take integration tests and load testing, usually, to figure out
what your capacity is.

Drop patterns help us identify failures quickly, they help us stop the bleeding,
then, they help us move forward again, when that wound has been healed.

`drop` shows how to use the `select` to walk away from a channel operation if it
will immediately block.

A Denial-of-service attack is a great example. We get a bunch of requests coming
to our server. If we try to handle every single request, we are gonna implode.
We have to handle what we can and drop other requests.

```go
// drop: You are a manager and you hire a new employee. Your new employee
// doesn't know immediately what they are expected to do and waits for
// you to tell them what to do. You prepare the work and send it to them. The
// amount of time they wait is unknown because you need a guarantee that the
// work your sending is received by the employee. You won't wait for the
// employee to take the work if they are not ready to receive it. In that case
// you drop the work on the floor and try again with the next piece of work.
func drop() {
    // We can use buffer that are larger than 1. We have to measure what the
    // buffer should be. It cannot be random.
    const cap = 5
    // Make a buffered channel of 5. This channel is our capacity.
    // Once we have at least 5 pending tasks waiting in our buffer here, we're
    // at capacity. We're not going to go to 6. At that point, we're just going
    // to drop the work on the floor. We're not going to take it. Taking anymore
    // is too much risk.
    ch := make(chan string, cap)

	go func() {
        // We are in the receive loop waiting for data to work on.
		for p := range ch {
			fmt.Println("employee : recv'd signal :", p)
		}
	}()

    // This will send the work to the channel.
    // E.g. maybe here we're reading some sort of network and as we pull data
    // out of the network, we send it into our buffer ch.
    // But once this buffer fills up, which means it blocks, the default case
    // comes in and drop things. We're not going to block.
	const work = 20
	for w := 0; w < work; w++ {
        // The select statement allows a single Goroutine, to handle multiple
        // channel operations, whether they're sends or receives, at exactly the
        // same time. So you can create event loops here. This is also how we're
        // going to do cancellation and deadlines, and a drop pattern is very
        // much like a cancellation.
		select {
		case ch <- "paper": // trying to send to channel
			fmt.Println("manager : sent signal :", w)
        default:
            // If the send is going to block, don't block, move on. Basically,
            // now going to decide what we're going to do with this data we read
            // off the network. Maybe send back a HTTP 500. We're not going to
            // create back pressure against this channel. This isn't necessarily
            // time out reduction, it's literally, capacity reduction.
			fmt.Println("manager : dropped data :", w)
		}
    }

    // The health of our system is measured by our capacity. If this is working
    // fast enough, then this ch buffer never gets full. But if something bad
    // happens, and the buffer gets full, we no longer block, we don't take
    // blocking latencies, we drop. Now, in order for this to work, we got to
    // be able to identify when we're full, without fully blocking.

	close(ch)
	fmt.Println("manager : sent shutdown signal")

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------------------")
}
```

## Cancellation Pattern

**Buffered channel: select and receive**

Cancellation deadlines are critical to software because a task or a request
cannot take forever. We will use the context package that's part of the standard
library for this, so I will use this code to show you some basic mechanics of
how the context package does things.

```go
// cancellation: You are a manager and you hire a new employee. Your new
// employee knows immediately what they are expected to do and starts their
// work. You sit waiting for the result of the employee's work. The amount
// of time you wait on the employee is unknown because you need a
// guarantee that the result sent by the employee is received by you. Except
// you are not willing to wait forever for the employee to finish their work.
// They have a specified amount of time and if they are not done, you don't
// wait and walk away.
func cancellation() {
    // We're using a buffered channel of 1, we're getting some delayed
    // guarantees here.
	ch := make(chan string, 1)

	go func() {
        // Simulate that this work could take from anywhere from 1 to 200 ms.
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
        ch <- "paper"
        // The key is is that this might be work that one day we want to cancel.
        // If we (this Goroutine) do the work, then we can't cancel it. We can't
        // because we're blocked waiting for it to get done.
        // This is where we need another Goroutine if we want to be able to cancel.
	}()

    // Here, cancel says that this work has 100 milliseconds of time to get done.
    // We're not willing to wait more than 150 milliseconds, if that happens,
    // we've got to move on.
	duration := 150 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

    // select allows us to handle multiple channel operations, both sends and
    // receives at the same time.
	select {
    // Here, we're performing 2 different receives on 2 different channels
    case d := <-ch:
        // Channel receive of the ch buffer
		fmt.Println("work complete", d)

    case <-ctx.Done():
        // Channel receive on a timer.
        // ctx is a Context and carries a cancellation signal.
        // ctx.Done returns a channel that's closed when work done on behalf of
        // this context should be canceled.
		fmt.Println("work cancelled")
	}

	time.Sleep(time.Second)
	fmt.Println("-------------------------------------------------------------")
}
```

## Goroutine leak

If we used an unbuffered channel `ch := make(chan string)`, we would have a very
common bug in this code.

One of the biggest bug we are going to have and potential memory problem is when
we write code like this and we don't give the Goroutine an opportunity to
terminate.

We are using an unbuffered channel and this Goroutine at some point, its
duration will finish and it will want to perform a send. But this is an
unbuffered channel. This send cannot be completed unless there is a
corresponding receive. What if this Goroutine times out and moves on? There is
no more corresponding receive. That send block forever. Therefore, we will have
a Goroutine leak, which means it will never be terminated.

Goroutine leaks eventually cause memory leak. But Go is really good at cleaning
things up, so sometimes these Goroutine leaks don't show themselves for hours,
days, or weeks because Go is very good at dealing with and minimizing memory.

The cleanest way to fix this bug is to use the buffered channel of 1. If this
send happens, we don't necessarily have the guarantee. We don't need it. We just
need to perform the signal then we can walk away. Therefore, either we get the
signal on the other side or we walk away. Even if we walk away, this send can
still be completed because there is room in the buffer for that send to happen.
