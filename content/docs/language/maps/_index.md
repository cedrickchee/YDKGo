---
title: Maps
weight: 8
---

# Maps

Maps provide a data structure that allow for the storage and management of key/value pair data.

- Maps provide a way to store and retrieve key/value pairs.
- Reading an absent key returns the zero value for the map's value type.
- Iterating over a map is always random.
- The map key must be a value that is comparable.
- Elements in a map are not addressable.
- Maps are a reference type.

## Declare, write, read, and delete

Declare and make a map that stores values of type user with a key of type string.

```go
users := make(map[string]user)
```

Add key/value pairs to the map.

```go
users["Roy"] = user{"Rob", "Roy"}
users["Ford"] = user{"Henry", "Ford"}
```

Iterate over the map printing each key and value.

```go
for key, value := range users {
    fmt.Println(key, value)
}
```

Delete the value at a specific key.

```go
delete(users, "Roy")
```

Reminder: if you use map as cache, memory leak if you don't remove key.

## Map literals

```go
// Declare and initialize the map with values.
users := map[string]user{
    "Roy":     {"Rob", "Roy"},
    "Ford":    {"Henry", "Ford"},
    "Mouse":   {"Mickey", "Mouse"},
    "Jackson": {"Michael", "Jackson"},
}
```

## Absent keys

If we need to check for the presence of a key we use a 2 variable assignment.
The 2nd variable is a bool.

If `found` is true, we will get a copy value of that type.
If found is false, `score` is still a value of type scores but is set to its
zero value.

```go
// Find the "anna" key.
score, found := scores["anna"]
```

## Map key restrictions

```go
// users defines a set of users.
type users []user

// Declare and make a map that uses a slice as the key.
u := make(map[users]int) // map key violations

// ./example3.go:22: invalid map key type users
```

Using this syntax, we can define a set of `users`.

This is a second way we can define `users`. We can use an existing type and use
it as a base for another type. These are two different types. There is no
relationship here. However, when we try use it as a key, like:
`u := make(map[users]int)` the complier says we cannot use that:
"invalid map key type users". The reason is: whatever we use for the key, the
value must be comparable. We have to use it in some sort of boolean expression
in order for the map to create a hash value for it.

## Sorting maps by key

```go
users := map[string]user{
    "Roy":     {"Rob", "Roy"},
    "Ford":    {"Henry", "Ford"},
    "Mouse":   {"Mickey", "Mouse"},
    "Jackson": {"Michael", "Jackson"},
}

// Pull the keys from the map.
var keys []string
for key := range users {
    keys = append(keys, key)
}

// Sort the keys alphabetically.
sort.Strings(keys)
```

Uses the standard library sort package.

Walk through a map by alphabetical key order. So, the output be not so random.

## Mechanical sympathy

Map is going to do its best to keep the data contiguous as well.

But think about that map really as one of those data structures when you want to
pinpoint and get to data quickly, based on a key.
