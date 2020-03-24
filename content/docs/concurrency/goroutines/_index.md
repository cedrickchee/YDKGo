---
title: Goroutines
weight: 1
---

# Goroutines

Goroutines are functions that are created and scheduled to be run independently by the Go scheduler. The Go scheduler is responsible for the management and execution of goroutines.

* Goroutines are functions that are scheduled to run independently.
* We must always maintain an account of running goroutines and shutdown cleanly.
* Concurrency is not parallelism.
	* Concurrency is about dealing with lots of things at once.
	* Parallelism is about doing lots of things at once.

## Scheduling in Go

Every time our Go's program starts up, it looks to see how many cores are
available. Then it creates a logical processor.

### OS Scheduler

The OS scheduler is considered a preemptive scheduler. It runs down there in the
kernel. Its job is to look at all the threads that are in runnable states and
gives them the opportunity to run on some cores. These algorithms are fairly
complex: waiting, bouncing threads, keeping memory of threads, caching, and so
on. The OS is doing all of that for us. The algorithm is really smart when it
comes to multicore processor. Go doesn't want to reinvent the wheel. It wants
to sit on top of the OS and leverage it.

The OS is still responsible for OS threads, scheduling OS threads efficiently.
If we have a 2 core machine and a thousands threads that the OS has to schedule,
that's a lot of work.

#### Context Switching

A context switch on some OS thread is expensive when the OS have no clues of
what that thread is doing. It has to save all the possible states in order to be
able to restore that to exactly the way it was. If there are fewer threads, each
thread can get more time to be rescheduled. If there are more threads, each
thread has less time over a long period of time.

#### Less Is More

"Less is more" is a really big concept here when we start to write concurrent
software. We want to leverage the preemptive scheduler. So the Go's scheduler,
the logical processor actually runs in user mode, the mode our application is
running at. Because of that, we have to call the Go's scheduler a cooperating
scheduler. What brilliant here is the runtime that coordinating the operation.
It still looks and feels as a preemptive scheduler up in user land. We will see
how "less is more" concept gets to present itself and we get to do a lot more
work with less. Our goal needs to be how much work we get done with the less
number of threads.

#### Find The Balance

Think about this in a simple way because processors are complex: hyperthreading,
multiple threads per core, clock cycle. We can only execute one OS thread at a
time on any given core. If we only have 1 core, only 1 thread can be executed at
a time. Anytime we have more threads in runnable states than we have cores, we
are creating load, latency and we are getting more work done as we want. There
needs to be this balance because not every thread is necessarily will be active
at the same time. It all comes down to determining, understanding the workload
for the software that we are writing.

### Go Scheduler

When our Go program comes up, it has to see how many cores that available. Let's
say it found 1. It is going to create a logical processor ("P") for that core.
Again, the OS is scheduling things around OS threads. Every "P" is assigned an OS
Thread (M). The 'M' stands for machine. This Thread is still managed by the OS
and the OS is still responsible for placing the Thread on a Core for
execution (run our code).

The Linux scheduler has a run queue. Threads are placed in run queue in certain
cores or family of cores and those are constantly bounded as threads are running.
Go will do the same thing. Go has its run queue as well. There are 2 different
run queues in the Go scheduler: the Global Run Queue (GRQ) and the Local Run
Queue (LRQ). Each P is given a LRQ.

### Goroutine

Every Go program is also given an initial Goroutine ("G"), which is the path of
execution for a Go program. Threads are paths of execution. That path of
execution needs to be scheduled. In Go, every function or method can be created
to be a Goroutine, can become an independent path of execution that can be
scheduled to run on some OS threads against some cores.

When we start our Go program, the first thing runtime will do is creating a
Goroutine and putting that in some main LRQ for some "P". In our case, we only
have 1 "P" here so we can imagine that Goroutine is attached to "P".

A Goroutine, just like thread, can be in one of three major states: sleeping,
executing or in runnable state asking to wait for some time to execute on the
hardware. When the runtime creates a Goroutine, it will placed in "P" and
multiplex on this thread. Remember that it's the OS that taking the thread,
scheduling it, placing it on some core and doing execution. So Go's scheduler
will take all the code related to that Goroutine's path of execution, place it
on a thread, tell the OS that this thread is in runnable state and can we
execute it. If the answer is yes, the OS starts to execute on some cores there
in the hardware.

As the main Goroutine runs, it might want to create more paths of execution,
more Goroutines. When that happens, those Goroutines might find themselves
initially in the GRQ. These would be Goroutines that are in runnable state but
haven't been assigned to some "P"s yet. Eventually, they would end up in the
LRQ where they're saying they would like some time to execute.

This queue does not necessarily follow First-In-First-Out protocol. We have to
understand that everything here is non-deterministic, just like the OS scheduler.
We cannot predict what the scheduler will do when all things are equal. It will
make sure there is a balance. Until we get into orchestration, till we learn how
to coordinate these execution of these Goroutines, there is no predictability.

Here is the mental model of our example.

```
// GRQ

   m
   |
+-----+         LRQ
|  P  | ----------
+-----+          |
   |             G1
  Gm             |
                 G2
```

We have "Gm" executing on for this thread for this "P", and we are creating 2
more Goroutines "G1" and "G2". Because this is a cooperating scheduler, that
means that these Goroutines have to cooperate to be scheduled, to be swapped
context switch on this OS thread "M".

There are 4 major places in our code where the scheduler has the opportunity to
make a scheduling decision.
- The keyword `go` that we are going to create Goroutines. That is also an
opportunity for the scheduler to rebalance when it has multiple "P".
- A system call. These system calls tend to happen all the time already.
- A channel operation because there is mutex (blocking call) that we will learn
later.
- Garbage collection.

Back to the example, says the scheduler might decide "Gm" has enough time to run,
it will put "Gm" back to the run queue and allow G1 to run on that "M". We are
now having context switch.

```
   M
   |
+-----+     LRQ
|  P  | ------
+-----+      |
   |         Gm
   G1        |
             G2
```

Let's say "G1" decides to open up a file. Opening up a file can take microsecond
or 10 milliseconds. We don't really know. If we allow this Goroutine to block
this OS thread while we open up that file, we are not getting more work done.
In this scenario here, having a single "P", we are single threaded software app.
All Goroutines only execute on the "M" attached to this "P". What happen is this
Goroutine will block this "M" for potential a long time. We are basically be
stalled while we still have works that need to get done. So the scheduler is not
going to allow that to happen, What actually happen is that the scheduler will
detach that "M" and "G1". It will bring a new "M", say "M2", then decide what
"G" from the run queue should run next, say "G2".

```
        M2
        |
M    +-----+     LRQ
|    |  P  | ------
G1   +-----+      |
        |         Gm
        G2
```

We now have 2 threads in a single threaded program. From our perspective, we are
still single threading because the code that we are writing, the code associated
with any "G" can only run against this "P" and this "M". However, we don't know
at any given time what "M" we are running on. "M" can get swapped out but we are
still single threaded.

Eventually, "G1" will come back, the file will be opened. The scheduler will
take this "G1" and put it back to the run queue so we can be executed against on
this "P" for some "M" ("M2" in this case). "M" get placed on the side for later
use. We are still maintaining these 2 threads. The whole process can happen
again.

```
        M2
        |
M    +-----+     LRQ
     |  P  | ------
      -----       |
        |         Gm
        G2        |
                  G1
```

It's a really brilliant system of trying to leverage this thread to its fullest
capability by doing more on 1 thread. It do so much on this thread we don't
need another.

There is something called a Network poller. It will do all the low level
asynchronous networking stuff. Our "G", if it is going to do anything like that,
it might be moved out to the Network poller and then brought back in. From our
perspective, here is what we have to remember: The code that we are writing
always run on some "P" against some "M". Depending on how many "P" we have,
that's how many threads variables for us to run.

Concurrency is about managing a lot of thing at once. This is what the scheduler
is doing. It manages the execution of these 3 Goroutines against this one "M"
for this "P". Only 1 Goroutine can be executed at a single time.

If we want to do something in parallel, which means doing a lot of things at
once, then we would need to create another "P" that has another "M", say "M3".

```
   M3              M2
   |               |
+-----+          +-----+     LRQ
|  P  | ------   |  P  | ------
+-----+      |   +-----+      |
   |         Gx     |         Gm
   Gx        |      G2        |
             Gx               G1
```

Both are scheduled by the OS. So now we can have 2 Goroutines running at the
same time in parallel.

#### Practical example

We have a multiple threaded software. The program launched 2 threads. Even if
both threads end up on the same core, each want to pass a message to each other.
What has to happen from the OS point of view?
We have to wait for thread 1 to get scheduled and placed on some cores — a
context switch ("CTX") has to happen here. While that's happening, thread is
asleep so it's not running at all. From thread 1, we send a message over and
want to wait to get a message back. In order to do that, there is another
context switch needs to be happened because we can put a different thread on
that core (?). We are waiting for the OS to schedule thread 2 so we are going to
get another context switch, waking up and running, processing the message and
sending the message back. On every single message that we are passing back and
forth, thread going from executable state to runnable state to asleep state.
This will cost a lot of context switches to occur.

```
T1                         T2
| CTX      message     CTX |
| -----------------------> |
| CTX                    | |
| CTX      message     CTX |
| <----------------------- |
| CTX                  CTX |
```

Let's see what happen when we are using Goroutines, even on a single core.
"G1" wants to send a message to "G2" and we perform a context switch. However,
the context here is user's space switch. "G1" can be taken out of the thread and
"G2" can be put on the thread. From the OS point of view, this thread never go
to sleep. This thread is always executing and never needed to be context
switched out. It is the Go's scheduler that keeps the Goroutines context
switched.

```
             M
             |
          +-----+
          |  P  |
          +-----+
G1                         G2
| CTX      message     CTX |
| -----------------------> |
| CTX                    | |
| CTX      message     CTX |
| <----------------------- |
| CTX                  CTX |
```

If a "P" for some "M" here has no work to do, there is no "G", the runtime
scheduler will try to spin that "M" for a little bit to keep it hot on the core.
Because if that thread goes cold, the OS will pull it off the core and put
something else on. So it just spin a little bit to see if there will be another
"G" comes in to get some work done.

This is how the scheduler work underneath. We have a "P", attached to thread "M".
The OS will do the scheduling. We don't need any more OS threads than cores we
have. If we have more threads than cores we have, all we do is putting load on
the OS. We allow the Go's scheduler to make decisions on our Goroutines, keeping
the least number of threads we need and hot all the time if we have work. The
Go's scheduler will look and feel preemptive even though we are calling a
cooperating scheduler.

However, let's not think about how the scheduler work. Think the following way
makes it easier for future development.
Every single "G", every Goroutine that is in runnable state, is running at the
same time.

## Creating Goroutines

How we manage concurrency in Go? Understand that Go routines are chaos. I like
to think of Go routines as children. You know that children are chaos.
So, really what we're will learn about, for the rest of this section, is what I
say, a parenting class.

### Language Mechanics

One of the most important thing that we must do from day one is to write
software that can startup and shutdown cleanly. This is very very important.

[Sample program](example1/example1.go) to show how to create goroutines and how
the scheduler behaves.

```go
func init() {
	// Allocate one logical processor for the scheduler to use.
	runtime.GOMAXPROCS(1)
}
```

`init` calls a function from the runtime package called `GOMAXPROCS`. This is
also an environment variable, which is why is all capitalized. Prior to 1.5,
when our Go program came up for the first time, it came up with just a
single "P", regardless of how many cores. The improvement that we made to the
garbage collector and scheduler changed all that.

```go
// wg is used to manage concurrency.
var wg sync.WaitGroup
```

`wg` is set to its zero value. This is one of the very special types in Go that
are usable in its zero value state. It is also called Asynchronous Counting
Semaphore. It has 3 methods: `Add`, `Done` and `Wait`.
`n` number of Goroutines can call this method at the same time and it's all get
serialized.
- Add keeps a count of how many Goroutines out there.
- Done decrements that count because some Goroutines are about to terminated.
- Wait holds the program until that count goes back down to zero.

```go
wg.Add(2)
```

We are creating 2 Goroutines.
We rather call `Add(1)` and call it over and over again to increment by 1. If we
don't how many Goroutines that we are going to create, that is a smell.

```go
go func() {
    lowercase()
    wg.Done()
}()
```

Create a Goroutine from the `lowercase` function using anonymous function.
We have a function decoration here with no name and being called by the `()` in
the end. We are declaring and calling this function right here, inside of main.
The big thing here is the keyword `go` in front of `func()`. We don't execute
this function right now in series here. Go schedules that function to be a "G",
say "G1", and load in some LRQ for our "P". This is our first "G". Remember, we
want to think that every "G" that is in runnable state is running at the same
time. Even though we have a single "P", even though we have a single thread,
we don't care. We are having 2 Goroutines running at the same time: `main` and
"G1".

```go
// Create a goroutine from the uppercase function.
go func() {
    uppercase()
    wg.Done()
}()
```

We are doing it again. We are now having 3 Goroutines running at the same time.

```go
fmt.Println("Waiting To Finish")
wg.Wait()
```

Wait for the Goroutines to finish.
This is holding `main` from terminating because when the `main` terminates, our
program terminates, regardless of what any other Goroutine is doing. There is a
golden rule here: We are not allowed to create a Goroutine unless we can tell
when and how it terminates. `Wait` allows us to hold the program until the 2
other Goroutines report that they are done. It will wait, count from 2 to 0.
When it reaches 0, the scheduler will wake up the `main` Goroutine again and
allow it to be terminated.

**Sequence**

We call the `uppercase` after `lowercase` but Go's scheduler chooses to call the
`lowercase` first. Remember we are running on a single thread so there is only
one Goroutine is executed at a given time here. We can't see that we are running
concurrently that the `uppercase` runs before the `lowercase`. Everything starts
and completes cleanly.

**What if we forget to hold `Wait`?**

We would see no output of `uppercase` and `lowercase`. This is pretty much a
data race. It's a race to see the program terminates before the scheduler stops
it and schedules another Goroutine to run. By not waiting, these Goroutine never
get a chance to execute at all.

**What if we forget to call `Done`?**

Deadlock!
This is a very special thing in Go. When the runtime determines that all the
Goroutines are there can no longer move forward, it will panic.

#### Goroutine Time Slicing

How the Go's scheduler, even though it is a cooperating scheduler (not
preemptive), it looks and feel preemptive because the runtime scheduler is
making all the decisions for us. It is not coming for us.

The [program](example2/example2.go) below will show us a context switch and how we can predict when the
context switch is going to happen. It is using the same pattern that we've seen
in the last program. The only difference is the `printHashes` function.

```go
func init() {
	// Allocate one logical processor for the scheduler to use.
	runtime.GOMAXPROCS(1)
}

func main() {
	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Create Goroutines")

	// Create the first goroutine and manage its lifecycle here.
	go func() {
		printHashes("A")
		wg.Done()
	}()

	// Create the second goroutine and manage its lifecycle here.
	go func() {
		printHashes("B")
		wg.Done()
	}()

	// Wait for the goroutines to finish.
	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("Terminating Program")
}

// printHashes calculates the sha1 hash for a range of
// numbers and prints each in hex encoding.
func printHashes(prefix string) {
	// print each has from 1 to 10. Change this to 50000 and
	// see how the scheduler behaves.
	for i := 1; i <= 50000; i++ {
		// Convert i to a string.
		num := strconv.Itoa(i)

		// Calculate hash for string num.
		sum := sha1.Sum([]byte(num))

		// Print prefix: 5-digit-number: hex encoded hash
		fmt.Printf("%s: %05d: %x\n", prefix, i, sum)
	}

	fmt.Println("Completed", prefix)
}
```

`printHashes` is not special. It just requires a little bit more time to
complete. When we run the program, what we will see are context switches at some
point for some particular number. We cannot predict when the context switch
happen. That's why we say the Go's scheduler looks and feels very preemptive
even though it is a cooperating scheduler.

### Goroutine and Parallelism

This [program](example3/example3.go) show how Goroutines run in parallel. We are
going to have 2 "P" with 2 "M", and 2 Gorouines running in parallel on each "M".
This is still the same program that we are starting with. The only difference is
that we are getting rid of the `lowercase` and `uppercase` function and putting
their code directly inside Go's anonymous functions.

Looking at the output, we can see a mix of uppercase of lowercase characters.
These Goroutines are running in parallel now.

```go
func init() {
	// Allocate two logical processors for the scheduler to use.
	runtime.GOMAXPROCS(2)
}

func main() {
	// wg is used to wait for the program to finish.
	// Add a count of two, one for each goroutine.
	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Start Goroutines")

	// Declare an anonymous function and create a goroutine.
	go func() {

		// Display the alphabet three times.
		for count := 0; count < 3; count++ {
			for r := 'a'; r <= 'z'; r++ {
				fmt.Printf("%c ", r)
			}
		}

		// Tell main we are done.
		wg.Done()
	}()

	// Declare an anonymous function and create a goroutine.
	go func() {

		// Display the alphabet three times.
		for count := 0; count < 3; count++ {
			for r := 'A'; r <= 'Z'; r++ {
				fmt.Printf("%c ", r)
			}
		}

		// Tell main we are done.
		wg.Done()
	}()

	// Wait for the goroutines to finish.
	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\nTerminating Program")
}
```
