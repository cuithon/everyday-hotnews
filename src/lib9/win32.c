// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include <u.h>
#include <libc.h>

int fork()
{
	return -1;
}

int p9rfork(int flags)
{
	return -1;
}

Waitmsg *p9wait()
{
	return 0;
}

int p9waitpid()
{
	return -1;
}
