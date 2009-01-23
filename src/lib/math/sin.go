// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package math

import "math"

func sinus(arg float64, quad int) float64 {
	// Coefficients are #3370 from Hart & Cheney (18.80D).
	const
	(
		P0	=  .1357884097877375669092680e8;
		P1	= -.4942908100902844161158627e7;
		P2	=  .4401030535375266501944918e6;
		P3	= -.1384727249982452873054457e5;
		P4	=  .1459688406665768722226959e3;
		Q0	=  .8644558652922534429915149e7;
		Q1	=  .4081792252343299749395779e6;
		Q2	=  .9463096101538208180571257e4;
		Q3	=  .1326534908786136358911494e3;
	)
	x := arg;
	if(x < 0) {
		x = -x;
		quad = quad+2;
	}
	x = x * (2/Pi);	/* underflow? */
	var y float64;
	if x > 32764 {
		var e float64;
		e, y = Modf(x);
		e = e + float64(quad);
		temp1, f := Modf(0.25*e);
		quad = int(e - 4*f);
	} else {
		k := int32(x);
		y = x - float64(k);
		quad = (quad + int(k)) & 3;
	}

	if quad&1 != 0 {
		y = 1-y;
	}
	if quad > 1 {
		y = -y;
	}

	yy := y*y;
	temp1 := ((((P4*yy+P3)*yy+P2)*yy+P1)*yy+P0)*y;
	temp2 := ((((yy+Q3)*yy+Q2)*yy+Q1)*yy+Q0);
	return temp1/temp2;
}

func Cos(arg float64) float64 {
	if arg < 0 {
		arg = -arg;
	}
	return sinus(arg, 1);
}

func Sin(arg float64) float64 {
	return sinus(arg, 0);
}
