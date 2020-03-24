---
title: Slices
weight: 7
---

# Slices

Slices are an incredibly important data structure in Go. They form the basis for how we manage and manipulate data in a flexible, performant and dynamic way. It is incredibly important for all Go programmers to learn how to uses slices.

- Slices are like dynamic arrays with special and built-in functionality.
- There is a difference between a slices length and capacity and they each service a purpose.
- Slices allow for multiple "views" of the same underlying array.
- Slices can grow through the use of the built-in function append.

## Declare and length

```go
slice1 := make([]string, 5)
slice1[0] = "Apple"
slice1[1] = "Orange"
slice1[2] = "Banana"
slice1[3] = "Grape"
slice1[4] = "Plum"
```

Create a slice with a length of 5 elements.

`make` is a special built-in function that only works with slice, map and channel.
`make` creates a slice that has an array of 5 strings behind it.
We are getting back a 3 word data structure:
- first word points to the backing array
- second word is length
- third word is capacity

```
// -------     -------------------------------
// |  *  | --> | nil | nil | nil | nil | nil |
// -------     |-----|-----|-----|-----|-----|
// |  5  |     |  0  |  0  |  0  |  0  |  0  |
// -------     -------------------------------
// |  5  |
// -------
```

### Length vs Capacity

Length is the number of elements from this pointer position we have access
to (read and write).

Capacity is the total number of elements from this pointer position that exist
in the backing array.

Syntactic sugar -> looks like array.
It also have the same cost that we've seen in array.
One thing to be mindful about: there is no value in the bracket `[]string`
inside the `make` function. With that in mind, we can constantly notice that we
are dealing with a slice, not array.

## Reference types

```go
fruits := make([]string, 5, 8)
fruits[0] = "Apple"
fruits[1] = "Orange"
fruits[2] = "Banana"
fruits[3] = "Grape"
fruits[4] = "Plum"
```

Create a slice with a length of 5 elements and a capacity of 8.

`make` allows us to adjust the capacity directly on construction of this
initialization. What we end up having now is a 3 word data structure where the
first word points to an array of 8 elements, length is 5 and capacity is 8.

```
// -------     -------------------------------------------------
// |  *  | --> | nil | nil | nil | nil | nil | nil | nil | nil |
// -------     |-----|-----|-----|-----|-----|-----|-----|-----|
// |  5  |     |  0  |  0  |  0  |  0  |  0  |  0  |  0  |  0  |
// -------     -------------------------------------------------
// |  8  |
// -------
```

It means that I can read and write to the first 5 elements and I have
3 elements of capacity that I can leverage later.

### Appending Slices

The idea of appending: making slice a dynamic data structure.

```go
// Declare a nil slice of strings, set to its zero value.
var data []string
```

This is a 3 word data structure: first one points to nil, second and last are zero.

**Empty literal construction for a slice**

What if I do `data := string{}`? Is it the same? No, because data in this case
is not set to its zero value.

This is why we always use var for zero value because not every type when we
create an empty literal we have its zero value in return. What actually happen
here is that we have a slice but it has a pointer (as opposed to nil).

This is consider an empty slice, not a nil slice.
There is a semantic between a nil slice and an empty slice.

Remember: when it comes to a reference type, any time a reference type is set
to its zero value, we consider it to be nil, interfaces, channels, maps, slices,
functions.

If we pass a nil slice to a marshal function, we get back a string that said
null but when we pass an empty slice, we get an empty JSON document.
But where does that pointer point to? It is an empty struct, which we will
review later.

```go
var data []string

// Capture the capacity of the slice.
lastCap := cap(data)

// Append ~100k strings to the slice.
for record := 1; record <= 1e5; record++ {

    // Use the built-in function append to add to the slice.
    value := fmt.Sprintf("Rec: %d", record)
    data = append(data, value)

    // When the capacity of the slice changes, display the changes.
    if lastCap != cap(data) {

        // Calculate the percent of change.
        capChg := float64(cap(data)-lastCap) / float64(lastCap) * 100

        // Save the new values for capacity.
        lastCap = cap(data)

        // Display the results.
        fmt.Printf("Addr[%p]\tIndex[%d]\t\tCap[%d - %2.f%%]\n",
            &data[0],
            record,
            cap(data),
            capChg)
    }
}
```

`append` allows us to add value to a slice, making the data structure dynamic,
yet still allows us to use that contiguous block of memory that gives us the
predictable access pattern from mechanical sympathy. The append call is working
with value semantic. Notice that append does mutate. But we're not using
pointers. We are not sharing this slice but appending to it and returning a
new copy of it. The slice gets to stay on the stack, not heap.

**What is memory leak in Go**

A memory leak in Go is when you maintain a reference to a value in the heap, and
that reference never goes away. This is complicated because you can't instrument
for memory leak when it's reference based. Who is to say that reference is or
isn't supposed to be there at any given time. So, if you think you have a
memory leak, which is the only way to look at the GC trace. Is the memory going
up on every garbage collection? And if it is, we have our memory leak.

Let's keep going here. Every time append runs, it checks the length and capacity.
If it is the same, it means that we have no room. append creates a new backing
array, double its size, copy the old value back in and append the new value.
It mutates its copy on its stack frame and return us a copy. We replace our
slice with the new copy. If it is not the same, it means that we have extra
elements of capacity we can use. Now we can bring these extra capacity into the
length and no copy is being made. This is very efficient.

```go
// Outputs:
// Addr[0xc000010200]      Index[1]                Cap[1 - +Inf%]
// Addr[0xc00000c080]      Index[2]                Cap[2 - 100%]
// Addr[0xc000064080]      Index[3]                Cap[4 - 100%]
// Addr[0xc00007e000]      Index[5]                Cap[8 - 100%]
// Addr[0xc000080000]      Index[9]                Cap[16 - 100%]
// Addr[0xc00007c200]      Index[17]               Cap[32 - 100%]
// Addr[0xc000082000]      Index[33]               Cap[64 - 100%]
// Addr[0xc000084000]      Index[65]               Cap[128 - 100%]
// Addr[0xc000079000]      Index[129]              Cap[256 - 100%]
// Addr[0xc000086000]      Index[257]              Cap[512 - 100%]
// Addr[0xc00008a000]      Index[513]              Cap[1024 - 100%]
// Addr[0xc000090000]      Index[1025]             Cap[1280 - 25%]
// Addr[0xc00009a000]      Index[1281]             Cap[1704 - 33%]
// Addr[0xc0000b2000]      Index[1705]             Cap[2560 - 50%]
// Addr[0xc0000c0000]      Index[2561]             Cap[3584 - 40%]
// Addr[0xc0000d4000]      Index[3585]             Cap[4608 - 29%]
// Addr[0xc0000ec000]      Index[4609]             Cap[6144 - 33%]
// Addr[0xc00010e000]      Index[6145]             Cap[7680 - 25%]
// Addr[0xc000134000]      Index[7681]             Cap[9728 - 27%]
// Addr[0xc000166000]      Index[9729]             Cap[12288 - 26%]
// Addr[0xc0001a6000]      Index[12289]            Cap[15360 - 25%]
// Addr[0xc0001f4000]      Index[15361]            Cap[19456 - 27%]
// Addr[0xc000258000]      Index[19457]            Cap[24576 - 26%]
// Addr[0xc0002d6000]      Index[24577]            Cap[30720 - 25%]
// Addr[0xc000372000]      Index[30721]            Cap[38400 - 25%]
// Addr[0xc000434000]      Index[38401]            Cap[48128 - 25%]
// Addr[0xc00053c000]      Index[48129]            Cap[60416 - 26%]
// Addr[0xc000628000]      Index[60417]            Cap[75776 - 25%]
// Addr[0xc000750000]      Index[75777]            Cap[94720 - 25%]
// Addr[0xc0008c2000]      Index[94721]            Cap[118784 - 25%]
```

Looking at the last column in the output, when the backing array is 1000
elements or less, it doubles the size of the backing array for growth. Once we
pass 1000 elements, growth rate moves to 25%.

### Taking Slices of Slices

```go
slice2 := slice1[2:4] // slicing syntax
```

Take a slice of `slice1`. We want just indexes 2 and 3.
The length of `slice2` is 2 and capacity is 6.
Parameters are [starting_index : (starting_index + length)]

```go
// Outputs:
// Length[5] Capacity[8]
// [0] 0xc00007c000 Apple
// [1] 0xc00007c010 Orange
// [2] 0xc00007c020 Banana
// [3] 0xc00007c030 Grape
// [4] 0xc00007c040 Plum
// Length[2] Capacity[6]
// [0] 0xc00007c020 Banana
// [1] 0xc00007c030 Grape
```

By looking at the output, we can see that they are sharing the same backing array.
Thes slice headers get to stay on the stack when we use these value semantics.
Only the backing array that needed to be on the heap.

**Side effects**

```go
// Change the value of the index 0 of slice2.
slice2[0] = "CHANGED"
```

When we change the value of the index 0 of slice2, who are going to see this
change? The answer is both. We have to always to aware that we are modifying an
existing slice. We have to be aware who are using it, who is sharing that
backing array.

How about `slice2 := append(slice2, "CHANGED")`?
Similar problem will occur with `append` if the length and capacity is not the
same.
Instead of changing `slice2` at index 0, we call `append` on `slice2`. Since the
length of `slice2` is 2, capacity is 6 at the moment, we have extra rooms for
modification. We go and change the element at index 3 of `slice2`, which is
index 4 of `slice2`. That is very dangerous.

So, what if the length and capacity is the same? Instead of making `slice2`
capacity 6, we set it to 2 by adding another parameter to the slicing syntax
like this: `slice2 := slice1[2:4:4]`.
When `append` looks at this slice and see that the length and capacity is the same, it wouldn't
bring in the element at index 4 of slice1. It would detach. `slice2` will have
a length of 2 and capacity of 2, still share the same backing array.
On the call to `append`, length and capacity will be different. The addresses
are also different. This is called 3 index slice. This new slice will get its
own backing array and we don't affect anything at all to our original slice.

**Copy a slice**

```go
// Make a new slice big enough to hold elements of slice 1 and copy the
// values over using the builtin copy function.
slice3 := make([]string, len(slice1))
copy(slice3, slice1)
```

`copy` only works with string and slice only.

### Slices and References

[Sample program](example5/example5.go) to show how one needs to be careful when appending to a slice
when you have a reference to an element.

```go
type user struct {
	likes int
}

// Declare a slice users with 3 values.
users := make([]user, 3)

// Share the user at index 1.
shareUser := &users[1]

// Add a like for the user that was shared.
shareUser.likes++

// Display the number of likes for all users.
for i := range users {
    fmt.Printf("User: %d Likes: %d\n", i, users[i].likes)
}

// Append a new value to the slice.
// This line of code raises a red flag.
// users is a slice with length 3, capacity 3. Since the length and capacity is
// the same, we're now going to have to create a new backing array.
// Our new backing array will be doubling in size, and our length and capacities
// are going to change. append then copy values over. users nows points to diffrent
// memory block and has a length of 4, capacity of 6.
users = append(users, user{})

// We continue increment likes for the user that was shared.
shareUser.likes++

// Notice the last like has not been recorded.
// When we change the value of the second element of the slice, it is not change
// because it points to the old slice. Everytime we read it, we will get the
// wrong value.

// By displaying the number of likes for all users, we can see that we are in trouble.
fmt.Println("*************************")
for i := range users {
    fmt.Printf("User: %d Likes: %d\n", i, users[i].likes)
}
```

In this case, we kind of have a memory leak in a sense that this memory can't
get released because of that pointer. And you might think this would never
happen, and we know some of the best Go developers on the planet who have
created bugs like this, not thinking how `append`, when length and capacity are
the same, is going to make a copy of the current data, which is the point of
truth, and now make this the point of truth, and yet we have these pointers now
working against the old data structures.

This is a side effect, these are the nasty bugs that are so hard to find, and so
anytime we're working with pointer semantics, that's great, it's gonna give us
levels of efficiency, right, we have to be careful there, but we also have to
make sure that we're very clean with data updates, like with slices, and that
our mutations are not going to cause problems, or the mutations are happening
in the wrong place.

#### Strings in Go

Strings in Go are UTF-8 based.
If we use different encoding scheme, we might have a problem.

[Sample program](example6/example6.go).

What's interesting about UTF-8 is that it's a three layer character set.
You've got bytes at the bottom, in the middle you have what are call code points.
And a code point is a 32-bit or 4-byte value. And then, after code points,
you have characters.

A code point is anywhere from one to four bytes. A character is anywhere from
one to multiple code points. You have this, kind of like, n-tiered type of
character set.

```go
// Declare a literal string with both Chinese and English characters.
s := "世界 means world"
```

This string actually is going to be 18 bytes. Why is that?
For each Chinese character, we need 3 byte for each one.
Because UTF-8 is built on 3 layers: bytes, code point and character. From Go
perspective, string are just bytes. That is what we are storing.

In our example, the first 3 bytes represents a single code point that represents
that single character. We can have anywhere from 1 to 4 bytes representing a
code point (a code point is a 32 bit value) and anywhere from 1 to multiple code
points can actually represent a character. To keep it simple, we only have 3
bytes representing 1 code point representing 1 character. So we can read `s` as
3 bytes, 3 bytes, 1 byte, 1 byte, ... (since there are only 2 Chinese characters
in the first place, the rest are English).

```go
// buf is an array, it's not a slice.
// utf8.UTFMax is a constant, which represents the max number of bytes you need for code point.
// UTFMax is 4 -- up to 4 bytes per encoded rune.
// So buf is an array of four bytes.
var buf [utf8.UTFMax]byte
```

Maximum number of bytes we need to represent any code point is 4.

```go
// Wait, we can range over string?

for i, r := range s {

    // Capture the number of bytes for this rune.
    rl := utf8.RuneLen(r)

    // Calculate the slice offset for the bytes associated with this rune.
    si := i + rl

    // Copy of rune from the string to our buffer.
    copy(buf[:], s[i:si]) // we're slicing the string, s here.

    // Display the details.
    fmt.Printf("%2d: %q; codepoint: %#6x; encoded bytes: %#v\n", i, r, r, buf[:rl])
}
```

When we are ranging over a string, are we doing it byte by byte or code point by
code point or character by character?
The answer is code point by code point.
On the first iteration, `i` is 0. On the next one, `i` is 3 because we are
moving to the next code point. Then `i` is 6.

`r` represents the type rune. Rune really isn't a type in Go. Rune is its own
type. It is an alias for int32 type. In fact, similar to type byte we are using,
it is just an alias for uint8.

`copy` is a built-in function and only works with slices and string.

We want to go through every code point and copy them into our array `buf`, and
display them on the screen.

`buf[:]`: the syntax in Go allows us to apply by slicing syntax to an array.

> Every array is just a slice waiting to happen. — favorite sayings in Go

That slicing syntax will create a new slice value using the backing array `buf` as
our storage, and setting length and capacity to 4 in our slice header.
All of them are on the stack. There is no allocation here.

### Range Mechanics

[Sample program](example8/example8.go) to show how the `for` range has both
value and pointer semantics.

Using the value semantic form of the `for` range.

```go
friends := []string{"Annie", "Betty", "Charley", "Doug", "Edward"}
for _, v := range friends {
    friends = friends[:2]
    fmt.Printf("v[%s]\n", v)
}

// Outputs:
// v[Annie]
// v[Betty]
// v[Charley]
// v[Doug]
// v[Edward]
```

Using the pointer semantic form of the `for` range.

```go
friends = []string{"Annie", "Betty", "Charley", "Doug", "Edward"}
for i := range friends {
    friends = friends[:2]
    fmt.Printf("v[%s]\n", friends[i])
}

// Outputs:
// v[Annie]
// v[Betty]
// panic: runtime error: index out of range [2] with length 2

// goroutine 1 [running]:
// main.main()
//         /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/language/slices/example8/example8.go:24 +0x1fd
// exit status 2
```
