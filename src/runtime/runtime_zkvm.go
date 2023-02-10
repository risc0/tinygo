//go:build tinygo.zkvm
// +build tinygo.zkvm

package runtime

/*
#include <stdlib.h>
__attribute__((always_inline))
int sys_cycle_count() {
    int result;
	__asm__("li t0, 2\n\t" // software ecall
	        "li a7, 4\n\t" // SYS_CYCLE_COUNT
	        "ecall\n\t"
	        "mv %[result], x10\n\t"
			:[result]"=r"(result)
			: // no input
			: "t0", "a0", "a1", "a7"// no clobber
			);
	return result;
};
*/
import "C"
import (
	"device/riscv"
	"unsafe"
)

const GOARCH = "zkvm"

// The bitness of the CPU (e.g. 8, 32, 64).
const TargetBits = 32

const deferExtraRegs = 0

type timeUnit uint32

var timestamp timeUnit

const baremetal = false

//go:extern _heap_start
var heapStartSymbol [0]byte

//go:extern _heap_end
var heapEndSymbol [0]byte

//go:extern _globals_start
var globalsStartSymbol [0]byte

//go:extern _globals_end
var globalsEndSymbol [0]byte

//go:extern _stack_top
var stackTopSymbol [0]byte

var (
	heapStart    = uintptr(unsafe.Pointer(&heapStartSymbol))
	heapEnd      = uintptr(unsafe.Pointer(&heapEndSymbol))
	globalsStart = uintptr(unsafe.Pointer(&globalsStartSymbol))
	globalsEnd   = uintptr(unsafe.Pointer(&globalsEndSymbol))
	stackTop     = uintptr(unsafe.Pointer(&stackTopSymbol))
)

func growHeap() bool {
	// On zkvm, there is no way the heap can be grown.
	return false
}

//go:linkname procPin sync/atomic.runtime_procPin
func procPin() {
}

//go:linkname procUnpin sync/atomic.runtime_procUnpin
func procUnpin() {
}

//export main
func main() {
	preinit()
	run()
	exit(0)
}

//go:extern __bss_begin
var __bss_begin [0]byte

//go:extern __bss_end
var __bss_end [0]byte

func preinit() {
	// Initialize .bss: zero-initialized global variables.
	ptr := unsafe.Pointer(&__bss_begin)
	for ptr != unsafe.Pointer(&__bss_end) {
		*(*uint32)(ptr) = 0
		ptr = unsafe.Pointer(uintptr(ptr) + 4)
	}
}

func ticksToNanoseconds(ticks timeUnit) int64 {
	return int64(ticks)
}

func nanosecondsToTicks(ns int64) timeUnit {
	return timeUnit(ns)
}

func sleepTicks(d timeUnit) {
	// TODO
	timestamp += d
}

func putchar(c byte) {
	// Do the move first because the parameter c is stored in a0, which gets overwritten for the syscall
	riscv.AsmFull("mv a1, {value}",
		map[string]interface{}{"value": uintptr(unsafe.Pointer(&c))})
	riscv.Asm("li t0, 2") // software ecall
	riscv.Asm("li a7, 2") // SYS_IO
	riscv.Asm("li a0, 1") // stdout
	// we're only printing a single character. Set the buffer length = 1
	riscv.Asm("li a2, 1")
	riscv.Asm("ecall")
	return
}

func abort() {
	exit(1)
}

//go:linkname exit syscall.Exit
func exit(code int) {
	riscv.Asm("li t0, 0")
	riscv.Asm("ecall")
}

//go:linkname now time.now
func now() (sec int64, nsec int32, mono int64) {
	mono = nanotime()
	sec = mono / (1000 * 1000 * 1000)
	nsec = int32((mono) - sec*(1000*1000*1000))
	return
}

// zkVM has no notion of time the closest thing to time is the current cycle count
func ticks() timeUnit {
	return timeUnit(C.sys_cycle_count())
}

// Align on word boundary.
func align(ptr uintptr) uintptr {
	return (ptr + 3) &^ 3
}

func getCurrentStackPointer() uintptr {
	return uintptr(stacksave())
}
