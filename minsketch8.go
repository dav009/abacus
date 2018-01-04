package abacus

import "sync"
import (
	"github.com/spaolacci/murmur3"
	"unsafe"
)

type CountTypeLog8 uint8
const MaxLog8 = ^(CountTypeLog8(0))

type SketchLog8 struct {
	Width uint32
	Depth uint32
	Count [][]CountTypeLog8
	mutex sync.RWMutex
}

func sizeOfCellLog8() uintptr{
	var a CountTypeLog8
	return unsafe.Sizeof(a)
}

func NewSketchLog8(width, depth uint32) (sk *SketchLog8) {
	sk = &SketchLog8{
		Width: width,
		Depth: depth,
		Count: make([][]CountTypeLog8, depth),
	}
	for i := uint32(0); i < depth; i++ {
		sk.Count[i] = make([]CountTypeLog8, width)
	}
	return sk
}

func (sk *SketchLog8) Incr(dat []byte) (min CountTypeLog8) {
	return sk.Add(dat, 1)
}

func (sk *SketchLog8) positions(dat []byte) (pos []uint32) {
	// reference: https://github.com/addthis/stream-lib/blob/master/src/main/java/com/clearspring/analytics/stream/membership/Filter.java
	hash1 := murmur3.Sum32WithSeed(dat, 0)
	hash2 := murmur3.Sum32WithSeed(dat, hash1)
	pos = make([]uint32, sk.Depth)
	for i := uint32(0); i < sk.Depth; i++ {
		pos[i] = (hash1 + i*hash2) % sk.Width
	}
	return pos
}

func (sk *SketchLog8) Add(dat []byte, cnt CountTypeLog8) (min CountTypeLog8) {
	pos := sk.positions(dat)
	min = sk.query(pos)

	min += cnt

	sk.mutex.Lock()
	for i := uint32(0); i < sk.Depth; i++ {
		v := sk.Count[i][pos[i]]
		if v < min {
			sk.Count[i][pos[i]] = min
		}
	}
	sk.mutex.Unlock()

	return min
}

func (sk *SketchLog8) Query(dat []byte) (min CountTypeLog8) {
	pos := sk.positions(dat)
	return sk.query(pos)
}

func (sk *SketchLog8) query(pos []uint32) (min CountTypeLog8) {
	min = MaxLog8

	sk.mutex.RLock()
	for i := uint32(0); i < sk.Depth; i++ {
		v := sk.Count[i][pos[i]]
		if min > v {
			min = v
		}
	}
	sk.mutex.RUnlock()

	return min
}


