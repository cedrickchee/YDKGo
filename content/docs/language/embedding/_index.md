---
title: Embedding
weight: 11
---

# Embedding

## Declaring Fields

[Sample program](example1/example1.go) to show how what we are doing is NOT
embedding a type but just using a type as a field.

```go
import "fmt"

// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements a method notifies users
// of different events.
func (u *user) notify() {
	fmt.Printf("Sending user email To %s<%s>\n", u.name, u.email)
}

// admin represents an admin user with privileges.
type admin struct {
	person user // NOT Embedding
	level  string
}
```

`person user` is not embedding. All we do here just create a person field based
on that other concrete type named `user`.

```go
func main() {
	// Create an admin user
	ad := admin{
		person: user{
			name:  "john smith",
			email: "john@yahoo.com",
		},
		level: "super",
	}

    // We can access fields methods.
	ad.person.notify()
}
```

Create an `admin` user using struct literal. Since `person` also has struct
type, we use another literal to initialize it.

We call `notify` through the `person` field through the `admin` type value.

## Embedding types

[Sample program](example2/example2.go) to show how to embed a type into another
type and the relationship between the inner and outer type.

```go
// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements a method notifies users
// of different events.
func (u *user) notify() {
	fmt.Printf("Sending user email To %s<%s>\n",
		u.name,
		u.email)
}

// admin represents an admin user with privileges.
type admin struct {
	user  // Embedded Type
	level string
}
```

Notice that we don't use the field `person` here anymore. We are now embedding a
value of type `user` inside value of type `admin`. This is an
inner-type-outer-type relationship where `user` is the inner type and `admin`
is the outer type.

### Inner type promotion

What special about embedding in Go is that we have inner type promotion
mechanism. In other words, anything related to the inner type can be promoted up
to the outer type. It will mean more in the construction below.

```go
func main() {
	// Create an admin user.
	ad := admin{
		user: user{
			name:  "john smith",
			email: "john@yahoo.com",
		},
		level: "super",
    }
}
```

We are now constructing outer type `admin` and inner type `user`. This inner
type value now looks like a field, but it is not a field. We can access it
through the type name like a field. We are initializing the inner value through
the struct literal of `user`.

```go
// We can access the inner type's method directly.
ad.user.notify()

// The inner type's method is promoted.
ad.notify()
```

Because of inner type promotion, we can access the `notify` method directly
through the outer type. Therefore, the output will be the same.

Inner type promotion is giving us a sense of type reuse, but will use it for
much more than that.

## Embedded types and interfaces

[Sample program](example3/example3.go) to show how embedded types work with
interfaces.

```go
// notifier is an interface that defined notification
// type behavior.
type notifier interface {
	notify()
}

// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements a method notifies users of different events using a
// pointer receiver.
func (u *user) notify() {
	fmt.Printf("Sending user email To %s<%s>\n",
		u.name,
		u.email)
}

// admin represents an admin user with privileges.
type admin struct {
	user
	level string
}

func main() {
	// Create an admin user.
	ad := admin{
		user: user{
			name:  "john smith",
			email: "john@yahoo.com",
		},
		level: "super",
	}

	// Send the admin user a notification.
	// The embedded inner type's implementation of the
	// interface is "promoted" to the outer type.
	sendNotification(&ad)
}
```

We are passing the address of outer type value `&ad`. Because of inner type
promotion, the outer type now implements all the same contract as the inner type.

Embedding does not create a subtyping relationship. This isn't base derived
class stuff like you see in OOP languages. `admin` is admin and `user`is still
`user`. The behavior that inner type value uses, the outer type exposes it as
well. It means that outer type value can implement the same interface/same
contract as the inner type.

We are getting type reuse. We are not mixing or sharing state but extending the
behavior up to the outer type.

```go
// sendNotification accepts values that implement the notifier
// interface and sends notifications.
func sendNotification(n notifier) {
	n.notify()
}
```

We have our polymorphic function here. Because of inner type promotion, this
`admin` value now satisfies all the same interfaces.

## Outer and inner type interface implementations

[Sample program](example4/example4.go) to show what happens when the outer and
inner type implement the same interface.

```go
// notifier is an interface that defined notification type behavior.
type notifier interface {
	notify()
}

// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements a method notifies users of different events.
func (u *user) notify() {
	fmt.Printf("Sending user email To %s<%s>\n",
		u.name,
		u.email)
}

// admin represents an admin user with privileges.
type admin struct {
	user
	level string
}

// notify implements a method notifies admins of different events.
func (a *admin) notify() {
	fmt.Printf("Sending admin Email To %s<%s>\n",
		a.name,
		a.email)
}
```

We now have two different implementations of `notifier` interface, one for the
inner type, one for the outer type. Because the outer type now implements that
interface, the inner type promotion doesn't happen. We have overwritten through
the outer type anything that inner type provides to us.

```go
func main() {
	// Create an admin user.
	ad := admin{
		user: user{
			name:  "john smith",
			email: "john@yahoo.com",
		},
		level: "super",
	}

	// Send the admin user a notification.
    // The embedded inner type's implementation of the interface is NOT
    // "promoted" to the outer type.
	sendNotification(&ad)

	// We can access the inner type's method directly.
	ad.user.notify()

	// The inner type's method is NOT promoted.
	ad.notify()
}

// sendNotification accepts values that implement the notifier interface and
// sends notifications.
func sendNotification(n notifier) {
	n.notify()
}
```
