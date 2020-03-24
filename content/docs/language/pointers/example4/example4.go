// Copyright 2014 Ardan Studios
//
// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to teach the mechanics of escape analysis.
package main

// user represents a user in the system.
type user struct {
	name  string
	email string
}

// main is the entry point for the application.
func main() {
	u1 := createUserV1()
	u2 := createUserV2()

	println("u1", &u1, "u2", u2)
}

// createUserV1 creates a user value and passed a copy back to the caller.
//
// shows how the variable does not escape.
// Since we know the size of the user value at compiled time, the complier
// will put this on a stack frame.
//go:noinline
func createUserV1() user {
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V1", &u)

	return u
}

// createUserV2 creates a user value and shares the value with the caller.
//
// shows how the variable escape.
// This looks almost identical to the createUserV1 function.
// It creates a value of type user and initialize it. It seems like we are doing
// the same here. However, there is one subtle difference: we do not return the
// value itself but the address of u. That is the value that is being passed
// back up the call stack. We are using pointer semantic.
//
// You might think about what we have after this call is: main has a pointer to
// a value that is on a stack frame below. If this is the case, then we are in
// trouble. Once we come back up the call stack, this memory is there but it is
// reusable again. It is no longer valid. Anytime now main makes a function
// call, we need to allocate the frame and initialize it.
//
// Think about zero value for a second here. It is enable to us to initialize every stack frame that
// we take. Stack are self cleaning. We clean our stack on the way down. Every time we make a
// function call, zero value, initialization, we are cleaning those stack frames. We leave that
// memory on the way up because we don't know if we need that again.
//
// Back to the example, it is bad because it looks like we take the address of user value, pass it
// back up to the call stack and we now have a pointer which is about to get erased. Therefore, it
// is not what will happen.
//
// What actually going to happen is the idea of escape analysis.
// Because of line "return &u", this value cannot be put inside the stack frame for this function
// so we have to put it out on the heap.
// Escape analysis decides what stay on stack and what not.
// In the createUserV1 function, because we are passing the copy of the value itself, it is safe to
// keep these things on the stack. But when we SHARE something above the call stack like this,
// escape analysis said this memory is no longer be valid when we get back to main, we must put it
// out there on the heap. main will end up having a pointer to the heap.
// In fact, this allocation happens immediately on the heap. createUserV2 is gonna have a pointer
// to the heap. But u is gonna base on value semantic.
//go:noinline
func createUserV2() *user {
	u := user{
		name:  "Bill",
		email: "bill@ardanlabs.com",
	}

	println("V2", &u)

	return &u
}

// Outputs:
// $ go build example4.go
// V1 0xc00003a6d8
// V2 0xc00007c000
// u1 0xc00003a730 u2 0xc00007c000

/*
// See escape analysis and inlining decisions.
$ go build -gcflags -m=2
# github.com/ardanlabs/gotraining/topics/go/language/pointers/example4
./example4.go:24:6: cannot inline createUserV1: marked go:noinline
./example4.go:38:6: cannot inline createUserV2: marked go:noinline
./example4.go:14:6: cannot inline main: non-leaf function
./example4.go:30:16: createUserV1 &u does not escape
./example4.go:46:9: &u escapes to heap
./example4.go:46:9: 	from ~r0 (return) at ./example4.go:46:2
./example4.go:39:2: moved to heap: u
./example4.go:44:16: createUserV2 &u does not escape
./example4.go:18:16: main &u1 does not escape
./example4.go:18:27: main &u2 does not escape

// See the intermediate representation phase before
// generating the actual arch-specific assembly.
$ go build -gcflags -S
0x0021 00033 (/.../example4.go:15)	CALL	"".createUserV1(SB)
0x0026 00038 (/.../example4.go:15)	MOVQ	(SP), AX
0x002a 00042 (/.../example4.go:15)	MOVQ	8(SP), CX
0x002f 00047 (/.../example4.go:15)	MOVQ	16(SP), DX
0x0034 00052 (/.../example4.go:15)	MOVQ	24(SP), BX
0x0039 00057 (/.../example4.go:15)	MOVQ	AX, "".u1+40(SP)
0x003e 00062 (/.../example4.go:15)	MOVQ	CX, "".u1+48(SP)
0x0043 00067 (/.../example4.go:15)	MOVQ	DX, "".u1+56(SP)
0x0048 00072 (/.../example4.go:15)	MOVQ	BX, "".u1+64(SP)
0x004d 00077 (/.../example4.go:16)	PCDATA	$0, $1

// See bounds checking decisions.
$ go build -gcflags="-d=ssa/check_bce/debug=1"

// See the actual machine representation by using
// the disasembler.
$ go tool objdump -s main.main example4
TEXT main.main(SB) /home/cedric/m/dev/work/repo/experiments/go/ultimate-go/language/pointers/example4/example4.go
  example4.go:16        0x4525b0                64488b0c25f8ffffff      MOVQ FS:0xfffffff8, CX
  example4.go:16        0x4525b9                483b6110                CMPQ 0x10(CX), SP
  example4.go:16        0x4525bd                0f86af000000            JBE 0x452672
  example4.go:16        0x4525c3                4883ec50                SUBQ $0x50, SP
  example4.go:16        0x4525c7                48896c2448              MOVQ BP, 0x48(SP)
  example4.go:16        0x4525cc                488d6c2448              LEAQ 0x48(SP), BP
  example4.go:17        0x4525d1                e8aa000000              CALL main.createUserV1(SB)
  example4.go:17        0x4525d6                488b0424                MOVQ 0(SP), AX
  example4.go:17        0x4525da                488b4c2408              MOVQ 0x8(SP), CX
  example4.go:17        0x4525df                488b542410              MOVQ 0x10(SP), DX
  example4.go:17        0x4525e4                488b5c2418              MOVQ 0x18(SP), BX
  example4.go:17        0x4525e9                4889442428              MOVQ AX, 0x28(SP)
  example4.go:17        0x4525ee                48894c2430              MOVQ CX, 0x30(SP)
  example4.go:17        0x4525f3                4889542438              MOVQ DX, 0x38(SP)
  example4.go:17        0x4525f8                48895c2440              MOVQ BX, 0x40(SP)
  example4.go:18        0x4525fd                e85e010000              CALL main.createUserV2(SB)
  example4.go:18        0x452602                488b0424                MOVQ 0(SP), AX
  example4.go:18        0x452606                4889442420              MOVQ AX, 0x20(SP)
  example4.go:20        0x45260b                e8c030fdff              CALL runtime.printlock(SB)
  example4.go:20        0x452610                488d05a4100200          LEAQ 0x210a4(IP), AX
  example4.go:20        0x452617                48890424                MOVQ AX, 0(SP)
  example4.go:20        0x45261b                48c744240803000000      MOVQ $0x3, 0x8(SP)
  example4.go:20        0x452624                e8e739fdff              CALL runtime.printstring(SB)
  example4.go:20        0x452629                488d442428              LEAQ 0x28(SP), AX
  example4.go:20        0x45262e                48890424                MOVQ AX, 0(SP)
  example4.go:20        0x452632                e89939fdff              CALL runtime.printpointer(SB)
  example4.go:20        0x452637                488d05a4100200          LEAQ 0x210a4(IP), AX
  example4.go:20        0x45263e                48890424                MOVQ AX, 0(SP)
  example4.go:20        0x452642                48c744240804000000      MOVQ $0x4, 0x8(SP)
  example4.go:20        0x45264b                e8c039fdff              CALL runtime.printstring(SB)
  example4.go:20        0x452650                488b442420              MOVQ 0x20(SP), AX
  example4.go:20        0x452655                48890424                MOVQ AX, 0(SP)
  example4.go:20        0x452659                e87239fdff              CALL runtime.printpointer(SB)
  example4.go:20        0x45265e                e8fd32fdff              CALL runtime.printnl(SB)
  example4.go:20        0x452663                e8e830fdff              CALL runtime.printunlock(SB)
  example4.go:21        0x452668                488b6c2448              MOVQ 0x48(SP), BP
  example4.go:21        0x45266d                4883c450                ADDQ $0x50, SP
  example4.go:21        0x452671                c3                      RET
  example4.go:16        0x452672                e8897affff              CALL runtime.morestack_noctxt(SB)
  example4.go:16        0x452677                e934ffffff              JMP main.main(SB)

// See a list of the symbols in an artifact with
// annotations and size.
$ go tool nm example4
481ca8 r internal/cpu.xgetbv.args_stackmap
4c20b0 D main..inittask
4c20b0 D main..inittask
452680 T main.createUserV1
481e20 R main.createUserV1.stkobj
452760 T main.createUserV2
4525b0 T main.main
481e40 R main.main.stkobj
482fc0 r masks
402250 t memeqbody
4061f0 T runtime.(*TypeAssertionError).Error
*/
