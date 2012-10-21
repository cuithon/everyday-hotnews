// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Runtime type representation; master is type.go
 *
 * The *Types here correspond 1-1 to type.go's *Type's, but are
 * prefixed with an extra header of 2 pointers, corresponding to the
 * interface{} structure, which itself is called type Type again on
 * the Go side.
 */

typedef struct CommonType CommonType;
typedef struct UncommonType UncommonType;
typedef struct InterfaceType InterfaceType;
typedef struct Method Method;
typedef struct IMethod IMethod;
typedef struct SliceType SliceType;
typedef struct FuncType FuncType;
typedef struct PtrType PtrType;

// Needs to be in sync with typekind.h/CommonSize
struct CommonType
{
	uintptr size;
	uint32 hash;
	uint8 _unused;
	uint8 align;
	uint8 fieldAlign;
	uint8 kind;
	Alg *alg;
	void *gc;
	String *string;
	UncommonType *x;
	Type *ptrto;
};

struct Method
{
	String *name;
	String *pkgPath;
	Type	*mtyp;
	Type *typ;
	void (*ifn)(void);
	void (*tfn)(void);
};

struct UncommonType
{
	String *name;
	String *pkgPath;
	Slice mhdr;
	Method m[];
};

struct Type
{
	void *type;	// interface{} value
	void *ptr;
	CommonType;
};

struct IMethod
{
	String *name;
	String *pkgPath;
	Type *type;
};

struct InterfaceType
{
	Type;
	Slice mhdr;
	IMethod m[];
};

struct MapType
{
	Type;
	Type *key;
	Type *elem;
};

struct ChanType
{
	Type;
	Type *elem;
	uintptr dir;
};

struct SliceType
{
	Type;
	Type *elem;
};

struct FuncType
{
	Type;
	bool dotdotdot;
	Slice in;
	Slice out;
};

struct PtrType
{
	Type;
	Type *elem;
};
