---
title: Arrays
weight: 7
---

# Arrays

## Part 1 - Mechanical sympathy

All data structures Go gives you is just the array, slice which looks like an array, and maps.

A look at a piece of [code](benchmarks/caching/caching.go) to show why Data Oriented Design matters. How data layouts matter more to performance than algorithm efficiency.

```sh
$ go test -run none -bench . -benchtime 3s
Elements in the link list 16777216
Elements in the matrix 16777216
goos: linux
goarch: amd64
pkg: github.com/cedrickchee/ultimate-go/language/arrays/benchmarks/caching
BenchmarkLinkListTraverse-4           93          36718229 ns/op
BenchmarkColumnTraverse-4             19         184652302 ns/op
BenchmarkRowTraverse-4               165          21787691 ns/op
PASS
ok      github.com/cedrickchee/ultimate-go/language/arrays/benchmarks/caching   21.069s
```

Notice the drastic change in the column traversal, we're almost talking like 3 ms, from our perspective, that's like eternity. And then the row traversal got a little faster, but again, it's very very consistent. Our link list and our row traversals are consistent, and our column traversal is all over the place. What is going on here?

### CPU cache

Animation that came from a great video from Scott Meyers where he talks about CPU caches and why they matter.
- [CPU Caches and Why You Care (18:50-20:30)](https://youtu.be/WDIkqP4JbkE?t=1129)
- [CPU Caches and Why You Care (44:36-45:40)](https://youtu.be/WDIkqP4JbkE?t=2676)

That's the relative speed of fetching memory, accessing memory, from the different caches.

Next is, [processor cache hierarchies](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/arrays/README.md#cache-hierarchies).

### Cache line

A cache line, historically, is 64 bytes of memory. You can think of all of the memory that's laid out on that machine, all of our virtual memory, and from a caching system perspective, it's the cache line.

- How can we create a situation where the cache line can be inside L1 or L2 before we need it?

This is our job, this is something now that becomes our responsibility, what we have to do is **write code that creates predictable access patterns to memory** if we want to be mechanically sympathetic with the hardware. If performance matters, then what we have to do is be much more efficient with how data gets in to the processor, not get the processors to run at higher clock speed. Predictable access patterns to memory is everything.

- How do you create a predictable access pattern to memory?

Here is the simplest way to do it. If you **allocate a contiguous block of memory**, and you walk through that memory on a predictable stride, well guess what, the prefetchers, which is little software sitting inside the processor, the **prefetchers can pick up on that data access**, and start bringing in those cache lines way ahead of when you need them. The prefetchers are everything, and we must be sympathetic with them.

### Array

Now, the cleanest and easiest way to create this predictable access pattern is to use an array, an array gives us the ability to allocate a contiguous block of memory, and when you define an array, you define it based on an element size, right, it's an array of string, it's an array of int, there's our predictable stride, every index is a predictable stride from the other.

Prefetchers love arrays, the hardware loves arrays, it really is the most important data structure there is today from the hardware. We almost don't even care what your algorithms are, an array is going to beat it.

Now there are times when maybe you're dealing with data that is so large that a linear traversal isn't going to be more efficient, but overall, if you're dealing with data, small data, that array and those predictable access patterns are going to beat out performance every time on these traversals.

### Slice

The array is not the most important data structure in Go. The slice is the most important data structure in Go.

And this isn't cause the slice uses an array underneath, technically, slice is really a vector, and if you've watched any C++ videos over the last five years of talking about performance, you will hear a lot about vectors, use vectors. Why, because, just like the slice, we're gonna be using arrays behind the scenes, we're gonna be doing those linear iterations, and we're gonna be creating predictable access patterns to memory that the prefetchers are going to pick up on.

### TLB cache

There's another cache in the hardware called the TLB.

The TLB is a very special cache that the operating system (OS) is gonna be managing. And what it does is it creates a cache of virtual addresses to OS page and physical memory locations. In other words, your program is working in virtual memory, it thinks it's got real memory, but it's not, it's given a full virtual memory, because the OS gives it that level of abstraction.

TLB cache misses:

- Cache misses can result in TLB cache misses as well.
- A TLB miss could be really, really deadly.

**Back to the piece of [code](benchmarks/caching/caching.go)**

So, if we go back to our results now, I want you to notice something here. We should be able to understand now why we see what see. Look at row traversal. Row traversal was not only the fastest, it was also incredibly consistent, and why is that? That is because, when we start walking through the matrix row by row, we're walking it down cache line by connected cache line. That row major traversal is creating a predictable access pattern to memory, and the prefetchers are coming in and reducing all of the main memory latency cost, we're gonna get the best performance I can get on my machine through row major traversal.

But why is column traversal so slow and inconsistent? I played a small game here with the matrix, I made the matrix large enough so this element and the next element not only were not really in a predictable stride, let's say, but made sure that those two elements ended up on two different OS pages. Yes, they were very far away from each other, and on different OS pages. I've got a ton of problems, this is basically pure random memory access, when we're doing column traversal. This is why our results were so inconsistent and so slow.

The link list sits somewhere in between. We're probably getting cache line misses, because this data is not guaranteed to be on a predictable stride, but we're probably getting this data all on the same page.

### Wrap-up

So now when we come back and we look at what Go's given us, arrays, slices, and maps, it all starts to make sense. We don't have link lists, and stacks, and queues, and these things, because they are not mechanically sympathetic with the hardware when your model is a real machine.

I will never say anything negative about Java, the JVM is an amazing piece of engineering, because it takes these mechanical sympathy issues and deals with them for you.

We don't have a virtual machine here. You have the real machine, which means you are more responsible for your data layouts, you're more responsible for your traversals, but that doesn't mean you have to get overwhelmed.

What it means is, if you use the slice, which is the most important data structure in Go, if you use it for the majority of your data needs, you are inherently creating predictable access patterns to memory, you are inherently being mechanically sympathetic, and by the way, the Go map also is constantly creating data underneath it that is contiguous. Think about it, your stacks, contiguous memory. Your slice, a vector, contiguous memory. A slice of values is our first choice, until it is not reasonable or practical to do so.

## Part 2 - Semantics

### Declare and initialize

```go
var strings [5]string
```

Declare an array of five strings that is initialized to its zero value.

Remember: a string is a 2 word data structure: a pointer and a length.

Since this array is set to its zero value, every string in this array is also
set to its zero value (empty string), which means that each string has the
first word pointed to nil and second word is 0.

```
40 bytes of memory [5 * 8 bytes]. 2 word data structure:

         -------------------------------
pointer  | nil | nil | nil | nil | nil |
         -------------------------------
length   |  0  |  0  |  0  |  0  |  0  |
         -------------------------------
```

```go
fruits[0] = "Apple" // a literal string named Apple
```

At index 0, a string now has a pointer to a backing array of bytes (characters
in string) and its length is 5.

**What is the cost?**

The cost of this assignment is the cost of copying 2 bytes.
We have two string values that have pointers to the same backing array of bytes.
Therefore, the cost of this assignment is just 2 words.

```
        -------------------------------
    --- |  *  | nil | nil | nil | nil |
    |   -------------------------------
    |   |  5  |  0  |  0  |  0  |  0  |
    |   -------------------------------
    |
    ---------------------
                        |
                        V
-------       ---------------------
|  *  |  ---> | A | p | p | l | e | (1)
-------       ---------------------
|  5  |
-------
```

### Iterate over the array

```go
for i, fruit := range strings {
    fmt.Println(i, fruit)
}
```

Using `range`, not only we can get the index but also a copy of the value in the
array. `fruit` is now a string value; its scope is within the for statement.
In the first iteration, we have the word "Apple". It is a string that has the
first word also points to (1) and the second word is 5. So we now have 3
different string value all sharing the same backing array.

What are we passing to the `Println` function?

We are using value semantic here. We are not sharing our string value.
`Println` is getting its own copy, its own string value. It means when we get
to the `Println` call, there are now 4 string values all sharing the same
backing array.

We don't want to take an address of a string. We know the size of a string
ahead of time.
- it has the ability to be on the stack.
- not creating allocation.
- not causing pressure on the GC.
- the string has been designed to leverage value mechanic, to stay on the stack,
out of the way of creating garbage.
- the only thing that has to be on the heap, if anything is the backing array,
which is the one thing that being shared.

### Contiguous memory allocations

```go
// Declare an array of 5 strings initialized with values.
friends := [5]string{"Annie", "Betty", "Charley", "Doug", "Edward"}

for i, v := range friends {
    fmt.Printf("Value[%s]\tAddress[%p] IndexAddr[%p]\n", v, &v, &friends[i])
}

// Prints:
// Value[Annie]    Address[0xc00009a040] IndexAddr[0xc00009c000]
// Value[Betty]    Address[0xc00009a040] IndexAddr[0xc00009c010]
// Value[Charley]  Address[0xc00009a040] IndexAddr[0xc00009c020]
// Value[Doug]     Address[0xc00009a040] IndexAddr[0xc00009c030]
// Value[Edward]   Address[0xc00009a040] IndexAddr[0xc00009c040]
```

Iterate over the array displaying the value and address of each element.
By looking at the output of this `Printf` function, we can see that this array
is truly a contiguous block of memory. We know a string is 2 word and depending
on computer architecture, it will have x byte. The distance between two
consecutive `IndexAddr` is exactly x byte. `v` is its own variable on the stack
and it has the same address every single time.


### Different type arrays

Arrays of different sizes are not of the same type.

```go
// Declare an array of 5 integers that is initialized to its zero value.
var five [5]int

// Declare an array of 4 integers that is initialized with some values.
four := [4]int{10, 20, 30, 40}

// Assign one array to the other
five = four
```

When we try to assign `four` to `five`, the compiler says that "cannot use
four (type [4]int) as type [5]int in assignment". This cannot happen because
they have different types (size and representation). The size of an array makes
up its type name: `[4]int` versus [5]int. Just like what we've seen with
pointer. The `*` in `*int` is not an operator but part of the type name.

Unsurprisingly, all array has known size at compiled time.
