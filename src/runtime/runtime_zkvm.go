//go:build tinygo.zkvm
// +build tinygo.zkvm

package runtime

import (
	"device/riscv"
	"unsafe"
)

type timeUnit int64

var timestamp timeUnit

const GOARCH = "zkvm"
const TargetBits = 32

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

// Align on word boundary.
func align(ptr uintptr) uintptr {
	return (ptr + 3) &^ 3
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

func ticks() timeUnit {
	return timestamp
}

func putchar(c byte) {
	// Do the move first because the parameter c is stored in a0, which gets overwritten for the syscall
	riscv.AsmFull("mv a1, {value}",
		map[string]interface{}{"value": uintptr(unsafe.Pointer(&c))})
	riscv.Asm("li t0, 2") // software ecall
	riscv.Asm("li a7, 2") // SYS_IO
	riscv.Asm("li a0, 1") // stdout
	// We're only printing a single character. Set the buffer length = 1
	riscv.Asm("li a2, 1")
	riscv.Asm("ecall")
	return
}

func getchar() byte {
	// TODO
	return 0
}

func buffered() int {
	// TODO
	return 0
}

func abort() {
	exit(1)
}

func exit(code int) {
	riscv.Asm("li t0, 0")
	riscv.Asm("ecall")
}
