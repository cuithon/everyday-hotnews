// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Parse "zoneinfo" time zone file.
// This is a fairly standard file format used on OS X, Linux, BSD, Sun, and others.
// See tzfile(5), http://en.wikipedia.org/wiki/Zoneinfo,
// and ftp://munnari.oz.au/pub/oldtz/

package time

import (
	"once";
	"os"
)

const (
	MaxFileSize = 8192;	// actual files are closer to 1K
	HeaderSize = 4+16+4*7
)

export var (
	BadZoneinfo = os.NewError("time: malformed zoneinfo");
	NoZoneinfo = os.NewError("time: unknown time zone")
)

// Simple I/O interface to binary blob of data.
type Data struct {
	p *[]byte
}

func (d *Data) Read(n int) *[]byte {
	if d.p == nil || len(d.p) < n {
		d.p = nil
		return nil
	}
	p := d.p[0:n];
	d.p = d.p[n:len(d.p)]
	return p
}

func (d *Data) Big4() (n uint32, ok bool) {
	p := d.Read(4)
	if p == nil {
		return 0, false
	}
	return uint32(p[0]) << 24 | uint32(p[1]) << 16 | uint32(p[2]) << 8 | uint32(p[3]), true
}

func (d *Data) Byte() (n byte, ok bool) {
	p := d.Read(1)
	if p == nil {
		return 0, false
	}
	return p[0], true
}


// Make a string by stopping at the first NUL
func ByteString(p *[]byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

// Parsed representation
type Zone struct {
	utcoff int;
	isdst bool;
	name string
}

type Zonetime struct {
	time int32;		// transition time, in seconds since 1970 GMT
	zone *Zone;		// the zone that goes into effect at that time
	isstd, isutc bool;	// ignored - no idea what these mean
}

func ParseZoneinfo(bytes *[]byte) (zt *[]Zonetime, err *os.Error) {
	data1 := Data(bytes);
	data := &data1;

	// 4-byte magic "TZif"
	if magic := data.Read(4); magic == nil || string(magic) != "TZif" {
		return nil, BadZoneinfo
	}

	// 1-byte version, then 15 bytes of padding
	var p *[]byte
	if p = data.Read(16); p == nil || p[0] != 0 && p[0] != '2' {
		return nil, BadZoneinfo
	}
	vers := p[0]

	// six big-endian 32-bit integers:
	//	number of UTC/local indicators
	//	number of standard/wall indicators
	//	number of leap seconds
	//	number of transition times
	//	number of local time zones
	//	number of characters of time zone abbrev strings
	const (
		NUTCLocal = iota;
		NStdWall;
		NLeap;
		NTime;
		NZone;
		NChar
	)
	var n [6]int
	for i := 0; i < 6; i++ {
		nn, ok := data.Big4()
		if !ok {
			return nil, BadZoneinfo
		}
		n[i] = int(nn);
	}

	// Transition times.
	txtimes1 := Data(data.Read(n[NTime]*4));
	txtimes := &txtimes1;

	// Time zone indices for transition times.
	txzones := data.Read(n[NTime]);

	// Zone info structures
	zonedata1 := Data(data.Read(n[NZone]*6));
	zonedata := &zonedata1;

	// Time zone abbreviations.
	abbrev := data.Read(n[NChar]);

	// Leap-second time pairs
	leapdata1 := Data(data.Read(n[NLeap]*8));
	leapdata := &leapdata1;

	// Whether tx times associated with local time types
	// are specified as standard time or wall time.
	isstd := data.Read(n[NStdWall]);

	// Whether tx times associated with local time types
	// are specified as UTC or local time.
	isutc := data.Read(n[NUTCLocal]);

	if data.p == nil {	// ran out of data
		return nil, BadZoneinfo
	}

	// If version == 2, the entire file repeats, this time using
	// 8-byte ints for txtimes and leap seconds.
	// We won't need those until 2106.

	// Now we can build up a useful data structure.
	// First the zone information.
	//	utcoff[4] isdst[1] nameindex[1]
	zone := new([]Zone, n[NZone])
	for i := 0; i < len(zone); i++ {
		var ok bool;
		var n uint32
		if n, ok = zonedata.Big4(); !ok {
			return nil, BadZoneinfo
		}
		zone[i].utcoff = int(n)
		var b byte
		if b, ok = zonedata.Byte(); !ok {
			return nil, BadZoneinfo
		}
		zone[i].isdst = b != 0
		if b, ok = zonedata.Byte(); !ok || int(b) >= len(abbrev) {
			return nil, BadZoneinfo
		}
		zone[i].name = ByteString(abbrev[b:len(abbrev)])
	}

	// Now the transition time info.
	zt = new([]Zonetime, n[NTime])
	for i := 0; i < len(zt); i++ {
		var ok bool;
		var n uint32
		if n, ok = txtimes.Big4(); !ok {
			return nil, BadZoneinfo
		}
		zt[i].time = int32(n)
		if int(txzones[i]) >= len(zone) {
			return nil, BadZoneinfo
		}
		zt[i].zone = &zone[txzones[i]];
		if i < len(isstd) {
			zt[i].isstd = isstd[i] != 0
		}
		if i < len(isutc) {
			zt[i].isutc = isutc[i] != 0
		}
	}
	return zt, nil
}

func ReadFile(name string, max int) (p *[]byte, err *os.Error) {
	fd, e := os.Open(name, os.O_RDONLY, 0)
	if e != nil {
		return nil, e
	}
	p = new([]byte, max+1)[0:0];
	n := 0
	for len(p) < max {
		nn, e := fd.Read(p[n:cap(p)])
		if e != nil {
			fd.Close()
			return nil, e
		}
		if nn == 0 {
			fd.Close()
			return p, nil
		}
		p = p[0:n+nn]
	}
	fd.Close()
	return nil, BadZoneinfo	// too long
}


func ReadZoneinfoFile(name string) (tx *[]Zonetime, err *os.Error) {
	data, e := ReadFile(name, MaxFileSize)
	if e != nil {
		return nil, e
	}
	tx, err = ParseZoneinfo(data)
	return tx, err
}

var zones *[]Zonetime
var zoneerr *os.Error

func SetupZone() {
	// TODO: /etc/localtime is the default time zone info
	// for the system, but libc allows setting an environment
	// variable in order to direct reading a different file
	// (in /usr/share/zoneinfo).  We should check that 
	// environment variable.
	zones, zoneerr = ReadZoneinfoFile("/etc/localtime");
}

export func LookupTimezone(sec int64) (zone string, offset int, err *os.Error) {
	once.Do(&SetupZone)
	if zoneerr != nil || zones == nil || len(zones) == 0 {
		return "GMT", 0, zoneerr
	}

	// Binary search for entry with largest time <= sec
	tz := zones;
	for len(tz) > 1 {
		m := len(tz)/2;
		if sec < int64(tz[m].time) {
			tz = tz[0:m]
		} else {
			tz = tz[m:len(tz)]
		}
	}
	z := tz[0].zone
	return z.name, z.utcoff, nil
}
