---
title: Pointers
weight: 3
---

# Pointers

This is one of the most important sections in the class. It help us learn how
we can look at the impact that our code is having on the machine.

### Everything in Go is pass by value

All of the code you're writing at some point gets into machine code, and the
OS's job is to choose a path of execution, a thread to execute those
instructions one after the other. What's important here is the data structures
now. There's 3 areas of memory:
- data segment: usually reserved for your global variables, your read-only values.
- stacks: a data structure that every thread is given. At the OS level, your
stack is a contiguous block of memory and usually it's allocated to be 1MB.
- heaps

The diagram:
- "M": an operating system thread, is a path of execution, and it has a stack,
and it needs that stack in order to do its job.
- "G": Goroutine, which is our path of execution at Go's level that contains
instructions that needed to be executed by the machine. You can also think of
Goroutine as a lightweight thread.
"G"s are very much like "M"s, we could almost, at this point, say that they're
the same but this is above the OS.

### Goroutine

When the program starts up, the Go runtime creates a Goroutine.

By the time the goroutine that was created for this Go program wants to execute
main, it's already executed a bunch of code from the run time.

Any time a goroutine makes a function call, what it's going to do is take some
memory off the stack. We call this a frame of memory. It will slice a frame
of memory off the stack.

The stack memory in Go starts out at 2K. It is very small. It can change over
time. The growing direction of the stack is downward.

### Stack frame

Every function is given a stack frame, memory execution of a function.
The size of every stack frame is known at compiled time. No value can be
placed on a stack unless the compiler knows its size ahead of time. If we don't
know the size of something at compiled time, it has to be on the heap.

Remember, we have the concept of zero value. Zero value enables us to initialize
every stack frame that we take. Stacks are self cleaning. We clean our stack on
the way down. Every time we make a function, zero value initialization cleaning
stack frame. Memory below the active frame no longer has integrity because it's
going to be reused. We leave that memory on the way up because we don't know if
we would need that again.

### Program boundaries

Every time you make a function call, we're crossing over a program boundary.
We can also have a boundary between goroutines when we will discuss it later.

## Pass by value

Passed by value means we make copies and we store copies. Frame allows the
goroutine to mutate memory without any cause of side effects throughout
the program.

### Value and pointer semantics behavior

If you want to write code in Go that is optimized for correctness, that you can
read and understand the impact of things, then your value and your pointer
semantics are everything.

## Sharing data

Pointer semantics serve one purpose and that is to share our piece of data
across a program boundary.

## Escape analysis

Value semantics have the benefit of being able to mutate memory and isolation
within our own sandboxes, but it has the cost of inefficiency. We will have
lost a copy of the data as we cross over these program boundaries.
Pointer semantics, however, fix the efficiency issue.

If we balance our value and our pointer semantics properly, leveraging the
aspects of the language helping the cognitive load over memory management, it's
going to be a lot better for us.

### Factory functions

We don't have constructors in Go. We don't want that. It hides cost, but what we
do have is what we call factory functions.

Factory function is a function that creates a value, initializes it for use,
and returns it back to the caller. This is great for readability, it doesn't
hide cost, we can read it, and lends to simplicity in terms of construction.

### The ampersand operator

It is very powerful from a readability standpoint. Ampersand means sharing.

### Static code analysis

The compiler is able to perform static code analysis called escape analysis.
Escape analysis determine whether a value gets to be placed on the stack, or
it escapes to the heap.

Our first priority is that a value stays on the stack. This is because that
memory is already there. It's very very fast to leverage the stack. Also stacks
are self-cleaning, which means that the garbage collector doesn't even get
involved.

An allocation in Go is when an escape analysis determines that a value cannot
be constructed on the stack, but has to be constructed on the heap.

### Sharing tells us everything

The way escape analysis works is it doesn't care about construction.
Construction in Go tells you nothing. What tells us everything is how a
value is shared.

### Mixing semantics

Example of clever code: during the construction, I am telling the compiler I don't want you to
be a value of type `user`. I want it to be a pointer to the value that we are
constructing. This is nightmare.

We are using pointer semantics during construction, even though we're creating
a variable, and now we've made this code much harder to read, and we're really
also mixing semantics as we go along the way. Anytime you mix semantics we're
going to have a problem.

Make sure that we're using the right semantics and semantic consistency
all of the time.

### Escape analysis report

When you use the `gcflags` on the `go build` calls, what you will get is
not a binary, but we will get the escape analysis report. This report tells
us why something is allocating.

## Stack growth

There's another part of allocation in Go, and that is that if the compiler
doesn't know the size of a value at compile time, it must immediately construct
it on the heap. Frames are not dynamic, so if the compiler doesn't know the
size of something at compile time, it cannot place it on the stack.
The compiler knows the size of a lot of things at compile time. It knows:
- struct types
- built in types

Sometimes you might have things like collections that are based on, their size
is based on a variable, which gives the compilers no idea what the size of
that is.

Go stack is 2K, being very small. What happens when you've got a Go routine
that's making lots of function calls and eventually it runs out of stack space?

Get a new stack.

Basically, imagine that we had our stack, we had some value there, and imagine
we were even sharing this value as we move down the call stack. Eventually, we
run out of stack space.

### Go contiguous stacks

What it's going to do is allocate a larger stack, 25% larger than the original
one, and then, what it's got to do is copy all these frames back over, in this
case, these pointers are relative, so they're very fast to fix. But, basically,
Go routine, during the function call, going to take a little latency hit on
creating the larger stack, copying those frames over, and readjusting any
of these pointers.

This isn't something that's going to happen all of the time. 2K is usually more
than enough for our stacks, because you usually don't go more than even like
10 function calls deep, you don't. There's other optimizations the compiler can
do to keep these frames very small.

When this happen, values on your stack can potentially be moving around.
This is a whole new world.

Code example:
```go
// Number of elements to grow each stack frame.
// Run with 1 and then with 1024.
const size = 1

// main is the entry point for the application.
func main() {
	s := "HELLO"
	stackCopy(&s, 0, [size]int{}) // we know the size of an array at compile time
}

// stackCopy is a recursive function.
// It calls itself over and over and over again, constantly sharing this
// string down the call stack, increasing the size of the stack.
func stackCopy(s *string, c int, a [size]int) {
	println(c, s, *s)

	c++
	if c == 10 {
		return
	}

	stackCopy(s, c, a)
}
```

The side effect is, since a value can move in memory that's on the stack,
this actually creates an interesting constraint for us in Go.
What this means is, is no stack can have a pointer to another stack.
Imagine, we had all of these stacks all over the place, hundreds of thousands
of Go routines with pointers to each other's stacks. That would be total chaos
if one stack had to grow.

### Local pointers

Since our stacks can move, it means that the only pointers to a stack would be
local pointers. Only that stack memory is for the Go routine. Stack memory
cannot be shared across Go routines.

The heap basically now is used for:
- any value that's going to be shared across Go routine boundaries
- any value that casting on the frame because there's an integrity issue
- any value where we don't know the size at compile time.

What we care about is not that our code is the fastest it can be,
we care about is it fast enough.

## Garbage Collection (GC)

### The design of the Go GC

Go 1.10: It's call a tri-color mark and sweep concurrent collector.
It's not a compacting garbage collector, memory on our heap does not move
around, which is getting interesting because memory on our stacks potentially
are. Once an allocation is made on the heap, it stays there until it gets
swept away.

### Pacing algorithm

Everything begins and ends with the pacing algorithm.
The GC has an algorithm and the pacing algorithm.

The pacing algorithms trying to do is balance 3 things. How do we maintain the
smallest heap size run at a reasonable pace where the stop the word latency time
is under a 100 microseconds and were able to even leverage less than or up to
25% of your available CPU capacity.

### CPU capacity cost

Where could the 25% come from? The garbage collector uses Go heap as well. Go is
written in Go the runtime and the compiler.

### Different types of garbage collectors

Each one has their own thing where the GC maybe run at a high level of
performance to get done quickly. Go is about the lower latency and we all just
run together and we do things that are very constant and consistent pace.

### Heap size

Diagram:
The size of your heap and the live heap.
Live heap contains maybe a map of a caching system.

We're trying to maintain the smallest possible at a reasonable pace, so the
stopped the world (STW) latency maintains itself in a 100 us or less.

As your program's running, the live heap is moving close to the top of the heap.
At some point, if it gets close enough, we have to bring it back down.

We can't let the live heap get all the way to the top of the heap because if we
will want to run it concurrently, that means that we would blow by it and
anytime the live heap passed beyond the scope of the size of the heap, there's
one configuration option in Go called `GOGC` and the default is 100 and that
means we will have 100% growth on that heap when the live heap has to go
by it.

Chart:
Shows the different areas of the garbage collector and where some of that
STW time is. During GC we have a very quick STW and that's to turn the right
barrier on. The write barrier one should be really quick.

### Write barrier

The idea of the write barrier is that these Go routines that are running
essentially need to report in what they're doing.

From Go 1.10 and before the only way to stop a process or to bring it to that
safe point is to wait for a Go routine to make a function call.
Scheduling happens during function calls, this because we have a cooperative
scheduler, not a preemptive scheduler.

### The heap is a very large graph

We have two things. We have our stacks, and our stacks have frames and in
some cases these frames are going to have values that point to values on the
heap. You will have other values here and some values can point to
other values.

From a tri-color perspective, we turn all of all of these values,
the stacks and these objects. They all start out as white. Iterate over this
entire tree that when we're done, all we have left are black values or white
values. Anything that's black has to stay in memory because there's a reference
to it from a stack.

### Balance value in pointer semantics

We have to leverage value semantics to their fullest extent and know when to
use the pointer semantics, understand the costs and the benefits of these things
and try to reduce the amount of allocations our program is having.

You will not write zero allocation software. We're not trying to prevent it.
We're trying to reduce it. Less is always more and if there's less work for
the GC to do, this is all going to happen much faster.

A larger heap doesn't necessarily mean better performance because it just means
that when the live heap gets to the top, it's got that much more work to get
to the bottom. So we don't really play with configuration here. We let the
pacing algorithm do it.

We write software that's consistent in allocations. We want to reduce
allocations, let the profiler but benchmarking tell us where there are wasted
or unnecessary allocations. Reduce that and just learn how to maintain great
balance between our value and our pointer semantics as we write code.
