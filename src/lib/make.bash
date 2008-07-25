# Copyright 2009 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

#!/bin/bash

rm -f *.6
for i in fmt.go flag.go container/vector.go
do
	6g $i
done
mv *.6 $GOROOT/pkg
