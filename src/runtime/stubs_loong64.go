// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build loong64

package runtime

// Called from assembly only; declared for go vet.
func load_g()
func save_g()

// getcallerfp returns the address of the frame pointer in the callers frame or 0 if not implemented.
func getcallerfp() uintptr { return 0 }
