// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build loong64

package runtime

// Called from assembly only; declared for go vet.
func load_g()
func save_g()

// Used by reflectcall and the reflect package.
//
// Spills/loads arguments in registers to/from an internal/abi.RegArgs
// respectively. Does not follow the Go ABI.
func spillArgs()
func unspillArgs()

// getfp returns the frame pointer register of its caller or 0 if not implemented.
// TODO: Make this a compiler intrinsic
func getfp() uintptr { return 0 }
