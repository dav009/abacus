package abacus

import "sync"
import (
	"github.com/spaolacci/murmur3"
	"unsafe"
)

type CountType uint32
const Max = ^(CountType(0))

type Sketch struct {
	Width uint32
	Depth uint32
	Count [][]CountType
	mutex sync.RWMutex
}

func sizeOfCell() uintptr{
	var a CountType
	return unsafe.Sizeof(a)
}

func NewSketch(width, depth uint32) (sk *Sketch) {
	sk = &Sketch{
		Width: width,
		Depth: depth,
		Count: make([][]CountType, depth),
	}
	for i := uint32(0); i < depth; i++ {
		sk.Count[i] = make([]CountType, width)
	}
	return sk
}

func (sk *Sketch) Incr(dat []byte) (min CountType) {
	return sk.Add(dat, 1)
}

func (sk *Sketch) positions(dat []byte) (pos []uint32) {
	// reference: https://github.com/addthis/stream-lib/blob/master/src/main/java/com/clearspring/analytics/stream/membership/Filter.java
	hash1 := murmur3.Sum32WithSeed(dat, 0)
	hash2 := murmur3.Sum32WithSeed(dat, hash1)
	pos = make([]uint32, sk.Depth)
	for i := uint32(0); i < sk.Depth; i++ {
		pos[i] = (hash1 + i*hash2) % sk.Width
	}
	return pos
}

func (sk *Sketch) Add(dat []byte, cnt CountType) (min CountType) {
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

func (sk *Sketch) Query(dat []byte) (min CountType) {
	pos := sk.positions(dat)
	return sk.query(pos)
}

func (sk *Sketch) query(pos []uint32) (min CountType) {
	min = Max

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