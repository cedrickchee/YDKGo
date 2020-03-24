---
title: Exporting
weight: 12
---

# Exporting

Exporting is kind of like the **encapsulation** piece of Go.

In object oriented programming, we are used to public, private, protected type
of access mechanism. Go is different. Everything in Go is about package.

Package is a self-contained unit of code. Every folder in our source tree is a
self-contained user code. We will get deeper into that in package oriented
design section.

The idea is: anything that is named in a given package can be exported or
accessible through other packages or unexported or not accessible through
other packages.

## Declare and access exported identifiers

[Sample program](example1/example1.go).

```go
import (
	"fmt"

    // This is a relative path to a physical location on our disk - relative to GOPATH
	"github.com/cedrickchee/ultimate-go/language/exporting/example1/counters"
)

func main() {
	// Create a variable of the exported type and initialize the value to 10.
	counter := counters.AlertCounter(10)

	fmt.Printf("Counter: %d\n", counter)
}
```

[Sample package](example1/counters/counters.go).

```go
// Package counters provides alert counter support.
package counters

// AlertCounter is an exported named type that
// contains an integer counter for alerts.
type AlertCounter int
```

## Declare unexported identifiers and restrictions

[Sample package](example2/counters/counters.go).

```go
// Package counters provides alert counter support.
package counters

// alertCounter is an unexported named type that
// contains an integer counter for alerts.
type alertCounter int
```

[Sample program](example2/example2.go).

This time, when we create a variable of the unexported type and initialize the
value to 10: `counter := counters.alertCounter(10)`, the compiler will say:
- cannot refer to unexported name counters.alertCounter
- undefined: counters.alertCounter

## Access values of unexported identifiers

[Sample package](example3/counters/counters.go).

```go
// Package counters provides alert counter support.
package counters

// alertCounter is an unexported named type that contains an integer counter
// for alerts.
type alertCounter int

// New creates and returns values of the unexported type alertCounter.
func New(value int) alertCounter {
	return alertCounter(value)
}
```

Declare an exported function called `New` - a factory function that knows how
to create and initialize the value of an unexported type. It returns an
unexported value of `alertCounter`.

[Sample program](example3/example3.go).

```go
func main() {
	// Create a variable of the unexported type using the exported New function
	// from the package counters.
	counter := counters.New(10)

	fmt.Printf("Counter: %d\n", counter)
}
```

The compiler is OK with this because exporting and unexporting is not about
the value like private and public mechanism, it is about the identifier itself.
Please don't do this, since there is no encapsulation here. We can just make
the type exported.

You should be using your Go linter. Your Go linter's going to tell you this,
"Hey, don't do this because it's annoying to use," and it's absolutely right.

## Unexported struct type fields

We saw package-level encapsulation, now let's look at type-level encapsulation.

[Sample package](example4/users/users.go).

```go
// Package users provides support for user management.
package users

// User represents information about a user.
type User struct {
	Name string
	ID   int

	password string
}
```

Exported type `User` represents information about a user. It has 2 exported
fields: `Name` and `ID` and 1 unexported field: `password`.

[Sample program](example4/example4.go).

```go
func main() {
	// Create a value of type User from the users package using struct literal.
	u := users.User{
		Name: "Chole",
		ID:   10,

		password: "xxxx",
    }

	fmt.Printf("User: %#v\n", u)
}
```

But, since `password` is unexported, we get compiler error: "unknown users.User
field 'password' in struct literal".

## Unexported embedded types

[Sample package](example5/users/users.go).

```go
// Package users provides support for user management.
package users

// user represents information about a user.
type user struct {
	Name string
	ID   int
}

// Manager represents information about a manager.
type Manager struct {
	Title string

	user
}
```

`user` is an unexported type with 2 exported fields.
`Manager` is an exported type embedded the unexported field `user`.

[Sample program](example5/example5.go).

```go
// Create a value of type Manager from the users package.
u := users.Manager{
    Title: "Dev Manager",
}
```

During construction, we are only able to initialize the exported field `Title`.
We cannot access the embedded type directly.

```go
// Set the exported fields from the unexported user inner type.
u.Name = "Chole"
u.ID = 10
```

However, once we have the `Manager` value, the exported fields from that
unexported type are accessible. We're really not getting any encapsulation here
because all of the fields got promoted up anyway.

We'd be making user capital `U` if I saw this code, but it is very, very common
in Go to have this situation where the type is un-exported and the fields are
exported. Example, marshaling and un-marshaling code in Go only will respect
fields that are exported.
