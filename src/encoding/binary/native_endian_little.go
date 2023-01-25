// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build 386 || amd64 || amd64p32 || alpha || arm || arm64 || loong64 || mipsle || mips64le || mips64p32le || nios2 || ppc64le || riscv || riscv64 || sh || wasm

package binary

type nativeEndian struct {
	littleEndian
}

// NativeEndian is the native-endian implementation of ByteOrder and AppendByteOrder.
var NativeEndian nativeEndian
