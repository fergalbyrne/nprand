/*
 Ogma Toolkit (OTK)
 Copyright (c) 2016 Ogma Intelligent Systems Corp. All rights reserved.
*/

package nprand

import (
	_ "fmt"
	"math/rand"
)

type MT struct {
	// Create two length 624 array to store the state of the generators
	index, lp uint32
	Key       [624]uint32
}

const (
	N          = 624
	M          = 397
	MATRIX_A   = 0x9908b0df
	UPPER_MASK = 0x80000000
	LOWER_MASK = 0x7fffffff
)

var GlobalMT rand.Source

// New Mersenne Twister
func NewMT() rand.Source {
	m := MT{index: 0}
	m.Seed(1)
	return &m
}

// Initialize the generator from a seed
func (m *MT) Seed(theseed int64) {
	// Lower 32bit
	seed := uint32(theseed)
	for m.lp = 0; m.lp < 624; m.lp++ {
		m.Key[m.lp] = seed
		seed = (uint32(1812433253)*(seed^(seed>>30)) + m.lp + 1) & 0xffffffff
		//m.arrayL[m.lp] = 0x6c078965*(m.arrayL[m.lp-1]^(m.arrayL[m.lp-1]>>30)) + m.lp
	}
}

func (m *MT) random_int32() (y uint32) {
	if m.lp == 624 {
		i := 0
		//fmt.Printf("tempering\n%v", m.Key)
		for i = 0; i < N-M; i++ {
			y = (m.Key[i] & UPPER_MASK) | (m.Key[i+1] & LOWER_MASK)
			m.Key[i] = m.Key[i+M] ^ (y >> 1) ^ (-(y & 1) & MATRIX_A)
		}
		for ; i < N-1; i++ {
			y = (m.Key[i] & UPPER_MASK) | (m.Key[i+1] & LOWER_MASK)
			m.Key[i] = m.Key[i+(M-N)] ^ (y >> 1) ^ (-(y & 1) & MATRIX_A)
		}
		y = (m.Key[N-1] & UPPER_MASK) | (m.Key[0] & LOWER_MASK)
		m.Key[N-1] = m.Key[M-1] ^ (y >> 1) ^ (-(y & 1) & MATRIX_A)

		m.lp = 0
		//fmt.Printf("tempered\n%v", m.Key)
	}
	y = m.Key[m.lp]
	m.lp++

	/* Tempering */
	y = y ^ (y >> 11)
	y = y ^ (y<<7)&0x9d2c5680
	y = y ^ (y<<15)&0xefc60000
	y = y ^ (y >> 18)

	return
}

func (m *MT) RandomInt32(low, high int32) int32 {
	rng := uint32(high - low)
	off := uint32(int32(low))
	return int32(m.RandomUint32(off, rng))
}

func (m *MT) RandomUint32(off, rng uint32) uint32 {
	val := rng
	mask := rng

	if rng <= 0 {
		//_ = m.random_int32()
		return off
	}

	/* Smallest bit mask >= max */
	mask = mask | (mask >> 1)
	mask = mask | (mask >> 2)
	mask = mask | (mask >> 4)
	mask = mask | (mask >> 8)
	mask = mask | (mask >> 16)

	val = m.random_int32() & mask
	for val > rng {
		val = m.random_int32() & mask
	}
	return off + val
}

func (m *MT) Int63() int64 {
	upper := uint64(m.random_int32()) << 32
	lower := uint64(m.random_int32())
	return int64(upper|lower) & 0x7FFFFFFFFFFFFFFF
}

func (m *MT) Float64() float64 {
again:
	f := float64(m.Int63()) / (1 << 63)
	if f == 1 {
		goto again // resample; this branch is taken O(never)
	}
	return f
}

func init() {
	GlobalMT = NewMT()
}
