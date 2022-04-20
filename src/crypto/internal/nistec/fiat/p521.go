// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by generate.go. DO NOT EDIT.

package fiat

import (
	"crypto/subtle"
	"errors"
)

// P521Element is an integer modulo 2^521 - 1.
//
// The zero value is a valid zero element.
type P521Element struct {
	// Values are represented internally always in the Montgomery domain, and
	// converted in Bytes and SetBytes.
	x p521MontgomeryDomainFieldElement
}

const p521ElementLen = 66

type p521UntypedFieldElement = [9]uint64

// One sets e = 1, and returns e.
func (e *P521Element) One() *P521Element {
	p521SetOne(&e.x)
	return e
}

// Equal returns 1 if e == t, and zero otherwise.
func (e *P521Element) Equal(t *P521Element) int {
	eBytes := e.Bytes()
	tBytes := t.Bytes()
	return subtle.ConstantTimeCompare(eBytes, tBytes)
}

// IsZero returns 1 if e == 0, and zero otherwise.
func (e *P521Element) IsZero() int {
	zero := make([]byte, p521ElementLen)
	eBytes := e.Bytes()
	return subtle.ConstantTimeCompare(eBytes, zero)
}

// Set sets e = t, and returns e.
func (e *P521Element) Set(t *P521Element) *P521Element {
	e.x = t.x
	return e
}

// Bytes returns the 66-byte big-endian encoding of e.
func (e *P521Element) Bytes() []byte {
	// This function is outlined to make the allocations inline in the caller
	// rather than happen on the heap.
	var out [p521ElementLen]byte
	return e.bytes(&out)
}

func (e *P521Element) bytes(out *[p521ElementLen]byte) []byte {
	var tmp p521NonMontgomeryDomainFieldElement
	p521FromMontgomery(&tmp, &e.x)
	p521ToBytes(out, (*p521UntypedFieldElement)(&tmp))
	p521InvertEndianness(out[:])
	return out[:]
}

// SetBytes sets e = v, where v is a big-endian 66-byte encoding, and returns e.
// If v is not 66 bytes or it encodes a value higher than 2^521 - 1,
// SetBytes returns nil and an error, and e is unchanged.
func (e *P521Element) SetBytes(v []byte) (*P521Element, error) {
	if len(v) != p521ElementLen {
		return nil, errors.New("invalid P521Element encoding")
	}

	// Check for non-canonical encodings (p + k, 2p + k, etc.) by comparing to
	// the encoding of -1 mod p, so p - 1, the highest canonical encoding.
	var minusOneEncoding = new(P521Element).Sub(
		new(P521Element), new(P521Element).One()).Bytes()
	for i := range v {
		if v[i] < minusOneEncoding[i] {
			break
		}
		if v[i] > minusOneEncoding[i] {
			return nil, errors.New("invalid P521Element encoding")
		}
	}

	var in [p521ElementLen]byte
	copy(in[:], v)
	p521InvertEndianness(in[:])
	var tmp p521NonMontgomeryDomainFieldElement
	p521FromBytes((*p521UntypedFieldElement)(&tmp), &in)
	p521ToMontgomery(&e.x, &tmp)
	return e, nil
}

// Add sets e = t1 + t2, and returns e.
func (e *P521Element) Add(t1, t2 *P521Element) *P521Element {
	p521Add(&e.x, &t1.x, &t2.x)
	return e
}

// Sub sets e = t1 - t2, and returns e.
func (e *P521Element) Sub(t1, t2 *P521Element) *P521Element {
	p521Sub(&e.x, &t1.x, &t2.x)
	return e
}

// Mul sets e = t1 * t2, and returns e.
func (e *P521Element) Mul(t1, t2 *P521Element) *P521Element {
	p521Mul(&e.x, &t1.x, &t2.x)
	return e
}

// Square sets e = t * t, and returns e.
func (e *P521Element) Square(t *P521Element) *P521Element {
	p521Square(&e.x, &t.x)
	return e
}

// Select sets v to a if cond == 1, and to b if cond == 0.
func (v *P521Element) Select(a, b *P521Element, cond int) *P521Element {
	p521Selectznz((*p521UntypedFieldElement)(&v.x), p521Uint1(cond),
		(*p521UntypedFieldElement)(&b.x), (*p521UntypedFieldElement)(&a.x))
	return v
}

func p521InvertEndianness(v []byte) {
	for i := 0; i < len(v)/2; i++ {
		v[i], v[len(v)-1-i] = v[len(v)-1-i], v[i]
	}
}
