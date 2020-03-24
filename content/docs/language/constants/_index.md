---
title: Constants
weight: 4
---

# Constants

## Constants are not variables

One of the really interesting things about constants, for me is that they only exist at compile time. Constants are something that really have a much different feel and flavor to than your common variables.

Most of the time we think about constants as being read-only variables, and it's absolutely not the case in Go.

## Two types of constants

There's two types of constants, constants of a kind and constants of a type.

```go
// Untyped Constants.
const ui = 12345    // kind: integer
const uf = 3.141592 // kind: floating-point

// Typed Constants.
const ti int = 12345        // type: int
const tf float64 = 3.141592 // type: float64
```

The difference between constants of a kind and of a type, are that constants of a kind can be implicitly converted by the compiler.

## Explicit and implicit conversions

**Kind Promotion** will tell you how things promote so floats promote over ints and types always promote over kind.

```go
// Variable answer will of type float64.
var answer = 3 * 0.333 // KindFloat(3) * KindFloat(0.333)
```

We're dealing with 256 bits of precision when we're dealing with constants of a kind, and when we now convert back to a variable, we're moving that down to a 64 bit level of precision. There will be some precision loss, but, remember that floating points already are already not precise, IEEE754 binary decimals.

```go
// Constant third will be of kind floating point.
const third = 1 / 3.0 // KindFloat(1) / KindFloat(3.0)
```

In the old days we used to call what we considered constants of a kind to be exact, they were like these very exact numbers because we had such high levels of precision that they were exact. So, we would look at third truly as 1/3 even though, eventually it turns to 56 bits but there be some precision loss.

```go
// Constant zero will be of kind integer.
const zero = 1 / 3 // KindInt(1) / KindInt(3)
```

But then on above code, you could see that there's no promotion going on. One is of kind int, three is of kind int, we do the division, we end up with zero, because that's what's that's going to be, everything stays within the kind integer.

**Promoto from Kind to Type**

```go
// This is an example of constant arithmetic between typed and
// untyped constants.
const one int8 = 1
const two = 2 * one // int8(2) * int8(1)
```

`two` ends up being a constant of type `int8`.

**Parallel type system**

```go
const (
	// Max integer value on 64 bit architecture.
	maxInt = 9223372036854775807

	// Much larger value than int64.
	bigger = 9223372036854775808543522345

	// Will NOT compile
	// Compiler: "constant 9223372036854775808543522345 overflows int64"
	// biggerInt int64 = 9223372036854775808543522345
)
```

**Practical use of constants**

See the power of constants and their use in the standard library.

```go
type Duration int64

const (
        Nanosecond  Duration = 1
        Microsecond          = 1000 * Nanosecond
        Millisecond          = 1000 * Microsecond
        Second               = 1000 * Millisecond
        Minute               = 60 * Second
        Hour                 = 60 * Minute
)
```

This is a second way to declare a type here in Go, what I would say is, the name type `Duration` is based, based on `int64`. This is **not an alias**, we have really two distinct named types here. We're just using `int64` as our base information or our base memory model for `Duration`. And I only want to do these types of things when the new type has its own representation and meaning, and it does here in the time package.

`Duration` represents time. Doesn't represent an `int64`, it represents nanoseconds of time. We look at these constants because this is a real practical and clean way of how this idea of type and kind work together in constants.

**Reasons why we can't have enumerations in Go**

We don't want to create types as aliases to get compiler protection when they're based on let's say, those built-in types, which is where constants are allowed to be. Constants can only be based on the built-in types because again, they only exist at compile time.

## iota

```go
const (
    A2 = iota // 0 : Start at 0
    B2        // 1 : Increment by 1
    C2        // 2 : Increment by 1
)
```

We only have to assign iota one time to the very first constant in the block and we will automatically, for free, get the incremental.

iota is a very powerful mechanism if you're creating a set of constants that are going to have some unique IDs and it just kind of let's the language set all that up for you.

**Wrap-up**

- Two types of constants, constants of a kind and constants of a type.
- Literal values in Go are constants of a kind, they're unnamed constants.
- Constants of a kind can be implicitly converted by the compiler.
- Constants of a kind can have up to 256 bits of precision.
