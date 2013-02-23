// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains the printf-checker.

package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
	"unicode/utf8"
)

var printfuncs = flag.String("printfuncs", "", "comma-separated list of print function names to check")

// printfList records the formatted-print functions. The value is the location
// of the format parameter. Names are lower-cased so the lookup is
// case insensitive.
var printfList = map[string]int{
	"errorf":  0,
	"fatalf":  0,
	"fprintf": 1,
	"panicf":  0,
	"printf":  0,
	"sprintf": 0,
}

// printList records the unformatted-print functions. The value is the location
// of the first parameter to be printed.  Names are lower-cased so the lookup is
// case insensitive.
var printList = map[string]int{
	"error":  0,
	"fatal":  0,
	"fprint": 1, "fprintln": 1,
	"panic": 0, "panicln": 0,
	"print": 0, "println": 0,
	"sprint": 0, "sprintln": 0,
}

// checkCall triggers the print-specific checks if the call invokes a print function.
func (f *File) checkFmtPrintfCall(call *ast.CallExpr, Name string) {
	if !vet("printf") {
		return
	}
	name := strings.ToLower(Name)
	if skip, ok := printfList[name]; ok {
		f.checkPrintf(call, Name, skip)
		return
	}
	if skip, ok := printList[name]; ok {
		f.checkPrint(call, Name, skip)
		return
	}
}

// literal returns the literal value represented by the expression, or nil if it is not a literal.
func (f *File) literal(value ast.Expr) *ast.BasicLit {
	switch v := value.(type) {
	case *ast.BasicLit:
		return v
	case *ast.ParenExpr:
		return f.literal(v.X)
	case *ast.BinaryExpr:
		if v.Op != token.ADD {
			break
		}
		litX := f.literal(v.X)
		litY := f.literal(v.Y)
		if litX != nil && litY != nil {
			lit := *litX
			x, errX := strconv.Unquote(litX.Value)
			y, errY := strconv.Unquote(litY.Value)
			if errX == nil && errY == nil {
				return &ast.BasicLit{
					ValuePos: lit.ValuePos,
					Kind:     lit.Kind,
					Value:    strconv.Quote(x + y),
				}
			}
		}
	case *ast.Ident:
		// See if it's a constant or initial value (we can't tell the difference).
		if v.Obj == nil || v.Obj.Decl == nil {
			return nil
		}
		valueSpec, ok := v.Obj.Decl.(*ast.ValueSpec)
		if ok && len(valueSpec.Names) == len(valueSpec.Values) {
			// Find the index in the list of names
			var i int
			for i = 0; i < len(valueSpec.Names); i++ {
				if valueSpec.Names[i].Name == v.Name {
					if lit, ok := valueSpec.Values[i].(*ast.BasicLit); ok {
						return lit
					}
					return nil
				}
			}
		}
	}
	return nil
}

// checkPrintf checks a call to a formatted print routine such as Printf.
// call.Args[formatIndex] is (well, should be) the format argument.
func (f *File) checkPrintf(call *ast.CallExpr, name string, formatIndex int) {
	if formatIndex >= len(call.Args) {
		return
	}
	lit := f.literal(call.Args[formatIndex])
	if lit == nil {
		if *verbose {
			f.Warn(call.Pos(), "can't check non-literal format in call to", name)
		}
		return
	}
	if lit.Kind != token.STRING {
		f.Badf(call.Pos(), "literal %v not a string in call to", lit.Value, name)
	}
	format, err := strconv.Unquote(lit.Value)
	if err != nil {
		// Shouldn't happen if parser returned no errors, but be safe.
		f.Badf(call.Pos(), "invalid quoted string literal")
	}
	firstArg := formatIndex + 1 // Arguments are immediately after format string.
	if !strings.Contains(format, "%") {
		if len(call.Args) > firstArg {
			f.Badf(call.Pos(), "no formatting directive in %s call", name)
		}
		return
	}
	// Hard part: check formats against args.
	argNum := firstArg
	for i, w := 0, 0; i < len(format); i += w {
		w = 1
		if format[i] == '%' {
			verb, flags, nbytes, nargs := f.parsePrintfVerb(call, format[i:])
			w = nbytes
			if verb == '%' { // "%%" does nothing interesting.
				continue
			}
			// If we've run out of args, print after loop will pick that up.
			if argNum+nargs <= len(call.Args) {
				f.checkPrintfArg(call, verb, flags, argNum, nargs)
			}
			argNum += nargs
		}
	}
	// TODO: Dotdotdot is hard.
	if call.Ellipsis.IsValid() && argNum != len(call.Args) {
		return
	}
	if argNum != len(call.Args) {
		expect := argNum - firstArg
		numArgs := len(call.Args) - firstArg
		f.Badf(call.Pos(), "wrong number of args for format in %s call: %d needed but %d args", name, expect, numArgs)
	}
}

// parsePrintfVerb returns the verb that begins the format string, along with its flags,
// the number of bytes to advance the format to step past the verb, and number of
// arguments it consumes.
func (f *File) parsePrintfVerb(call *ast.CallExpr, format string) (verb rune, flags []byte, nbytes, nargs int) {
	// There's guaranteed a percent sign.
	flags = make([]byte, 0, 5)
	nbytes = 1
	end := len(format)
	// There may be flags.
FlagLoop:
	for nbytes < end {
		switch format[nbytes] {
		case '#', '0', '+', '-', ' ':
			flags = append(flags, format[nbytes])
			nbytes++
		default:
			break FlagLoop
		}
	}
	getNum := func() {
		if nbytes < end && format[nbytes] == '*' {
			nbytes++
			nargs++
		} else {
			for nbytes < end && '0' <= format[nbytes] && format[nbytes] <= '9' {
				nbytes++
			}
		}
	}
	// There may be a width.
	getNum()
	// If there's a period, there may be a precision.
	if nbytes < end && format[nbytes] == '.' {
		flags = append(flags, '.') // Treat precision as a flag.
		nbytes++
		getNum()
	}
	// Now a verb.
	c, w := utf8.DecodeRuneInString(format[nbytes:])
	nbytes += w
	verb = c
	if c != '%' {
		nargs++
	}
	return
}

// printfArgType encodes the types of expressions a printf verb accepts. It is a bitmask.
type printfArgType int

const (
	argBool printfArgType = 1 << iota
	argInt
	argRune
	argString
	argFloat
	argPointer
	anyType printfArgType = ^0
)

type printVerb struct {
	verb  rune
	flags string // known flags are all ASCII
	typ   printfArgType
}

// Common flag sets for printf verbs.
const (
	numFlag      = " -+.0"
	sharpNumFlag = " -+.0#"
	allFlags     = " -+.0#"
)

// printVerbs identifies which flags are known to printf for each verb.
// TODO: A type that implements Formatter may do what it wants, and vet
// will complain incorrectly.
var printVerbs = []printVerb{
	// '-' is a width modifier, always valid.
	// '.' is a precision for float, max width for strings.
	// '+' is required sign for numbers, Go format for %v.
	// '#' is alternate format for several verbs.
	// ' ' is spacer for numbers
	{'b', numFlag, argInt},
	{'c', "-", argRune | argInt},
	{'d', numFlag, argInt},
	{'e', numFlag, argFloat},
	{'E', numFlag, argFloat},
	{'f', numFlag, argFloat},
	{'F', numFlag, argFloat},
	{'g', numFlag, argFloat},
	{'G', numFlag, argFloat},
	{'o', sharpNumFlag, argInt},
	{'p', "-#", argPointer},
	{'q', " -+.0#", argRune | argInt | argString},
	{'s', " -+.0", argString},
	{'t', "-", argBool},
	{'T', "-", anyType},
	{'U', "-#", argRune | argInt},
	{'v', allFlags, anyType},
	{'x', sharpNumFlag, argRune | argInt | argString},
	{'X', sharpNumFlag, argRune | argInt | argString},
}

const printfVerbs = "bcdeEfFgGopqstTvxUX"

func (f *File) checkPrintfArg(call *ast.CallExpr, verb rune, flags []byte, argNum, nargs int) {
	// Linear scan is fast enough for a small list.
	for _, v := range printVerbs {
		if v.verb == verb {
			for _, flag := range flags {
				if !strings.ContainsRune(v.flags, rune(flag)) {
					f.Badf(call.Pos(), "unrecognized printf flag for verb %q: %q", verb, flag)
					return
				}
			}
			// Verb is good. If nargs>1, we have something like %.*s and all but the final
			// arg must be integer.
			for i := 0; i < nargs-1; i++ {
				if !f.matchArgType(argInt, call.Args[argNum+i]) {
					f.Badf(call.Pos(), "arg for * in printf format not of type int")
				}
			}
			for _, v := range printVerbs {
				if v.verb == verb {
					if !f.matchArgType(v.typ, call.Args[argNum+nargs-1]) {
						f.Badf(call.Pos(), "arg for printf verb %%%c of wrong type", verb)
					}
					break
				}
			}
			return
		}
	}
	f.Badf(call.Pos(), "unrecognized printf verb %q", verb)
}

func (f *File) matchArgType(t printfArgType, arg ast.Expr) bool {
	if f.pkg == nil {
		return true // Don't know; assume OK.
	}
	// TODO: for now, we can only test builtin types and untyped constants.
	typ := f.pkg.types[arg]
	if typ == nil {
		return true
	}
	basic, ok := typ.(*types.Basic)
	if !ok {
		return true
	}
	switch basic.Kind {
	case types.Bool:
		return t&argBool != 0
	case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
		fallthrough
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
		return t&argInt != 0
	case types.Float32, types.Float64, types.Complex64, types.Complex128:
		return t&argFloat != 0
	case types.String:
		return t&argString != 0
	case types.UnsafePointer:
		return t&argPointer != 0
	case types.UntypedBool:
		return t&argBool != 0
	case types.UntypedComplex:
		return t&argFloat != 0
	case types.UntypedFloat:
		// If it's integral, we can use an int format.
		switch f.pkg.values[arg].(type) {
		case int, int8, int16, int32, int64:
			return t&(argInt|argFloat) != 0
		case uint, uint8, uint16, uint32, uint64:
			return t&(argInt|argFloat) != 0
		}
		return t&argFloat != 0
	case types.UntypedInt:
		return t&(argInt|argFloat) != 0 // You might say Printf("%g", 1234)
	case types.UntypedRune:
		return t&(argInt|argRune) != 0
	case types.UntypedString:
		return t&argString != 0
	case types.UntypedNil:
		return t&argPointer != 0 // TODO?
	}
	return false
}

// checkPrint checks a call to an unformatted print routine such as Println.
// call.Args[firstArg] is the first argument to be printed.
func (f *File) checkPrint(call *ast.CallExpr, name string, firstArg int) {
	isLn := strings.HasSuffix(name, "ln")
	isF := strings.HasPrefix(name, "F")
	args := call.Args
	// check for Println(os.Stderr, ...)
	if firstArg == 0 && !isF && len(args) > 0 {
		if sel, ok := args[0].(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok {
				if x.Name == "os" && strings.HasPrefix(sel.Sel.Name, "Std") {
					f.Warnf(call.Pos(), "first argument to %s is %s.%s", name, x.Name, sel.Sel.Name)
				}
			}
		}
	}
	if len(args) <= firstArg {
		// If we have a call to a method called Error that satisfies the Error interface,
		// then it's ok. Otherwise it's something like (*T).Error from the testing package
		// and we need to check it.
		if name == "Error" && f.pkg != nil && f.isErrorMethodCall(call) {
			return
		}
		// If it's an Error call now, it's probably for printing errors.
		if !isLn {
			// Check the signature to be sure: there are niladic functions called "error".
			if f.pkg == nil || firstArg != 0 || f.numArgsInSignature(call) != firstArg {
				f.Badf(call.Pos(), "no args in %s call", name)
			}
		}
		return
	}
	arg := args[firstArg]
	if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		if strings.Contains(lit.Value, "%") {
			f.Badf(call.Pos(), "possible formatting directive in %s call", name)
		}
	}
	if isLn {
		// The last item, if a string, should not have a newline.
		arg = args[len(call.Args)-1]
		if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if strings.HasSuffix(lit.Value, `\n"`) {
				f.Badf(call.Pos(), "%s call ends with newline", name)
			}
		}
	}
}

// numArgsInSignature tells how many formal arguments the function type
// being called has. Assumes type checking is on (f.pkg != nil).
func (f *File) numArgsInSignature(call *ast.CallExpr) int {
	// Check the type of the function or method declaration
	typ := f.pkg.types[call.Fun]
	if typ == nil {
		return 0
	}
	// The type must be a signature, but be sure for safety.
	sig, ok := typ.(*types.Signature)
	if !ok {
		return 0
	}
	return len(sig.Params)
}

// isErrorMethodCall reports whether the call is of a method with signature
//	func Error() error
// where "error" is the universe's error type. We know the method is called "Error"
// and f.pkg is set.
func (f *File) isErrorMethodCall(call *ast.CallExpr) bool {
	// Is it a selector expression? Otherwise it's a function call, not a method call.
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	// The package is type-checked, so if there are no arguments, we're done.
	if len(call.Args) > 0 {
		return false
	}
	// Check the type of the method declaration
	typ := f.pkg.types[sel]
	if typ == nil {
		return false
	}
	// The type must be a signature, but be sure for safety.
	sig, ok := typ.(*types.Signature)
	if !ok {
		return false
	}
	// There must be a receiver for it to be a method call. Otherwise it is
	// a function, not something that satisfies the error interface.
	if sig.Recv == nil {
		return false
	}
	// There must be no arguments. Already verified by type checking, but be thorough.
	if len(sig.Params) > 0 {
		return false
	}
	// Finally the real questions.
	// There must be one result.
	if len(sig.Results) != 1 {
		return false
	}
	// It must have return type "string" from the universe.
	result := sig.Results[0].Type
	if types.IsIdentical(result, types.Typ[types.String]) {
		return true
	}
	return true
}

// Error methods that do not satisfy the Error interface and should be checked.
type errorTest1 int

func (errorTest1) Error(...interface{}) string {
	return "hi"
}

type errorTest2 int // Analogous to testing's *T type.
func (errorTest2) Error(...interface{}) {
}

type errorTest3 int

func (errorTest3) Error() { // No return value.
}

type errorTest4 int

func (errorTest4) Error() int { // Different return type.
	return 3
}

type errorTest5 int

func (errorTest5) error() { // niladic; don't complain if no args (was bug)
}

// This function never executes, but it serves as a simple test for the program.
// Test with make test.
func BadFunctionUsedInTests() {
	var b bool
	var i int
	var r rune
	var s string
	var x float64
	var p *int
	// Some good format/argtypes
	fmt.Printf("")
	fmt.Printf("%b %b", 3, i)
	fmt.Printf("%c %c %c %c", 3, i, 'x', r)
	fmt.Printf("%d %d", 3, i)
	fmt.Printf("%e %e %e", 3, 3e9, x)
	fmt.Printf("%E %E %E", 3, 3e9, x)
	fmt.Printf("%f %f %f", 3, 3e9, x)
	fmt.Printf("%F %F %F", 3, 3e9, x)
	fmt.Printf("%g %g %g", 3, 3e9, x)
	fmt.Printf("%G %G %G", 3, 3e9, x)
	fmt.Printf("%o %o", 3, i)
	fmt.Printf("%p %p", p, nil)
	fmt.Printf("%q %q %q %q", 3, i, 'x', r)
	fmt.Printf("%s %s", "hi", s)
	fmt.Printf("%t %t", true, b)
	fmt.Printf("%T %T", 3, i)
	fmt.Printf("%U %U", 3, i)
	fmt.Printf("%v %v", 3, i)
	fmt.Printf("%x %x %x %x", 3, i, "hi", s)
	fmt.Printf("%X %X %X %X", 3, i, "hi", s)
	fmt.Printf("%.*s %d %g", 3, "hi", 23, 2.3)
	// Some bad format/argTypes
	fmt.Printf("%b", 2.3)                      // ERROR "arg for printf verb %b of wrong type"
	fmt.Printf("%c", 2.3)                      // ERROR "arg for printf verb %c of wrong type"
	fmt.Printf("%d", 2.3)                      // ERROR "arg for printf verb %d of wrong type"
	fmt.Printf("%e", "hi")                     // ERROR "arg for printf verb %e of wrong type"
	fmt.Printf("%E", true)                     // ERROR "arg for printf verb %E of wrong type"
	fmt.Printf("%f", "hi")                     // ERROR "arg for printf verb %f of wrong type"
	fmt.Printf("%F", 'x')                      // ERROR "arg for printf verb %F of wrong type"
	fmt.Printf("%g", "hi")                     // ERROR "arg for printf verb %g of wrong type"
	fmt.Printf("%G", i)                        // ERROR "arg for printf verb %G of wrong type"
	fmt.Printf("%o", x)                        // ERROR "arg for printf verb %o of wrong type"
	fmt.Printf("%p", 23)                       // ERROR "arg for printf verb %p of wrong type"
	fmt.Printf("%q", x)                        // ERROR "arg for printf verb %q of wrong type"
	fmt.Printf("%s", b)                        // ERROR "arg for printf verb %s of wrong type"
	fmt.Printf("%t", 23)                       // ERROR "arg for printf verb %t of wrong type"
	fmt.Printf("%U", x)                        // ERROR "arg for printf verb %U of wrong type"
	fmt.Printf("%x", nil)                      // ERROR "arg for printf verb %x of wrong type"
	fmt.Printf("%X", 2.3)                      // ERROR "arg for printf verb %X of wrong type"
	fmt.Printf("%.*s %d %g", 3, "hi", 23, 'x') // ERROR "arg for printf verb %g of wrong type"
	// TODO
	fmt.Println()                      // not an error
	fmt.Println("%s", "hi")            // ERROR "possible formatting directive in Println call"
	fmt.Printf("%s", "hi", 3)          // ERROR "wrong number of args for format in Printf call"
	fmt.Printf("%"+("s"), "hi", 3)     // ERROR "wrong number of args for format in Printf call"
	fmt.Printf("%s%%%d", "hi", 3)      // correct
	fmt.Printf("%08s", "woo")          // correct
	fmt.Printf("% 8s", "woo")          // correct
	fmt.Printf("%.*d", 3, 3)           // correct
	fmt.Printf("%.*d", 3, 3, 3)        // ERROR "wrong number of args for format in Printf call"
	fmt.Printf("%.*d", "hi", 3)        // ERROR "arg for \* in printf format not of type int"
	fmt.Printf("%.*d", i, 3)           // correct
	fmt.Printf("%.*d", s, 3)           // ERROR "arg for \* in printf format not of type int"
	fmt.Printf("%q %q", multi()...)    // ok
	fmt.Printf("%#q", `blah`)          // ok
	printf("now is the time", "buddy") // ERROR "no formatting directive"
	Printf("now is the time", "buddy") // ERROR "no formatting directive"
	Printf("hi")                       // ok
	const format = "%s %s\n"
	Printf(format, "hi", "there")
	Printf(format, "hi") // ERROR "wrong number of args for format in Printf call"
	f := new(File)
	f.Warn(0, "%s", "hello", 3)  // ERROR "possible formatting directive in Warn call"
	f.Warnf(0, "%s", "hello", 3) // ERROR "wrong number of args for format in Warnf call"
	f.Warnf(0, "%r", "hello")    // ERROR "unrecognized printf verb"
	f.Warnf(0, "%#s", "hello")   // ERROR "unrecognized printf flag"
	// Something that satisfies the error interface.
	var e error
	fmt.Println(e.Error()) // ok
	// Something that looks like an error interface but isn't, such as the (*T).Error method
	// in the testing package.
	var et1 errorTest1
	fmt.Println(et1.Error())        // ERROR "no args in Error call"
	fmt.Println(et1.Error("hi"))    // ok
	fmt.Println(et1.Error("%d", 3)) // ERROR "possible formatting directive in Error call"
	var et2 errorTest2
	et2.Error()        // ERROR "no args in Error call"
	et2.Error("hi")    // ok, not an error method.
	et2.Error("%d", 3) // ERROR "possible formatting directive in Error call"
	var et3 errorTest3
	et3.Error() // ok, not an error method.
	var et4 errorTest4
	et4.Error() // ok, not an error method.
	var et5 errorTest5
	et5.error() // ok, not an error method.
}

// printf is used by the test.
func printf(format string, args ...interface{}) {
	panic("don't call - testing only")
}

// multi is used by the test.
func multi() []interface{} {
	panic("don't call - testing only")
}
