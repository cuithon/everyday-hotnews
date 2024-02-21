// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains various functionality that is
// different between go/types and types2. Factoring
// out this code allows more of the rest of the code
// to be shared.

package types

import (
	"go/ast"
	"go/token"
)

const isTypes2 = false

// cmpPos compares the positions p and q and returns a result r as follows:
//
// r <  0: p is before q
// r == 0: p and q are the same position (but may not be identical)
// r >  0: p is after q
//
// If p and q are in different files, p is before q if the filename
// of p sorts lexicographically before the filename of q.
func cmpPos(p, q token.Pos) int { return int(p - q) }

// hasDots reports whether the last argument in the call is followed by ...
func hasDots(call *ast.CallExpr) bool { return call.Ellipsis.IsValid() }

// dddErrPos returns the positioner for reporting an invalid ... use in a call.
func dddErrPos(call *ast.CallExpr) positioner { return atPos(call.Ellipsis) }

// argErrPos returns positioner for reportign an invalid argument count.
func argErrPos(call *ast.CallExpr) positioner { return inNode(call, call.Rparen) }

// endPos returns the position of the first character immediately after node n.
func endPos(n ast.Node) token.Pos { return n.End() }
