---
title: Composition
weight: 1
---

# Interface and Composition Design

## Grouping Types

### Grouping by State

[Sample program](grouping/example1/example1.go).

```go
// Animal contains all the base fields for animals.
type Animal struct {
	Name     string
	IsMammal bool
}

// Speak provides generic behavior for all animals and how they speak.
// This is completely useless because animals themselves cannot speak. This
// cannot apply to all animals.
// There's no value in it. These things need tests, which can't really be
// accurate. It increase our lines of code, increase our bug potential, increase
// our test cases, and it adds zero value to the software.
func (a *Animal) Speak() {
	fmt.Printf(
		"UGH! My name is %s, it is %t I am a mammal\n",
		a.Name,
		a.IsMammal,
	)
}

// Dog contains everything an Animal is but specific
// attributes that only a Dog has.
type Dog struct {
	Animal
	PackFactor int
}

// Speak knows how to speak like a dog.
func (d *Dog) Speak() {
	fmt.Printf(
		"Woof! My name is %s, it is %t I am a mammal with a pack factor of %d.\n",
		d.Name,
		d.IsMammal,
		d.PackFactor,
	)
}

// Cat contains everything an Animal is but specific
// attributes that only a Cat has.
type Cat struct {
	Animal
	ClimbFactor int
}

// Speak knows how to speak like a cat.
func (c *Cat) Speak() {
	fmt.Printf(
		"Meow! My name is %s, it is %t I am a mammal with a climb factor of %d.\n",
		c.Name,
		c.IsMammal,
		c.ClimbFactor,
	)
}

func main() {
    // This code will not compile.
    // Here, we try to group the Cat and Dog based on the fact that they are Animals.
    // We are trying to leverage subtyping in Go. However, Go doesn't have it.
	animals := []Animal{
		// Create a Dog by initializing its Animal parts and then its specific Dog attributes.
		Dog{
			Animal: Animal{
				Name:     "Fido",
				IsMammal: true,
			},
			PackFactor: 5,
		},

		// Create a Cat by initializing its Animal parts and then its specific Cat attributes.
		Cat{
			Animal: Animal{
				Name:     "Milo",
				IsMammal: true,
			},
			ClimbFactor: 4,
		},
	}

	// Have the Animals speak.
	for _, animal := range animals {
		animal.Speak()
	}
}
```

The problem: Why is it that I cannot group Dogs and Cats by what they are, which
is an Animal?

The very first things you've got to realize is Go goes very strong on their type
system. The fact that you've embedded Animal into Dog and Cat means zero, it
means nothing.

There is no subtyping, there is no subclassing in Go. All types are their own
and the concepts of base and derived types do not exist in Go. This pattern does
not provide a good design principle in a Go program.

Go doesn't encourage us to group types by common DNA. We need to stop designing
APIs around this idea that types have a common DNA because if we only focus on
who we are, it is very limiting on who can we group with. Subtyping doesn't
promote diversity. We lock types in a very small subset that can be grouped
with. But when we focus on behavior, we open up entire world to us.

`Animal` type is an example of type pollution.

Believe it or not, a little copying and pasting can go a very long way.
We've been taught these concepts of DRY, Do not Repeat Yourself. But I think the
cost of DRY in Go is worse than maybe some of these other languages you've
worked on. Coupling is always going to be a bigger cost than a little
duplication.

Then, how do we correct this program? We will talk about this next.

### Grouping By Behavior

Stop thinking about who you are (concrete base type). Stop thinking about what
cats and dogs are and start focusing on what cats and dogs do.
It's convention over configuration.

This is an [example](grouping/example2/example2.go) of using composition and
interfaces.

Let's remove the `Animal` type. Let's just bring in an interface, the `Speaker`
interface, that's behavior.

We will group common types by their behavior and not by their state.

```go
// Speaker provide a common behavior for all concrete types
// to follow if they want to be a part of this group. This
// is a contract for these concrete types to follow.
type Speaker interface {
	Speak()
}

// Dog contains everything a Dog needs.
// We have copied those common fields into Dog. This is going to make code more
// readable, easier to debug, and easier to test. Now, Dog is more precise.
// These benefits outweight that cost of DRY.
type Dog struct {
	Name       string
	IsMammal   bool
	PackFactor int
}

// Speak knows how to speak like a dog.
// This makes a Dog now part of a group of concrete types that know how to speak.
func (d *Dog) Speak() {
	fmt.Printf(
		"Woof! My name is %s, it is %t I am a mammal with a pack factor of %d.\n",
		d.Name,
		d.IsMammal,
		d.PackFactor,
	)
}

// Cat contains everything a Cat needs.
type Cat struct {
	Name        string
	IsMammal    bool
	ClimbFactor int
}

// Speak knows how to speak like a cat.
// This makes a Cat now part of a group of concrete types that know how to speak.
func (c *Cat) Speak() {
	fmt.Printf(
		"Meow! My name is %s, it is %t I am a mammal with a climb factor of %d.\n",
		c.Name,
		c.IsMammal,
		c.ClimbFactor,
	)
}
```

What brilliant about Go is that it doesn't have to be configured ahead of time.
The compiler automatically identifies interface and behaviors at compile time.
It means that we can write code today that compliant with any interface that
exists today or tomorrow. It doesn't matter where that is declared because the
compiler can do this on the fly.

```go
func main() {
	// Create a list of Animals that know how to speak.
	speakers := []Speaker{
		// Create a Dog by initializing Dog attributes.
		Dog{
			Name:       "Fido",
			IsMammal:   true,
			PackFactor: 5,
		},

		// Create a Cat by initializing Cat attributes.
		Cat{
			Name:        "Milo",
			IsMammal:    true,
			ClimbFactor: 4,
		},
	}

	// Have the Speakers speak.
	for _, spkr := range speakers {
		spkr.Speak()
	}
}

// =============================================================================

// NOTES:

// Here are some guidelines around declaring types:
// 	* Declare types that represent something new or unique.
// 	* Validate that a value of any type is created or used on its own.
// 	* Embed types to reuse existing behaviors you need to satisfy.
// 	* Question types that are an alias or abstraction for an existing type.
// 	* Question types whose sole purpose is to share common state.
//
// Practical examples:
//  - Bad:
//        type handle int // type aliasing
//        func foo(h handle)
//  - Ok:
//        type Duration int64 // Duration doesn't represents integer. It represents nanoseconds of time.
```

## Decoupling

Let's really look at how composition works in Go. Guidelines around doing this
and writing code in Go.

It's probably a little bit against, what you're used to doing. A lot of you have been taught, to start with the interface, start with the behavior, try to figure out what those contracts are. I don't want you to do that. That is guessing.

Remember, the problem is solved in the concrete, not in the abstract. So we need a concrete influenced, implementation solution first, in order to know how to decouple. This is going to simplify your code, and allow you to focus on what's important, which is the problem, and remember, the problem is the data.

Before we get into this example. Things that I see, that are missing, and one of these things, is the idea of done. How do you know when you're done with a piece of code?

First, make sure that you have your unit tests and you got to have some level of test coverage.

Second, is the code decoupled from the change we expect to happen? Now, we might decide that we know what the decoupling is required, but we don't need it today, because we don't have multiple implementations with something. We don't have to need today. Remember, you're writing code for today, we're designing and architecting for tomorrow. So, deciding to decouple, may or may not happen immediately, but I want to ask the question, do we know what has to be decoupled, and do we want to do that? And I can do this, and I can ask this question, because decoupling is part two of everything I do. It is a refactoring. We solve problems in the concrete first, we refactor into the decoupling.

Things, that I see on that developers have problems with: they don't understand, how to create a **layered API**.

This is really going to help you with a codebase, instead of you trying to solve every problem, in a few functions. What we really want to do, is find a layered approach, then initially just three layers, is all we're going to need.

We can have what I call our **primitive layer**, and this layer knows how to do one thing, and one thing very, very well. We're always focused on what is that one thing, and one thing very, very well. We write out the code for this layer, we write our unit tests for this layer right, and we write this layer, so it is testable.

Now, when I say **testable**, this doesn't mean it interfaces, remember, I'm working in the concrete, interfaces, decoupling happens through a refactoring, and yet I have to already start writing unit tests. Unit tests, code that is testable, usually means that the data, that we're passing in is reproducible, and the data coming out, okay, is testable. The problem you're trying to solve is about the data. So, when we're think about unit testing here, I want you to think about, can I write a function in the primitive layer, where I can give it a set of data, get a set of data back out, and we can validate it and test it.

Your next layer, I would call the lower Level, the **lower level layer**, and that's the layer that sits on top of the primitive API, that does maybe some raw things, a little higher level than the primitive layer, and again we're going to write those unit tests, and this layer to, should be testable in its own right. We should be able to test everything coming in, and everything going out, and a lot of time, this layer and these layers, are probably unexported, they could be, but sometimes it's nice to export the Lower Level, the Lower Layer API, because it's very usable if you write it correctly, and it gives the developer, the users, maybe more control over the things they need.

Then you're going to have your, **high level API**. This is where you're trying to do as much as you can for the user, to make their life better. You get it sitting on top of the Lower Level, which is sitting on top of the Primitive Level. Since every one of these levels have been coded, and unit tested, by the time you get up here, we're really just thinking about, an ease of use for the developer, and we're going to be able to hopefully unit test this as well. Sometimes these higher level functions, require more kinda integration kinds of tests, but we want to think about, how we test every layer all the way up. Eventually, your higher level tests might be able to replace some of your lower level, and primitive unit tests, because they cover those test cases for you, and so there are times, where I'm even writing unit tests, that I know aren't going to exist forever, because as I move up these levels of API, I'm getting the same code coverage on the tests that I've written here, and that's okay.

I don't want you to be worried, about throwing code away right, refactoring is about making something work, and looking at how we make it better, how we make it more readable, how do we do things, and the best days are when we actually remove code, so don't get hung up that I might be writing some unit tests here, that eventually we're going to delete, because we're going to get some better coverage here, and writing those tests right away, because I want to make sure these layers are working.

**Prototyping**

TL;DR: Prototyping or writing proof of concept and solving problem in the concrete
first is important. This allow us to ask ourselves: what can change?, what
change is coming? so we can start decoupling and refactor. Refactoring need to
become a part of the development cycle.

This idea of the concrete implementation could mean one of two things. It could be a prototype and I really believe in prototype oriented design because I don't want you guessing, I want you knowing. And so, it could be a prototype. It could be a prototype that we could even put in production right away. Could be a production piece of coding even if it's in the concrete. The idea today is we have to be able to get code into production faster, a code that has integrity, a code that we can really start solving problems because, for me, technical debt is when you were already working on a piece of code and it never leaves your laptop, it never gets into production. We have to be better at all of this.

Here is **the problem:** we have a system called Xenia that has a database. There is
another system called Pillar, which is a web server with some front-end that
consume it. It has a database too. Our goal is to move the Xenia's data into
Pillar's system.

How long will it take? How do we know when a piece of code is done so we can
move on the next piece of code? If you are a technical manager, how do you know
whether your debt is "wasting effort" or "not putting enough effort"?

Being done has 2 parts: one is test coverage, 80% in general and 100% on the
happy path. Second is about changes. By asking what can change, from technical
perspective and business perspective, we make sure that we refactor the code to
be able to handle that change.

One example is, we can give you a concrete version in 2 days but we need 2 weeks
to be able to refactor this code to deal with the change that we know it's coming.

The plan is to solve one problem at a time. Don't be overwhelmed by everything.
Write a little code, write some tests and refactor. Write layer of APIs that
work on top of each other, knowing that each layer is a strong foundation to
the next.

Do not pay too much attention in the implementation detail.
It's the mechanics here that are important. We are optimizing for correctness,
not performance. We can always go back if it doesn't perform well enough to
speed things up.

So let's look at the base prototype code that we end up putting in production, that solves the problem at hand.

### Struct Composition

[Sample program](decoupling/example1/example1.go)

```go
// Sample program demonstrating struct composition.
package main

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"
)

// The first problem that we have to solve is that we need a software that run
// on a timer. It need to connect to Xenia, read that database, identify all the
// data we haven't moved and pull it in.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// =============================================================================

// Data is the structure of the data we are copying.
// For simplicity, just pretend it is a string data.
type Data struct {
	Line string
}

// =============================================================================

// Xenia is a system we need to pull data from.
type Xenia struct {
	Host    string
	Timeout time.Duration
}

// Pull knows how to pull data out of Xenia.
// We could do func (*Xenia) Pull() (*Data, error) that return the data and error.
// However, this would cost an allocation on every call and we don't want that.
// Using the function below, we know data is a struct type and its size ahead
// of time. Therefore they could be on the stack.
func (*Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF

	case 5:
		return errors.New("Error reading data from Xenia")

	default:
		d.Line = "Data"
		fmt.Println("In:", d.Line)
		return nil
	}
}

// Pillar is a system we need to store data into.
type Pillar struct {
	Host    string
	Timeout time.Duration
}

// Store knows how to store data into Pillar.
// We are using pointer semantics for consistency.
func (*Pillar) Store(d *Data) error {
	fmt.Println("Out:", d.Line)
	return nil
}

// =============================================================================

// System wraps Xenia and Pillar together into a single system.
// We have the API based on Xenia and Pillar. We want to build another API on
// top of this and use it as a foundation. One way is to have a type that have
// the behavior of being able to pull and store. We can do that through
// composition. System is based on the embedded value of Xenia and Pillar. And
// because of inner type promotion, System know how to pull and store.
type System struct {
	Xenia
	Pillar
}

// =============================================================================

// pull knows how to pull bulks of data from Xenia, leveraging the foundation
// that we have built. We don't need to add method to System to do this. There
// is no state inside System that we want the System to maintain. Instead, we
// want the System to understand the behavior. Functions are a great way of
// writing API because functions can be more readable than any method can.
// We always want to start with an idea of writing API from the package level
// with functions. When we write a function, all the input must be passed in.
// When we use a method, its signature doesn't indicate any level, what field or
// state that we are using on that value that we use to make the call.
func pull(x *Xenia, data []Data) (int, error) {
	for i := range data {
		if err := x.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store bulks of data into Pillar.
// Similar to the function above.
// We might wonder if it is efficient. However, we are optimizing for
// correctness, not performance. When it is done, we will test it. If it is not
// fast enough, we will add more complexities to make it run faster.
func store(p *Pillar, data []Data) (int, error) {
	for i := range data {
		if err := p.Store(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// Copy knows how to pull and store data from the System.
// Now we can call the pull and store functions, passing Xenia and Pillar through.
func Copy(sys *System, batch int) error {
	data := make([]Data, batch)

	for {
		i, err := pull(&sys.Xenia, data)
		if i > 0 {
			if _, err := store(&sys.Pillar, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =============================================================================

func main() {
	sys := System{
		Xenia: Xenia{
			Host:    "localhost:8000",
			Timeout: time.Second,
		},
		Pillar: Pillar{
			Host:    "localhost:9000",
			Timeout: time.Second,
		},
	}

	if err := Copy(&sys, 3); err != io.EOF {
		fmt.Println(err)
	}
}
```

### Decoupling With Interface

Now that we have that concrete implementation, we can decouple this code from
change.

[Sample program](decoupling/example2/example2.go)

```go
func init() {
	rand.Seed(time.Now().UnixNano())
}

// =============================================================================

// Data is the structure of the data we are copying.
type Data struct {
	Line string
}

// =============================================================================

// Puller declares behavior for pulling data.
type Puller interface {
	Pull(d *Data) error
}

// Storer declares behavior for storing data.
type Storer interface {
	Store(d *Data) error
}

// =============================================================================

// Xenia is a system we need to pull data from.
type Xenia struct {
	Host    string
	Timeout time.Duration
}

// Pull knows how to pull data out of Xenia.
func (*Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF

	case 5:
		return errors.New("Error reading data from Xenia")

	default:
		d.Line = "Data"
		fmt.Println("In:", d.Line)
		return nil
	}
}

// Pillar is a system we need to store data into.
type Pillar struct {
	Host    string
	Timeout time.Duration
}

// Store knows how to store data into Pillar.
func (*Pillar) Store(d *Data) error {
	fmt.Println("Out:", d.Line)
	return nil
}

// =============================================================================

// System wraps Xenia and Pillar together into a single system.
type System struct {
	Xenia
	Pillar
}

// =============================================================================

// pull knows how to pull bulks of data from any Puller.
func pull(p Puller, data []Data) (int, error) {
	for i := range data {
		if err := p.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store bulks of data from any Storer.
func store(s Storer, data []Data) (int, error) {
	for i := range data {
		if err := s.Store(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// Copy knows how to pull and store data from the System.
func Copy(sys *System, batch int) error {
	data := make([]Data, batch)

	for {
		i, err := pull(&sys.Xenia, data)
		if i > 0 {
			if _, err := store(&sys.Pillar, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =============================================================================

func main() {
	sys := System{
		Xenia: Xenia{
			Host:    "localhost:8000",
			Timeout: time.Second,
		},
		Pillar: Pillar{
			Host:    "localhost:9000",
			Timeout: time.Second,
		},
	}

	if err := Copy(&sys, 3); err != io.EOF {
		fmt.Println(err)
	}
}
```

By looking at the API (functions), we need to decouple the API from the concrete
implementation. The decoupling that we do must get all the way down into
initialization. To do this right, the only piece of code that we need to change
is initialization. Everything else should be able to act on the behavior that
these types are going to provide.

`pull` is based on the concrete. It only knows how to work on Xenia. However,
if we are able to decouple `pull` to use any system that know how to pull data,
we can get the highest level of decoupling. Since the algorithm we have is
already efficient, we don't need to add another level of generalization and
destroy the work we did in the concrete. Same thing with `store`.

It is nice to work from the concrete up. When we do this, not only we are
solving problem efficiently and reducing technical debt but the contracts, they
come to us. We already know what the contract is for pulling/storing data.
We already validate that and this is what we need.

Let's just decouple these 2 functions and add 2 interfaces.
The `Puller` interface knows how to pull and the `Storer` knows how to store.
Xenia already implemented the `Puller` interface and Pillar already implemented
the `Storer` interface. Now we can come into pull/store, decouple this function
from the concrete. Instead of passing Xenial and Pillar, we pass in the `Puller`
and `Storer`. The algorithm doesn't change. All we doing is now calling
pull/store indirectly through the interface value.

`Copy` also doesn't have to change because Xenia/Pillar already implemented
the interfaces. However, we are not done because `Copy` is still bounded to the
concrete. `Copy` can only work with pointer of type system. We need to decouple
`Copy` so we can have a decoupled system that knows how to pull and store.
We will do it in the next file.

### Interface Composition

[Sample program](decoupling/example3/example3.go).

```go
func init() {
	rand.Seed(time.Now().UnixNano())
}

// =============================================================================

// Data is the structure of the data we are copying.
type Data struct {
	Line string
}

// =============================================================================

// Puller declares behavior for pulling data.
type Puller interface {
	Pull(d *Data) error
}

// Storer declares behavior for storing data.
type Storer interface {
	Store(d *Data) error
}

// PullStorer declares behavior for both pulling and storing.
type PullStorer interface {
	Puller
	Storer
}

// =============================================================================

// Xenia is a system we need to pull data from.
type Xenia struct {
	Host    string
	Timeout time.Duration
}

// Pull knows how to pull data out of Xenia.
func (*Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF

	case 5:
		return errors.New("Error reading data from Xenia")

	default:
		d.Line = "Data"
		fmt.Println("In:", d.Line)
		return nil
	}
}

// Pillar is a system we need to store data into.
type Pillar struct {
	Host    string
	Timeout time.Duration
}

// Store knows how to store data into Pillar.
func (*Pillar) Store(d *Data) error {
	fmt.Println("Out:", d.Line)
	return nil
}

// =============================================================================

// System wraps Xenia and Pillar together into a single system.
type System struct {
	Xenia
	Pillar
}

// =============================================================================

// pull knows how to pull bulks of data from any Puller.
func pull(p Puller, data []Data) (int, error) {
	for i := range data {
		if err := p.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store bulks of data from any Storer.
func store(s Storer, data []Data) (int, error) {
	for i := range data {
		if err := s.Store(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// Copy knows how to pull and store data from any System.
func Copy(ps PullStorer, batch int) error {
	data := make([]Data, batch)

	for {
		i, err := pull(ps, data)
		if i > 0 {
			if _, err := store(ps, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =============================================================================

func main() {
	sys := System{
		Xenia: Xenia{
			Host:    "localhost:8000",
			Timeout: time.Second,
		},
		Pillar: Pillar{
			Host:    "localhost:9000",
			Timeout: time.Second,
		},
	}

	if err := Copy(&sys, 3); err != io.EOF {
		fmt.Println(err)
	}
}
```

Let's just add another interface. Let's use interface composition to do this.
`PullStorer` has both behaviors: `Puller` and `Storer`. Any concrete type that
implement both pull and store is a `PullStorer`. `System` is a `PullStorer`
because it is embedded of these 2 types, `Xenia` and `Pillar`. Now we just need
to go into `Copy`, replace the `System` pointer with `PullStorer` and no other
code need to change.

Looking closely at `Copy`, there is something that could potentially confuse us.
We are passing the `PullStorer` interface value directly into `pull` and `store`
respectively. If we look into `pull` and `store`, they don't want a `PullStorer`.
One want a `Puller` and one want a `Storer`. Why does the compiler allow us to
pass a value of different type value while it didn't allow us to do that before?
This is because Go has what is called: **implicit interface conversion**.
This is possible because:
- All interface values have the exact same model (implementation details).
- If the type information is clear, the concrete type that exists in one
interface has enough behaviors for another interface. It is true that any
concrete type that is stored inside of a `PullStorer` must also implement the
`Storer` and `Puller`.

Let's walk through this code. In the main function, we are creating a value of
our `System` type. As we know, our `System` type value is based on the embedding
of two concrete types: `Xenia` and `Pillar`, where `Xenia` knows how to `pull`
and `Pillar` knows how to `store`. Because of inner type promotion, `System`
knows how to `pull` and `store` both inherently. We are passing the address of
our `System` to `Copy`. `Copy` then creates the `PullStorer` interface.
The first word is a `System` pointer and the second word point to the original
value. This interface now knows how to `pull` and `store`. When we call `pull`
off of `ps`, we call `pull` off of `System`, which eventually call pull off
of `Xenia`. Here is the kicker: the implicit interface conversion.
We can pass the interface value `ps` to `pull` because the compiler knows that
any concrete type stored inside the `PullStorer` must also implement `Puller`.
We end up with another interface called `Puller`. Because the memory models are
the same for all interfaces, we just copy those 2 words so they are all sharing
the same interface type. Now when we call `pull` off of `Puller`, we call `pull`
off of `System`. Similar to `Storer`. All using value semantic for the
interface value and pointer semantic to share.

```
    ps                          system
+---------+             +-------------------+
|         +--+pull      | +-------+         |
| *System |             | | Xenia +--+pull  +--+pull
|         +--+store     | |       |         |
|         |             | +-------+         +--+store
+---------+             | +-------+         |
|         |             | |Pillar +--+store |
|    *    +------------>+ |       |         |
|         |             | +-------+         |
|         |             +---------------+---+
+---------+                             ^
                                        |
     p                  s               |
+---------+        +---------+          |
|         +--+pull |         +--+store  |
| *System |        | *System |          |
|         |        |         |          |
|         |        |         |          |
+---------+        +---------+          |
|         |        |         |          |
|    *    +--------+    *    +----------+
|         |        |         |
|         |        |         |
+---------+        +---------+
```

Our `System` type is still concrete system type because it is still based on
two concrete types, `Xenial` and `Pillar`. If we have another `System`,
say `Bob`, we have to change in type `System` struct.
This is not scalable.

### Decoupling With Interface Composition

[Sample code](decoupling/example4/example4.go).


```go
func init() {
	rand.Seed(time.Now().UnixNano())
}

// =============================================================================

// Data is the structure of the data we are copying.
type Data struct {
	Line string
}

// =============================================================================

// Puller declares behavior for pulling data.
type Puller interface {
	Pull(d *Data) error
}

// Storer declares behavior for storing data.
type Storer interface {
	Store(d *Data) error
}

// PullStorer declares behavior for both pulling and storing.
type PullStorer interface {
	Puller
	Storer
}

// =============================================================================

// Xenia is a system we need to pull data from.
type Xenia struct {
	Host    string
	Timeout time.Duration
}

// Pull knows how to pull data out of Xenia.
func (*Xenia) Pull(d *Data) error {
	switch rand.Intn(10) {
	case 1, 9:
		return io.EOF

	case 5:
		return errors.New("Error reading data from Xenia")

	default:
		d.Line = "Data"
		fmt.Println("In:", d.Line)
		return nil
	}
}

// Pillar is a system we need to store data into.
type Pillar struct {
	Host    string
	Timeout time.Duration
}

// Store knows how to store data into Pillar.
func (*Pillar) Store(d *Data) error {
	fmt.Println("Out:", d.Line)
	return nil
}

// =============================================================================

// System wraps Pullers and Stores together into a single system.
type System struct {
	Puller
	Storer
}

// =============================================================================

// pull knows how to pull bulks of data from any Puller.
func pull(p Puller, data []Data) (int, error) {
	for i := range data {
		if err := p.Pull(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// store knows how to store bulks of data from any Storer.
func store(s Storer, data []Data) (int, error) {
	for i := range data {
		if err := s.Store(&data[i]); err != nil {
			return i, err
		}
	}

	return len(data), nil
}

// Copy knows how to pull and store data from any System.
func Copy(ps PullStorer, batch int) error {
	data := make([]Data, batch)

	for {
		i, err := pull(ps, data)
		if i > 0 {
			if _, err := store(ps, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =============================================================================

func main() {
	sys := System{
		Puller: &Xenia{
			Host:    "localhost:8000",
			Timeout: time.Second,
		},
		Storer: &Pillar{
			Host:    "localhost:9000",
			Timeout: time.Second,
		},
	}

	if err := Copy(&sys, 3); err != io.EOF {
		fmt.Println(err)
	}
}
```

We change our concrete type `System`. Instead of using 2 concrete types `Xenia`
and `Pillar`, we use 2 interface types `Puller` and `Storer`. Our concrete type
`System` where we can have concrete behaviors is now based on the embedding of 2
interface types. It means that we can inject any data, not based on the common
DNA but on the data that providing the capability, the behavior that we need.
Now our code can be fully decouple because any value that implements the `Puller`
interface can be stored inside the `System` (same with `Storer` interface). We
can create multiple `System`s and that data can be passed in `Copy`.
We don't need method here. We just need one function that accept data and its
behavior will change based on the data we put in.

Now `System` is not based on `Xenia` and `Pillar` anymore. It is based on 2
interfaces, one that stores `Xenia` and one that stores `Pillar`. We get the
extra layer of decoupling. If the `System` change, no big deal. We replace the
`System` as we need to during the program startup.

We solve this problem. We put this in production. Every single refactoring that
we did went into production before we did the next one. We keep minimizing
technical debt.

### Remove Interface Pollution

Code readability review: since I'm not going to have more than one `System` type, this interface is now
pollution. I can go back and work in the concrete, because I'm never going to
have a second implementation of `PullStorer`, which means that I can go ahead
and remove code.

[Sample code](decoupling/example5/example5.go)

```go
// PullStorer declares behavior for both pulling and storing.
// type PullStorer interface {
// 	Puller
// 	Storer
// }

// Copy knows how to pull and store data from any System.
func Copy(sys *System, batch int) error {
	data := make([]Data, batch)

	for {
		i, err := pull(sys, data)
		if i > 0 {
			if _, err := store(sys, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =============================================================================

func main() {
	sys := System{
		Puller: &Xenia{
			Host:    "localhost:8000",
			Timeout: time.Second,
		},
		Storer: &Pillar{
			Host:    "localhost:9000",
			Timeout: time.Second,
		},
	}

	if err := Copy(&sys, 3); err != io.EOF {
		fmt.Println(err)
	}
}
```

### More Precise API

Another code readability review: This idea of data injection is cool. But it's
clever because this API isn't as precise as it could be.

`System` is really hiding the cost of initialization. We don't really know what
a `System` is when we look at `Copy`. We don't realize we need these 2
behaviors. But now with this new API definition we've reduced fraud, we're going
to make this API easier to use, easier to test, because it's exactly saying,
give me concrete data that implements `Puller`, give me concrete data that
implements `Storer`.

[Sample program](decoupling/example6/example6.go)

```go
// Copy knows how to pull and store data from any System.
func Copy(p Puller, s Storer, batch int) error {
	data := make([]Data, batch)

	for {
		i, err := pull(p, data)
		if i > 0 {
			if _, err := store(s, data[:i]); err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}
}

// =============================================================================

func main() {
	x := Xenia{
		Host:    "localhost:8000",
		Timeout: time.Second,
	}

	p := Pillar{
		Host:    "localhost:9000",
		Timeout: time.Second,
	}

    // Share Xenia (&x) and share Pillar (&p) with Copy.
	if err := Copy(&x, &p, 3); err != io.EOF {
		fmt.Println(err)
	}
}
```

## Conversion and Assertions

Let's go explore a little bit deeper the idea that we're passing concrete data
across these program boundaries when we're using interface values.

### Interface Conversions

[Sample code](assertions/example1/example1.go).

```go
// Mover provides support for moving things.
type Mover interface {
	Move()
}

// Locker provides support for locking and unlocking things.
type Locker interface {
	Lock()
	Unlock()
}

// MoveLocker provides support for moving and locking things.
type MoveLocker interface {
	Mover
	Locker
}

// bike represents a concrete type for the example.
type bike struct{}

// Move can change the position of a bike.
func (bike) Move() {
	fmt.Println("Moving the bike")
}

// Lock prevents a bike from moving.
func (bike) Lock() {
	fmt.Println("Locking the bike")
}

// Unlock allows a bike to be moved.
func (bike) Unlock() {
	fmt.Println("Unlocking the bike")
}

func main() {
	// Declare variables of the MoveLocker and Mover interfaces set to their
	// zero value.
	var ml MoveLocker
	var m Mover

	// Create a value of type bike and assign the value to the MoveLocker
	// interface value.
	ml = bike{}

	// An interface value of type MoveLocker can be implicitly converted into
	// a value of type Mover. They both declare a method named move.
	m = ml

	//    ml                          m
	// +------+                    +------+
	// | bike |        bike        | bike |
	// +------+      +------+      +------+
	// |  *   | ---> |      | <--- |  *   |
	// +------+      +------+      +------+

    // However, we can't go in the other direction, like so:
    ml = m
    // The compiler will say:
	// prog.go:68: cannot use m (type Mover) as type MoveLocker in assignment:
    //	   Mover does not implement MoveLocker (missing Lock method)

    // **************
    // Type assertion
    // **************
	// Interface type Mover does not declare methods named lock and unlock.
	// Therefore, the compiler can't perform an implicit conversion to assign
	// a value of interface type Mover to an interface value of type MoveLocker.
	// It is irrelevant that the concrete type value of type bike that is stored
	// inside of the Mover interface value implements the MoveLocker interface.

	// We can perform a type assertion at runtime to support the assignment.

	// Perform a type assertion against the Mover interface value to access
	// a COPY of the concrete type value of type bike that was stored inside
	// of it. Then assign the COPY of the concrete type to the MoveLocker
    // interface.

    // This is the syntax for type assertion.
    // We are taking the interface value itself, dot (bike). We are using bike
    // as an parameter. If m is not nil and there is a bike inside of m, we will
    // get a copy of it since we are using value semantic. Or else, a panic
    // occurs. b is having a copy of bike value.
    b := m.(bike)

    // We can prevent panic when type assertion breaks by destructuring
    // the boolean value that represents type assertion result
    b, ok := m.(bike)
    fmt.Println("Does m has value of bike?:", ok)

	ml = b

	// It's important to note that the type assertion syntax provides a way
	// to state what type of value is stored inside the interface. This is
	// more powerful from a language and readability standpoint, than using
	// a casting syntax, like in other languages.
}
```

## Runtime Type Assertions

[Code sample](assertions/example2/example2.go).

```go
// car represents something you drive.
type car struct{}

// String implements the fmt.Stringer interface.
func (car) String() string {
	return "Vroom!"
}

// cloud represents somewhere you store information.
type cloud struct{}

// String implements the fmt.Stringer interface.
func (cloud) String() string {
	return "Big Data!"
}

func main() {
	// Seed the number random generator.
	rand.Seed(time.Now().UnixNano())

	// Create a slice of the Stringer interface values.

	//  ---------------------
	// |   car   |   cloud   |
	//  ---------------------
	// |    *    |     *     |
	//  ---------------------
	//      A          A
	//      |          |
	//     car       cloud
	//    -----      -----
	//   |     |    |     |
	//    -----      -----

	mvs := []fmt.Stringer{
		car{},
		cloud{},
	}

	// Let's run this experiment ten times.
	for i := 0; i < 10; i++ {
		// Choose a random number from 0 to 1.
		rn := rand.Intn(2)

        // Perform a type assertion that we have a concrete type
        // of cloud in the interface value we randomly chose.
		// This shows us that this checking is at runtime, not compile time.
		if v, ok := mvs[rn].(cloud); ok {
			fmt.Println("Got Lucky:", v)
			continue
		}

		fmt.Println("Got Unlucky")
	}
}
```

We have to guarantee that variable in question (x in `x.(T)`) can always be
asserted correctly as `T` type.

Or else, we wouldn't want to use that `ok` variable because we want it to panic
if there is an integrity issue. We must shut it down immediately if that happens
if we cannot recover from a panic and guarantee that we are back at 100%
integrity, the software has to be restarted. Shutting down means you have to
call `log.Fatal`, `os.exit`, or `panic` for stack trace.
When we use type assertion, we need to understand when it is okay that whatever
we are asking for is not there.

Important note: if the type assertion is causing us to call the concrete value
out, that should raise a big flag. We are using interface to maintain a level of
decoupling and now we are using type assertion to go back to the concrete.
When we are in the concrete, we are putting our codes in the situation where
cascading changes can cause widespread refactoring. What we want with interface
is the opposite, internal changes minimize cascading changes.

## Interface Pollution

It comes from the fact that people are designing software from the interface
first down instead of concrete type up.

Why are we using an interface here?

Myth #1: We are using interfaces because we have to use interfaces.
Answer: No. We don't have to use interfaces. We use it when it is practical and
reasonable to do so. Even though they are wonderful, there is a cost of using
interfaces: a level of indirection and potential allocation when we store
concrete type inside of them. Unless the cost of that is worth whatever
decoupling we are getting, we shouldn't be using interfaces.

Myth #2: We need to be able to test our code so we need to use interfaces.
Answer: No. We must design our API that are usable for user application
developer first, not our test.

[Sample code](pollution/example1/example1.go).

Below is an example that creates interface pollution by improperly using an
interface when one is not needed.

```go
// Server defines a contract for tcp servers.
type Server interface {
	Start() error
	Stop() error
	Wait() error
}
```

This is a little bit of smell because this is some sort of APIs that going to be
exposed to user and already that is a lot of behaviors brought in a generic
interface.

```go
// server is our Server implementation.
type server struct {
	host string

	// PRETEND THERE ARE MORE FIELDS.
}

// NewServer returns an interface value of type Server with a server implementation.
func NewServer(host string) Server {
	// SMELL - Storing an unexported type pointer in the interface.
	return &server{host}
}
```

Here is the factory function. It immediately starts to smell even worse. It is
returning the interface value. It is not that functions and interfaces cannot
return interface values. They can. But normally, that should raise a flag.
The concrete type is the data that has the behavior and the interface
normally should be used as accepting the input to the data, not necessarily
going out.

```go
// Start allows the server to begin to accept requests.
func (s *server) Start() error {
	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Stop shuts the server down.
func (s *server) Stop() error {
	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Wait prevents the server from accepting new connections.
func (s *server) Wait() error {
	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}
```

This code here couldn't care less nor would it change if `srv` was the concrete
type, not the interface. The interface is not providing any level of support
whatsoever. There is no decoupling here that is happening. It is not giving us
anything special here. All is doing is causing us another level of indirection.

```go
func main() {
	// Create a new Server.
	srv := NewServer("localhost")

	// Use the API.
	srv.Start()
	srv.Stop()
	srv.Wait()
}

// =============================================================================

// NOTES:

// Smells:
//  * The package declares an interface that matches the entire API of its own concrete type.
//  * The interface is exported but the concrete type is unexported.
//  * The factory function returns the interface value with the unexported concrete type value inside.
//  * The interface can be removed and nothing changes for the user of the API.
//  * The interface is not decoupling the API from change.
```

### Remove Interface Pollution

This is basically just removing the improper interface usage from previous
example.

[Sample code](pollution/example2/example2.go)

```go
// Server is our Server implementation.
type Server struct {
	host string

	// PRETEND THERE ARE MORE FIELDS.
}

// NewServer returns an interface value of type Server with a server implementation.
func NewServer(host string) *Server {
	// SMELL - Storing an unexported type pointer in the interface.
	return &Server{host}
}

// Start allows the server to begin to accept requests.
func (s *Server) Start() error {
	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Stop shuts the server down.
func (s *Server) Stop() error {
	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Wait prevents the server from accepting new connections.
func (s *Server) Wait() error {
	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

func main() {
	// Create a new Server.
	srv := NewServer("localhost")

	// Use the API.
	srv.Start()
	srv.Stop()
	srv.Wait()
}
```

Here are some guidelines around interface pollution:
* Use an interface:
  * When users of the API need to provide an implementation detail.
  * When API’s have multiple implementations that need to be maintained.
  * When parts of the API that can change have been identified and require decoupling.
* Question an interface:
  * When its only purpose is for writing testable API’s (write usable API’s first).
  * When it’s not providing support for the API to decouple from change.
  * When it's not clear how the interface makes the code better.

## Mocking

We don't want to be using interfaces anytime we have to think about mocking.

It is important to mock things.
Most things over the network can be mocked in our test. However, mocking our
database is a different story because it is too complex. This is where Docker
can come in and simplify our code by allowing us to launch our database while
running our tests and have that clean database for everything we do.

Every API only need to focus on its test. We no longer have to worry about the
application user or user over API test. We used to worry about: if we don't have
that interface, the user who use our API can't write test. That is gone.
The example below will demonstrate the reason.

### Package to Mock

[Sample package](mocking/example1/pubsub/pubsub.go)

```go
// Package pubsub simulates a package that provides publication/subscription
// type services.
package pubsub

// PubSub provides access to a queue system.
type PubSub struct {
	host string

	// PRETEND THERE ARE MORE FIELDS.
}

// New creates a pubsub value for use.
func New(host string) *PubSub {
	ps := PubSub{
		host: host,
	}

	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.

	return &ps
}

// Publish sends the data for the specified key.
func (ps *PubSub) Publish(key string, v interface{}) error {
	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}

// Subscribe sets up an request to receive messages for the specified key.
func (ps *PubSub) Subscribe(key string) error {
	// PRETEND THERE IS A SPECIFIC IMPLEMENTATION.
	return nil
}
```

Imagine we are working at a company that decides to incorporate Go as a part of
its stack. They have their internal pubsub system that all applications are
supposed to used. Maybe they are doing event sourcing and there is a single
pubsub platform they are using that is not going to be replaced. They need the
pubsub API for Go that they can start building services that connect into this
event source. So what can change? Can the event source change? If the answer is
no, then it immediately tells us that we don't need to use interfaces. We can
built the entire API in the concrete, which we would do it first anyway. We then
write tests to make sure everything work.

A couple days later, they come to us with a problem. They have to write tests
and they cannot hit the pubsub system directly when my test run so they need to
mock that out. They want us to give them an interface. However, we don't need an
interface because our API doesn't need an interface. They need an interface,
not us. They need to decouple from the pubsub system, not us. They can do any
decoupling they want because this is Go.

[Sample client program](mocking/example1/example1.go) to show how you can
personally mock concrete types when you need to for your own packages or tests.

```go
// publisher is an interface to allow this package to mock the pubsub
// package support.
type publisher interface {
	Publish(key string, v interface{}) error
	Subscribe(key string) error
}
```

When we are writing our applications, declare our own interface that map out all
the APIs call we need for the APIs. The concrete types APIs in the previous files
satisfy it out of the box. We can write the entire application with mocking
decoupling from concrete implementations.

```go
// mock is a concrete type to help support the mocking of the pubsub package.
type mock struct{}

// Publish implements the publisher interface for the mock.
func (m *mock) Publish(key string, v interface{}) error {
	// ADD YOUR MOCK FOR THE PUBLISH CALL.
	return nil
}

// Subscribe implements the publisher interface for the mock.
func (m *mock) Subscribe(key string) error {
	// ADD YOUR MOCK FOR THE SUBSCRIBE CALL.
	return nil
}

func main() {
	// Create a slice of publisher interface values. Assign
	// the address of a pubsub.PubSub value and the address of
	// a mock value.
	pubs := []publisher{
		pubsub.New("localhost"),
		&mock{},
	}

	// Range over the interface value to see how the publisher
	// interface provides the level of decoupling the user needs.
	// The pubsub package did not need to provide the interface type.
	for _, p := range pubs {
		p.Publish("key", "value")
		p.Subscribe("key")
	}
}
```
