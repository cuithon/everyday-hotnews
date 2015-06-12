// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains some utility functions to help define Funcs for testing.
// As an example, the following func
//
//   b1:
//     v1 = Arg <mem> [.mem]
//     Plain -> b2
//   b2:
//     Exit v1
//   b3:
//     v2 = Const <bool> [true]
//     If v2 -> b3 b2
//
// can be defined as
//
//   fun := Fun("entry",
//       Bloc("entry",
//           Valu("mem", OpArg, TypeMem, 0, ".mem"),
//           Goto("exit")),
//       Bloc("exit",
//           Exit("mem")),
//       Bloc("deadblock",
//          Valu("deadval", OpConst, TypeBool, 0, true),
//          If("deadval", "deadblock", "exit")))
//
// and the Blocks or Values used in the Func can be accessed
// like this:
//   fun.blocks["entry"] or fun.values["deadval"]

package ssa

// TODO(matloob): Choose better names for Fun, Bloc, Goto, etc.
// TODO(matloob): Write a parser for the Func disassembly. Maybe
//                the parser can be used instead of Fun.

import (
	"log"
	"reflect"
	"testing"
)

// Compare two Funcs for equivalence. Their CFGs must be isomorphic,
// and their values must correspond.
// Requires that values and predecessors are in the same order, even
// though Funcs could be equivalent when they are not.
// TODO(matloob): Allow values and predecessors to be in different
// orders if the CFG are otherwise equivalent.
func Equiv(f, g *Func) bool {
	valcor := make(map[*Value]*Value)
	var checkVal func(fv, gv *Value) bool
	checkVal = func(fv, gv *Value) bool {
		if fv == nil && gv == nil {
			return true
		}
		if valcor[fv] == nil && valcor[gv] == nil {
			valcor[fv] = gv
			valcor[gv] = fv
			// Ignore ids. Ops and Types are compared for equality.
			// TODO(matloob): Make sure types are canonical and can
			// be compared for equality.
			if fv.Op != gv.Op || fv.Type != gv.Type || fv.AuxInt != gv.AuxInt {
				return false
			}
			if !reflect.DeepEqual(fv.Aux, gv.Aux) {
				// This makes the assumption that aux values can be compared
				// using DeepEqual.
				// TODO(matloob): Aux values may be *gc.Sym pointers in the near
				// future. Make sure they are canonical.
				return false
			}
			if len(fv.Args) != len(gv.Args) {
				return false
			}
			for i := range fv.Args {
				if !checkVal(fv.Args[i], gv.Args[i]) {
					return false
				}
			}
		}
		return valcor[fv] == gv && valcor[gv] == fv
	}
	blkcor := make(map[*Block]*Block)
	var checkBlk func(fb, gb *Block) bool
	checkBlk = func(fb, gb *Block) bool {
		if blkcor[fb] == nil && blkcor[gb] == nil {
			blkcor[fb] = gb
			blkcor[gb] = fb
			// ignore ids
			if fb.Kind != gb.Kind {
				return false
			}
			if len(fb.Values) != len(gb.Values) {
				return false
			}
			for i := range fb.Values {
				if !checkVal(fb.Values[i], gb.Values[i]) {
					return false
				}
			}
			if len(fb.Succs) != len(gb.Succs) {
				return false
			}
			for i := range fb.Succs {
				if !checkBlk(fb.Succs[i], gb.Succs[i]) {
					return false
				}
			}
			if len(fb.Preds) != len(gb.Preds) {
				return false
			}
			for i := range fb.Preds {
				if !checkBlk(fb.Preds[i], gb.Preds[i]) {
					return false
				}
			}
			return true

		}
		return blkcor[fb] == gb && blkcor[gb] == fb
	}

	return checkBlk(f.Entry, g.Entry)
}

// fun is the return type of Fun. It contains the created func
// itself as well as indexes from block and value names into the
// corresponding Blocks and Values.
type fun struct {
	f      *Func
	blocks map[string]*Block
	values map[string]*Value
}

// Fun takes the name of an entry bloc and a series of Bloc calls, and
// returns a fun containing the composed Func. entry must be a name
// supplied to one of the Bloc functions. Each of the bloc names and
// valu names should be unique across the Fun.
func Fun(c *Config, entry string, blocs ...bloc) fun {
	f := new(Func)
	f.Config = c
	blocks := make(map[string]*Block)
	values := make(map[string]*Value)
	// Create all the blocks and values.
	for _, bloc := range blocs {
		b := f.NewBlock(bloc.control.kind)
		blocks[bloc.name] = b
		for _, valu := range bloc.valus {
			// args are filled in the second pass.
			values[valu.name] = b.NewValue0IA(0, valu.op, valu.t, valu.auxint, valu.aux)
		}
	}
	// Connect the blocks together and specify control values.
	f.Entry = blocks[entry]
	for _, bloc := range blocs {
		b := blocks[bloc.name]
		c := bloc.control
		// Specify control values.
		if c.control != "" {
			cval, ok := values[c.control]
			if !ok {
				log.Panicf("control value for block %s missing", bloc.name)
			}
			b.Control = cval
		}
		// Fill in args.
		for _, valu := range bloc.valus {
			v := values[valu.name]
			for _, arg := range valu.args {
				a, ok := values[arg]
				if !ok {
					log.Panicf("arg %s missing for value %s in block %s",
						arg, valu.name, bloc.name)
				}
				v.AddArg(a)
			}
		}
		// Connect to successors.
		for _, succ := range c.succs {
			addEdge(b, blocks[succ])
		}
	}
	return fun{f, blocks, values}
}

// Bloc defines a block for Fun. The bloc name should be unique
// across the containing Fun. entries should consist of calls to valu,
// as well as one call to Goto, If, or Exit to specify the block kind.
func Bloc(name string, entries ...interface{}) bloc {
	b := bloc{}
	b.name = name
	seenCtrl := false
	for _, e := range entries {
		switch v := e.(type) {
		case ctrl:
			// there should be exactly one Ctrl entry.
			if seenCtrl {
				log.Panicf("already seen control for block %s", name)
			}
			b.control = v
			seenCtrl = true
		case valu:
			b.valus = append(b.valus, v)
		}
	}
	if !seenCtrl {
		log.Panicf("block %s doesn't have control", b.name)
	}
	return b
}

// Valu defines a value in a block.
func Valu(name string, op Op, t Type, auxint int64, aux interface{}, args ...string) valu {
	return valu{name, op, t, auxint, aux, args}
}

// Goto specifies that this is a BlockPlain and names the single successor.
// TODO(matloob): choose a better name.
func Goto(succ string) ctrl {
	return ctrl{BlockPlain, "", []string{succ}}
}

// If specifies a BlockIf.
func If(cond, sub, alt string) ctrl {
	return ctrl{BlockIf, cond, []string{sub, alt}}
}

// Exit specifies a BlockExit.
func Exit(arg string) ctrl {
	return ctrl{BlockExit, arg, []string{}}
}

// bloc, ctrl, and valu are internal structures used by Bloc, Valu, Goto,
// If, and Exit to help define blocks.

type bloc struct {
	name    string
	control ctrl
	valus   []valu
}

type ctrl struct {
	kind    BlockKind
	control string
	succs   []string
}

type valu struct {
	name   string
	op     Op
	t      Type
	auxint int64
	aux    interface{}
	args   []string
}

func addEdge(b, c *Block) {
	b.Succs = append(b.Succs, c)
	c.Preds = append(c.Preds, b)
}

func TestArgs(t *testing.T) {
	c := NewConfig("amd64", DummyFrontend{})
	fun := Fun(c, "entry",
		Bloc("entry",
			Valu("a", OpConst, TypeInt64, 14, nil),
			Valu("b", OpConst, TypeInt64, 26, nil),
			Valu("sum", OpAdd, TypeInt64, 0, nil, "a", "b"),
			Valu("mem", OpArg, TypeMem, 0, ".mem"),
			Goto("exit")),
		Bloc("exit",
			Exit("mem")))
	sum := fun.values["sum"]
	for i, name := range []string{"a", "b"} {
		if sum.Args[i] != fun.values[name] {
			t.Errorf("arg %d for sum is incorrect: want %s, got %s",
				i, sum.Args[i], fun.values[name])
		}
	}
}

func TestEquiv(t *testing.T) {
	c := NewConfig("amd64", DummyFrontend{})
	equivalentCases := []struct{ f, g fun }{
		// simple case
		{
			Fun(c, "entry",
				Bloc("entry",
					Valu("a", OpConst, TypeInt64, 14, nil),
					Valu("b", OpConst, TypeInt64, 26, nil),
					Valu("sum", OpAdd, TypeInt64, 0, nil, "a", "b"),
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Goto("exit")),
				Bloc("exit",
					Exit("mem"))),
			Fun(c, "entry",
				Bloc("entry",
					Valu("a", OpConst, TypeInt64, 14, nil),
					Valu("b", OpConst, TypeInt64, 26, nil),
					Valu("sum", OpAdd, TypeInt64, 0, nil, "a", "b"),
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Goto("exit")),
				Bloc("exit",
					Exit("mem"))),
		},
		// block order changed
		{
			Fun(c, "entry",
				Bloc("entry",
					Valu("a", OpConst, TypeInt64, 14, nil),
					Valu("b", OpConst, TypeInt64, 26, nil),
					Valu("sum", OpAdd, TypeInt64, 0, nil, "a", "b"),
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Goto("exit")),
				Bloc("exit",
					Exit("mem"))),
			Fun(c, "entry",
				Bloc("exit",
					Exit("mem")),
				Bloc("entry",
					Valu("a", OpConst, TypeInt64, 14, nil),
					Valu("b", OpConst, TypeInt64, 26, nil),
					Valu("sum", OpAdd, TypeInt64, 0, nil, "a", "b"),
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Goto("exit"))),
		},
	}
	for _, c := range equivalentCases {
		if !Equiv(c.f.f, c.g.f) {
			t.Error("expected equivalence. Func definitions:")
			t.Error(c.f.f)
			t.Error(c.g.f)
		}
	}

	differentCases := []struct{ f, g fun }{
		// different shape
		{
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Goto("exit")),
				Bloc("exit",
					Exit("mem"))),
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Exit("mem"))),
		},
		// value order changed
		{
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Valu("b", OpConst, TypeInt64, 26, nil),
					Valu("a", OpConst, TypeInt64, 14, nil),
					Exit("mem"))),
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Valu("a", OpConst, TypeInt64, 14, nil),
					Valu("b", OpConst, TypeInt64, 26, nil),
					Exit("mem"))),
		},
		// value auxint different
		{
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Valu("a", OpConst, TypeInt64, 14, nil),
					Exit("mem"))),
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Valu("a", OpConst, TypeInt64, 26, nil),
					Exit("mem"))),
		},
		// value aux different
		{
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Valu("a", OpConst, TypeInt64, 0, 14),
					Exit("mem"))),
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Valu("a", OpConst, TypeInt64, 0, 26),
					Exit("mem"))),
		},
		// value args different
		{
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Valu("a", OpConst, TypeInt64, 14, nil),
					Valu("b", OpConst, TypeInt64, 26, nil),
					Valu("sum", OpAdd, TypeInt64, 0, nil, "a", "b"),
					Exit("mem"))),
			Fun(c, "entry",
				Bloc("entry",
					Valu("mem", OpArg, TypeMem, 0, ".mem"),
					Valu("a", OpConst, TypeInt64, 0, nil),
					Valu("b", OpConst, TypeInt64, 14, nil),
					Valu("sum", OpAdd, TypeInt64, 0, nil, "b", "a"),
					Exit("mem"))),
		},
	}
	for _, c := range differentCases {
		if Equiv(c.f.f, c.g.f) {
			t.Error("expected difference. Func definitions:")
			t.Error(c.f.f)
			t.Error(c.g.f)
		}
	}
}

// opcodeMap returns a map from opcode to the number of times that opcode
// appears in the function.
func opcodeMap(f *Func) map[Op]int {
	m := map[Op]int{}
	for _, b := range f.Blocks {
		for _, v := range b.Values {
			m[v.Op]++
		}
	}
	return m
}

// opcodeCounts checks that the number of opcodes listed in m agree with the
// number of opcodes that appear in the function.
func checkOpcodeCounts(t *testing.T, f *Func, m map[Op]int) {
	n := opcodeMap(f)
	for op, cnt := range m {
		if n[op] != cnt {
			t.Errorf("%s appears %d times, want %d times", op, n[op], cnt)
		}
	}
}
