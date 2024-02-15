// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !cmd_go_bootstrap

package main

import "golang.org/x/telemetry"

var TelemetryStart = func() {
	telemetry.Start(telemetry.Config{Upload: true})
}
