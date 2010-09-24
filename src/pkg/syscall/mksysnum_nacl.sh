#!/usr/bin/perl
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

my $command = "mksysnum_nacl.sh ". join(' ', @ARGV);

print <<EOF;
// $command
// MACHINE GENERATED BY THE ABOVE COMMAND; DO NOT EDIT

package syscall

const(
EOF

while(<>){
	if(/^#define NACL_sys_(\w+)\s+([0-9]+)/){
		my $name = "SYS_$1";
		my $num = $2;
		$name =~ y/a-z/A-Z/;
		print "	$name = $num;\n";
	}
}

print <<EOF;
)

EOF
