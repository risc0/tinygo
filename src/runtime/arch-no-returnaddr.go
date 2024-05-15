//go:build avr || tinygo.wasm || tinygo.zkvm
// +build avr tinygo.wasm tinygo.zkvm

package runtime

import "unsafe"

const hasReturnAddr = false

func returnAddress(level uint32) unsafe.Pointer {
	return nil
}
