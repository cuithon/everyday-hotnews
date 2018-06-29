#!/usr/bin/env perl
# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
# Generate system call table for FreeBSD from master list
# (for example, /usr/src/sys/kern/syscalls.master).

use strict;

my $command = "mksysnum_freebsd.pl " . join(' ', @ARGV);

print <<EOF;
// $command
// Code generated by the command above; DO NOT EDIT.

package syscall

const (
EOF

while(<>){
	if(/^([0-9]+)\s+\S+\s+STD\s+({ \S+\s+(\w+).*)$/){
		my $num = $1;
		my $proto = $2;
		my $name = "SYS_$3";
		$name =~ y/a-z/A-Z/;

		# There are multiple entries for enosys and nosys, so comment them out.
		if($name =~ /^SYS_E?NOSYS$/){
			$name = "// $name";
		}
		if($name eq 'SYS_SYS_EXIT'){
			$name = 'SYS_EXIT';
		}
		if($name =~ /^SYS_CAP_+/ || $name =~ /^SYS___CAP_+/){
			next
		}

		print "	$name = $num;  // $proto\n";

		# We keep Capsicum syscall numbers for FreeBSD
		# 9-STABLE here because we are not sure whether they
		# are mature and stable.
		if($num == 513){
			print " SYS_CAP_NEW = 514 // { int cap_new(int fd, uint64_t rights); }\n";
			print " SYS_CAP_GETRIGHTS = 515 // { int cap_getrights(int fd, \\\n";
			print " SYS_CAP_ENTER = 516 // { int cap_enter(void); }\n";
			print " SYS_CAP_GETMODE = 517 // { int cap_getmode(u_int *modep); }\n";
		}
	}
}

print <<EOF;
)
EOF
