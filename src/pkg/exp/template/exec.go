// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package template

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// state represents the state of an execution. It's not part of the
// template so that multiple executions of the same template
// can execute in parallel.
type state struct {
	tmpl *Template
	wr   io.Writer
	set  *Set
	line int // line number for errors
}

// errorf formats the error and terminates processing.
func (s *state) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("template: %s:%d: %s", s.tmpl.name, s.line, format)
	panic(fmt.Errorf(format, args...))
}

// error terminates processing.
func (s *state) error(err os.Error) {
	s.errorf("%s", err)
}

// Execute applies a parsed template to the specified data object,
// writing the output to wr.
func (t *Template) Execute(wr io.Writer, data interface{}) os.Error {
	return t.ExecuteInSet(wr, data, nil)
}

// ExecuteInSet applies a parsed template to the specified data object,
// writing the output to wr. Nested template invocations will be resolved
// from the specified set.
func (t *Template) ExecuteInSet(wr io.Writer, data interface{}, set *Set) (err os.Error) {
	defer t.recover(&err)
	state := &state{
		tmpl: t,
		wr:   wr,
		set:  set,
		line: 1,
	}
	if t.root == nil {
		state.errorf("must be parsed before execution")
	}
	state.walk(reflect.ValueOf(data), t.root)
	return
}

// Walk functions step through the major pieces of the template structure,
// generating output as they go.
func (s *state) walk(data reflect.Value, n node) {
	switch n := n.(type) {
	case *actionNode:
		s.line = n.line
		s.printValue(n, s.evalPipeline(data, n.pipeline))
	case *listNode:
		for _, node := range n.nodes {
			s.walk(data, node)
		}
	case *ifNode:
		s.walkIfOrWith(nodeIf, data, n.pipeline, n.list, n.elseList)
	case *rangeNode:
		s.walkRange(data, n)
	case *textNode:
		if _, err := s.wr.Write(n.text); err != nil {
			s.error(err)
		}
	case *templateNode:
		s.walkTemplate(data, n)
	case *withNode:
		s.walkIfOrWith(nodeWith, data, n.pipeline, n.list, n.elseList)
	default:
		s.errorf("unknown node: %s", n)
	}
}

// walkIfOrWith walks an 'if' or 'with' node. The two control structures
// are identical in behavior except that 'with' sets dot.
func (s *state) walkIfOrWith(typ nodeType, data reflect.Value, pipe []*commandNode, list, elseList *listNode) {
	val := s.evalPipeline(data, pipe)
	truth := false
	switch val.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		truth = val.Len() > 0
	case reflect.Bool:
		truth = val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		truth = val.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		truth = val.Uint() != 0
	case reflect.Float32, reflect.Float64:
		truth = val.Float() != 0
	case reflect.Complex64, reflect.Complex128:
		truth = val.Complex() != 0
	case reflect.Chan, reflect.Func, reflect.Ptr:
		truth = !val.IsNil()
	default:
		s.errorf("if/with can't use value of type %T", val.Interface())
	}
	if truth {
		if typ == nodeWith {
			data = val
		}
		s.walk(data, list)
	} else if elseList != nil {
		s.walk(data, elseList)
	}
}

func (s *state) walkRange(data reflect.Value, r *rangeNode) {
	val := s.evalPipeline(data, r.pipeline)
	switch val.Kind() {
	case reflect.Array, reflect.Slice:
		if val.Len() == 0 {
			break
		}
		for i := 0; i < val.Len(); i++ {
			s.walk(val.Index(i), r.list)
		}
		return
	case reflect.Map:
		if val.Len() == 0 {
			break
		}
		for _, key := range val.MapKeys() {
			s.walk(val.MapIndex(key), r.list)
		}
		return
	default:
		s.errorf("range can't iterate over value of type %T", val.Interface())
	}
	if r.elseList != nil {
		s.walk(data, r.elseList)
	}
}

func (s *state) walkTemplate(data reflect.Value, t *templateNode) {
	name := s.evalArg(data, reflect.TypeOf("string"), t.name).String()
	if s.set == nil {
		s.errorf("no set defined in which to invoke template named %q", name)
	}
	tmpl := s.set.tmpl[name]
	if tmpl == nil {
		s.errorf("template %q not in set", name)
	}
	data = s.evalPipeline(data, t.pipeline)
	newState := *s
	newState.tmpl = tmpl
	newState.walk(data, tmpl.root)
}

// Eval functions evaluate pipelines, commands, and their elements and extract
// values from the data structure by examining fields, calling methods, and so on.
// The printing of those values happens only through walk functions.

func (s *state) evalPipeline(data reflect.Value, pipe []*commandNode) reflect.Value {
	value := reflect.Value{}
	for _, cmd := range pipe {
		value = s.evalCommand(data, cmd, value) // previous value is this one's final arg.
	}
	return value
}

func (s *state) evalCommand(data reflect.Value, cmd *commandNode, final reflect.Value) reflect.Value {
	firstWord := cmd.args[0]
	if field, ok := firstWord.(*fieldNode); ok {
		return s.evalFieldNode(data, field, cmd.args, final)
	}
	if len(cmd.args) > 1 || final.IsValid() {
		// TODO: functions
		s.errorf("can't give argument to non-method %s", cmd.args[0])
	}
	switch word := cmd.args[0].(type) {
	case *dotNode:
		return data
	case *boolNode:
		return reflect.ValueOf(word.true)
	case *numberNode:
		// These are ideal constants but we don't know the type
		// and we have no context.  (If it was a method argument,
		// we'd know what we need.) The syntax guides us to some extent.
		switch {
		case word.isComplex:
			return reflect.ValueOf(word.complex128) // incontrovertible.
		case word.isFloat && strings.IndexAny(word.text, ".eE") >= 0:
			return reflect.ValueOf(word.float64)
		case word.isInt:
			return reflect.ValueOf(word.int64)
		case word.isUint:
			return reflect.ValueOf(word.uint64)
		}
	case *stringNode:
		return reflect.ValueOf(word.text)
	default:
		s.errorf("can't handle command %q", firstWord)
	}
	panic("not reached")
}

func (s *state) evalFieldNode(data reflect.Value, field *fieldNode, args []node, final reflect.Value) reflect.Value {
	// Up to the last entry, it must be a field.
	n := len(field.ident)
	for i := 0; i < n-1; i++ {
		data = s.evalField(data, field.ident[i])
	}
	// Now it can be a field or method and if a method, gets arguments.
	return s.evalMethodOrField(data, field.ident[n-1], args, final)
}

func (s *state) evalField(data reflect.Value, fieldName string) reflect.Value {
	for data.Kind() == reflect.Ptr {
		data = reflect.Indirect(data)
	}
	switch data.Kind() {
	case reflect.Struct:
		// Is it a field?
		field := data.FieldByName(fieldName)
		// TODO: look higher up the tree if we can't find it here. Also unexported fields
		// might succeed higher up, as map keys.
		if field.IsValid() && field.Type().PkgPath() == "" { // valid and exported
			return field
		}
		s.errorf("%s has no field %s", data.Type(), fieldName)
	default:
		s.errorf("can't evaluate field %s of type %s", fieldName, data.Type())
	}
	panic("not reached")
}

func (s *state) evalMethodOrField(data reflect.Value, fieldName string, args []node, final reflect.Value) reflect.Value {
	ptr := data
	for data.Kind() == reflect.Ptr {
		ptr, data = data, reflect.Indirect(data)
	}
	// Is it a method? We use the pointer because it has value methods too.
	if method, ok := ptr.Type().MethodByName(fieldName); ok {
		return s.evalMethod(ptr, method, args, final)
	}
	if len(args) > 1 || final.IsValid() {
		s.errorf("%s is not a method but has arguments", fieldName)
	}
	switch data.Kind() {
	case reflect.Struct:
		return s.evalField(data, fieldName)
	default:
		s.errorf("can't handle evaluation of field %s of type %s", fieldName, data.Type())
	}
	panic("not reached")
}

var (
	osErrorType = reflect.TypeOf(new(os.Error)).Elem()
)

func (s *state) evalMethod(v reflect.Value, method reflect.Method, args []node, final reflect.Value) reflect.Value {
	typ := method.Type
	fun := method.Func
	numIn := len(args)
	if final.IsValid() {
		numIn++
	}
	if !typ.IsVariadic() && numIn < typ.NumIn()-1 || !typ.IsVariadic() && numIn != typ.NumIn() {
		s.errorf("wrong number of args for %s: want %d got %d", method.Name, typ.NumIn(), len(args))
	}
	// We allow methods with 1 result or 2 results where the second is an os.Error.
	switch {
	case typ.NumOut() == 1:
	case typ.NumOut() == 2 && typ.Out(1) == osErrorType:
	default:
		s.errorf("can't handle multiple results from method %q", method.Name)
	}
	// Build the arg list.
	argv := make([]reflect.Value, numIn)
	// First arg is the receiver.
	argv[0] = v
	// Others must be evaluated.
	for i := 1; i < len(args); i++ {
		argv[i] = s.evalArg(v, typ.In(i), args[i])
	}
	// Add final value if necessary.
	if final.IsValid() {
		argv[len(args)] = final
	}
	result := fun.Call(argv)
	// If we have an os.Error that is not nil, stop execution and return that error to the caller.
	if len(result) == 2 && !result[1].IsNil() {
		s.error(result[1].Interface().(os.Error))
	}
	return result[0]
}

func (s *state) evalArg(data reflect.Value, typ reflect.Type, n node) reflect.Value {
	if field, ok := n.(*fieldNode); ok {
		value := s.evalFieldNode(data, field, []node{n}, reflect.Value{})
		if !value.Type().AssignableTo(typ) {
			s.errorf("wrong type for value; expected %s; got %s", typ, value.Type())
		}
		return value
	}
	switch typ.Kind() {
	case reflect.Bool:
		return s.evalBool(data, typ, n)
	case reflect.String:
		return s.evalString(data, typ, n)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return s.evalInteger(data, typ, n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return s.evalUnsignedInteger(data, typ, n)
	case reflect.Float32, reflect.Float64:
		return s.evalFloat(data, typ, n)
	case reflect.Complex64, reflect.Complex128:
		return s.evalComplex(data, typ, n)
	}
	s.errorf("can't handle node %s for method arg of type %s", n, typ)
	panic("not reached")
}

func (s *state) evalBool(v reflect.Value, typ reflect.Type, n node) reflect.Value {
	if n, ok := n.(*boolNode); ok {
		value := reflect.New(typ).Elem()
		value.SetBool(n.true)
		return value
	}
	s.errorf("expected bool; found %s", n)
	panic("not reached")
}

func (s *state) evalString(v reflect.Value, typ reflect.Type, n node) reflect.Value {
	if n, ok := n.(*stringNode); ok {
		value := reflect.New(typ).Elem()
		value.SetString(n.text)
		return value
	}
	s.errorf("expected string; found %s", n)
	panic("not reached")
}

func (s *state) evalInteger(v reflect.Value, typ reflect.Type, n node) reflect.Value {
	if n, ok := n.(*numberNode); ok && n.isInt {
		value := reflect.New(typ).Elem()
		value.SetInt(n.int64)
		return value
	}
	s.errorf("expected integer; found %s", n)
	panic("not reached")
}

func (s *state) evalUnsignedInteger(v reflect.Value, typ reflect.Type, n node) reflect.Value {
	if n, ok := n.(*numberNode); ok && n.isUint {
		value := reflect.New(typ).Elem()
		value.SetUint(n.uint64)
		return value
	}
	s.errorf("expected unsigned integer; found %s", n)
	panic("not reached")
}

func (s *state) evalFloat(v reflect.Value, typ reflect.Type, n node) reflect.Value {
	if n, ok := n.(*numberNode); ok && n.isFloat {
		value := reflect.New(typ).Elem()
		value.SetFloat(n.float64)
		return value
	}
	s.errorf("expected float; found %s", n)
	panic("not reached")
}

func (s *state) evalComplex(v reflect.Value, typ reflect.Type, n node) reflect.Value {
	if n, ok := n.(*numberNode); ok && n.isComplex {
		value := reflect.New(typ).Elem()
		value.SetComplex(n.complex128)
		return value
	}
	s.errorf("expected complex; found %s", n)
	panic("not reached")
}

// printValue writes the textual representation of the value to the output of
// the template.
func (s *state) printValue(n node, v reflect.Value) {
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			s.errorf("%s: nil value", n)
		}
	case reflect.Chan, reflect.Func, reflect.Interface:
		s.errorf("can't print %s of type %s", n, v.Type())
	}
	fmt.Fprint(s.wr, v.Interface())
}
