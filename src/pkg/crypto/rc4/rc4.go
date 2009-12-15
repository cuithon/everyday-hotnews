// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This package implements RC4 encryption, as defined in Bruce Schneier's
// Applied Cryptography.
package rc4

// BUG(agl): RC4 is in common use but has design weaknesses that make
// it a poor choice for new protocols.

import (
	"os"
	"strconv"
)

// A Cipher is an instance of RC4 using a particular key.
type Cipher struct {
	s    [256]byte
	i, j uint8
}

type KeySizeError int

func (k KeySizeError) String() string {
	return "crypto/rc4: invalid key size " + strconv.Itoa(int(k))
}

// NewCipher creates and returns a new Cipher.  The key argument should be the
// RC4 key, at least 1 byte and at most 256 bytes.
func NewCipher(key []byte) (*Cipher, os.Error) {
	k := len(key)
	if k < 1 || k > 256 {
		return nil, KeySizeError(k)
	}
	var c Cipher
	for i := 0; i < 256; i++ {
		c.s[i] = uint8(i)
	}
	var j uint8 = 0
	for i := 0; i < 256; i++ {
		j += c.s[i] + key[i%k]
		c.s[i], c.s[j] = c.s[j], c.s[i]
	}
	return &c, nil
}

// XORKeyStream will XOR each byte of the given buffer with a byte of the
// generated keystream.
func (c *Cipher) XORKeyStream(buf []byte) {
	for i := range buf {
		c.i += 1
		c.j += c.s[c.i]
		c.s[c.i], c.s[c.j] = c.s[c.j], c.s[c.i]
		buf[i] ^= c.s[c.s[c.i]+c.s[c.j]]
	}
}

// Reset zeros the key data so that it will no longer appear in the
// process's memory.
func (c *Cipher) Reset() {
	for i := range c.s {
		c.s[i] = 0
	}
	c.i, c.j = 0, 0
}
