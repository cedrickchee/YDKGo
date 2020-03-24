---
title: Variables
weight: 1
---

# Variables

## Type

- Type is everything. Type is life.
- Our basic unit of memory that we will be working with in our programming model is a byte.
- Type provides us two pieces of information, size, how much memory, how many bytes of memory we need to read or write at any given time, and its representation, what it represents.
- Go has the built-in types: numerics, string, bool.
- When you look at the name `float64`, it's really interesting because it tells us both parts of the type information, `float64`, 64 tells us that it's an eight byte, 64 bit value that's giving us the cost in terms of memory footprint, and float tells us it's an IEEE 754 binary decimal.
- When we declare a type without being very specific, such as `uint` or `int`, it gets mapped based on the architecture we are building the code against. On a 64-bit OS, `int` will map to `int64`. Similarly, on a 32 bit OS, it becomes `int32`.
- The word size is the number of bytes in a word, which matches our address size. For example, in 64-bit architecture, the word size is 64 bit (8 bytes), address size is 64 bit then our integer should be 64 bit.

## Zero value concept

- Zero value's very very important, and it's an integrity play in Go.
- Zero value is all memory that we allocate gets initialized at least to its zero value state.
- Every single value we create must be initialized. If we don't specify it, it will be set to the zero value. The entire allocation of memory, we reset that bit to 0.
```go
Type            Initialized Value
Boolean         false
Integer         0
Floating Point  0
Complex         0i
String          "" (empty string)
Pointer         nil
```

- Strings are a series of `uint8` types.
- A string is a two word data structure: first word represents a pointer to a backing array, the second word represents its length. If it is a zero value then the first word is nil, the second word is 0.

## Declare and initialize

- `var` is that kind of readability marker that we are declaring and initializing to zero value, and `var` gives us zero value 100% of the time in this language.
```go
var a int
var b string
```

- Short variable declaration operator as a productivity operator — we can define and initialize at the same time.
```go
cc := 3.14159
```

## Conversion versus casting

- Go doesn't have casting, it has conversion.
- Instead of telling a compiler to pretend to have some more bytes, we have to allocate more memory.
```go
// Specify type and perform a conversion.
aaa := int32(10)
```
