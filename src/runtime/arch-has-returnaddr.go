//go:build !(avr || tinygo.wasm || tinygo.zkvm)
// +build !avr,!tinygo.wasm,!tinygo.zkvm

package runtime

import "unsafe"

const hasReturnAddr = true

//export llvm.returnaddress
func returnAddress(level uint32) unsafe.Pointer
