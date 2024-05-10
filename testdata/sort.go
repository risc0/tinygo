package main

/*
#include <stdint.h>
#cgo LDFLAGS: /Users/austinabell/development/github.com/risc0/risc0/target/riscv32im-risc0-zkvm-elf/release/librisc0_zkvm_platform.a
void sys_halt(uint8_t exit_code, uint32_t* initial_sha_state);
void sys_write(uint32_t fd, uint8_t* byte_ptr, int len);
void sys_sha_buffer(uint32_t* out_state, uint32_t* in_state, uint8_t* buf, uint32_t count);
*/
import "C"
import (
	"encoding/binary"
	"sort"
	"unsafe"
)

// sort.Slice implicitly uses reflect.Swapper

func strings() {
	data := []string{"aaaa", "cccc", "bbb", "fff", "ggg"}
	sort.Slice(data, func(i, j int) bool {
		return data[i] > data[j]
	})
	println("strings")
	for _, d := range data {
		println(d)
	}
}

func int64s() {
	sd := []int64{1, 6, 3, 2, 1923, 123, -123, -29, 3, 0, 1}
	sort.Slice(sd, func(i, j int) bool {
		return sd[i] > sd[j]
	})
	println("int64s")
	for _, d := range sd {
		println(d)
	}

	ud := []uint64{1, 6, 3, 2, 1923, 123, 29, 3, 0, 1}
	sort.Slice(ud, func(i, j int) bool {
		return ud[i] > ud[j]
	})
	println("uint64s")
	for _, d := range ud {
		println(d)
	}
}

func int32s() {
	sd := []int32{1, 6, 3, 2, 1923, 123, -123, -29, 3, 0, 1}
	sort.Slice(sd, func(i, j int) bool {
		return sd[i] > sd[j]
	})
	println("int32s")
	for _, d := range sd {
		println(d)
	}

	ud := []uint32{1, 6, 3, 2, 1923, 123, 29, 3, 0, 1}
	sort.Slice(ud, func(i, j int) bool {
		return ud[i] > ud[j]
	})
	println("uint32s")
	for _, d := range ud {
		println(d)
	}
}

func int16s() {
	sd := []int16{1, 6, 3, 2, 1923, 123, -123, -29, 3, 0, 1}
	sort.Slice(sd, func(i, j int) bool {
		return sd[i] > sd[j]
	})
	println("int16s")
	for _, d := range sd {
		println(d)
	}

	ud := []uint16{1, 6, 3, 2, 1923, 123, 29, 3, 0, 1}
	sort.Slice(ud, func(i, j int) bool {
		return ud[i] > ud[j]
	})
	println("uint16s")
	for _, d := range ud {
		println(d)
	}
}

func int8s() {
	sd := []int8{1, 6, 3, 2, 123, -123, -29, 3, 0, 1}
	sort.Slice(sd, func(i, j int) bool {
		return sd[i] > sd[j]
	})
	println("int8s")
	for _, d := range sd {
		println(d)
	}

	ud := []uint8{1, 6, 3, 2, 123, 29, 3, 0, 1}
	sort.Slice(ud, func(i, j int) bool {
		return ud[i] > ud[j]
	})
	println("uint8s")
	for _, d := range ud {
		println(d)
	}
}

func ints() {
	sd := []int{1, 6, 3, 2, 123, -123, -29, 3, 0, 1}
	sort.Slice(sd, func(i, j int) bool {
		return sd[i] > sd[j]
	})
	println("ints")
	for _, d := range sd {
		println(d)
	}

	ud := []uint{1, 6, 3, 2, 123, 29, 3, 0, 1}
	sort.Slice(ud, func(i, j int) bool {
		return ud[i] > ud[j]
	})
	println("uints")
	for _, d := range ud {
		println(d)
	}
}

func structs() {
	type s struct {
		name string
		a    uint64
		b    uint32
		c    uint16
		d    int
		e    *struct {
			aa uint16
			bb int
		}
	}

	data := []s{
		{
			name: "struct 1",
			d:    100,
			e: &struct {
				aa uint16
				bb int
			}{aa: 123, bb: -10},
		},
		{
			name: "struct 2",
			d:    1,
			e: &struct {
				aa uint16
				bb int
			}{aa: 15, bb: 10},
		},
		{
			name: "struct 3",
			d:    10,
			e: &struct {
				aa uint16
				bb int
			}{aa: 31, bb: -1030},
		},
		{
			name: "struct 4",
			e: &struct {
				aa uint16
				bb int
			}{},
		},
	}
	sort.Slice(data, func(i, j int) bool {
		di := data[i]
		dj := data[j]
		return di.d*di.e.bb > dj.d*dj.e.bb
	})
	println("structs")
	for _, d := range data {
		println(d.name)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func shaBuffer(initialState [8]uint32, bytes []byte) [8]uint32 {
	var outState [8]uint32
	// TODO doesn't really handle alignment
	padLen := computeU32sNeeded(len(bytes))

	// Create buffer of u32s to ensure word aligned
	padBuf := make([]uint32, padLen)
	if !(len(bytes) <= len(padBuf)*4) {
		panic("bytes too long")
	}
	padBufu8 := unsafe.Slice((*byte)(unsafe.Pointer(&padBuf[0])), 4*len(padBuf))

	// Copy whole bytes
	length := int(min(len(bytes), padLen*4))
	copy(padBufu8[:length], bytes[:length])

	// Add END marker since this is always with a trailer
	padBufu8[length] = 0x80

	// Add trailer with number of bits written. This needs to be big endian.
	bitsTrailer := 8 * uint32(len(bytes))

	// Swap bits to BE
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, bitsTrailer)
	bitsTrailer = uint32(binary.LittleEndian.Uint32(b))

	padBuf[padLen-1] = bitsTrailer

	numBlocks := padLen / 16
	C.sys_sha_buffer(&outState[0], &initialState[0], &padBufu8[0], uint32(numBlocks))
	return outState
}

func taggedStruct(tag string, digests [][8]uint32) [8]uint32 {
	tagDigest := shaBuffer(sha256InitState(), []byte(tag))

	// Create buffer for all bytes to hash
	allBytes := make([]byte, 0, len(tagDigest)*4+len(digests)*32+2)

	for _, word := range tagDigest {
		allBytes = binary.LittleEndian.AppendUint32(allBytes, word)
	}

	for _, digest := range digests {
		for _, word := range digest {
			allBytes = binary.LittleEndian.AppendUint32(allBytes, word)
		}
	}

	// TODO doesn't handle extending with data, not needed for this use case yet

	allBytes = binary.LittleEndian.AppendUint16(allBytes, uint16(len(digests)))

	return shaBuffer(sha256InitState(), allBytes)
}

func sha256InitState() [8]uint32 {
	// BE conversion of the sha256 init state
	return [8]uint32{
		1743128938,
		2242799547,
		1928556092,
		989155237,
		2136084049,
		2355627419,
		2883158815,
		432922715,
	}
}

func main() {
	strings()
	int64s()
	int32s()
	int16s()
	int8s()
	ints()
	structs()

	var assumptionsDigestState = [8]uint32{0, 0, 0, 0, 0, 0, 0, 0}

	outputBytes := [4]byte{0,1,2,3}

	// No output, but finalize still calls sha on empty bytes it seems
	journalDigest := shaBuffer(sha256InitState(), outputBytes[:])

	C.sys_write(3 /* journal */, &outputBytes[0], C.int(len(outputBytes)))

	output := taggedStruct("risc0.Output", [][8]uint32{journalDigest, assumptionsDigestState})

	C.sys_halt(0, &output[0])
}

func computeU32sNeeded(lenBytes int) int {
	const WordSize = 4
	const BlockWords = 16
	// Add one byte for end marker
	nWords := alignUp(lenBytes+1, WordSize) / WordSize
	// Add two words for length at end (even though we only
	// use one of them, being a 32-bit architecture)
	nWords += 2

	return alignUp(nWords, BlockWords)
}

func alignUp(addr, al int) int {
	return (addr + al - 1) & ^(al - 1)
}
