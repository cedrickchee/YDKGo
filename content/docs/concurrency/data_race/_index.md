---
title: "Data Races"
weight: 2
---

# Data Races

A data race is when two or more goroutines attempt to read and write to the same resource at the same time. Race conditions can create bugs that appear totally random or can never surface as they corrupt data. Atomic functions and mutexes are a way to synchronize the access of shared resources between goroutines.

* Goroutines need to be coordinated and synchronized.
* When two or more goroutines attempt to access the same resource, we have a data race.
* Atomic functions and mutexes can provide the support we need.

We will learn important topics around synchronization and orchestration, but we
will focus on synchronization. These are the two things you have to worry about
as a multi-threaded software developer, synchronization and orchestration, and
data races are one of these places where it is the nastiest bugs you ever want to
have.

## Cache Coherency and False Sharing

A data race is is when you've got at least two paths of execution, like two Go routines, accessing the same memory location at the same time, where one is doing a read, and the other is at least doing a write. I mean they both could me doing writes too, and that would be a data race. You cannot be mutating memory.

It's two paths of execution at the same time, we will have data corruption. This is where synchronization comes in.

The best way to think about synchronization is if you went to Starbucks and you got in line because you want to get some coffee. Now you're in line, you're waiting for your turn to get up to the counter. Anytime Goroutines have to get in line, that is a synchronization issue. But once you get to the counter and you start talking to the person at the register, you now have an orchestration issue. We're having a conversation, we're exchanging money, there's data going back and forth. That is orchestration.

**Hardware level**

We've gotta appreciate this when we start writing multi-threaded software because it can really hurt our performance. Our cacheing systems, though they're helping us try to reduce latency to main memory access, you can actually be thrashing memory as well if we're not careful.

**Cache coherency problem**

Cache coherency problems where all we're doing is thrashing through memory because copies of that data is being leveraged, re-modified, across all these cores.

**False sharing**

False sharing occurs when the cores are reading and writing to different
addresses on the same cache line. Even though they are not sharing data, the
caches act like they are.

False sharing occurs when you don't really have a synchronization problem, but
we still have the cache coherency issue.

## Data Race and Race Detection

As soon as we add another Goroutine to our program, we add a huge amount of
complexity. We can't always let the Goroutine run stateless. There has to be
coordination. There are, in fact, 2 things that we can do with multithread
software.
1. We either have to synchronize access to share state like that `WaitGroup` is
done with `Add`, `Done` and `Wait`.
2. Or we have to coordinate these Goroutines to behave in a predictable or
responsible manner.

Up until the use of channel, we have to use atomic function, mutex, to do both.
The channel gives us a simple way to do orchestration. However, in many cases,
using atomic function, mutex, and synchronizing access to shared state is the
best way to go.

Atomic instructions are the fastest way to go because deep down in memory, Go is
synchronizaing 4-8 bytes at a time.
Mutexes are the next fastest. Channels are very slow because not only they are
mutexes, there are all data structure and logic that go with them.

Data races is when we have multiple Goroutines trying to access the same memory
location.
For example, in the simplest case, we have a integer that is a counter. We have
2 Goroutines that want to read and write to that variable at the same time.
If they are actually doing it at the same time, they are going to trash each
other read and write. Therefore, this type of synchronizing access to the shared
state has to be coordinated.

The problem with data races is that they always appear random.

[Sample program](example1/example1.go) to show how to create race conditions in
our programs. We don't want to do this.

```go
import (
	"fmt"
	"sync"
)

// counter is a variable incremented by all goroutines.
var counter int

func main() {
	// Number of goroutines to use.
	const grs = 2

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(grs)

	// Create two goroutines.
	for g := 0; g < grs; g++ {
		go func() {
			for i := 0; i < 2; i++ {
				// Capture the value of Counter.
				value := counter

				// Increment our local value of Counter.
				value++

				fmt.Println(value)

				// Store the value back into Counter.
				counter = value
			}

			wg.Done()
		}()
	}

	// Wait for the goroutines to finish.
	wg.Wait()
	fmt.Println("Final Counter:", counter)
}
```

`for g := 0; g < grs; g++ {`
This loop twice: perform a read to a local counter, increase by 1, write it back
to the shared state. Every time we run the program, the output should be 4.
The data races that we have here is that: at any given time, both Goroutines
could be reading and writing at the same time. However, we are very lucky in
this case. What we are seeing is that, each Goroutine is executing the 3
statements atomically completely by accident every time this code run.
If we put the line `fmt.Println(value)`, this will trigger the data race to
happen. Once we read the value out of that shared state, we are will force the
context switch. Then we come back, we are not getting 4 as frequent.

The print statement is a system call. The scheduler going "oh, you want to make
a system call? Well we're going to now move you from a running state into a
waiting state." We just got, context switch.

### Race detector tool

To identify race condition:

```sh
$ go run -race <file_name>

$ ./example1
==================
WARNING: DATA RACE
Read at 0x0000005fa138 by goroutine 8:
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/example1/example1.go:35 +0x47

Previous write at 0x0000005fa138 by goroutine 7:
  main.main.func1()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/example1/example1.go:41 +0x63

Goroutine 8 (running) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/example1/example1.go:31 +0xab

Goroutine 7 (finished) created at:
  main.main()
      /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/concurrency/data_race/example1/example1.go:31 +0xab
==================
Final Counter: 4
Found 1 data race(s)
```

The general rule is this; if the race detector finds a race, you have a race;
fix it. But if the race detector doesn't find a race, it doesn't mean you don't
have one, it just hasn't found it yet.

I've been in environments where, on my Mac, the race detector doesn't find
anything because the scheduler isn't creating enough, let's say, wacky
concurrency or context switches. Then the code moves to CI on some Linux
platform and CI keeps failing on the race detection.

You really should be running your tests with the race detector on, especially if
you're writing multi-threaded software.

## Synchronization with Atomic Functions

Start with the atomic instructions to clean this code up.

[Sample program](example2/example2.go).

```go
import (
	"fmt"
	"sync"
	"sync/atomic"
)

// counter is a variable incremented by all goroutines.
var counter int64
```

Notice that it's not just an `int` but `int64`. When we use the atomic package
we need precision based integers. We can't just use `int` anymore. We have to
specify specifically are we using 64 bit or 32 bit integers.

```go
func main() {
	// Number of goroutines to use.
	const grs = 2

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(grs)

	// Create two goroutines.
	for g := 0; g < grs; g++ {
		go func() {
			for i := 0; i < 2; i++ {
				atomic.AddInt64(&counter, 1)
			}

			wg.Done()
		}()
	}

	// Wait for the goroutines to finish.
	wg.Wait()

	// Display the final value.
	fmt.Println("Final Counter:", counter)
}
```

`atomic.AddInt64(&counter, 1)` safely add one to counter.
Add the atomic functions that we have take an address as the first parameter and
that is being synchronized, no matter many Goroutines they are. If we call one
of these function on the same location, they will get serialized. This is the
fastest way to serialization. We can run this program all day long and still
get 4 every time.

## Synchronization with Mutexes

Mutex allows us to create a block of code, or multiple lines of code, and treat
that as one atomic operation.

Mutex allows us to have the API like the `WaitGroup` where any Goroutine can
execute one at a time.

[Sample program](example3/example3.go).

```go
import (
	"fmt"
	"sync"
)

// counter is a variable incremented by all goroutines.
var counter int
```

```go
// mutex is used to define a critical section of code.
var mutex sync.Mutex
```

I like to think of mutexes as creating rooms in your code where all Goroutines
have to go through. However, only one Goroutine can go at a time. The scheduler
will decide who can get in and which one is next.

The scheduler kind of acts like a bouncer, just like if you went to a bar.

We cannot determine what the scheduler will do. Hopefully, it will be fair.
Just because one Goroutine got to the door before another, it doesn't mean that
Goroutine will get to the end first. Nothing here is predictable.

The key here is, once a Goroutine is allowed in, it must report that it's out.
All the Goroutines come in will ask for a lock and unlock when it leave for
other one to get in.

Two different functions can use the same mutex which means only one Goroutine
can execute any of given functions at a time.

```go
func main() {
	// Number of goroutines to use.
	const grs = 2

	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(grs)

	// Create two goroutines.
	for g := 0; g < grs; g++ {
		go func() {
			for i := 0; i < 2; i++ {
				// Only allow one goroutine through this critical section at a time.
				// Creating these artificial curly brackets gives readability.
				// We don't have to do this but it is highly recommended.
				// The Lock and Unlock function must always be together in line of sight.
				mutex.Lock()
				{
					// Capture the value of counter.
					value := counter

					// Increment our local value of counter.
					value++

					// Store the value back into counter.
					counter = value
				}
				mutex.Unlock()
				// Release the lock and allow any waiting goroutine through.
			}

			wg.Done()
		}()
	}

	// Wait for the goroutines to finish.
	wg.Wait()
	fmt.Printf("Final Counter: %d\n", counter)
}
```

Rule: The same function that calls `Lock()` must call `Unlock()`. If you miss a
`Lock()` or an `Unlock()`, you're will have a deadlock situation. All these
Goroutines are stuck because they can't get in. We want the `Lock()` and the
`Unlock()` to always be together.

## Read/Write Mutex

What read/write mutex does is, it allows us to have multiple reads across that
one write.

Sometimes we have a shared resource where we want many Goroutines reading it.
Every now and then, one Goroutine can come in and make change to the resource.
When that happens, everybody has to stop reading. It doesn't make sense to
synchronize reads in this type of scenario because we are just adding latency
to our program for no reason.


[Sample code](example4/example4.go).

```go
// data is a slice that will be shared.
var data []string
```

```go
// rwMutex is used to define a critical section of code.
var rwMutex sync.RWMutex
```

RWMutex is a little bit slower than Mutex but we are optimizing for correctness
first so we don't care about that for now.

```go
// Number of reads occurring at ay given time.
var readCount int64
```

As soon as we see `int64` here, we should start thinking about using atomic
instruction.

```go
// init is called prior to main.
func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	// wg is used to manage concurrency.
	var wg sync.WaitGroup
	wg.Add(1)

	// Create a writer goroutine that performs 10 different writes.
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			writer(i)
		}
		wg.Done()
	}()

	// Create eight reader goroutines that runs forever.
	for i := 0; i < 8; i++ {
		go func(id int) {
			for {
				reader(id)
			}
		}(i)
	}

	// Wait for the write goroutine to finish.
	wg.Wait()
	fmt.Println("Program Complete")
}

// writer adds a new string to the slice in random intervals.
func writer(i int) {
	// Only allow one goroutine to read/write to the slice at a time.
	rwMutex.Lock()
	{
		// Capture the current read count.
		// Keep this safe though we can due without this call.
		// We want to make sure that no other Goroutines are reading. The value of rc should always
		// be 0 when this code run.
		rc := atomic.LoadInt64(&readCount)

		// Perform some work since we have a full lock.
		fmt.Printf("****> : Performing Write : RCount[%d]\n", rc)
		data = append(data, fmt.Sprintf("String: %d", i))
	}
	rwMutex.Unlock()
	// Release the lock.
}

// reader wakes up and iterates over the data slice.
func reader(id int) {
	// Any goroutine can read when no write operation is taking place.
	// RLock has the corresponding RUnlock.
	rwMutex.RLock()
	{
		// Increment the read count value by 1.
		rc := atomic.AddInt64(&readCount, 1)

		// Perform some read work and display values.
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
		fmt.Printf("%d : Performing Read : Length[%d] RCount[%d]\n", id, len(data), rc)

		// Decrement the read count value by 1.
		atomic.AddInt64(&readCount, -1)
	}
	rwMutex.RUnlock()
	// Release the read lock.
}
```

There's a real cost to mutexes and atomic functions, and that's latency. Latency
can be good when we have to coordinate orchestrating. But, if we can reduce
latency using Read/Write Mutex, life is better.

If we are using mutex, make sure that we get in and out of mutex as fast as
possible. Don't do anything extra. Sometimes just reading the shared state into
a local variable is all we need to do. The less operation we can perform on the
mutex, the better. We then reduce the latency to the bare minimum.

## Map Data Race

Most people in Go don't realize is that accessing a map is not inherently
synchronous. We don't get that for free. You don't get anything for free when it
comes to synchronization and orchestration. You're responsible for it.

Because so many people have had data races with maps, the language now has build
into the run time data race detection for map access. Go cares about integrity
over everything. There's a cost to integrity, which is performance. But they've
really made sure that this little bit of race detection with map access isn't
going to hurt you.

Let's take a look at how that works with this [program](example5/example5.go).

```go
// scores holds values incremented by multiple goroutines.
var scores = make(map[string]int)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 1000; i++ {
			scores["A"]++
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			scores["B"]++
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("Final scores:", scores)
}
```

I've created a global variable game called `scores`. You shouldn't be doing this
in production code.

There are 2 Goroutines trying to write to the same piece of data technically at
the same time. This is a concurrent write. This is a data race.

Now prior to the version of Go we're running on here, this code would have ran
in production and you would have had data corruption everywhere and we wouldn't
have known it until after the fact.

## Interface-Based Race Condition

[Sample program](advanced/example1/example1.go) to show a more complicated race
condition using an interface value. This produces a read to an interface value after
a partial write.

This code doesn't blow up.

The data races in this case was two writes that happened at the same time that
were not synchronized. Line 78 and line 67.
