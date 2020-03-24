---
title: Testing
weight: 1
---

# Testing

## Basic Unit Testing

Unit tests are integrated into the Go programming language and into the tooling.
There are third-party packages out there like GoConvey, but I want to avoid
using third-party stuff, just because it adds extra dependencies and things to
your code.

> A unit test is a test of behavior whose success or failure is wholly
> determined by the correctness of the test and the correctness of the unit
> under test. — Kevin Henney

So, one thing we've to validate and make sure is that our tests are adding value,
and also that we're not over unit testing or under unit testing.

If you want to write tests, then your files have to be in the
format `<filename>_test.go`. Otherwise, the testing tool will not find the tests.

These files aren't going to be built and compiled into your final binaries.

Test files should be in the same package as your code. We might also want to
have a folder called test for more than unit test, say integration test.

The package name can be the name only or name_test.
If we go with name_test, it allows us to make sure these tests work with the
package. The only reason that we don't want to do this is when we have a
function or method that is unexported. However, if we don't use name_test, it
will raise a red flag because if we cannot test the exported API to get the
coverage for unexported API then we know are missing something. Thus, 9 out of
10, this is what we want.

My testing philosophies: always make sure that those unexported APIs are very
testable. Make sure the exported APIs are very usable. So, it's a balance.

[Sample test](example1/example1_test.go).

```go
package example1

import (
	"net/http"
	"testing"
)
```

Import Go testing package.

```go
const succeed = "\u2713"
const failed = "\u2717"
```

These constants, which are the Unicode characters for x and check. They gives us
visualization.

```go
// TestDownload validates the http Get function can download content.
func TestDownload(t *testing.T) {
	url := "https://www.ardanlabs.com/blog/index.xml"
	statusCode := 200

	t.Log("Given the need to test downloading content.")
	{
		t.Logf("\tTest 0:\tWhen checking %q for status code %d", url, statusCode)
		{
			// Make sure this test code is really as close to how you would
			// write production code.
			resp, err := http.Get(url)
			// Most importantly, if a function fails, like there's an error, you
			// need to check your error here. We don't want to see blank
			// identifier being used in your test code.
			if err != nil {
				t.Fatalf("\t%s\tShould be able to make the Get call : %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to make the Get call.", succeed)

			defer resp.Body.Close()

			if resp.StatusCode == statusCode {
				t.Logf("\t%s\tShould receive a %d status code.", succeed, statusCode)
			} else {
				t.Errorf("\t%s\tShould receive a %d status code : %d", failed, statusCode, resp.StatusCode)
			}
		}
	}
}
```

`TestBasic` validates the `http.Get` function can download content.
Every test will be associated with test function. It starts with the word `Test`
and the first word after `Test` must be capitalized. This function is exported,
or the testing tool will not find it. It uses a `testing.T` pointer as its
parameter.

When writing test, we want to focus on usability first. We must write it the
same way as we would write it in production. We also want the verbosity of
tests. I like the verbosity of tests because when they run in the CI, I'd like
to have a lot of information. I'm really about not hiding information.

So we have 3 different methods of `t`: `Log` or `Logf`, `Fatal` or `Fatalf`,
`Error` or `Error f`. That is the core APIs for testing.

- `Log`: Write documentation out into the log output.
- `Error`: Write documentation and also say that this test is failed but we are
continue moving forward to execute code in this test.
- `Fatal`: Similarly, document that this test is failed but we are done. We move
on to the next test function.

The whole idea of a unit test is to validate that something is working.
Normally, we will validate that some API that you've written, whether it's
exported or unexported, is behaving, that if I give it this input, it produces
this sort of output.

The idea of Given, When, Should:

- `Given`: Why are we writing this test?
- `When`: What data are we using for this test?
- `Should`: What should happen or not happen?

We are also using the artificial block `{ }` between a long `Log` function.
They help with readability.

`go test` command: the testing tool will find that function and the test.

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example1
$ go test
PASS
ok      github.com/cedrickchee/ultimate-go/testing/tests/example1       1.143s
```

We can also say `go test -v` for verbosity, we will get a full output of the
logging.

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example1
$ go test -v
=== RUN   TestDownload
--- PASS: TestDownload (2.10s)
    example1_test.go:22: Given the need to test downloading content.
    example1_test.go:24:        Test 0: When checking "https://www.ardanlabs.com/blog/index.xml" for status code 200
    example1_test.go:30:        ✓       Should be able to make the Get call.
    example1_test.go:35:        ✓       Should receive a 200 status code.
PASS
ok      github.com/cedrickchee/ultimate-go/testing/tests/example1       2.100s
```

Suppose that we have a lot of test functions, `go test -run TestBasic` will
only run the `TestBasic` function. The `run` flag uses a regular expression to
be able to filter the test functions that you just want to run.

## Table Unit Testing

Table tests are a way of writing unit test when you have a base code and you
want to throw a lot of different input and expected output around it.

It set up a data structure of input to expected output. This way we don't need
a separate function for each one of these. We just have 1 test function. As we
go along, we just add more to the table.

[Sample test](example2/example2_test.go).

```go
tests := []struct {
	url        string
	statusCode int
}{
	{"https://www.ardanlabs.com/blog/index.xml", http.StatusOK},
	{"http://rss.cnn.com/rss/cnn_topstorie.rss", http.StatusNotFound},
}
```

This table is a slice of anonymous struct type. It is the URL we will call and
`statusCode` are what we expect.

```go
t.Log("Given the need to test downloading different content.")
{
	for i, tt := range tests {
		t.Logf("\tTest: %d\tWhen checking %q for status code %d", i, tt.url, tt.statusCode)
		{
			resp, err := http.Get(tt.url)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to make the Get call : %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to make the Get call.", succeed)

			defer resp.Body.Close()

			if resp.StatusCode == tt.statusCode {
				t.Logf("\t%s\tShould receive a %d status code.", succeed, tt.statusCode)
			} else {
				t.Errorf("\t%s\tShould receive a %d status code : %v", failed, tt.statusCode, resp.StatusCode)
			}
		}
	}
}
```

Let's see the output:

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example2
$ go test -v
=== RUN   TestDownload
--- PASS: TestDownload (7.63s)
    example2_test.go:28: Given the need to test downloading different content.
    example2_test.go:31:        Test: 0 When checking "https://www.ardanlabs.com/blog/index.xml" for status code 200
    example2_test.go:37:        ✓       Should be able to make the Get call.
    example2_test.go:42:        ✓       Should receive a 200 status code.
    example2_test.go:31:        Test: 1 When checking "http://rss.cnn.com/rss/cnn_topstorie.rss" for status code 404
    example2_test.go:37:        ✓       Should be able to make the Get call.
    example2_test.go:42:        ✓       Should receive a 404 status code.
PASS
ok      github.com/cedrickchee/ultimate-go/testing/tests/example2       7.634s
```

## Mocking Server

The last two tests we created, they were fantastic. But the problem with those 2
tests. They require a connection to the outside world. We're hitting a live
sever. We cannot assume that we always have access to the resources we need.

Therefore, mocking becomes an important part of testing in many cases.
I always want to be very careful about mocking, I don't want to mock databases,
I want to use Docker for that. Certain systems, where it's critical that we know
that we're talking to live systems or binary protocols. But a lot of things can
be mocked, and these tests that we did run, could be mocked.

I know that web servers work. So I don't care about the HTTP protocols, here.
What I care about is that if I send this request and I get back the response,
that I can process it. I will show you how you can mock a HTTP GET call, already
built into Go's standard library and the language.

[Sample test](example3/example3_test.go)

```go
import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)
```

Go gives us this package called httptest. We're going to use this to mock these
GET calls internally.

```go
const succeed = "\u2713"
const failed = "\u2717"

// feed is mocking the XML document we expect to receive.
// Notice that we are using ` instead of " so we can reserve special characters.
var feed = `<?xml version="1.0" encoding="UTF-8"?>
<rss>
<channel>
    <title>Going Go Programming</title>
    <description>Golang : https://github.com/goinggo</description>
    <link>http://www.goinggo.net/</link>
    <item>
        <pubDate>Sun, 15 Mar 2015 15:04:00 +0000</pubDate>
        <title>Object Oriented Programming Mechanics</title>
        <description>Go is an object oriented language.</description>
        <link>http://www.goinggo.net/2015/03/object-oriented</link>
    </item>
</channel>
</rss>`

// Item defines the fields associated with the item tag in the buoy RSS document.
type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
}

// Channel defines the fields associated with the channel tag in the buoy RSS
// document.
type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Items       []Item   `xml:"item"`
}

// Document defines the fields associated with the buoy RSS document.
type Document struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
	URI     string
}
```

I've created some structs here, and with the XML tag, so we can pull this data
when we mock the call and get this raw string back. We can unmarshal it into our
real data structures here and do some actual testing.

```go
func mockServer() *httptest.Server {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprintln(w, feed)
	}

	return httptest.NewServer(http.HandlerFunc(f))
}
```

`mockServer` returns a pointer of type `httptest.Server` to handle the mock get
call. This mock function calls `NewServer` function that will stand up a web
server for us automatically. All we have to give `NewServer` is a function of
the `Handler` type, which is `f`.

`f` creates an anonymous function with the signature of `ResponseWriter` and
`Request`. This is the core signature of everything related to http in Go.
`ResponseWriter` is an interface that allows us to write the response out.
Normally when we get this interface value, there is already a concrete type
value stored inside of it that support what we are doing. `Request` is a
concrete type that we will get with the request.

This is how it will work.
We will get a mock server started by making `NewServer` call. When the request
comes into it, execute `f`. Therefore, `f` is doing the entire mock.
We will send `200` down the line, set the header to XML and use `Fprintln` to
take the `ResponseWriter` interface value and feeding with the raw string we
defined above.

```go
// TestDownload validates the http Get function can download content and
// the content can be unmarshaled and clean.
func TestDownload(t *testing.T) {
	statusCode := http.StatusOK

	// Call the mock sever and defer close to shut it down cleanly.
	server := mockServer()
	defer server.Close()

	// Now, it's just the matter of using server value to know what URL we need
	// to use to run this mock. From the http.Get point of view, it is making an
	// URL call. It has no idea that it's hitting the mock server. We have
	// mocked out a perfect response.
	t.Log("Given the need to test downloading content.")
	{
		t.Logf("\tTest 0:\tWhen checking %q for status code %d", server.URL, statusCode)
		{
			resp, err := http.Get(server.URL)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to make the Get call : %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to make the Get call.", succeed)

			defer resp.Body.Close()

			if resp.StatusCode != statusCode {
				t.Fatalf("\t%s\tShould receive a %d status code : %v", failed, statusCode, resp.StatusCode)
			}
			t.Logf("\t%s\tShould receive a %d status code.", succeed, statusCode)

			// When we get the response back, we are unmarshaling it from XML to
			// our struct type and do some extra validation with that as we go.
			var d Document
			if err := xml.NewDecoder(resp.Body).Decode(&d); err != nil {
				t.Fatalf("\t%s\tShould be able to unmarshal the response : %v", failed, err)
			}
			t.Logf("\t%s\tShould be able to unmarshal the response.", succeed)

			if len(d.Channel.Items) == 1 {
				t.Logf("\t%s\tShould have 1 item in the feed.", succeed)
			} else {
				t.Errorf("\t%s\tShould have 1 item in the feed : %d", failed, len(d.Channel.Items))
			}
		}
	}
}
```

Run test:

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example3
$ go test -v
=== RUN   TestDownload
--- PASS: TestDownload (0.00s)
    example3_test.go:82: Given the need to test downloading content.
    example3_test.go:84:        Test 0: When checking "http://127.0.0.1:33826" for status code 200
    example3_test.go:90:        ✓       Should be able to make the Get call.
    example3_test.go:97:        ✓       Should receive a 200 status code.
    example3_test.go:103:       ✓       Should be able to unmarshal the response.
    example3_test.go:106:       ✓       Should have 1 item in the feed.
PASS
ok      github.com/cedrickchee/ultimate-go/testing/tests/example3       0.004s
```

I want you to see how fast this test ran. It runs much faster because we know
that it didn't leave my machine. We also see that we have the localhost IP
address on a port.

## Testing Internal Endpoints

A lot of us are building web APIs using the HTTPS protocols, and we would like
to test it as well without manually having to stand up our own server.

The Go standard library HTTP test package supports this.

Below is a [basic code to build a non-production level web API](example4/example4.go).

```go
package main

import (
	"log"
	"net/http"

	// Import handler package that has a set of routes that we will work with.
	"github.com/ardanlabs/gotraining/topics/go/testing/tests/example4/handlers"
)

func main() {
	handlers.Routes()

	log.Println("listener : Started : Listening on: http://localhost:4000")
	http.ListenAndServe(":4000", nil)
}
```

[handlers code](example4/handlers/handlers.go).

```go
package handlers

import (
	"encoding/json"
	"net/http"
)

// Routes sets the routes for the web service.
func Routes() {
	http.HandleFunc("/sendjson", SendJSON)
}
```

`Routes` has 1 route call sendjson. When that route is executed, it will call
the `SendJSON` function.

```go
// SendJSON returns a simple JSON document.
func SendJSON(rw http.ResponseWriter, r *http.Request) {
	u := struct {
		Name  string
		Email string
	}{
		Name:  "Bill",
		Email: "bill@ardanlabs.com",
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(200)
	json.NewEncoder(rw).Encode(&u)
}
```

`SendJSON` has the same signature that we had before using `ResponseWriter` and
`Request`. That's how we write handlers in Go. Then we define our literal type,
create an anonymous struct, initialize it and unmarshall it into JSON and send
it down the wire.

Build and run it:

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example4
$ go build

$ ./example4
2020/03/19 15:17:01 listener : Started : Listening on: http://localhost:4000
```

cURL that URL:

```sh
$ curl -X POST http://localhost:4000/sendjson
{"Name":"Bill","Email":"bill@ardanlabs.com"}
```

You've just hit that endpoint.

### Internal test

Below is how to test the execution of an internal endpoint without having to
stand up the server.

[Example test for handlers](example4/handlers/handlers_test.go).

```go
package handlers_test
```

We are using `handlers_test` for package name because we want to make sure we
only touch the exported API.

```go
func init() {
	handlers.Routes()
}
```

The `init` function is really important. It binds your routes to HTTP handler
function. Not only your app can bind them, but your tests can bind them as well.
If we forget to do this then nothing will work.

```go
// TestSendJSON testing the sendjson internal endpoint.
func TestSendJSON(t *testing.T) {
	url := "/sendjson"
	statusCode := 200

	t.Log("Given the need to test the SendJSON endpoint.")
	{
		// Create a nil request body GET for the URL.
		r := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()
		http.DefaultServeMux.ServeHTTP(w, r)

		t.Logf("\tTest 0:\tWhen checking %q for status code %d", url, statusCode)
		{
			if w.Code != 200 {
				t.Fatalf("\t%s\tShould receive a status code of %d for the response. Received[%d].", failed, statusCode, w.Code)
			}
			t.Logf("\t%s\tShould receive a status code of %d for the response.", succeed, statusCode)

			// If we got the 200, we try to unmarshal and validate it.
			var u struct {
				Name  string
				Email string
			}

			if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
				t.Fatalf("\t%s\tShould be able to decode the response.", failed)
			}
			t.Logf("\t%s\tShould be able to decode the response.", succeed)

			if u.Name == "Bill" {
				t.Logf("\t%s\tShould have \"Bill\" for Name in the response.", succeed)
			} else {
				t.Errorf("\t%s\tShould have \"Bill\" for Name in the response : %q", failed, u.Name)
			}

			if u.Email == "bill@ardanlabs.com" {
				t.Logf("\t%s\tShould have \"bill@ardanlabs.com\" for Email in the response.", succeed)
			} else {
				t.Errorf("\t%s\tShould have \"bill@ardanlabs.com\" for Email in the response : %q", failed, u.Email)
			}
		}
	}
}
```

In order to mock this call, we don't need the network. What we need to do is
create a request and run it through the Mux so we will bypass the network call
together, run the request directly through the Mux to test the route and the
handler.

`w := httptest.NewRecorder()` gives us a pointer to its concrete type called
`ResponseRecorder` that already implemented the `ResponseWriter` interface.

`http.DefaultServeMux.ServeHTTP(w, r)`  asks for a `ResponseWriter` and a
`Request`. This call will perform the Mux and call that handler to test it
without network. When his call comes back, the recorder value `w` has the result
of the entire execution. Now we can use that to validate.

Run test:

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example4/handlers
$ go test -run TestSendJSON -v
=== RUN   TestSendJSON
--- PASS: TestSendJSON (0.00s)
    handlers_test.go:32: Given the need to test the SendJSON endpoint.
    handlers_test.go:38:        Test 0: When checking "/sendjson" for status code 200
    handlers_test.go:43:        ✓       Should receive a status code of 200 for the response.
    handlers_test.go:53:        ✓       Should be able to decode the response.
    handlers_test.go:56:        ✓       Should have "Bill" for Name in the response.
    handlers_test.go:62:        ✓       Should have "bill@ardanlabs.com" for Email in the response.
PASS
ok      github.com/cedrickchee/ultimate-go/testing/tests/example4/handlers      0.003s
```

### Example test

This is another type of test in Go. Examples are both documentations and tests.

[handlers_example_test.go](example4/handlers/handlers_example_test.go).

```go
package handlers_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

// ExampleSendJSON provides a basic example example.
func ExampleSendJSON() {
	r := httptest.NewRequest("GET", "/sendjson", nil)
	w := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(w, r)

	var u struct {
		Name  string
		Email string
	}

	if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
		log.Println("ERROR:", err)
	}

	fmt.Println(u)
	// Output:
	// {Bill bill@ardanlabs.com}
}
```

Naming convention `ExampleSendJSON`: Notice that we are binding that `Example`
to our `SendJSON` function.

### Godoc

If we execute `godoc -http :3000`, Go will generate for us a server that present
the documentation of our code. The interface will looks familiar like the
official golang.org interface, but then inside "Packages" section are our local
packages.

### Example is also a test

Example functions are a little bit more concrete in term of showing people how
to use our API. More interestingly, Examples are not only for documentation but
they can also be tests. For them to be tests, we need to add a comment at the
end of the functions: one is Output and one is expected output. If we change the
expected output to be something wrong then, the complier will tell us when we
run the test. Example:

```go
fmt.Println(u)
// Output:
// {Bill bill@ardanlabs.com}
```

So anything in this example writes to standard out, we can now validate.

So, let's go and test that.

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example4/handlers
$ go test -run ExampleSendJSON
PASS
ok      github.com/cedrickchee/ultimate-go/testing/tests/example4/handlers      0.003s
```

Example tests are powerful. They give users examples how to use the API and
validate that the APIs and examples are working.

## Sub Tests

Sub test let us streamline our test functions, filters out command-line level
large tests into smaller sub tests.

Sub tests help when we do table tests because one of the things that we have
with table tests is the ability to do this data-driven testing. Sub tests let
us filter a piece of data at the command line-level.

[Sample test](example5/example5_test.go).

```go
// TestDownload validates the http Get function can download content and
// handles different status conditions properly.
func TestDownload(t *testing.T) {
	// Data for our standard table tests.
	// name field will give us the ability to name each piece of data and then
	// filter things on the command-line.
	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{"statusok", "https://www.ardanlabs.com/blog/index.xml", http.StatusOK},
		{"statusnotfound", "http://rss.cnn.com/rss/cnn_topstorie.rss", http.StatusNotFound},
	}

	t.Log("Given the need to test downloading different content.")
	{
		// Range over our table but this time, create an anonymous function (tf)
		// that takes a testing T parameter. This is a test function inside a test
		// function. What nice about it is that we will have a new function for
		// each set of data that we have in our table. Therefore, we will end up
		// with 2 different functions here.
		for i, tt := range tests {
			tf := func(t *testing.T) {
				t.Logf("\tTest: %d\tWhen checking %q for status code %d", i, tt.url, tt.statusCode)
				{
					resp, err := http.Get(tt.url)
					if err != nil {
						t.Fatalf("\t%s\tShould be able to make the Get call : %v", failed, err)
					}
					t.Logf("\t%s\tShould be able to make the Get call.", succeed)

					defer resp.Body.Close()

					if resp.StatusCode == tt.statusCode {
						t.Logf("\t%s\tShould receive a %d status code.", succeed, tt.statusCode)
					} else {
						t.Errorf("\t%s\tShould receive a %d status code : %v", failed, tt.statusCode, resp.StatusCode)
					}
				}
			}

			t.Run(tt.name, tf)
		}
	}
}
```

Run this test:

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example5
$ go test -run TestDownload -v
=== RUN   TestDownload
=== RUN   TestDownload/statusok
=== RUN   TestDownload/statusnotfound
--- PASS: TestDownload (3.62s)
    example5_test.go:34: Given the need to test downloading different content.
    --- PASS: TestDownload/statusok (3.06s)
        example5_test.go:38:    Test: 0 When checking "https://www.ardanlabs.com/blog/index.xml" for status code 200
        example5_test.go:44:    ✓       Should be able to make the Get call.
        example5_test.go:49:    ✓       Should receive a 200 status code.
    --- PASS: TestDownload/statusnotfound (0.56s)
        example5_test.go:38:    Test: 1 When checking "http://rss.cnn.com/rss/cnn_topstorie.rss" for status code 404
        example5_test.go:44:    ✓       Should be able to make the Get call.
        example5_test.go:49:    ✓       Should receive a 404 status code.
PASS
ok      github.com/cedrickchee/ultimate-go/testing/tests/example5       3.622s
```

It executed it in series, one after the other.

Now, let's say I just wanted to run "statusok". I didn't have to go into my code
and comment out any particular line in my table.

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example5
$ go test -run TestDownload/statusok -v
=== RUN   TestDownload
=== RUN   TestDownload/statusok
--- PASS: TestDownload (1.37s)
    example5_test.go:34: Given the need to test downloading different content.
    --- PASS: TestDownload/statusok (1.37s)
        example5_test.go:38:    Test: 0 When checking "https://www.ardanlabs.com/blog/index.xml" for status code 200
        example5_test.go:44:    ✓       Should be able to make the Get call.
        example5_test.go:49:    ✓       Should receive a 200 status code.
PASS
ok      github.com/cedrickchee/ultimate-go/testing/tests/example5       1.372s
```

We can run all of this data in parallel when it's reasonable and practical to
do so.

```go
// TestParallelize validates the http Get function can download content and
// handles different status conditions properly but runs the tests in parallel.
func TestParallelize(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		statusCode int
	}{
		{"statusok", "https://www.goinggo.net/post/index.xml", http.StatusOK},
		{"statusnotfound", "http://rss.cnn.com/rss/cnn_topstorie.rss", http.StatusNotFound},
	}

	t.Log("Given the need to test downloading different content.")
	{
		for i, tt := range tests {
			tf := func(t *testing.T) {
				t.Parallel()

				t.Logf("\tTest: %d\tWhen checking %q for status code %d", i, tt.url, tt.statusCode)
				{
					resp, err := http.Get(tt.url)
					if err != nil {
						t.Fatalf("\t%s\tShould be able to make the Get call : %v", failed, err)
					}
					t.Logf("\t%s\tShould be able to make the Get call.", succeed)

					defer resp.Body.Close()

					if resp.StatusCode == tt.statusCode {
						t.Logf("\t%s\tShould receive a %d status code.", succeed, tt.statusCode)
					} else {
						t.Errorf("\t%s\tShould receive a %d status code : %v", failed, tt.statusCode, resp.StatusCode)
					}
				}
			}

			t.Run(tt.name, tf)
		}
	}
}
```

The only difference here is that we call `t.Parallel()` function inside each of
these individual sub test function. This tells the testing tool to run each as a
separate or independent Goroutine.

See what happens when we run this:

```sh
~/m/dev/work/repo/experiments/go/ultimate-go/testing/tests/example5
$ go test -run TestParallelize -v
=== RUN   TestParallelize
=== RUN   TestParallelize/statusok
=== PAUSE TestParallelize/statusok
=== RUN   TestParallelize/statusnotfound
=== PAUSE TestParallelize/statusnotfound
=== CONT  TestParallelize/statusok
=== CONT  TestParallelize/statusnotfound
--- PASS: TestParallelize (0.00s)
    example5_test.go:73: Given the need to test downloading different content.
    --- PASS: TestParallelize/statusok (0.48s)
        example5_test.go:79:    Test: 1 When checking "http://rss.cnn.com/rss/cnn_topstorie.rss" for status code 404
        example5_test.go:85:    ✓       Should be able to make the Get call.
        example5_test.go:90:    ✓       Should receive a 404 status code.
    --- PASS: TestParallelize/statusnotfound (0.48s)
        example5_test.go:79:    Test: 1 When checking "http://rss.cnn.com/rss/cnn_topstorie.rss" for status code 404
        example5_test.go:85:    ✓       Should be able to make the Get call.
        example5_test.go:90:    ✓       Should receive a 404 status code.
PASS
ok      github.com/cedrickchee/ultimate-go/testing/tests/example5       0.481s
```

Look how fast it ran. It wasn't a second anymore.

Sub tests are fantastic when we've got these data-driven tests, because it allow
us to not just be able to isolate on the command line one given piece of data
over the other, but run all of this stuff in parallel, which means our unit
tests can run faster.

## Code Coverage

In one of the earlier sections, I told you, "how do you know when you're done
with a piece of code?" One of the things we said was code coverage.

### Basic code coverage commands

Look at the coverage of all test cases:

```sh
$ go test -cover
PASS
coverage: 100.0% of statements
ok      github.com/cedrickchee/ultimate-go/testing/tests/example4/handlers      0.003s
```

Get the cover profile and write it out to `c.out`:

```sh
$ go test -coverprofile c.out
```

I've generated a `c.out` file with the profiling information for the cover
report. But how do I look at this coverage? Well, go tool cover. Get an HTML
representation in the browser:

```sh
$ go tool cover -html=c.out
```

Making sure your tests cover as much of your code as possible is critical.
Make sure you're at that 70 to 80% code coverage.
Go's testing tool allows you to create a profile for the code that is executed
during all the tests and see a visual of what is and is not covered.
