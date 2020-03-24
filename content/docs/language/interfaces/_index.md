---
title: Interfaces
weight: 10
---

# Interfaces

Interfaces provide a way to declare types that define only behavior. This behavior can be implemented by concrete types, such as struct or named types, via methods. When a concrete type implements the set of methods for an interface, values of the concrete type can be assigned to variables of the interface type. Then method calls against the interface value actually call into the equivalent method of the concrete value. Since any concrete type can implement any interface, method calls against an interface value are polymorphic in nature.

* The method set for a value, only includes methods implemented with a value receiver.
* The method set for a pointer, includes methods implemented with both pointer and value receivers.
* Methods declared with a pointer receiver, only implement the interface with pointer values.
* Methods declared with a value receiver, implement the interface with both a value and pointer receiver.
* The rules of method sets apply to interface types.
* Interfaces are reference types, don't share with a pointer.
* This is how we create polymorphic behavior in Go.

## Part 1: Polymorphism

"Polymorphism means that you write a certain program and it behaves differently depending on the data that it operates on." - Tom Kurtz (inventor of BASIC)

[Sample program](example1/example1.go) to show how polymorphic behavior with
interfaces.

```go
// reader is an interface that defines the act of reading data.
type reader interface {
	read(b []byte) (int, error) // 1
}
```

`interface` is technically a **valueless type**. This interface doesn't declare
state. They only define a method set of behavior. They defines a contract of
behavior. Through that contract of behavior, we have polymorphism.

It is a 2 word data structure that has 2 pointers.

When we say `var r reader`, we would have a `nil` value interface because
interface is a reference type.

Interface types are not real. `r` is not real, there's nothing concrete about an
interface type.

Interfaces are about behavior. We don't want to see interfaces that describe
things: user/animal/car interface. These are not interfaces. Those are things,
that's your concrete data.

We could have written this API a little bit differently.
Technically, I could have said, how many bytes do you want me to read and I will
return that slice of byte and an error, like so:

```go
type reader interface {
    read(i int) ([]byte, error)` // 2
}
```

Why do we choose the other one instead?
Every time we call (2), it will cost an allocation because the method would have to
allocate a slice of some unknown type and share it back up the call stack.
The method would have to allocate a slice of some unknown type and share it back
up the call stack. The backing array for that slice has to be an allocation.
But if we stick with (1), the caller is allocating a slice. Even the backing
array for that is ended up on a heap, it is just 1 allocation. We can call this
10000 times and it is still 1 allocation.

```go
// ************
// Relationship
// ************

// We store concrete type values inside interfaces.
type file struct {
	name string
}

// read implements the reader interface for a file.
func (file) read(b []byte) (int, error) {
	s := "<rss><channel><title>Going Go Programming</title></channel></rss>"
	copy(b, s)
	return len(s), nil
}
```

**Convention over configuration**

We do not configure an interface to the concrete type like you might see in
other languages.

We just have to declare the method and the compiler at compile time can identify
interface compliance, satisfaction.

### Concrete type vs interface type

A concrete type is any type that can have a method. Only user defined type can
have a method. Method allows a piece of data to expose capabilities, primarily
around interfaces. `file` defines a system file. It is a concrete type because
it has the method `read` below. It is identical to the method in the `reader`
interface. Because of this, we can say the concrete type `file` implements the
`reader` interface using a value receiver. There is no fancy syntax.
The complier can automatically recognize the implementation here.

```go
// pipe defines a named pipe network connection.
type pipe struct {
	name string
}

// read implements the reader interface for a network connection.
func (pipe) read(b []byte) (int, error) {
	s := `{name: "hoanh", title: "developer"}`
	copy(b, s)
	return len(s), nil
}
```

This is the second concrete type that uses a value receiver. We now have two
different pieces of data, both exposing the reader's contract and implementation
for this contract.

### Polymorphic function

```go
// retrieve can read any device and process the data.
func retrieve(r reader) error {
	data := make([]byte, 100)

	len, err := r.read(data)
	if err != nil {
		return err
	}

	fmt.Println(string(data[:len]))
	return nil
}
```
This is called a polymorphic function.

The parameter is being used here is the `reader` type. But it is valueless.
**Interface types are valueless**. What does it mean?
This function will accept any data that implement the `reader` contract.
This function knows nothing about the concrete and it is completely decoupled.
It is the highest level of decoupling we can get.
The algorithm is still efficient and compact. All we have is a level of
indirection to the concrete type data values in order to be able to execute
the algorithm.

```go
func main() {
	// Create two values one of type file and one of type pipe.
	f := file{"data.json"}
	p := pipe{"cfg_service"}

	// Call the retrieve function for each concrete type.
	retrieve(f)
	retrieve(p)
}
```

Here we are passing the value itself, which means a copy of `f` going to pass
across the program boundary.

The compiler will ask: Does this `file` value implement the `reader` interface?
The answer is yes because there is a method there using the value receiver that
implements its contract.

```
//       reader           iTable
//   -------------      ----------
//   |           |      |  file  |
//   |     *     | -->  ----------
//   |           |      |   *    | --> code
//   -------------      ----------
//   |           |      ----------
//   |     *     | -->  | f copy | --> read()
//   |           |      ----------
//   -------------
```

The second word of the interface value will store its own copy of `f`.
The first word points to a special data structure that we call the iTable.

The iTable serves 2 purposes:
- The first part describes the type of value being stored. In our case, it is
the `file` value.
- The second part gives us a matrix of function pointers so we can actually
execute the right method when we call that through the interface.

```go
retrieve(f)
```

When we do a `read` against the interface, we can do an iTable lookup, find
that `read` associated with this type, then call that value against the `read`
method - the concrete implementation of read for this type of value.

```
//       reader           iTable
//   -------------      ----------
//   |           |      |  pipe  |
//   |     *     | -->  ----------
//   |           |      |   *    | --> code
//   -------------      ----------
//   |           |      ----------
//   |     *     | -->  | p copy | --> read()
//   |           |      ----------
//   -------------
```

Similar with `p`. Now the first word of `reader` interface points to `pipe`,
not `file` and the second word points to a copy of `pipe` value.

```go
retrieve(p)
```

The behavior changes because the data changes.

Note: moving forward, for simplicity, instead of drawing the a pointer pointing
to iTable, we only draw *pipe, like so:

```
//  -------
// | *pipe |
//  -------
// |   *   |  --> p copy
//  ------
```

## Part 2: Method Sets and Address of Value

### Method Sets

[Sample program](example2/example2.go) to show how to understand method sets.

```go
// notifier is an interface that defines notification
// type behavior.
type notifier interface {
	notify()
}

// user defines a user in the program.
type user struct {
	name  string
	email string
}

// notify implements the notifier interface with a pointer receiver.
// Interface via pointer receiver
func (u *user) notify() {
	fmt.Printf("Sending User Email To %s<%s>\n",
		u.name,
		u.email)
}

func main() {

	// Create a value of type User and send a notification.
	u := user{"Bill", "bill@email.com"}

	// Values of type user do not implement the interface because pointer
	// receivers don't belong to the method set of a value.

	sendNotification(u)

	// ./example1.go:36: cannot use u (type user) as type notifier in argument to sendNotification:
	//   user does not implement notifier (notify method has pointer receiver)
}

// This is our polymorphic function.
// sendNotification accepts values that implement the notifier
// interface and sends notifications.
// This is again saying: I will accept any value or pointer that implement the notifier interface.
// I will call that behavior against the interface itself.
func sendNotification(n notifier) {
	n.notify()
}
```

`sendNotification(u)` call polymorphic function but passing `u` using value
semantic. However, the compiler doesn't allow it: "cannot use u (type user) as
type notifier in argument to sendNotification: user does not implement
notifier (notify method has pointer receiver)". This is setting up for an
integrity issue.

In the specification, there are sets of rules around the
concepts of **method sets**. What we are doing is against these rules.

What are the rules?
- For any value of a given type `T`, only those methods implemented with a value
receiver belong to the method sets of that type.
- For any value of a given type `*T` (pointer of a given type), both value
receiver and pointer receiver methods belong to the method sets of that type.

In other words, if we are working with a pointer of some types, all the methods
that has been declared are associated with that pointer. But if we are working
with a value of some types, only those methods that operated on value semantic
can be applied.

In the previous lesson about method, we are calling them before without any
problem. That is true. When we are dealing with method, method call against the
concrete values themselves, Go can adjust to make the call. However, we are not
trying to call a method here. We are trying to store a concrete type value
inside the interface. For that to happen, that value must satisfy the contract.

The question now becomes: Why can't pointer receiver be associated with the
method sets for value? What is the integrity issue here that doesn't allow us to
use pointer semantic for value of type `T`?

It is not 100% guarantee that any value that can satisfy the interface has an
address. We can never call a pointer receiver because if that value doesn't have
an address, it is not shareable. For example:

[Sample program](example3/example3.go) to show how you can't always get the
address of a value.

```go
// Declare a type named duration that is based on an integer
type duration int

// Declare a method name notify using a pointer receiver.
// This type now implements the notifier interface using a pointer receiver.
func (d *duration) notify() {
    fmt.Println("Sending Notification in", *d)
}

// Take a value 42, convert it to type duration and try to call the notify method.
// Here are what the compiler says:
// - "cannot call pointer method on duration(42)"
// - "cannot take the address of duration(42)"
func main() {
    duration(42).notify()
}
```

Why can't we get the address? Because 42 is not stored in a variable. It is
still literal value that we don't know ahead the type. Yet it still does
implement the notifier interface.

Come back to our example, when we get the error, we know that we are mixing
semantics. `u` implements the interface using a pointer receiver and now we are
trying to work with copy of that value, instead of trying to share it.
It is not consistent.

The lesson:
- If we implement interface using pointer receiver, we must use pointer semantic.
- If we implement interface using value receiver, we then have the ability to
use value semantic and pointer semantic. However, for consistency, we want to
use value semantic most of the time, unless we are doing something like
Unmarshal function.

To fix the issue, instead of passing value `u`,
we must pass the address of `u` (`&u`). We create a `user` value and pass the
address of that, which means the interface now has a pointer of type `user` and
we get to point to the original value.

```
// ---------
// | *User |
// ---------
// |   *   | --> original user value
// ---------
```

## Part 3: Storage by Value

[Sample program](example4/example4.go) to show how the concrete value assigned
to the interface is what is stored inside the interface.

```go
// printer displays information.
type printer interface {
	print()
}

// user defines a user in the program.
type user struct {
	name string
}

// print displays the user's name.
func (u user) print() {
	fmt.Printf("User Name: %s\n", u.name)
}

func main() {
	// Create values of type user and admin.
	u := user{"Bill"}

	// Add the values and pointers to the slice of
    // printer interface values.
    //
	//   index 0   index 1
	// ---------------------
	// |   User  |  *User  |
	// ---------------------
	// |    *    |    *    |
	// ---------------------
	//      A         A
	//      |         |
	//     copy    original
	entities := []printer{
        // Store a copy of the user value in the interface value.
        // When we do that, the interface value has its own copy of the value.
		u, // Changes to the original value will not be seen.

        // Store a copy of the address of the user value in the interface value.
        // When we do that, the interface value has its own copy of the address.
		&u, // Changes to the original value will be seen.
	}

	// Change the name field on the user value.
	u.name = "Bill_CHG"

	// Iterate over the slice of entities and call
	// print against the copied interface value.
	for _, e := range entities {
		e.print()
	}
}
```
