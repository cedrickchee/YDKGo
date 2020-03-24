---
title: Benchmarking
weight: 2
---

# Benchmarking

Go has support for testing the performance of your code.

## Basic Benchmarking

Let's look at the benchmarking mechanics and look at how we can write these and
run a couple of them right there on the command line.

Benchmark file's have to have `<file_name>_test.go` and use the `Benchmark`
functions like below. The goal is to know what perform better and what allocate
more or less between `Sprint` and `Sprintf`.

Our guess is that `Sprint` is gonna be better because it doesn't have any
overhead doing the formatting. However, this is not true. Remember we have to
optimize for correctness so we don't want to guess.

Since this code's going to be compiled, and the compiler is capable of throwing
**dead code** away. With the SSA backend, it is capable of identifying, like this
scenario, you call a function like `Sprint`, and it doesn't really mutate state.
It's using value semantics. It gets its own copy of the string, returns a new
string. So, if you don't capture the new string coming out of `Sprint` or
`Sprintf`, technically the compiler can identify that, let's not waste any CPU
cycles making this function call. Suddenly, you're trying to do a benchmark like
this to figure out what's faster. It's blazing fast for both, because you didn't
capture the output.

[Sample benchmark](basic/basic_test.go).

```go
import (
	"fmt"
	"testing"
)

var gs string

func BenchmarkSprint(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}

	gs = s
}
```

`BenchmarkSprint` tests the performance of using `Sprint`.
All the code we want to benchmark need to be inside the `b.N` `for` loop.
The first time the tool call it, `b.N` is equal to 1. It will keep increasing
the value of `N` and run long enough based on our bench time.
`fmt.Sprint` returns a value and we want to capture this value so it doesn't
look like dead code. We assign it to the global variable `gs`.

```go
func BenchmarkSprintf(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello")
	}

	gs = s
}
```

`BenchmarkSprint` tests the performance of using `Sprintf`.

Run all the benchmarks:

```sh
go test -run none -bench . -benchtime 3s -benchmem
```

We're usig `go test` like we did before. But since there are no functions that
start with test, I don't want to waste the tools time looking for them. That's
why I use `-run none`. `benchmem` flag tell it to look at memory allocations.

Sample output:

```sh
goos: linux
goarch: amd64
pkg: github.com/cedrickchee/ultimate-go/testing/benchmarks/basic
BenchmarkSprint-4       33619288               103 ns/op               5 B/op          1 allocs/op
BenchmarkSprintf-4      43777660                81.0 ns/op             5 B/op          1 allocs/op
PASS
ok      github.com/cedrickchee/ultimate-go/testing/benchmarks/basic     9.914s
```

column 2 (e.g. 33619288): is showing us how many times we executed that code.
column 3 (e.g. 5 B/op): 5 bytes allocated over one object.
column 4 (e.g. 1 allocs/op): one object on the heap worth 5 bytes of memory.

`Sprintf` ran faster, at 81.0 ns/op with the same memory allocation. In other
words, our guess was wrong. `Sprintf` is actually faster than `Sprint`.

We don't want to guess about performance.

## Sub Benchmarks

Like sub test, we can also do sub benchmark.

[Sample benchmark](sub/sub_test.go).

```go
// BenchmarkSprint tests all the Sprint related benchmarks as sub benchmarks.
func BenchmarkSprintSub(b *testing.B) {
	b.Run("none", benchSprint)
	b.Run("format", benchSprintf)
}

// benchSprint tests the performance of using Sprint.
func benchSprint(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}

	gs = s
}

// benchSprintf tests the performance of using Sprintf.
func benchSprintf(b *testing.B) {
	var s string

	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello")
	}

	gs = s
}
```

I can do the go test again.

```sh
go test -run none -bench . -benchtime 3s -benchmem
```

Sample output:

```sh
goos: linux
goarch: amd64
pkg: github.com/cedrickchee/ultimate-go/testing/benchmarks/sub
BenchmarkSprint/none-4          27487278               110 ns/op               5 B/op          1 allocs/op
BenchmarkSprint/format-4        40980721                82.0 ns/op             5 B/op          1 allocs/op
PASS
ok      github.com/cedrickchee/ultimate-go/testing/benchmarks/sub       6.597s
```

More ways to run these benchmarks:

```sh
go test -run none -bench BenchmarkSprint/none -benchtime 3s -benchmem
go test -run none -bench BenchmarkSprint/format -benchtime 3s -benchmem
```

I would be very careful of running sub benchmarks in parallel. You've got to
make sure your machine is idle when you're doing all of this.

## Validate Benchmarks

One of the most important things you must do when benchmarking is validate that
the results are accurate.

There was this article the author said, "What I want to see is that if I do a
merge sort over a million integers and if I use a Goroutine every time I split
that list in half, I throw a different Goroutine at each half, how fast is
that?". The author wanted to know what was the fastest way to go and the
conclusion is, that only one Goroutine was the fastest way to do sorting.

The message is, just because you throw more goroutines at a problem doesn't
necessarily mean it's going to run faster, especially if you're not leveraging
mechanical sympathy.

I don't get it. So I went ahead and I wrote merge sort.

[Sample merge sort program](validate/validate_test.go).

```go
// n contains the data to sort.
var n []int

// Generate the numbers to sort.
func init() {
	for i := 0; i < 1000000; i++ {
		n = append(n, i)
	}
}

// Benchamark sort using one Goroutine
func BenchmarkSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		single(n)
	}
}

// Benchmark sort using unlimited number of Goroutines. Goroutine for every split.
func BenchmarkUnlimited(b *testing.B) {
	for i := 0; i < b.N; i++ {
		unlimited(n)
	}
}

// Benchmark sort using number of Goroutines for the number of CPU cores I have.
func BenchmarkNumCPU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numCPU(n, 0)
	}
}

// single uses a single goroutine to perform the merge sort.
func single(n []int) []int {
	// Once we have a list of one we can begin to merge values.
	if len(n) <= 1 {
		return n
	}

	// Split the list in half.
	i := len(n) / 2

	// Sort the left side.
	l := single(n[:i])

	// Sort the right side.
	r := single(n[i:])

	// Place things in order and merge ordered lists.
	return merge(l, r)
}

// unlimited uses a goroutine for every split to perform the merge sort.
func unlimited(n []int) []int {
	// Once we have a list of one we can begin to merge values.
	if len(n) <= 1 {
		return n
	}

	// Split the list in half.
	i := len(n) / 2

	// Maintain the ordered left and right side lists.
	var l, r []int

	// For each split we will have 2 goroutines.
	var wg sync.WaitGroup
	wg.Add(2)

	// Sort the left side concurrently.
	go func() {
		l = unlimited(n[:i])
		wg.Done()
	}()

	// Sort the right side concurrenyly.
	go func() {
		r = unlimited(n[i:])
		wg.Done()
	}()

	// Wait for the spliting to end.
	wg.Wait()

	// Place things in order and merge ordered lists.
	return merge(l, r)
}

// numCPU uses the same number of goroutines that we have cores
// to perform the merge sort.
func numCPU(n []int, lvl int) []int {
	// Once we have a list of one we can begin to merge values.
	if len(n) <= 1 {
		return n
	}

	// Split the list in half.
	i := len(n) / 2

	// Maintain the ordered left and right side lists.
	var l, r []int

	// Cacluate how many levels deep we can create goroutines.
	// On an 8 core machine we can keep creating goroutines until level 4.
	// 		Lvl 0		1  Lists		1  Goroutine
	//		Lvl 1		2  Lists		2  Goroutines
	//		Lvl 2		4  Lists		4  Goroutines
	//		Lvl 3		8  Lists		8  Goroutines
	//		Lvl 4		16 Lists		16 Goroutines

	// On 8 core machine this will produce the value of 3.
	maxLevel := int(math.Log2(float64(runtime.NumCPU())))

	// We don't need more goroutines then we have logical processors.
	if lvl <= maxLevel {
		lvl++

		// For each split we will have 2 goroutines.
		var wg sync.WaitGroup
		wg.Add(2)

		// Sort the left side concurrently.
		go func() {
			l = numCPU(n[:i], lvl)
			wg.Done()
		}()

		// Sort the right side concurrenyly.
		go func() {
			r = numCPU(n[i:], lvl)
			wg.Done()
		}()

		// Wait for the spliting to end.
		wg.Wait()

		// Place things in order and merge ordered lists.
		return merge(l, r)
	}

	// Sort the left and right side on this goroutine.
	l = numCPU(n[:i], lvl)
	r = numCPU(n[i:], lvl)

	// Place things in order and merge ordered lists.
	return merge(l, r)
}

// merge performs the merging to the two lists in proper order.
func merge(l, r []int) []int {
	// Declare the sorted return list with the proper capacity.
	ret := make([]int, 0, len(l)+len(r))

	// Compare the number of items required.
	for {
		switch {
		case len(l) == 0:
			// We appended everything in the left list so now append
			// everything contained in the right and return.
			return append(ret, r...)

		case len(r) == 0:
			// We appended everything in the right list so now append
			// everything contained in the left and return.
			return append(ret, l...)

		case l[0] <= r[0]:
			// First value in the left list is smaller than the
			// first value in the right so append the left value.
			ret = append(ret, l[0])

			// Slice that first value away.
			l = l[1:]

		default:
			// First value in the right list is smaller than the
			// first value in the left so append the right value.
			ret = append(ret, r[0])

			// Slice that first value away.
			r = r[1:]
		}
	}
}
```

I reconstructed everything that the author did and I'm gonna go ahead and run
this benchmark just like the author did in the article.

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/benchmarks/validate
$ go test -run none -bench . -benchtime 3s
```

Sample output:

```sh
goos: linux
goarch: amd64
pkg: github.com/cedrickchee/ultimate-go/testing/benchmarks/validate
BenchmarkSingle-4             26         133750947 ns/op
BenchmarkUnlimited-4           2        1593642608 ns/op
BenchmarkNumCPU-4             14         220086627 ns/op
PASS
ok      github.com/cedrickchee/ultimate-go/testing/benchmarks/validate  12.405s
```

Single Goroutine: it took 133 milliseconds for 26 times.
Unlimited Goroutine: it took 1.59 seconds for 2 times.
4 Goroutine: it took 220 milliseconds for 14 times.

This was the conclusion that the author had, as well, in the post.

I still don't believe it. You must validate your results. Then, I do the
following:

Isolate running NumCPU by itself.

```sh
BenchmarkNumCPU-4             24         161248624 ns/op
PASS
ok      github.com/cedrickchee/ultimate-go/testing/benchmarks/validate  4.054s
```

NumCPU is now faster. 59 milliseconds, significantly faster than a single
Goroutine, which is exactly what I expected to see. What happened here?
Machine must be idle when you're running these benchmarks. When we started
running NumCPU in this batch, the runtime in the machine was still busy cleaning
up the mess we made with the unlimited test.

I like the idea that throwing Goroutines at a problem doesn't mean it makes it
necessarily run any faster. But, when we can run things in parallel efficiently
with mechanical sympathies, it really should be faster than just using a single
Goroutine on one core and now we've really proven that.

## Advanced Performance

**Package Review**

- [Prediction](prediction/README.md)
- [Caching](caching/README.md)
- [False Sharing](falseshare/README.md)
