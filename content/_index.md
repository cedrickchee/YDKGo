---
title: Introduction
type: docs
bookToc: false
---

![Gopher](/img/gopher.png#center)

---

# What _Is_ Go?

You don't know Go, yet. Neither do I, not fully anyway. None of us do. But we can all start getting to know Go better.

# About the Book

Go programming language study notes turned book.

This book is inspired by ['You Don't Know JS Yet' (YDKJS)](https://github.com/getify/You-Dont-Know-JS) book series. YDKJS helped me understand JavaScript under the hood, after more than 8 years writing software in JS.

You Don't Know Go (YDKGo) book is based on [Ultimate Go training](https://www.ardanlabs.com/ultimate-go/), which is an intermediate-level class for engineers with some experience with Go trying to dig deeper into the language.

---

# Table of Contents

## [Lesson 1: Design Guidelines](/docs/guidelines) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/README.md#design-guidelines)]

- Philosophy
- [Prepare your mind](/docs/guidelines/_index.md#prepare-your-mind)
- [Productivity vs performance](/docs/guidelines/_index.md#productivity-versus-performance)
- [Correctness vs Performance](/docs/guidelines/_index.md#correctness-versus-performance)
- [Code reviews](/docs/guidelines/_index.md#code-reviews)

## [Language Mechanics](/docs/language)

<!--- Topics [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/courses/go/language/README.md)]-->

### Lesson 2: Language Syntax

- [Variables](/docs/language/variables) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/variables/README.md) | [code example](docs/language/variables/example1/example1.go) | [exercise 1 solution](/docs/language/variables/exercise1/exercise1.go)]
  - Built-in types
  - Zero value concept
  - Declare and initialize variables
  - Conversion vs casting
- [Struct Types](/docs/language/struct_types) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/struct_types/README.md) | [exercise 1 solution](/docs/language/struct_types/exercise1/exercise1.go)]
  - Declare, create and initialize struct types [[code example](/docs/language/struct_types/example1/example1.go)]
  - Anonymous struct types [[code example](/docs/language/struct_types/example2/example2.go)]
  - Named vs Unnamed types [[code example](/docs/language/struct_types/example3/example3.go)]
- [Pointers](/docs/language/pointers) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/pointers/README.md) | [exercise 1 solution](/docs/language/pointers/exercise1/exercise1.go) | [exercise 2 solution](/docs/language/pointers/exercise2/exercise2.go)]
  - Part 1 (Pass by Value) [[code example](/docs/language/pointers/example1/example1.go)]
  - Part 2 (Sharing Data) [[code example](/docs/language/pointers/example2/example2.go) | [code example](/docs/language/pointers/example3/example3.go)]
  - Part 3 (Escape Analysis) [[code example](/docs/language/pointers/example4/example4.go)]
  - Part 4 (Stack Growth) [[code example](/docs/language/pointers/example5/example5.go)]
  - Part 5 (Garbage Collection)
- [Constants](/docs/language/constants) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/constants/README.md) | [exercise 1 solution](/docs/language/constants/exercise1/exercise1.go)]
  - Declare and initialize constants [[code example](/docs/language/constants/example1/example1.go)]
  - Parallel type system (Kind) [[code example](/docs/language/constants/example2/example2.go)]
  - iota  [[code example](/docs/language/constants/example3/example3.go)]
  - Implicit conversion [[code example](/docs/language/constants/example4/example4.go)]
- [Functions](/docs/language/functions) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/functions/README.md) | [exercise 1 solution](/docs/language/functions/exercise1/exercise1.go)]
  - Return multiple values [[code example](/docs/language/functions/example1/example1.go)]
  - Blank identifier [[code example](/docs/language/functions/example2/example2.go)]
  - Redeclarations [[code example](/docs/language/functions/example3/example3.go)]
  - Anonymous Functions/Closures [[code example](/docs/language/functions/example4/example4.go)]
  - Advanced code review
    - Recover panics [[code example](/docs/language/functions/advanced/example1/example1.go)]

### Lesson 3: Data Structures

- [Data-Oriented Design](/docs/language/arrays/data_oriented_design.md)
  - Design guidelines [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/#data-oriented-design)]
- [Arrays](/docs/language/arrays/arrays.md) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/arrays/README.md) | [exercise 1 solution](/docs/language/arrays/exercise1/exercise1.go)]
  - Part 1 (Mechanical Sympathy)
  - Part 2 (Semantics)
    - Declare, initialize and iterate [[code example](/docs/language/arrays/example1/example1.go)]
    - Different type arrays [[code example](/docs/language/arrays/example2/example2.go)]
    - Contiguous memory allocations [[code example](/docs/language/arrays/example3/example3.go)]
    - Range mechanics [[code example](/docs/language/arrays/example4/example4.go)]
- [Slices](/docs/language/slices) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/slices) | [exercise 1 solution](/docs/language/slices/exercise1/exercise1.go)]
  - Part 1
    - Declare and Length [[code example](/docs/language/slices/example1/example1.go)]
    - Reference Types [[code example](/docs/language/slices/example2/example2.go)]
  - Part 2 (Appending Slices) [[code example](/docs/language/slices/example4/example4.go)]
  - Part 3 (Taking Slices of Slices) [[code example](/docs/language/slices/example3/example3.go)]
  - Part 4 (Slices and References) [[code example](/docs/language/slices/example5/example5.go)]
  - Part 5 (Strings and Slices) [[code example](/docs/language/slices/example6/example6.go)]
  - Part 6 (Range Mechanics) [[code example](/docs/language/slices/example8/example8.go)]
  - Part 7 (Variadic Functions) [[code example](/docs/language/slices/example7/example7.go)]
- [Maps](/docs/language/maps) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/maps/README.md) | [exercise 1 solution](/docs/language/maps/exercise1/exercise1.go)]
  - Declare, write, read, and delete [[code example](/docs/language/maps/example1/example1.go)]
  - Absent keys [[code example](/docs/language/maps/example2/example2.go)]
  - Map key restrictions [[code example](/docs/language/maps/example3/example3.go)]
  - Map literals and range [[code example](/docs/language/maps/example4/example4.go)]
  - Sorting maps by key [[code example](/docs/language/maps/example5/example5.go)]
  - Taking an element's address [[code example](/docs/language/maps/example6/example6.go)]
  - Maps are Reference Types [[code example](/docs/language/maps/example7/example7.go)]

### Lesson 4: Decoupling

- [Methods](/docs/language/methods) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/methods/README.md) | [exercise 1 solution](/docs/language/methods/exercise1/exercise1.go)]
  - Part 1 (Declare and Receiver Behavior) [[code example](/docs/language/methods/example1/example1.go)]
  - Part 2 (Value and Pointer Semantics) [[code example](/docs/language/methods/example5/example5.go)]
  - Part 3 (Function/Method Variables) [[code example](/docs/language/methods/example3/example3.go)]
  - Part 4 (Named Typed Methods) [[code example](/docs/language/methods/example2/example2.go)]
  - Part 5 (Function Types) [[code example](/docs/language/methods/example4/example4.go)]
- [Interfaces](/docs/language/interfaces) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/interfaces/README.md) | [exercise 1 solution](/docs/language/interfaces/exercise1/exercise1.go)]
  - Part 1 (Polymorphism) [[code example](/docs/language/interfaces/example1/example1.go)]
  - Part 2
    - Method Sets [[code example](/docs/language/interfaces/example2/example2.go)]
    - Address of Value [[code example](/docs/language/interfaces/example3/example3.go)]
  - Part 3 (Storage By Value) [[code example](/docs/language/interfaces/example4/example4.go)]
<!--  - Part 4 (Repetitive Code That Needs Polymorphism)  [[code example](/docs/language/interfaces/example0/example0.go)]
  - Part 5 (Type Assertions) [[code example](/docs/language/interfaces/example5/example5.go)]
  - Part 6 (Conditional Type Assertions) [[code example](/docs/language/interfaces/example6/example6.go)]
  - Part 7 (The Empty Interface and Type Switches) [[code example](/docs/language/interfaces/example7/example7.go)] -->
- [Embedding](/docs/language/embedding) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/embedding/README.md) | [exercise 1 solution](/docs/language/embedding/exercise1/exercise1.go)]
  - Declaring Fields [[code example](/docs/language/embedding/example1/example1.go)]
  - Embedding types [[code example](/docs/language/embedding/example2/example2.go)]
  - Embedded types and interfaces [[code example](/docs/language/embedding/example3/example3.go)]
  - Outer and inner type interface implementations [[code example](/docs/language/embedding/example4/example4.go)]
- [Exporting](/docs/language/exporting) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/language/exporting/README.md) | [exercise 1 solution](/docs/language/exporting/exercise1/exercise1.go)]
  - Declare and access exported identifiers - Pkg [[code example](/docs/language/exporting/example1/counters/counters.go)]
  - Declare and access exported identifiers - Main [[code example](/docs/language/exporting/example1/example1.go)]
  - Declare unexported identifiers and restrictions - Pkg [[code example](/docs/language/exporting/example2/counters/counters.go)]
  - Declare unexported identifiers and restrictions - Main [[code example](/docs/language/exporting/example2/example2.go)]
  - Access values of unexported identifiers - Pkg [[code example](/docs/language/exporting/example3/counters/counters.go)]
  - Access values of unexported identifiers - Main [[code example](/docs/language/exporting/example3/example3.go)]
  - Unexported struct type fields - Pkg [[code example](/docs/language/exporting/example4/users/users.go)]
  - Unexported struct type fields - Main [[code example](/docs/language/exporting/example4/example4.go)]
  - Unexported embedded types - Pkg [[code example](/docs/language/exporting/example5/users/users.go)]
  - Unexported embedded types - Main [[code example](/docs/language/exporting/example5/example5.go)]

## [Software Design](/docs/design) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/courses/go/design/README.md)]

### [Lesson 5: Composition](/docs/design/composition) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/design/composition/README.md) | [exercise 1 solution](/docs/design/composition/exercises/exercise1/exercise1.go)]
- Design Guidelines [[guide](https://github.com/ardanlabs/gotraining/tree/master/topics/go#interface-and-composition-design)]
- Grouping Types
  - Grouping By State [[code example](/docs/design/composition/grouping/example1/example1.go)]
  - Grouping By Behavior [[code example](/docs/design/composition/grouping/example2/example2.go)]
- Decoupling
  - Struct Composition [[code example](/docs/design/composition/decoupling/example1/example1.go)]
  - Decoupling With Interface [[code example](/docs/design/composition/decoupling/example2/example2.go)]
  - Interface Composition [[code example](/docs/design/composition/decoupling/example3/example3.go)]
  - Decoupling With Interface Composition [[code example](/docs/design/composition/decoupling/example4/example4.go)]
  - Remove Interface Pollution [[code example](/docs/design/composition/decoupling/example5/example5.go)]
  - More Precise API [[code example](/docs/design/composition/decoupling/example6/example6.go)]
- Conversion and Assertions
  - Interface Conversions [[code example](/docs/design/composition/assertions/example1/example1.go)]
  - Runtime Type Assertions [[code example](/docs/design/composition/assertions/example2/example2.go)]
  - Behavior Changes [[code example](/docs/design/composition/assertions/example3/example3.go)]
- Interface Pollution
  - Create Interface Pollution [[code example](/docs/design/composition/pollution/example1/example1.go)]
  - Remove Interface Pollution [[code example](/docs/design/composition/pollution/example2/example2.go)]
- Mocking
  - Package To Mock [[code example](/docs/design/composition/mocking/example1/pubsub/pubsub.go)]
  - Client [[code example](/docs/design/composition/mocking/example1/example1.go)]

### [Lesson 6: Error Handling](/docs/design/error_handling) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/design/error_handling/README.md) | exercise solution: [1](/docs/design/error_handling/exercise1/exercise1.go) | [2](/docs/design/error_handling/exercise2/exercise2.go)]
- Default Error Values [[code example](/docs/design/error_handling/example1/example1.go)]
- Error Variables [[code example](/docs/design/error_handling/example2/example2.go)]
- Type As Context [[code example](/docs/design/error_handling/example3/example3.go)]
- Behavior As Context [[code example](/docs/design/error_handling/example4/example4.go)]
- Find The Bug [[code example](/docs/design/error_handling/example5/example5.go) | [the reason](/docs/design/error_handling/example5/reason/reason.go)]
- Wrapping Errors [[code example](/docs/design/error_handling/example6/example6.go)]

### [Lesson 7: Packaging](/docs/design/packaging)
- Language Mechanics [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/design/packaging/README.md#language-mechanics)]
- Design Guidelines [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/design/packaging/README.md#design-philosophy)]
- Package-Oriented Design [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/design/packaging/README.md#package-oriented-design)]

## [Concurrency](/docs/concurrency) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/courses/go/concurrency/README.md)]

### [Lesson 8: Mechanics - Goroutines](/docs/concurrency/goroutines) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/concurrency/goroutines/README.md) | [exercise 1 solution](/docs/concurrency/goroutines/exercise1/exercise1.go)]
- Scheduling in Go
  - Part 1 (OS Scheduler) [[article](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part1.html)]
  - Part 2 (Go Scheduler) [[article](https://www.ardanlabs.com/blog/2018/08/scheduling-in-go-part2.html)]
  - Part 3 (Concurrency) [[article](https://www.ardanlabs.com/blog/2018/12/scheduling-in-go-part3.html)]
- Language Mechanics [[code example](/docs/concurrency/goroutines/example1/example1.go)]
- Goroutine Time Slicing [[code example](/docs/concurrency/goroutines/example2/example2.go)]
- Goroutine and Parallelism [[code example](/docs/concurrency/goroutines/example3/example3.go)]

### [Lesson 9: Mechanics - Data Races](/docs/concurrency/data_race) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/concurrency/data_race/README.md) | [exercise 1 solution](/docs/concurrency/data_race/exercise1/exercise1.go)]
- Data Race and Race Detection [[code example](/docs/concurrency/data_race/example1/example1.go)]
- Synchronization with Atomic Functions [[code example](/docs/concurrency/data_race/example2/example2.go)]
- Synchronization with Mutexes [[code example](/docs/concurrency/data_race/example3/example3.go)]
- Read/Write Mutex [[code example](/docs/concurrency/data_race/example4/example4.go)]
- Map Data Race [[code example](/docs/concurrency/data_race/example5/example5.go)]
- Interface-Based Race Condition [[code example](/docs/concurrency/data_race/advanced/example1/example1.go)]

### [Lesson 10: Mechanics - Channels](/docs/concurrency/channels) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/concurrency/channels/README.md) | exercise solution: [1](/docs/concurrency/channels/exercises/exercise1/exercise1.go) | [2](/docs/concurrency/channels/exercises/exercise2/exercise2.go) | [3](/docs/concurrency/channels/exercises/exercise3/exercise3.go) | [4](/docs/concurrency/channels/exercises/exercise4/exercise4.go)]
- Design Guidelines [[guide](https://github.com/ardanlabs/gotraining/tree/master/topics/go#channel-design)]
- Signaling Semantics
  - Language Mechanics
- Basic Patterns
  - Part 1 (Wait for Task) [[code example](/docs/concurrency/channels/example1/example1.go)]
  - Part 2 (Wait for Result) [[code example](/docs/concurrency/channels/example1/example1.go)]
  - Part 3 (Wait for Finished) [[code example](/docs/concurrency/channels/example1/example1.go)]
- Pooling Pattern [[code example](/docs/concurrency/channels/example1/example1.go)]
- Fan Out Pattern
  - Part 1 [[code example](/docs/concurrency/channels/example1/example1.go)]
  - Part 2 [[code example](/docs/concurrency/channels/example1/example1.go)]
- Drop Pattern [[code example](/docs/concurrency/channels/example1/example1.go)]
- Cancellation Pattern [[code example](/docs/concurrency/channels/example1/example1.go)]

### Lesson 11: Concurreny Patterns

- [Context](/docs/concurrency/context) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/packages/context/README.md) | [exercise 1 solution](/docs/concurrency/context/exercise1/exercise1.go)]
  - Store / Retrieve context values [[code example](/docs/concurrency/context/example1/example1.go)]
  - WithTimeout [[code example](/docs/concurrency/context/example4/example4.go)]
  - Request/Response Context Timeout [[code example](/docs/concurrency/context/example5/example5.go)]
  - WithCancel [[code example](/docs/concurrency/context/example2/example2.go)]
  - WithDeadline [[code example](/docs/concurrency/context/example3/example3.go)]
- Failure Detection [[code example](/docs/concurrency/patterns/advanced/main.go)]

## Testing and Profiling [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/courses/go/tooling/README.md)]

### [Lesson 12: Testing](/docs/testing/tests) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/testing/tests/README.md)]
- Basic Unit Testing [[code example](/docs/testing/tests/example1/example1_test.go)]
- Table Unit Testing [[code example](/docs/testing/tests/example2/example2_test.go)]
- Mocking Server [[code example](/docs/testing/tests/example3/example3_test.go)]
- Testing Internal Endpoints [[code example](/docs/testing/tests/example4/handlers/handlers_test.go)]
- Sub Tests [[code example](/docs/testing/tests/example5/example5_test.go)]
- Code Coverage

### [Lesson 13: Benchmarking](/docs/testing/benchmarks) [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/testing/benchmarks/README.md) | [exercise 1 solution](/docs/testing/benchmarks/exercises/exercise1/bench_test.go)]
- Basic Benchmarking [[code example](/docs/testing/benchmarks/basic/basic_test.go)]
- Sub Benchmarks [[code example](/docs/testing/benchmarks/sub/sub_test.go)]
- Validate Benchmarks [[code example](/docs/testing/benchmarks/validate/validate_test.go)]

### [Lesson 14: Profiling and Tracing](/docs/profiling)
- Profiling Guidelines [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/profiling/README.md)]
- Stack Traces [code example: [1](/docs/profiling/stack_trace/example1/example1.go) | [2](/docs/profiling/stack_trace/example2/example2.go) | [3](/docs/profiling/stack_trace/example3/example3.go)]
- Micro Level Optimization using Benchmarks [[code example](/docs/profiling/memcpu/stream_test.go)]
- Macro Level Optimization
  - Part 1: GODEBUG Tracing [code example: [1](/docs/profiling/project/main.go) | [2](/docs/profiling/godebug/godebug.go)]
  - Part 2: Memory Profiling [code example: [1](/docs/profiling/project/main.go) | [2](/docs/profiling/pprof/pprof.go)]
  - Part 3: Tooling Changes
  - Part 4: CPU Profiling [[code example](/docs/profiling/project/main.go)]
- Execution Tracing [[code example](/docs/profiling/trace/trace.go)]

## Extra Lesson

- Advanced Performance and Benchmarking [[guide](/docs/testing/benchmarks/_index.md#advanced-performance)]
- Fuzzing [[guide](https://github.com/ardanlabs/gotraining/blob/master/topics/go/testing/fuzzing/README.md)]

---

#### License

This book contains a variety of content; some developed by Cedric Chee, and some from third-parties. The third-party content is distributed under the license provided by those parties.

The content of this project itself is licensed under the [Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License](http://creativecommons.org/licenses/by-nc-sa/4.0/), and the underlying source code used to format and display that content is licensed under the [Apache License, Version 2.0](LICENSE).

The [Go gopher](http://blog.golang.org/gopher) was designed by [Renee French](http://reneefrench.blogspot.com/), is licensed under Creative Commons 3.0 Attributions.

[Gopher picture](https://github.com/MariaLetta/free-gophers-pack) by Maria Letta.
