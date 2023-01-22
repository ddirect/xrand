package xrand

import (
	"bytes"
	"encoding/binary"
	"time"

	"golang.org/x/exp/rand"
)

type Xrand struct {
	*rand.Rand
	bits     uint64
	bitsLeft int
}

func timeSeed() uint64 {
	return uint64(time.Now().UnixNano())
}

func New() *Xrand {
	return NewWithSeed(timeSeed())
}

func NewWithSeed(seed uint64) *Xrand {
	return &Xrand{rand.New(rand.NewSource(seed)), 0, 0}
}

func NewPair() (*Xrand, *Xrand) {
	seed := timeSeed()
	return NewWithSeed(seed), NewWithSeed(seed)
}

func NewSlice(size int) []*Xrand {
	seed := timeSeed()
	res := make([]*Xrand, size)
	for i := range res {
		res[i] = NewWithSeed(seed)
	}
	return res
}

// Returns a random int between min and max, included
func (r *Xrand) Uniform(min, max int) int {
	return min + r.Intn(max-min+1)
}

func (r *Xrand) Bool() bool {
	if r.bitsLeft == 0 {
		r.bits = r.Uint64()
		r.bitsLeft = 64
	}
	r.bitsLeft--
	res := (r.bits & 1) != 0
	r.bits >>= 1
	return res
}

func (r *Xrand) UniformFactory(min, max int) func() int {
	if min == max {
		return func() int {
			return min
		}
	} else {
		return func() int {
			return r.Uniform(min, max)
		}
	}
}

func (r *Xrand) Buffer(size int) []byte {
	b := make([]byte, size)
	r.Fill(b)
	return b
}

func (r *Xrand) Fill(b []byte) {
	for len(b) >= 8 {
		binary.LittleEndian.PutUint64(b, r.Uint64())
		b = b[8:]
	}
	if len(b) > 0 {
		r.Read(b)
	}
}

func (r *Xrand) Verify(b []byte) bool {
	for len(b) >= 8 {
		if binary.LittleEndian.Uint64(b) != r.Uint64() {
			return false
		}
		b = b[8:]
	}
	if len(b) > 0 {
		var bufstore [8]byte
		buf := bufstore[:len(b)]
		r.Read(buf)
		return bytes.Compare(b, buf) == 0
	}
	return true
}
