package abacus

import (
	"github.com/sasha-s/go-hll"
	"github.com/spaolacci/murmur3"
	"math/big"
)

type Abacus interface{
	Counts(key string) (CountType, error)
	Update(items []string) error
	Total() (*big.Int, error)
	Cardinality() (CountType, error)
}

func widthAndDepthFromSize(sizeMB uint) (uint32, uint32){
	width := uint64(uint64(sizeMB*1000000) / uint64( 2 * 8 * sizeOfCell() ))
	depth :=(uint64(sizeMB)*1000000) / (width * uint64(sizeOfCell()))
	return uint32(width), uint32(depth)
}

type memoryAbacus struct{
	MaxMemorySize uint
	s *Sketch
	h hll.HLL
	total *big.Int
}

func (a *memoryAbacus)  Counts(key string) (CountType, error) {
	return a.s.Query([]byte(key)), nil
}

func (a *memoryAbacus) Update(items []string) error {
	for  _, key := range items {
		a.s.Incr([]byte(key))
		a.h.Add(uint64(murmur3.Sum32([]byte(key))))
		a.total = a.total.Add(big.NewInt(1), a.total)
	}

  return nil
}

func (a *memoryAbacus) Total() (*big.Int, error){
	return a.total, nil
}

func (a *memoryAbacus) Cardinality() (CountType, error){
	return CountType(a.h.EstimateCardinality()),nil
}

func New(maxMemoryMB uint) memoryAbacus {
	w, d := widthAndDepthFromSize(maxMemoryMB)
	sketch := NewSketch(w, d)
	s, _ := hll.SizeByP(16)
	h := make(hll.HLL, s)
	a:= memoryAbacus{ MaxMemorySize: maxMemoryMB, s:sketch, h:h, total: big.NewInt(0)}
	return a
}


type memoryAbacusLog8 struct{
	MaxMemorySize uint
	s *SketchLog8
	h hll.HLL
	total *big.Int
}

func (a *memoryAbacusLog8)  Counts(key string) (CountType, error) {
	return CountType(a.s.Query([]byte(key))), nil
}

func (a *memoryAbacusLog8) Update(items []string) error {
	for  _, key := range items {
		a.s.Incr([]byte(key))
		a.h.Add(uint64(murmur3.Sum32([]byte(key))))
		a.total = a.total.Add(big.NewInt(1), a.total)
	}

  return nil
}

func (a *memoryAbacusLog8) Total() (*big.Int, error){
	return a.total, nil
}

func (a *memoryAbacusLog8) Cardinality() (CountType, error){
	return CountType(a.h.EstimateCardinality()),nil
}

func widthAndDepthFromSizeLog8(sizeMB uint) (uint32, uint32){
	width := uint64(uint64(sizeMB*1000000) / uint64( 2 * 8 * sizeOfCellLog8() ))
	depth :=(uint64(sizeMB)*1000000) / (width * uint64(sizeOfCellLog8()))
	return uint32(width), uint32(depth)
}

func NewAbacus8Log(maxMemoryMB uint) memoryAbacusLog8 {
	w, d := widthAndDepthFromSizeLog8(maxMemoryMB)
	sketch := NewSketchLog8(w, d)
	s, _ := hll.SizeByP(16)
	h := make(hll.HLL, s)
	a:= memoryAbacusLog8{ MaxMemorySize: maxMemoryMB, s:sketch, h:h, total: big.NewInt(0)}
	return a
}
