//go:build tinygo.zkvm
// +build tinygo.zkvm

package runtime

/*
#include <stdint.h>
#cgo LDFLAGS: /Users/austinabell/development/github.com/risc0/risc0/target/riscv32im-risc0-zkvm-elf/release/librisc0_zkvm_platform.a
void sys_write(uint32_t fd, char* c, int len);
void sys_halt(uint8_t exit_code);
*/
import "C"
import (
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
	c_char := C.char(c)
	C.sys_write(1 /* stdout */, &c_char, 1)

	return
}

func abort() {
	exit(1)
}

//go:linkname exit syscall.Exit
func exit(code int) {
	// TODO exiting this way will not currently be able to pull journal hash
	// var initialSHAState = [8]uint32{
	// 	0x6a09e667,
	// 	0xbb67ae85,
	// 	0x3c6ef372,
	// 	0xa54ff53a,
	// 	0x510e527f,
	// 	0x9b05688c,
	// 	0x1f83d9ab,
	// 	0x5be0cd19,
	// }

	// // Arch is little endian, but these values are expected as big endian for SHA, swap.
	// for i := 0; i < len(initialSHAState); i++ {
	// 	b := make([]byte, 8)
	// 	binary.BigEndian.PutUint32(b, initialSHAState[i])
	// 	initialSHAState[i] = uint32(binary.LittleEndian.Uint32(b))
	// }
	C.sys_halt(uint8(code))
}

//go:linkname now time.now
func now() (sec int64, nsec int32, mono int64) {
	mono = nanotime()
	sec = mono / (1000 * 1000 * 1000)
	nsec = int32((mono) - sec*(1000*1000*1000))
	return
}

//go:linkname cycle_count sys_cycle_count
func cycle_count() uint32

// zkVM has no notion of time the closest thing to time is the current cycle count
func ticks() timeUnit {
	return timeUnit(cycle_count())
}

// Align on word boundary.
func align(ptr uintptr) uintptr {
	return (ptr + 3) &^ 3
}

func getCurrentStackPointer() uintptr {
	return uintptr(stacksave())
}
