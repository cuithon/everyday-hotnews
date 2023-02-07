// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pprof

import (
	"bytes"
	"encoding/json"
	"fmt"
	"internal/abi"
	"internal/profile"
	"internal/testenv"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

// translateCPUProfile parses binary CPU profiling stack trace data
// generated by runtime.CPUProfile() into a profile struct.
// This is only used for testing. Real conversions stream the
// data into the profileBuilder as it becomes available.
//
// count is the number of records in data.
func translateCPUProfile(data []uint64, count int) (*profile.Profile, error) {
	var buf bytes.Buffer
	b := newProfileBuilder(&buf)
	tags := make([]unsafe.Pointer, count)
	if err := b.addCPUData(data, tags); err != nil {
		return nil, err
	}
	b.build()
	return profile.Parse(&buf)
}

// fmtJSON returns a pretty-printed JSON form for x.
// It works reasonably well for printing protocol-buffer
// data structures like profile.Profile.
func fmtJSON(x any) string {
	js, _ := json.MarshalIndent(x, "", "\t")
	return string(js)
}

func TestConvertCPUProfileEmpty(t *testing.T) {
	// A test server with mock cpu profile data.
	var buf bytes.Buffer

	b := []uint64{3, 0, 500} // empty profile at 500 Hz (2ms sample period)
	p, err := translateCPUProfile(b, 1)
	if err != nil {
		t.Fatalf("translateCPUProfile: %v", err)
	}
	if err := p.Write(&buf); err != nil {
		t.Fatalf("writing profile: %v", err)
	}

	p, err = profile.Parse(&buf)
	if err != nil {
		t.Fatalf("profile.Parse: %v", err)
	}

	// Expected PeriodType and SampleType.
	periodType := &profile.ValueType{Type: "cpu", Unit: "nanoseconds"}
	sampleType := []*profile.ValueType{
		{Type: "samples", Unit: "count"},
		{Type: "cpu", Unit: "nanoseconds"},
	}

	checkProfile(t, p, 2000*1000, periodType, sampleType, nil, "")
}

func f1() { f1() }
func f2() { f2() }

// testPCs returns two PCs and two corresponding memory mappings
// to use in test profiles.
func testPCs(t *testing.T) (addr1, addr2 uint64, map1, map2 *profile.Mapping) {
	switch runtime.GOOS {
	case "linux", "android", "netbsd":
		// Figure out two addresses from /proc/self/maps.
		mmap, err := os.ReadFile("/proc/self/maps")
		if err != nil {
			t.Fatal(err)
		}
		mprof := &profile.Profile{}
		if err = mprof.ParseMemoryMap(bytes.NewReader(mmap)); err != nil {
			t.Fatalf("parsing /proc/self/maps: %v", err)
		}
		if len(mprof.Mapping) < 2 {
			// It is possible for a binary to only have 1 executable
			// region of memory.
			t.Skipf("need 2 or more mappings, got %v", len(mprof.Mapping))
		}
		addr1 = mprof.Mapping[0].Start
		map1 = mprof.Mapping[0]
		map1.BuildID, _ = elfBuildID(map1.File)
		addr2 = mprof.Mapping[1].Start
		map2 = mprof.Mapping[1]
		map2.BuildID, _ = elfBuildID(map2.File)
	case "windows":
		addr1 = uint64(abi.FuncPCABIInternal(f1))
		addr2 = uint64(abi.FuncPCABIInternal(f2))

		exe, err := os.Executable()
		if err != nil {
			t.Fatal(err)
		}

		start, end, err := readMainModuleMapping()
		if err != nil {
			t.Fatal(err)
		}

		map1 = &profile.Mapping{
			ID:           1,
			Start:        start,
			Limit:        end,
			File:         exe,
			BuildID:      peBuildID(exe),
			HasFunctions: true,
		}
		map2 = &profile.Mapping{
			ID:           1,
			Start:        start,
			Limit:        end,
			File:         exe,
			BuildID:      peBuildID(exe),
			HasFunctions: true,
		}
	case "js":
		addr1 = uint64(abi.FuncPCABIInternal(f1))
		addr2 = uint64(abi.FuncPCABIInternal(f2))
	default:
		addr1 = uint64(abi.FuncPCABIInternal(f1))
		addr2 = uint64(abi.FuncPCABIInternal(f2))
		// Fake mapping - HasFunctions will be true because two PCs from Go
		// will be fully symbolized.
		fake := &profile.Mapping{ID: 1, HasFunctions: true}
		map1, map2 = fake, fake
	}
	return
}

func TestConvertCPUProfile(t *testing.T) {
	addr1, addr2, map1, map2 := testPCs(t)

	b := []uint64{
		3, 0, 500, // hz = 500
		5, 0, 10, uint64(addr1 + 1), uint64(addr1 + 2), // 10 samples in addr1
		5, 0, 40, uint64(addr2 + 1), uint64(addr2 + 2), // 40 samples in addr2
		5, 0, 10, uint64(addr1 + 1), uint64(addr1 + 2), // 10 samples in addr1
	}
	p, err := translateCPUProfile(b, 4)
	if err != nil {
		t.Fatalf("translating profile: %v", err)
	}
	period := int64(2000 * 1000)
	periodType := &profile.ValueType{Type: "cpu", Unit: "nanoseconds"}
	sampleType := []*profile.ValueType{
		{Type: "samples", Unit: "count"},
		{Type: "cpu", Unit: "nanoseconds"},
	}
	samples := []*profile.Sample{
		{Value: []int64{20, 20 * 2000 * 1000}, Location: []*profile.Location{
			{ID: 1, Mapping: map1, Address: addr1},
			{ID: 2, Mapping: map1, Address: addr1 + 1},
		}},
		{Value: []int64{40, 40 * 2000 * 1000}, Location: []*profile.Location{
			{ID: 3, Mapping: map2, Address: addr2},
			{ID: 4, Mapping: map2, Address: addr2 + 1},
		}},
	}
	checkProfile(t, p, period, periodType, sampleType, samples, "")
}

func checkProfile(t *testing.T, p *profile.Profile, period int64, periodType *profile.ValueType, sampleType []*profile.ValueType, samples []*profile.Sample, defaultSampleType string) {
	t.Helper()

	if p.Period != period {
		t.Errorf("p.Period = %d, want %d", p.Period, period)
	}
	if !reflect.DeepEqual(p.PeriodType, periodType) {
		t.Errorf("p.PeriodType = %v\nwant = %v", fmtJSON(p.PeriodType), fmtJSON(periodType))
	}
	if !reflect.DeepEqual(p.SampleType, sampleType) {
		t.Errorf("p.SampleType = %v\nwant = %v", fmtJSON(p.SampleType), fmtJSON(sampleType))
	}
	if defaultSampleType != p.DefaultSampleType {
		t.Errorf("p.DefaultSampleType = %v\nwant = %v", p.DefaultSampleType, defaultSampleType)
	}
	// Clear line info since it is not in the expected samples.
	// If we used f1 and f2 above, then the samples will have line info.
	for _, s := range p.Sample {
		for _, l := range s.Location {
			l.Line = nil
		}
	}
	if fmtJSON(p.Sample) != fmtJSON(samples) { // ignore unexported fields
		if len(p.Sample) == len(samples) {
			for i := range p.Sample {
				if !reflect.DeepEqual(p.Sample[i], samples[i]) {
					t.Errorf("sample %d = %v\nwant = %v\n", i, fmtJSON(p.Sample[i]), fmtJSON(samples[i]))
				}
			}
			if t.Failed() {
				t.FailNow()
			}
		}
		t.Fatalf("p.Sample = %v\nwant = %v", fmtJSON(p.Sample), fmtJSON(samples))
	}
}

var profSelfMapsTests = `
00400000-0040b000 r-xp 00000000 fc:01 787766                             /bin/cat
0060a000-0060b000 r--p 0000a000 fc:01 787766                             /bin/cat
0060b000-0060c000 rw-p 0000b000 fc:01 787766                             /bin/cat
014ab000-014cc000 rw-p 00000000 00:00 0                                  [heap]
7f7d76af8000-7f7d7797c000 r--p 00000000 fc:01 1318064                    /usr/lib/locale/locale-archive
7f7d7797c000-7f7d77b36000 r-xp 00000000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77b36000-7f7d77d36000 ---p 001ba000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d36000-7f7d77d3a000 r--p 001ba000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d3a000-7f7d77d3c000 rw-p 001be000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d3c000-7f7d77d41000 rw-p 00000000 00:00 0
7f7d77d41000-7f7d77d64000 r-xp 00000000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so
7f7d77f3f000-7f7d77f42000 rw-p 00000000 00:00 0
7f7d77f61000-7f7d77f63000 rw-p 00000000 00:00 0
7f7d77f63000-7f7d77f64000 r--p 00022000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so
7f7d77f64000-7f7d77f65000 rw-p 00023000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so
7f7d77f65000-7f7d77f66000 rw-p 00000000 00:00 0
7ffc342a2000-7ffc342c3000 rw-p 00000000 00:00 0                          [stack]
7ffc34343000-7ffc34345000 r-xp 00000000 00:00 0                          [vdso]
ffffffffff600000-ffffffffff601000 r-xp 00000090 00:00 0                  [vsyscall]
->
00400000 0040b000 00000000 /bin/cat
7f7d7797c000 7f7d77b36000 00000000 /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d41000 7f7d77d64000 00000000 /lib/x86_64-linux-gnu/ld-2.19.so
7ffc34343000 7ffc34345000 00000000 [vdso]
ffffffffff600000 ffffffffff601000 00000090 [vsyscall]

00400000-07000000 r-xp 00000000 00:00 0
07000000-07093000 r-xp 06c00000 00:2e 536754                             /path/to/gobench_server_main
07093000-0722d000 rw-p 06c92000 00:2e 536754                             /path/to/gobench_server_main
0722d000-07b21000 rw-p 00000000 00:00 0
c000000000-c000036000 rw-p 00000000 00:00 0
->
07000000 07093000 06c00000 /path/to/gobench_server_main
`

var profSelfMapsTestsWithDeleted = `
00400000-0040b000 r-xp 00000000 fc:01 787766                             /bin/cat (deleted)
0060a000-0060b000 r--p 0000a000 fc:01 787766                             /bin/cat (deleted)
0060b000-0060c000 rw-p 0000b000 fc:01 787766                             /bin/cat (deleted)
014ab000-014cc000 rw-p 00000000 00:00 0                                  [heap]
7f7d76af8000-7f7d7797c000 r--p 00000000 fc:01 1318064                    /usr/lib/locale/locale-archive
7f7d7797c000-7f7d77b36000 r-xp 00000000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77b36000-7f7d77d36000 ---p 001ba000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d36000-7f7d77d3a000 r--p 001ba000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d3a000-7f7d77d3c000 rw-p 001be000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d3c000-7f7d77d41000 rw-p 00000000 00:00 0
7f7d77d41000-7f7d77d64000 r-xp 00000000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so
7f7d77f3f000-7f7d77f42000 rw-p 00000000 00:00 0
7f7d77f61000-7f7d77f63000 rw-p 00000000 00:00 0
7f7d77f63000-7f7d77f64000 r--p 00022000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so
7f7d77f64000-7f7d77f65000 rw-p 00023000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so
7f7d77f65000-7f7d77f66000 rw-p 00000000 00:00 0
7ffc342a2000-7ffc342c3000 rw-p 00000000 00:00 0                          [stack]
7ffc34343000-7ffc34345000 r-xp 00000000 00:00 0                          [vdso]
ffffffffff600000-ffffffffff601000 r-xp 00000090 00:00 0                  [vsyscall]
->
00400000 0040b000 00000000 /bin/cat
7f7d7797c000 7f7d77b36000 00000000 /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d41000 7f7d77d64000 00000000 /lib/x86_64-linux-gnu/ld-2.19.so
7ffc34343000 7ffc34345000 00000000 [vdso]
ffffffffff600000 ffffffffff601000 00000090 [vsyscall]

00400000-0040b000 r-xp 00000000 fc:01 787766                             /bin/cat with space
0060a000-0060b000 r--p 0000a000 fc:01 787766                             /bin/cat with space
0060b000-0060c000 rw-p 0000b000 fc:01 787766                             /bin/cat with space
014ab000-014cc000 rw-p 00000000 00:00 0                                  [heap]
7f7d76af8000-7f7d7797c000 r--p 00000000 fc:01 1318064                    /usr/lib/locale/locale-archive
7f7d7797c000-7f7d77b36000 r-xp 00000000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77b36000-7f7d77d36000 ---p 001ba000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d36000-7f7d77d3a000 r--p 001ba000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d3a000-7f7d77d3c000 rw-p 001be000 fc:01 1180226                    /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d3c000-7f7d77d41000 rw-p 00000000 00:00 0
7f7d77d41000-7f7d77d64000 r-xp 00000000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so
7f7d77f3f000-7f7d77f42000 rw-p 00000000 00:00 0
7f7d77f61000-7f7d77f63000 rw-p 00000000 00:00 0
7f7d77f63000-7f7d77f64000 r--p 00022000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so
7f7d77f64000-7f7d77f65000 rw-p 00023000 fc:01 1180217                    /lib/x86_64-linux-gnu/ld-2.19.so
7f7d77f65000-7f7d77f66000 rw-p 00000000 00:00 0
7ffc342a2000-7ffc342c3000 rw-p 00000000 00:00 0                          [stack]
7ffc34343000-7ffc34345000 r-xp 00000000 00:00 0                          [vdso]
ffffffffff600000-ffffffffff601000 r-xp 00000090 00:00 0                  [vsyscall]
->
00400000 0040b000 00000000 /bin/cat with space
7f7d7797c000 7f7d77b36000 00000000 /lib/x86_64-linux-gnu/libc-2.19.so
7f7d77d41000 7f7d77d64000 00000000 /lib/x86_64-linux-gnu/ld-2.19.so
7ffc34343000 7ffc34345000 00000000 [vdso]
ffffffffff600000 ffffffffff601000 00000090 [vsyscall]
`

func TestProcSelfMaps(t *testing.T) {

	f := func(t *testing.T, input string) {
		for tx, tt := range strings.Split(input, "\n\n") {
			in, out, ok := strings.Cut(tt, "->\n")
			if !ok {
				t.Fatal("malformed test case")
			}
			if len(out) > 0 && out[len(out)-1] != '\n' {
				out += "\n"
			}
			var buf strings.Builder
			parseProcSelfMaps([]byte(in), func(lo, hi, offset uint64, file, buildID string) {
				fmt.Fprintf(&buf, "%08x %08x %08x %s\n", lo, hi, offset, file)
			})
			if buf.String() != out {
				t.Errorf("#%d: have:\n%s\nwant:\n%s\n%q\n%q", tx, buf.String(), out, buf.String(), out)
			}
		}
	}

	t.Run("Normal", func(t *testing.T) {
		f(t, profSelfMapsTests)
	})

	t.Run("WithDeletedFile", func(t *testing.T) {
		f(t, profSelfMapsTestsWithDeleted)
	})
}

// TestMapping checks the mapping section of CPU profiles
// has the HasFunctions field set correctly. If all PCs included
// in the samples are successfully symbolized, the corresponding
// mapping entry (in this test case, only one entry) should have
// its HasFunctions field set true.
// The test generates a CPU profile that includes PCs from C side
// that the runtime can't symbolize. See ./testdata/mappingtest.
func TestMapping(t *testing.T) {
	testenv.MustHaveGoRun(t)
	testenv.MustHaveCGO(t)

	prog := "./testdata/mappingtest/main.go"

	// GoOnly includes only Go symbols that runtime will symbolize.
	// Go+C includes C symbols that runtime will not symbolize.
	for _, traceback := range []string{"GoOnly", "Go+C"} {
		t.Run("traceback"+traceback, func(t *testing.T) {
			cmd := exec.Command(testenv.GoToolPath(t), "run", prog)
			if traceback != "GoOnly" {
				cmd.Env = append(os.Environ(), "SETCGOTRACEBACK=1")
			}
			cmd.Stderr = new(bytes.Buffer)

			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("failed to run the test program %q: %v\n%v", prog, err, cmd.Stderr)
			}

			prof, err := profile.Parse(bytes.NewReader(out))
			if err != nil {
				t.Fatalf("failed to parse the generated profile data: %v", err)
			}
			t.Logf("Profile: %s", prof)

			hit := make(map[*profile.Mapping]bool)
			miss := make(map[*profile.Mapping]bool)
			for _, loc := range prof.Location {
				if symbolized(loc) {
					hit[loc.Mapping] = true
				} else {
					miss[loc.Mapping] = true
				}
			}
			if len(miss) == 0 {
				t.Log("no location with missing symbol info was sampled")
			}

			for _, m := range prof.Mapping {
				if miss[m] && m.HasFunctions {
					t.Errorf("mapping %+v has HasFunctions=true, but contains locations with failed symbolization", m)
					continue
				}
				if !miss[m] && hit[m] && !m.HasFunctions {
					t.Errorf("mapping %+v has HasFunctions=false, but all referenced locations from this lapping were symbolized successfully", m)
					continue
				}
			}

			if traceback == "Go+C" {
				// The test code was arranged to have PCs from C and
				// they are not symbolized.
				// Check no Location containing those unsymbolized PCs contains multiple lines.
				for i, loc := range prof.Location {
					if !symbolized(loc) && len(loc.Line) > 1 {
						t.Errorf("Location[%d] contains unsymbolized PCs and multiple lines: %v", i, loc)
					}
				}
			}
		})
	}
}

func symbolized(loc *profile.Location) bool {
	if len(loc.Line) == 0 {
		return false
	}
	l := loc.Line[0]
	f := l.Function
	if l.Line == 0 || f == nil || f.Name == "" || f.Filename == "" {
		return false
	}
	return true
}

// TestFakeMapping tests if at least one mapping exists
// (including a fake mapping), and their HasFunctions bits
// are set correctly.
func TestFakeMapping(t *testing.T) {
	var buf bytes.Buffer
	if err := Lookup("heap").WriteTo(&buf, 0); err != nil {
		t.Fatalf("failed to write heap profile: %v", err)
	}
	prof, err := profile.Parse(&buf)
	if err != nil {
		t.Fatalf("failed to parse the generated profile data: %v", err)
	}
	t.Logf("Profile: %s", prof)
	if len(prof.Mapping) == 0 {
		t.Fatal("want profile with at least one mapping entry, got 0 mapping")
	}

	hit := make(map[*profile.Mapping]bool)
	miss := make(map[*profile.Mapping]bool)
	for _, loc := range prof.Location {
		if symbolized(loc) {
			hit[loc.Mapping] = true
		} else {
			miss[loc.Mapping] = true
		}
	}
	for _, m := range prof.Mapping {
		if miss[m] && m.HasFunctions {
			t.Errorf("mapping %+v has HasFunctions=true, but contains locations with failed symbolization", m)
			continue
		}
		if !miss[m] && hit[m] && !m.HasFunctions {
			t.Errorf("mapping %+v has HasFunctions=false, but all referenced locations from this lapping were symbolized successfully", m)
			continue
		}
	}
}

// Make sure the profiler can handle an empty stack trace.
// See issue 37967.
func TestEmptyStack(t *testing.T) {
	b := []uint64{
		3, 0, 500, // hz = 500
		3, 0, 10, // 10 samples with an empty stack trace
	}
	_, err := translateCPUProfile(b, 2)
	if err != nil {
		t.Fatalf("translating profile: %v", err)
	}
}
