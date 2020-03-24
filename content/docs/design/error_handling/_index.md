---
title: Error Handling
weight: 2
---

# Error Handling Design

The next step in learning design is error handling.

If you want to reduce the chances of legacy code, you've got to worry about
failure. Failure is always what we're coding against and error handling is
everything.

It's not about exception handling, because exception handling I think hides
context around errors.

## Default Error Values

Go errors are just values and they can be anything you need them to be. But from
an API design perspective, for me error handling is about showing respect to the
user of your API, giving your user enough context to make an informed decision
about the state of the application and giving them enough information to be able
to either recover or make a decision to shut down.

There's two ways to shut down an application in Go. You can go to OS, the OS
package, dot exit, and you can set a return code on that, that's the fastest way,
or you can call the built in function panic.

Now you'll choose one over the other depending on if you need a stack trace or
not. So if you need the stack trace you're gonna call panic, if you don't you
just call OS exit.

First, let's look at the language mechanic first on how the default error type
is implemented.

[Sample program](example1/example1.go).

```go
// http://golang.org/pkg/builtin/#error
type error interface {
	Error() string
}
```

This is pre-included in the language so it looks like an unexported type.
It has one active behavior, which is `Error` returned a string. Error handling
is decoupled because we are always working with error interface when we are
testing our code.

Go errors are just values.
We are going to valuate these through the decoupling of the
interface. Decoupling error handling means that cascading changes will bubble up
through the user application, causes cascading wide effect through the
code base. It's important that we leverage the interface here as much as we can.

```go
// http://golang.org/src/pkg/errors/errors.go
type errorString struct {
	s string
}
```

This is the default concrete type that comes from the error package. It is an
unexported type that has an unexported field. This gives us enough context to
make us form a decision. We have responsibility around error handling to give
the caller enough context to make them form a decision so they know how to
handle this situation.

```go
// http://golang.org/src/pkg/errors/errors.go
func (e *errorString) Error() string {
	return e.s
}
```

This is using a pointer receiver and returning a string. If the caller must call
this method and parse a string to see what is going on then we fail. This method
is only for logging information about the error.

```go
// http://golang.org/src/pkg/errors/errors.go
// New returns an error that formats as the given text.
func New(text string) error {
	return &errorString{text}
}
```

When we call `New`, what we are doing is creating `errorString` value, putting
some sort of string in there. Since we are returning the address of a concrete
type, the user will get an error interface value where the first word is
a `*errorString` and the second word points to the original value. We are going
to stay decoupled during the error handling.

```
//       error
// ----------------
// | *errorString |          errorString
// ----------------     ---------------------
// |      *       | --> |   "Bad Request"   |
// ----------------     ---------------------
```

```go
func main() {
	if err := webCall(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Life is good")
}

// webCall performs a web operation.
func webCall() error {
	return New("Bad Request")
}
```

This is a very typical way of error handling in Go. We are calling `webCall`
and return the error interface and store that in a variable. `nil` is a special
value in Go.

What `error != nil` actually means is that we are asking if there is a concrete
type value that is stored in error type interface. Because if error is
not `nil`, there is a concrete value stored inside. If it is the case, we've got
an error. Now do we handle the error, do we return the error up the call stack
for someone else to handle? We will talk about this latter.

## Error Variables

[Sample program](example2/example2.go) to show how to use error variables to
help the caller determine the exact error being returned.

```go
var (
	// ErrBadRequest is returned when there are problems with the request.
	ErrBadRequest = errors.New("Bad Request")

	// ErrPageMoved is returned when a 301/302 is returned.
	ErrPageMoved = errors.New("Page Moved")
)
```

These error variables are great when you have a function that can return more
than one error. We want these to be on the top of the source code file.

Naming convention: starting with `Err`.

They have to be exported because our user need to access to them. These are all
error interfaces that we have discussed in the last file, with variables tied to
them. The contexts for these errors are the variables themselves. This allows us
to continue using the default error type, that unexported type with unexported
field to maintain a level of decoupling through error handling.

## Type as Context

It is not always possible to be able to say the interface value itself will be
enough context. Sometimes, it requires more context. For example, a networking
problem can be really complicated. Error variables wouldn't work there.

Only when the error variables wouldn't work, we should go ahead and start
working with custom concrete type for the error.

[Sample program](example3/example3.go) to show how to implement a custom error
type based on the JSON package in the standard library.

```go
// An UnmarshalTypeError describes a JSON value that was not appropriate for
// a value of a specific Go type.
// Naming convention: The word "Error" ends at the name of the type.
type UnmarshalTypeError struct {
	Value string       // description of JSON value
	Type  reflect.Type // type of Go value it could not be assigned to
}
```

We've got these two user-defined, right, custom error types, part of the standard library's JSON package.


```go
// Error implements the error interface.
func (e *UnmarshalTypeError) Error() string {
	return "json: cannot unmarshal " + e.Value + " into Go value of type " + e.Type.String()
}
```

`UnmarshalTypeError` implements the error interface.
We are using pointer semantic. In the implementation, we are validating all the
fields are being used in the error message. If not, we have a problem. Because
why would you add a field to the custom error type and not displaying on your
log when this method would call. We only do this when we really need it.

```go
// An InvalidUnmarshalError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidUnmarshalError struct {
	Type reflect.Type
}
```

This concrete type is used when we don't pass the address of a value into
`Unmarshal` function.

```go
// Unmarshal simulates an unmarshal call that always fails.
func Unmarshal(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

	return &UnmarshalTypeError{"string", reflect.TypeOf(v)}
}
```

The one thing I wanna show you here is the second parameter. The second parameter's based on the empty interface. The empty interface tells us nothing, because any piece of data, any value or pointer, satisfies it, because it doesn't have to have any behavior.

I want to be very careful when we're using the empty interface. Don't use it to write generic APIs.

We should be using the empty interface when we have to pass data around, where that data can be hidden without problems, or in this case, where we're gonna use the reflect package, because what we wanna do is at runtime, or dynamically inspect the concrete data. This is great if you wanna do model validation.

We then return different error types depending on these.

```go
// user is a type for use in the Unmarshal call.
type user struct {
	Name int
}

func main() {
    var u user
    // err is error interface value
	err := Unmarshal([]byte(`{"name":"bill"}`), u) // Run with a value and pointer.
	if err != nil {
		switch e := err.(type) { // type assertion
		case *UnmarshalTypeError: // *UnmarshalTypeError is a concrete type
			fmt.Printf("UnmarshalTypeError: Value[%s] Type[%v]\n", e.Value, e.Type) // e will be a copy of that pointer
		case *InvalidUnmarshalError:
			fmt.Printf("InvalidUnmarshalError: Type[%v]\n", e.Type)
		default:
			fmt.Println(err)
		}
		return
	}

	fmt.Println("Name:", u.Name)
}
```

**Flaw when using type as context here**

This idea of type as context for error handling is kind of dangerous, because we're setting ourselves up for some potential cascading changes, throughout the code base.

But type as context can be very powerful when you need to move concrete data across program boundaries, where both sides need to work with the concrete data, itself, maintaining levels of decoupling. I'm just afraid, when we're using it here, with error handling, because I'd really rather be processing the error interface value directly, or at least maintain levels of decoupling.

So then, how can we maintain this idea of custom error types without moving into the concrete? Behavior as context.

## Behavior as Context

Behavior as context allows us now to use these custom error types, but to stay decoupled.

[Sample code](example4/example4.go) to show how to implement behavior as context.

```go
// client represents a single connection in the room.
type client struct {
	name   string
	reader *bufio.Reader
}

// TypeAsContext shows how to check multiple types of possible custom error
// types that can be returned from the net package.
func (c *client) TypeAsContext() {
	for {
        // reader is an interface value that we can imagine that we've
        // abstracted a network call.
        // When we make this call over the network, we should get a
        // line of data back or we might get an error.
        line, err := c.reader.ReadString('\n')
        // Check to see if there's a concrete value stored inside the error
        // interface
		if err != nil {
            // We are doing a type as context, and checking the different types
            // of concrete error values that we have in the net package.
            // There's lots of types of errors that can occur in the network,
            // and the net package tries to cover them all.
			switch e := err.(type) { // e is concrete value
            case *net.OpError:
                // Remember, errors are just values in Go. They can be anything
                // we need them to be, both in state and behavior.
				if !e.Temporary() {
                    // If it's temporary, then we know that we're still in a state
                    // of integrity and keep going. If it's not temporary, we've
                    // lost integrity. Maybe that listener has gone down, that
                    // socket has dropped. And now we have to make sure that we
                    // can recover.
					log.Println("Temporary: Client leaving chat")
					return
				}

			case *net.AddrError:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			case *net.DNSConfigError:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			default:
				if err == io.EOF {
					log.Println("EOF: Client leaving chat")
					return
				}

				log.Println("read-routine", err)
			}
		}

		fmt.Println(line)
	}
}
```

**net package and `OpError`**

Let's switch over to the [net package](https://golang.org/pkg/net/#OpError)
for a second.

`OpError` is the most common type of concrete type we might be using for error
handling in the net package.

Notice there is a naming convention for a custom error type, that's that it ends
in the word `Error`.

Let's see if we can find the [`Temporary`](https://golang.org/pkg/net/#OpError.Temporary)
method for `OpError`. Look how complex this method is.

This `Temporary` method returning to a fault is brilliant, and it's important,
because it allow us to simplify whether or not there is or is not an integrity
issue or something as complex as a networking issue.

So let's apply the ideas of behavior as context and clean up this code.

```go
// temporary is declared to test for the existence of the method coming
// from the net package.
type temporary interface {
	Temporary() bool
}

// BehaviorAsContext shows how to check for the behavior of an interface
// that can be returned from the net package.
func (c *client) BehaviorAsContext() {
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			switch e := err.(type) {
			case temporary:
				if !e.Temporary() {
					log.Println("Temporary: Client leaving chat")
					return
				}

			default:
				if err == io.EOF {
					log.Println("EOF: Client leaving chat")
					return
				}

				log.Println("read-routine", err)
			}
		}

		fmt.Println(line)
	}
}
```

We go from the 3 cases on the error where we're going from decoupling to that
concrete, and we move it all down into one case. What we're going to ask is
does that concrete data also implement my `temporary` interface? I am now
staying within a decoupled state. When I do that now, I go to just one case,
that's it. I don't care what the concrete data is. All I care about is that it
implements `Temporary`. And `Temporary` will tell me whether or not I have an
integrity issue.

This is beautiful. The fact an error value can have not just state but behavior,
can allow your call or your user to maintain decoupling in error handling.

> Thank to Go implicit conversion. We can maintain a level of decopling by
> creating an interface with methods or behaviors that we only want, and use it
> instead of concrete type for type assertion switch.

Here's a general rule that I want you to follow. If your custom error type can
have any one of these four methods, it doesn't have to have all of them, just
one, then I want the custom error type to be defined as unexported with
unexported fields, because if it's unexported with unexported fields, you're
forcing that your user can never, ever go from a decoupled state to that
concrete. Here are the four methods: temporary, time out, not found, and
not authorized.

Temporary covers a tremendous amount. Temporary is really kind of this blanket
statement that you have an integrity issue or you don't.

Now if you have a custom error type like the JSON package did, time out,
temporary, not found, not authorized. They don't work for those types.
Those types have to be exported.

## Find the Bug

This is a piece of code that was, I've retrofitted it from a bug that a friend
of mine, who is a very experienced Go developer a couple years ago produced.

[Sample program](example5/example5.go) to show see if you can find the bug.

```go
import "log"

// customError is our base and it's an empty struct because we're not going
// to have any state.
type customError struct{}

// Error implements the error interface using a pointer semantics, we should be
// with these struct types.
func (c *customError) Error() string {
	return "Find the bug."
}

// fail returns nil values for both return types.
func fail() ([]byte, *customError) {
	return nil, nil
}

func main() {
    // Declare an error interface variable set to its zero value (a nil interface)
    var err error
    // We call fail(), which does not fail, returning the result. fail() does not
    // fail, but we check to see if fail() failed.
    // As far as this code is concerned, fail() failed. How is it possible that
    // fail() is failing when we are returning nil, the zero value, for the failure?
    //
    //            err (zero value)
    //          ---------
    // pointer  |  nil  |
    //          ---------
    // value    |  nil  |
    //          ---------
    //
    // There is no concrete data stored inside of error. The value we're now
    // sticking in here absolutely is nil. That really hasn't changed.
    // The question you have to ask yourself is what type of nil?
    //
    // Remember nil takes on the type it needs.
    // If you look at the return, we're seeing a pointer to a custom error.
    //
    //            err (zero value)
    //          ------------------
    // pointer  |  *customError  |
    //          ------------------
    // value    |  nil           |
    //          ------------------
    //
    // We've actually stored a concrete piece of data inside the error interface
    // value. The value happens to be nil, but it's of the custom error type.
    // This is the bug.
    // The bug is that the developer didn't use the error interface
    // on the return, they used the custom error type. The custom error type is an
    // artifact. It is the custom or the concrete value we're using. Error handling
    // happens in the decoupled state.
	if _, err = fail(); err != nil {
		log.Fatal("Why did this fail?")
	}

	log.Println("No Error")
}

```

We fix this code by always returning the error interface value. We're not
returning a nil custom error, we're returning the nil interface value, which now
means that there is nothing stored inside the interface.

```go
func fail() ([]byte, error) {
	return nil, nil
}
```

Next, this will teach us this last part of the interface mechanics.

_TODO_

## Wrapping Errors

### Logging

One of the things I find that is lacking in a lot of Go code is logging
consistency. Error handling and logging are just this one thing, and bring them
together if we want any consistency in them.

We write applications that log a lot of things, and most of the time we're
logging as an insurance policy to be able to find bugs when errors occur. And I
did that for a long time.

The reality is that there's too much activity on our systems today. Our user
bases can grow to a million people almost overnight and so logging that much has
a huge significant cost.

A lot of times logging is going to create a large amount of allocations, which
is going to put a lot of pressure on your heap. Now that's not unique to Go.
So I want you to consider that logging is important, but we've got to constantly
balance signal to noise in the log because if you're writing logs, writing data
to your logs, that you end up never, ever reading or using, you're wasting CPU
cycles on something that you could've been doing actual real work.
And it goes beyond just the CPU cycles of your process. You're eating network
bandwidth, disk I/O, other complexities that go through the entire system.

So during development I really wanna make sure that we always have a good level
of signal in our logs and we're logging from a trace perspective the bare
minimum we need, but then we're logging the errors in a way that there's always
enough context if we want to take the time or need to take the time to look at it.

How do you make sure there's enough context in the log, both from a tracing
perspective, bare minimal, and then an error perspective and not duplicate
errors throughout a log, have a consistent pattern that we all can follow and
review during code reviews, where we're doing logging the same way and it's
not random?

Now I wanna show you a pattern using Dave Cheney's errors package.
Dave Cheney's error package is very, really nice, and it gives us a consistent
way to apply error handling, apply logging and have code consistency throughout,
minimizing again a lot of pain.

**Error handling**

When a piece of code decides to handle the error, it means that code is
responsible for logging it and logging in the full context of it.

It also means that code has to make a decision. Can we recover or not? If no,
then error handling means we shut down the app, either with the stack trace on
the panic call or OS exit. If we can recover, then that code has to recover the
application back to its correct state and keep it going.

[Sample program](example6/example6.go) to show how wrapping errors work.

```go
import (
	"fmt"

    // This is Dave Cheney's errors package.
	"github.com/pkg/errors"
)

// AppError represents a custom error type.
type AppError struct {
	State int
}

// secondCall makes a call to a third function and wraps any error.
func secondCall(i int) error {
	if err := thirdCall(); err != nil {
		return errors.Wrap(err, "secondCall->thirdCall()")
	}
	return nil
}

// thirdCall create an error value we will validate.
func thirdCall() error {
	return &AppError{99}
}
```

So here we call `thirdCall`, we get back an error, interface value. We ask is
there a concrete value stored inside the error interface? The answer is yes.

When the answer is yes, now the developer writing this code has to make a
choice. It's really boolean. Am I gonna handle the error here, or not?

If the answer is handling the error here, we deal with it.
But if the answer is no, I'm not gonna handle the error, then there's only one
thing you're allowed to do, and that is wrap the error with context.

We would prefer to handle the error, the lower in the call stack we handle the
error, the better opportunity you're gonna have for recovery.

In this case the developer has decided not to handle the error. They don't have
to worry about logging and recovery anymore. All they worry about is the wrap
call. The wrap call does two things. There will be call stack context and
user context.

```go
errors.Wrap(err, "secondCall->thirdCall()")
```

The call stack context is going to take a, where we are in the code at that
line of code that we're doing the wrap. We now know exactly where we are when
this error occurred.

And we're also able to add some user context. In this case they just indicating
that `secondCall` was calling `thirdCall`.

Now we take this error, we wrap it, and we send it back up the call stack.

`firstCall` now is involved. `firstCall` now has to decide am I gonna handle the
error? If the answer is no, then you only have one choice. We wrap the error
again with that context.

We go to `main`. `main` made a call to `firstCall`, it gets back the error.

```go
// Make the function call and validate the error.
if err := firstCall(10); err != nil {

    // Use type as context to determine cause.
    // Generic type assertions. The Cause function knows how to unwind this
    // down and get us back to the root error value. We're type asserting
    // against the root error value.
    switch v := errors.Cause(err).(type) {
    case *AppError:

        // We got our custom error type.
        fmt.Println("Custom App Error:", v.State)

    default:

        // We did not get any specific error type.
        fmt.Println("Default Error")
    }

    // Display the stack trace for the error.
    fmt.Println("\nStack Trace\n********************************")
    fmt.Printf("%+v\n", err)
    fmt.Println("\nNo Trace\n********************************")
    fmt.Printf("%v\n", err)
    // If you're using the standard library fmt or log packages, if you put a
    // %+v in the formatting, the full stack trace and context of this whole
    // thing gets logged to, in this case, standard out or standard error.
    // %v just gives you the user context. %+v gives you both.
}
```

This is in a fantastic pattern.

Don't forget that we have dashboards. We have metrics. I don't prefer writing
data into the logs. I'm not a big fan of structured logging because for me, logs
serve the purpose of being able to find and fix bugs, but that's me. I rather
use my metric systems and my dashboards for those data points, and try to tie
that stuff together.
