---
title: Struct Types
weight: 2
---

# Struct Types

## Declare and initialize

Declare a variable of type `example` set to its zero value.
```go
type example struct {
    flag    bool
    pi      float32
    counter int64
}

var e1 example
```

How much memory do we allocate for `example`?
- a `bool` is 1 byte
- `int16` is 2 bytes
- `float32` is 4 bytes

Putting together, we have 7 bytes. However, the actual answer is 8. That leads
us to a new concept of padding and alignment.

### Padding and alignment

The padding byte is sitting between the `bool` and the `int16`. The reason is
because of alignment.

The idea of alignment: It is more efficient for this piece of hardware to read
memory on its alignment boundary. We will take care of the alignment boundary
issues so the hardware people don't.

**Rule 1:**
Depending on the size a particular value, Go determines the alignment we need.
Every 2 bytes value must follow a 2 bytes boundary. Since the `bool` value is
only 1 byte and start at address 0, then the next `int16` must start on
address 2. The byte at address that get skipped over becomes a 1 byte padding.
Similarly, if it is a 4 bytes value then we will have a 3 bytes padding value.

**Rule 2:**
The largest field represents the padding for the entire struct. We need to
minimize the amount of padding as possible. Always lay out the field from
highest to smallest. This will push any padding down to the bottom.

Example:
```go
// In this case, the entire struct size has to follow a 8 bytes value because int64 is 8 bytes.
type example struct {
    counter int64
    pi      float32
    flag    bool
}
```

## Name type versus anonymous type

If we have two name type identical struct, we can't assign one to another.

For example, `example1` and `example2` are identical struct,
`var ex1 example1`, `var ex2 example2`.
`ex1 = ex2` is not allowed.
We have to explicitly say that `ex1 = example2(ex2)` by performing a conversion.

```go
// Declare a variable of an anonymous type and init using a struct literal.
e3 := struct {
    flag    bool
    counter int16
    pi      float32
}{
    flag:    true,
    counter: 10,
    pi:      3.141592,
}
```

However, if `ex` is a value of identical anonymous struct type (like `e3`
above), then it is possible to assign `ex1 = ex`

```go
var e4 example
e4 = e3 // OK
```
