---
title: Functions
weight: 5
---

# Functions

- Functions can return multiple values and most return an error value.
- The error value should always be checked as part of the programming logic.
- The blank identifier can be used to ignore return values.

## Return multiple values

```go
func retrieveUser(name string) (*user, error) {
}
```

During a code review, if I see a function that returns three or more values, it raises a big flag with me. I think the standard library does it a half to a dozen times tops. So it's not common to have a function that returns three or more.

It's common in Go, to return both some sort of value and an error.

Go doesn't have constructors. So, this idea of a function getting executed during construction doesn't exist. This is a good thing. In Go, we have factory functions. A function that is called that returns an initialized value, and without the function, that value could not have been initialized properly in any other way.

A lot of factory functions in Go do start with the name `new`, but I would consider `retrieveUser` here a factory function as well. It is creating a value of type `user`, it's sharing it back up the call stack after it gets initialized for use. So, a function like this is very common in Go.

```go
// Make a call to get the user in a json response.
r, err := getUser(name)
if err != nil {
    return nil, err
}

// Unmarshal the json document into a value of
// the user struct type.
var u user // create a value of type user
// &u (ampersand operator means sharing).
// we're sharing that user value down the call stack, so the Unmarshal function can
// read this document, initialize it.
err = json.Unmarshal([]byte(r), &u)
// then, we're sharing user value back up the call stack. This is going to
// create an allocation
return &u, err
```

This example leverage the readability of the ampersand to show sharing. It starts with that value mechanic. Start with the value, let the convention in your code determine what's happening there, don't be clever.

## Scope

```go
var u *user // global variable

func main() {
    var u *user // this variable shadow the global variable, u above.

	// Retrieve the user profile.
	u, err := retrieveUser("sally")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Display the user profile.
	fmt.Printf("%+v\n", *u)
}
```

`if` statements, `for` statements, `switch` statements, they all have curly brackets. They all define a block of scope, and scope flows down.

```go
// Retrieve the user profile.
if u, err := retrieveUser("sally"); err != nil {
    fmt.Println(err)
    return
}
```

These statements, they actually come with their own block of scope, something that we don't have in other languages.

**Idiom in Go**

We're looking at reducing the number of lines of code, we're looking at simplifying things. We don't have to pre-declare any variables. We want to declare variables when we need them in the smallest level of scope possible. This lends itself to readability. But even more, it lends itself to being able to use short variable names. This is a core idiom in Go. And we like this idiom because the variable name, the variable itself, is probably the most used identifier in every function. So, the longer that variable name is, the more noise it creates in a function.

**Basic guideline**

The farther away a variable is being declared from where it's being used, the longer the name has to be because the context is missing. But the closest level of scope that is being used, like that `if` statement, the shorter it can be. It can actually be a single letter variable. And what also allows this is that we have to be very clear about our function names.

I don't want function names that are this big. But this function's called `retrieveUser`. Is there any chance that `u` does not represent a `user`? And is there any chance that we don't know that? I think it's slim to none. Calling `u` `user` here would really add noise to this function. And we want the logic to pop out. We don't want that noise, we want a good signal ratio here in what we're doing. So, your function names become important, and since we're declaring `u` now within the smallest scope possible, we can keep that down to a single letter name.

## Blank identifier

```go
if _, err := updateUser(&u); err != nil {
    fmt.Println(err)
    return
}
```

This allows us to not declare a variable for a value when we're required to.

## Error handling

Another thing about **readability**. There is this idea around exception handling that said, by being able to treat an error as an exception, what we can do is inside our `try` statement, write our happy path, our positive path, happy path code. When everything is going right, this is what it is. And then inside the `catch`, as an exception, this is what happens when everything's going wrong. There was this idea that we would get readability by separating these two mediums. But we've already learned in the first section that if you take integrity seriously, you cannot separate your error handling from your code. You cannot be an exception. It has to be part of your code. With that said, how do you get to write code that has happy path, and then has that negative path, that is still readable?

This is how we do it in Go. Have as many returns as you need in a function for readability. We don't want to see an `else` clause. Unless it's a single statement, `else`s do not help for readability. Don't use an `else` to push code to the bottom of a function. What we want to do is follow this **happy path, negative path**.

## Anonymous functions and closures

Declare an anonymous function and call it.

```go
func() {
    fmt.Println("Direct:", n)
}()
```

Declare an anonymous function and assign it to a variable.

```go
f := func() {
    fmt.Println("Variable:", n)
}

// Call the anonymous function through the variable.
f()
```
