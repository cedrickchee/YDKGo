---
title: Methods
weight: 9
---

# Methods

Methods are functions that give data the ability to exhibit behavior.

- Methods are functions that declare a receiver variable.
- Receivers bind a method to a type and can use value or pointer semantics.
- Value semantics mean a copy of the value is passed across program boundaries.
- Pointer semantics mean a copy of the values address is passed across program boundaries.
- Stick to a single semantic for a given type and be consistent.

We will start talking about methods, but we really will start talking about is
decoupling. Up until now, we've been really talking about the concrete data.

```
//      --------------
//  ^   | Behavior   |
//  |   | Decoupling |
//  |   |------------|
//  |   | Concrete   |
//  |   | Data       |
//      --------------
```

## Declare and Receiver Behavior

Go has functions, but Go also has the concept of a method.
A method is a function, a function that has what we call a receiver.

[Sample program](example1/example1.go) to show how to declare methods and how
the Go compiler supports them.

```go
// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements a method with a value receiver: u of type user
func (u user) notify() {
	fmt.Printf("Sending User Email To %s<%s>\n",
		u.name,
		u.email)
}
```

In Go, a function is called a method if that function has receiver declared
within itself. It looks and feels like a parameter but it is exactly what it is.
Using the value receiver, the method operates on its own copy of the value that
is used to make the call.

```go
// changeEmail implements a method with a pointer receiver: u of type pointer user
func (u *user) changeEmail(email string) {
	u.email = email
}
```

Using the pointer reciever, the method operates on shared access.

These 2 methods above are just for studying the difference between a value
receiver and a pointer receiver. In production, we will have to ask ourself why
we choose to use inconsistent receiver's type.

Should I be using a value receiver or should I be using a pointer receiver?
A lot of Go developers get stuck and they don't know what to do. Too many Go
developers, especially early on, start doing this — if the method I am writing
has to mutate the data, I'll use a pointer receiver. I'll share the data in. But
if the method doesn't have to mutate, I'll use a value receiver. I'll let it
just have a copy. This is very bad, we cannot do this.

```go
// Values of type user can be used to call methods declared with both value
// and pointer receivers.
bill := user{"Bill", "bill@email.com"}
bill.changeEmail("bill@hotmail.com")
bill.notify()

// Pointers of type user can also be used to call methods declared with both
// value and pointer receiver.
joan := &user{"Joan", "joan@email.com"}
joan.changeEmail("joan@hotmail.com")
joan.notify()
```

`joan` in this example is a pointer that has the type `*user`. We are still
able to call `notify`. This is still correct. As long as we deal with the
type `user`, Go can adjust to make the call.

Behind the scene, we have something like `(*joan).notify()`. Go will take the
value that `joan` points to and make sure that `notify` leverages its value
semantic and works on its own copy.

Similarly, `bill` has the type `user` but still be able to call `changeEmail`.
Go will take the address of `bill` and do the rest for you: `(*bill).changeEmail()`.

## Value and Pointer Semantics

General guidelines on when to be using value semantics and when to be using
pointer semantics. There are exceptions to everything. Remember, semantic
consistency is everything.

1. If you are working with the built-in types, strings, numerics, and bool, you are
to be using value semantics. Including fields in a struct.
2. Value semantics for reference types (slice, map, channel, interface). (There's
one exception to this, however. A slice and a map, you may take the address of
a slice or a map only if you're sharing it down the call-stack and to a function
that's either named decode or un-marshall.)
3. User-defined types/struct types: you've got to choose, if you're defining
your own struct type, what semantic is going to be in play. If you're not sure
what to use then we're going to use those pointer semantics. Then if you're
absolutely sure that we can use value semantics, we want to use those value
semantics. They're really our first choice.

Take a look at some standard library [code](example5/example5.go). These is a
named type from the net package.

```go
type IP []byte // new type named IP which is based on slice of bytes.
type IPMask []byte
```

`IP` and `IPMask` are reference types. Since we use value semantics for
reference types, the implementation is using value semantics for both.

```go
func (ip IP) Mask(mask IPMask) IP {}
```

`Mask` is using a value receiver and returning a value of type `IP`. This method
is using value semantics for type `IP`.

```go
func ipEmptyString(ip IP) string {
	if len(ip) == 0 {
		return ""
	}
	return ip.String()
}
```

`ipEmptyString` accepts a value of type `IP` and returns a value of type string.
The function is using value semantics for type `IP`.

### Pointer semantic

We pulled this out of the standard library Time package.

```go
type Time struct {
	sec  int64
	nsec int32
	loc  *Location
}

func Now() Time {
	sec, nsec := now()
	return Time{sec + unixToInternal, nsec, Local}
}
```

Should `Time` use value or pointer semantics?
If you need to modify a time value should you mutate the value or create a
new one?

The best way to understand what semantic is going to be used is to look at the
factory function for type. It dictates the semantics that will be used.
In this example, the `Now` function returns a value of type `Time`.
It is making a copy of its `Time` value and passing it back up.
This means `Time` value can be on the stack. We should be using value semantic
all the way through.

```go
// Add is using a value receiver and returning a value of type Time. This
// method is using value semantics for Time.
func (t Time) Add(d Duration) Time {
	t.sec += int64(d / 1e9)
	nsec := int32(t.nsec) + int32(d%1e9)
	if nsec >= 1e9 {
		t.sec++
		nsec -= 1e9
	} else if nsec < 0 {
		t.sec--
		nsec += 1e9
	}
	t.nsec = nsec
	return t
}
```

Add is a mutation operation. If we go with the idea that we should be using
pointer semantic when we mutate something and value semantic when we don't,
then `Add` is implemented wrong. However, it has not been wrong because it is
the type that has to drive the semantic, not the implementation of the method.
The method must adhere to the semantic that we choose. `Add` is using a value
receiver and returning a value of type `Time`. It is mutating its local copy and
returning to us something new.

Here's another example:

```go
// div accepts a value of type Time and returns values of built-in types.
// The function is using value semantics for type Time.
func div(t Time, d Duration) (qmod2 int, r Duration) {
	// Code here
}
```

```go
// The only use pointer semantics for the `Time` api are these
// unmarshal related functions.
func (t *Time) UnmarshalBinary(data []byte) error {
func (t *Time) GobDecode(data []byte) error {
func (t *Time) UnmarshalJSON(data []byte) error {
func (t *Time) UnmarshalText(data []byte) error {
```

I told you there are exceptions to everything, and here's one of those exceptions.

When you are implementing or using functions that have un-marshall or decode in
it they require naturally to have pointer semantics. You need to share because
they're gonna be mutated. So here are four APIs in the Time package that are
switching semantics from value to pointer.

Here is a factory function for the file, for OS:

```go
// Factory functions dictate the semantics that will be used. The Open function
// returns a pointer of type File. This means we should be using pointer
// semantics and share File values.
func Open(name string) (file *File, err error) {
	return OpenFile(name, O_RDONLY, 0)
}
```

Notice that this factory function returns not a value of type file but a pointer
of type file. It means that pointer semantics are at play. It also means that
you do not have a right to make a copy of a value that a pointer points to when
it's been shared with you like that. Assume that it is dangerous to make copies
if something has been shared.

### Function/Method Variables

Methods are really just made up. They are not real. All we have is state and
all we have is functions. Methods give us this syntactic sugar, that this belief
system that a piece of data has behavior.

[Sample program](example3/example3.go) to show how to declare function variables.

```go
// data is a struct to bind methods to.
type data struct {
	name string
	age  int
}

type bill data
```

Does this new type `bill` have behavior like `data` does?

There is no behavior associated with `bill` data. All of the behavior is only
still associated with `data`. Why is this? One, because these methods are not
declared inside the type, like you would see in a class-based/OOP system,
they're declared outside the type. This is Go, separating state from behavior.
There is no OOP in GO. So, understand that even though `bill` is based
on `data`, the basing is based on its memory model, not on this idea of behavior.

### Methods are just functions

```go
// Declare a variable of type data.
d := data{
    name: "Bill",
}

fmt.Println("Proper Calls to Methods:")

// How we actually call methods in Go.
d.displayName()
d.setAge(45)

fmt.Println("\nWhat the Compiler is Doing:")

// This is what Go is doing underneath.
data.displayName(d)
(*data).setAge(&d, 45)
```

This is what Go is doing underneath.
When we call `d.displayName()`, the compiler will call `data.displayName`,
showing that we are using a value receiver of type `data`, and pass the `data`
in as the first parameter.
Taking a look at the function again: `func (d data) displayName()`, that
receiver is the parameter because it is truly a parameter. It is the first
parameter to a function that call `displayName`.
Similar to `d.setAge(45)`. Go is calling a function that based on the pointer
receiver and passing `data` to its parameters. We are adjusting to make the call
by taking the address of `d`.

### Function variable

```go
fmt.Println("\nCall Value Receiver Methods with Variable:")

f1 := d.displayName
```
Declare a function variable for the method bound to the `d` variable.
The function variable will get its own copy of `d` because the method is using
a value receiver.
`f1` is now a reference type: a pointer variable. We don't call the method here.
There is no `()` at the end of `displayName`.

```go
f1()
```

Call the method via the variable.

```
// We have a level of indirection (decoupling), two pointers to get to the code

//  -----
// |  *  | --> code
//  -----
// |  *  | --> copy of d
//  -----
```

`f1` is pointer and it points to a special 2 word data structure. The first word
points to the code for that method we want to execute, which is `displayName` in
this case. We cannot call `displayName` unless we have a value of type data.
So the second word is a pointer to the copy of data. `displayName` uses a value
receiver so it works on its own copy. When we make an assignment to `f1`, we are
having a copy of `d`.

```go
// Change the value of d.
d.name = "Joan"

// Call the method via the variable. We don't see the change.
f1()
```

When we change the value of `d` to "Joan", `f1` is not going to see the change.

However, if we do this again if `f2`, then we will see the change:

```go
fmt.Println("\nCall Pointer Receiver Method with Variable:")

// Declare a function variable for the method bound to the d variable.
// The function variable will get the address of d because the method
// is using a pointer receiver.
f2 := d.setAge

// Call the method via the variable.
f2(45)

// Change the value of d.
d.name = "Sammy"

// Call the method via the variable. We see the change.
f2(45)
```

```
// 2 word data structure

//  -----
// |  *  | --> code
//  -----
// |  *  | --> original d
//  -----
```

`f2` is also a pointer that has 2 word data structure. The first word points to
`setAge`, but the second words doesn't point to its copy any more, but to its
original.

## Escape analysis flaw

We have double indirection to get to the data, and the escape analysis algorithm
has a flaw in that, once we have indirection like this, the escape analysis
can't track whether or not this value can stay on a stack or not. In other
words, even though there's no reason for `d` to end up on the heap, it has to
now allocate.
