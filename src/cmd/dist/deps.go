// Code generated by mkdeps.bash; DO NOT EDIT.

package main

var builddeps = map[string][]string{

	"bufio": {
		"bytes",        // bufio
		"errors",       // bufio
		"io",           // bufio
		"unicode/utf8", // bufio
	},

	"bytes": {
		"errors",       // bytes
		"internal/cpu", // bytes
		"io",           // bytes
		"unicode",      // bytes
		"unicode/utf8", // bytes
	},

	"cmd/go": {
		"cmd/go/internal/base",     // cmd/go
		"cmd/go/internal/bug",      // cmd/go
		"cmd/go/internal/cfg",      // cmd/go
		"cmd/go/internal/clean",    // cmd/go
		"cmd/go/internal/doc",      // cmd/go
		"cmd/go/internal/envcmd",   // cmd/go
		"cmd/go/internal/fix",      // cmd/go
		"cmd/go/internal/fmtcmd",   // cmd/go
		"cmd/go/internal/generate", // cmd/go
		"cmd/go/internal/get",      // cmd/go
		"cmd/go/internal/help",     // cmd/go
		"cmd/go/internal/list",     // cmd/go
		"cmd/go/internal/run",      // cmd/go
		"cmd/go/internal/test",     // cmd/go
		"cmd/go/internal/tool",     // cmd/go
		"cmd/go/internal/version",  // cmd/go
		"cmd/go/internal/vet",      // cmd/go
		"cmd/go/internal/work",     // cmd/go
		"flag",                     // cmd/go
		"fmt",                      // cmd/go
		"log",                      // cmd/go
		"os",                       // cmd/go
		"path/filepath",            // cmd/go
		"runtime",                  // cmd/go
		"strings",                  // cmd/go
	},

	"cmd/go/internal/base": {
		"bytes",               // cmd/go/internal/base
		"cmd/go/internal/cfg", // cmd/go/internal/base
		"cmd/go/internal/str", // cmd/go/internal/base
		"errors",              // cmd/go/internal/base
		"flag",                // cmd/go/internal/base
		"fmt",                 // cmd/go/internal/base
		"go/build",            // cmd/go/internal/base
		"go/scanner",          // cmd/go/internal/base
		"log",                 // cmd/go/internal/base
		"os",                  // cmd/go/internal/base
		"os/exec",             // cmd/go/internal/base
		"os/signal",           // cmd/go/internal/base
		"path/filepath",       // cmd/go/internal/base
		"runtime",             // cmd/go/internal/base
		"strings",             // cmd/go/internal/base
		"sync",                // cmd/go/internal/base
		"syscall",             // cmd/go/internal/base
	},

	"cmd/go/internal/bug": {
		"bytes",                  // cmd/go/internal/bug
		"cmd/go/internal/base",   // cmd/go/internal/bug
		"cmd/go/internal/cfg",    // cmd/go/internal/bug
		"cmd/go/internal/envcmd", // cmd/go/internal/bug
		"cmd/go/internal/web",    // cmd/go/internal/bug
		"fmt",           // cmd/go/internal/bug
		"io",            // cmd/go/internal/bug
		"io/ioutil",     // cmd/go/internal/bug
		"os",            // cmd/go/internal/bug
		"os/exec",       // cmd/go/internal/bug
		"path/filepath", // cmd/go/internal/bug
		"regexp",        // cmd/go/internal/bug
		"runtime",       // cmd/go/internal/bug
		"strings",       // cmd/go/internal/bug
	},

	"cmd/go/internal/cfg": {
		"cmd/internal/objabi", // cmd/go/internal/cfg
		"fmt",           // cmd/go/internal/cfg
		"go/build",      // cmd/go/internal/cfg
		"os",            // cmd/go/internal/cfg
		"path/filepath", // cmd/go/internal/cfg
		"runtime",       // cmd/go/internal/cfg
	},

	"cmd/go/internal/clean": {
		"cmd/go/internal/base", // cmd/go/internal/clean
		"cmd/go/internal/cfg",  // cmd/go/internal/clean
		"cmd/go/internal/load", // cmd/go/internal/clean
		"cmd/go/internal/work", // cmd/go/internal/clean
		"fmt",           // cmd/go/internal/clean
		"io/ioutil",     // cmd/go/internal/clean
		"os",            // cmd/go/internal/clean
		"path/filepath", // cmd/go/internal/clean
		"strings",       // cmd/go/internal/clean
	},

	"cmd/go/internal/cmdflag": {
		"cmd/go/internal/base", // cmd/go/internal/cmdflag
		"flag",                 // cmd/go/internal/cmdflag
		"fmt",                  // cmd/go/internal/cmdflag
		"os",                   // cmd/go/internal/cmdflag
		"strconv",              // cmd/go/internal/cmdflag
		"strings",              // cmd/go/internal/cmdflag
	},

	"cmd/go/internal/doc": {
		"cmd/go/internal/base", // cmd/go/internal/doc
		"cmd/go/internal/cfg",  // cmd/go/internal/doc
	},

	"cmd/go/internal/envcmd": {
		"cmd/go/internal/base", // cmd/go/internal/envcmd
		"cmd/go/internal/cfg",  // cmd/go/internal/envcmd
		"cmd/go/internal/load", // cmd/go/internal/envcmd
		"cmd/go/internal/work", // cmd/go/internal/envcmd
		"encoding/json",        // cmd/go/internal/envcmd
		"fmt",                  // cmd/go/internal/envcmd
		"os",                   // cmd/go/internal/envcmd
		"runtime",              // cmd/go/internal/envcmd
		"strings",              // cmd/go/internal/envcmd
	},

	"cmd/go/internal/fix": {
		"cmd/go/internal/base", // cmd/go/internal/fix
		"cmd/go/internal/cfg",  // cmd/go/internal/fix
		"cmd/go/internal/load", // cmd/go/internal/fix
		"cmd/go/internal/str",  // cmd/go/internal/fix
	},

	"cmd/go/internal/fmtcmd": {
		"cmd/go/internal/base", // cmd/go/internal/fmtcmd
		"cmd/go/internal/cfg",  // cmd/go/internal/fmtcmd
		"cmd/go/internal/load", // cmd/go/internal/fmtcmd
		"cmd/go/internal/str",  // cmd/go/internal/fmtcmd
		"os",            // cmd/go/internal/fmtcmd
		"path/filepath", // cmd/go/internal/fmtcmd
		"runtime",       // cmd/go/internal/fmtcmd
		"sync",          // cmd/go/internal/fmtcmd
	},

	"cmd/go/internal/generate": {
		"bufio",                // cmd/go/internal/generate
		"bytes",                // cmd/go/internal/generate
		"cmd/go/internal/base", // cmd/go/internal/generate
		"cmd/go/internal/cfg",  // cmd/go/internal/generate
		"cmd/go/internal/load", // cmd/go/internal/generate
		"cmd/go/internal/work", // cmd/go/internal/generate
		"fmt",           // cmd/go/internal/generate
		"io",            // cmd/go/internal/generate
		"log",           // cmd/go/internal/generate
		"os",            // cmd/go/internal/generate
		"os/exec",       // cmd/go/internal/generate
		"path/filepath", // cmd/go/internal/generate
		"regexp",        // cmd/go/internal/generate
		"strconv",       // cmd/go/internal/generate
		"strings",       // cmd/go/internal/generate
	},

	"cmd/go/internal/get": {
		"bytes",                 // cmd/go/internal/get
		"cmd/go/internal/base",  // cmd/go/internal/get
		"cmd/go/internal/cfg",   // cmd/go/internal/get
		"cmd/go/internal/load",  // cmd/go/internal/get
		"cmd/go/internal/str",   // cmd/go/internal/get
		"cmd/go/internal/web",   // cmd/go/internal/get
		"cmd/go/internal/work",  // cmd/go/internal/get
		"encoding/json",         // cmd/go/internal/get
		"encoding/xml",          // cmd/go/internal/get
		"errors",                // cmd/go/internal/get
		"fmt",                   // cmd/go/internal/get
		"go/build",              // cmd/go/internal/get
		"internal/singleflight", // cmd/go/internal/get
		"io",            // cmd/go/internal/get
		"log",           // cmd/go/internal/get
		"net/url",       // cmd/go/internal/get
		"os",            // cmd/go/internal/get
		"os/exec",       // cmd/go/internal/get
		"path/filepath", // cmd/go/internal/get
		"regexp",        // cmd/go/internal/get
		"runtime",       // cmd/go/internal/get
		"strings",       // cmd/go/internal/get
		"sync",          // cmd/go/internal/get
	},

	"cmd/go/internal/help": {
		"bufio",                // cmd/go/internal/help
		"bytes",                // cmd/go/internal/help
		"cmd/go/internal/base", // cmd/go/internal/help
		"fmt",           // cmd/go/internal/help
		"io",            // cmd/go/internal/help
		"os",            // cmd/go/internal/help
		"strings",       // cmd/go/internal/help
		"text/template", // cmd/go/internal/help
		"unicode",       // cmd/go/internal/help
		"unicode/utf8",  // cmd/go/internal/help
	},

	"cmd/go/internal/list": {
		"bufio",                // cmd/go/internal/list
		"cmd/go/internal/base", // cmd/go/internal/list
		"cmd/go/internal/cfg",  // cmd/go/internal/list
		"cmd/go/internal/load", // cmd/go/internal/list
		"cmd/go/internal/work", // cmd/go/internal/list
		"encoding/json",        // cmd/go/internal/list
		"go/build",             // cmd/go/internal/list
		"io",                   // cmd/go/internal/list
		"os",                   // cmd/go/internal/list
		"strings",              // cmd/go/internal/list
		"text/template",        // cmd/go/internal/list
	},

	"cmd/go/internal/load": {
		"cmd/go/internal/base", // cmd/go/internal/load
		"cmd/go/internal/cfg",  // cmd/go/internal/load
		"cmd/go/internal/str",  // cmd/go/internal/load
		"cmd/internal/buildid", // cmd/go/internal/load
		"crypto/sha1",          // cmd/go/internal/load
		"fmt",                  // cmd/go/internal/load
		"go/build",             // cmd/go/internal/load
		"go/token",             // cmd/go/internal/load
		"io/ioutil",            // cmd/go/internal/load
		"log",                  // cmd/go/internal/load
		"os",                   // cmd/go/internal/load
		"path",                 // cmd/go/internal/load
		"path/filepath",        // cmd/go/internal/load
		"regexp",               // cmd/go/internal/load
		"runtime",              // cmd/go/internal/load
		"sort",                 // cmd/go/internal/load
		"strings",              // cmd/go/internal/load
		"unicode",              // cmd/go/internal/load
	},

	"cmd/go/internal/run": {
		"cmd/go/internal/base", // cmd/go/internal/run
		"cmd/go/internal/cfg",  // cmd/go/internal/run
		"cmd/go/internal/load", // cmd/go/internal/run
		"cmd/go/internal/str",  // cmd/go/internal/run
		"cmd/go/internal/work", // cmd/go/internal/run
		"fmt",     // cmd/go/internal/run
		"os",      // cmd/go/internal/run
		"strings", // cmd/go/internal/run
	},

	"cmd/go/internal/str": {
		"bytes",        // cmd/go/internal/str
		"fmt",          // cmd/go/internal/str
		"unicode",      // cmd/go/internal/str
		"unicode/utf8", // cmd/go/internal/str
	},

	"cmd/go/internal/test": {
		"bytes",                   // cmd/go/internal/test
		"cmd/go/internal/base",    // cmd/go/internal/test
		"cmd/go/internal/cfg",     // cmd/go/internal/test
		"cmd/go/internal/cmdflag", // cmd/go/internal/test
		"cmd/go/internal/load",    // cmd/go/internal/test
		"cmd/go/internal/str",     // cmd/go/internal/test
		"cmd/go/internal/work",    // cmd/go/internal/test
		"errors",                  // cmd/go/internal/test
		"flag",                    // cmd/go/internal/test
		"fmt",                     // cmd/go/internal/test
		"go/ast",                  // cmd/go/internal/test
		"go/build",                // cmd/go/internal/test
		"go/doc",                  // cmd/go/internal/test
		"go/parser",               // cmd/go/internal/test
		"go/token",                // cmd/go/internal/test
		"os",                      // cmd/go/internal/test
		"os/exec",                 // cmd/go/internal/test
		"path",                    // cmd/go/internal/test
		"path/filepath",           // cmd/go/internal/test
		"regexp",                  // cmd/go/internal/test
		"sort",                    // cmd/go/internal/test
		"strings",                 // cmd/go/internal/test
		"text/template",           // cmd/go/internal/test
		"time",                    // cmd/go/internal/test
		"unicode",                 // cmd/go/internal/test
		"unicode/utf8",            // cmd/go/internal/test
	},

	"cmd/go/internal/tool": {
		"cmd/go/internal/base", // cmd/go/internal/tool
		"cmd/go/internal/cfg",  // cmd/go/internal/tool
		"fmt",     // cmd/go/internal/tool
		"os",      // cmd/go/internal/tool
		"os/exec", // cmd/go/internal/tool
		"sort",    // cmd/go/internal/tool
		"strings", // cmd/go/internal/tool
	},

	"cmd/go/internal/version": {
		"cmd/go/internal/base", // cmd/go/internal/version
		"fmt",     // cmd/go/internal/version
		"runtime", // cmd/go/internal/version
	},

	"cmd/go/internal/vet": {
		"cmd/go/internal/base",    // cmd/go/internal/vet
		"cmd/go/internal/cfg",     // cmd/go/internal/vet
		"cmd/go/internal/cmdflag", // cmd/go/internal/vet
		"cmd/go/internal/load",    // cmd/go/internal/vet
		"cmd/go/internal/str",     // cmd/go/internal/vet
		"cmd/go/internal/work",    // cmd/go/internal/vet
		"flag",                    // cmd/go/internal/vet
		"fmt",                     // cmd/go/internal/vet
		"os",                      // cmd/go/internal/vet
		"path/filepath",           // cmd/go/internal/vet
		"strings",                 // cmd/go/internal/vet
	},

	"cmd/go/internal/web": {
		"errors", // cmd/go/internal/web
		"io",     // cmd/go/internal/web
	},

	"cmd/go/internal/work": {
		"bufio",                // cmd/go/internal/work
		"bytes",                // cmd/go/internal/work
		"cmd/go/internal/base", // cmd/go/internal/work
		"cmd/go/internal/cfg",  // cmd/go/internal/work
		"cmd/go/internal/load", // cmd/go/internal/work
		"cmd/go/internal/str",  // cmd/go/internal/work
		"cmd/internal/buildid", // cmd/go/internal/work
		"container/heap",       // cmd/go/internal/work
		"debug/elf",            // cmd/go/internal/work
		"errors",               // cmd/go/internal/work
		"flag",                 // cmd/go/internal/work
		"fmt",                  // cmd/go/internal/work
		"go/build",             // cmd/go/internal/work
		"io",                   // cmd/go/internal/work
		"io/ioutil",            // cmd/go/internal/work
		"log",                  // cmd/go/internal/work
		"os",                   // cmd/go/internal/work
		"os/exec",              // cmd/go/internal/work
		"path",                 // cmd/go/internal/work
		"path/filepath",        // cmd/go/internal/work
		"regexp",               // cmd/go/internal/work
		"runtime",              // cmd/go/internal/work
		"strconv",              // cmd/go/internal/work
		"strings",              // cmd/go/internal/work
		"sync",                 // cmd/go/internal/work
		"time",                 // cmd/go/internal/work
	},

	"cmd/internal/buildid": {
		"bytes",           // cmd/internal/buildid
		"crypto/sha256",   // cmd/internal/buildid
		"debug/elf",       // cmd/internal/buildid
		"debug/macho",     // cmd/internal/buildid
		"encoding/binary", // cmd/internal/buildid
		"fmt",             // cmd/internal/buildid
		"io",              // cmd/internal/buildid
		"os",              // cmd/internal/buildid
		"strconv",         // cmd/internal/buildid
	},

	"cmd/internal/objabi": {
		"flag",          // cmd/internal/objabi
		"fmt",           // cmd/internal/objabi
		"log",           // cmd/internal/objabi
		"os",            // cmd/internal/objabi
		"path/filepath", // cmd/internal/objabi
		"runtime",       // cmd/internal/objabi
		"strconv",       // cmd/internal/objabi
		"strings",       // cmd/internal/objabi
	},

	"compress/flate": {
		"bufio",     // compress/flate
		"fmt",       // compress/flate
		"io",        // compress/flate
		"math",      // compress/flate
		"math/bits", // compress/flate
		"sort",      // compress/flate
		"strconv",   // compress/flate
		"sync",      // compress/flate
	},

	"compress/zlib": {
		"bufio",          // compress/zlib
		"compress/flate", // compress/zlib
		"errors",         // compress/zlib
		"fmt",            // compress/zlib
		"hash",           // compress/zlib
		"hash/adler32",   // compress/zlib
		"io",             // compress/zlib
	},

	"container/heap": {
		"sort", // container/heap
	},

	"context": {
		"errors",  // context
		"fmt",     // context
		"reflect", // context
		"sync",    // context
		"time",    // context
	},

	"crypto": {
		"hash",    // crypto
		"io",      // crypto
		"strconv", // crypto
	},

	"crypto/sha1": {
		"crypto",       // crypto/sha1
		"hash",         // crypto/sha1
		"internal/cpu", // crypto/sha1
	},

	"crypto/sha256": {
		"crypto",       // crypto/sha256
		"hash",         // crypto/sha256
		"internal/cpu", // crypto/sha256
	},

	"debug/dwarf": {
		"encoding/binary", // debug/dwarf
		"errors",          // debug/dwarf
		"fmt",             // debug/dwarf
		"io",              // debug/dwarf
		"path",            // debug/dwarf
		"sort",            // debug/dwarf
		"strconv",         // debug/dwarf
		"strings",         // debug/dwarf
	},

	"debug/elf": {
		"bytes",           // debug/elf
		"compress/zlib",   // debug/elf
		"debug/dwarf",     // debug/elf
		"encoding/binary", // debug/elf
		"errors",          // debug/elf
		"fmt",             // debug/elf
		"io",              // debug/elf
		"os",              // debug/elf
		"strconv",         // debug/elf
		"strings",         // debug/elf
	},

	"debug/macho": {
		"bytes",           // debug/macho
		"debug/dwarf",     // debug/macho
		"encoding/binary", // debug/macho
		"fmt",             // debug/macho
		"io",              // debug/macho
		"os",              // debug/macho
		"strconv",         // debug/macho
	},

	"encoding": {
		"runtime", // encoding
	},

	"encoding/base64": {
		"encoding/binary", // encoding/base64
		"io",              // encoding/base64
		"strconv",         // encoding/base64
	},

	"encoding/binary": {
		"errors",  // encoding/binary
		"io",      // encoding/binary
		"math",    // encoding/binary
		"reflect", // encoding/binary
	},

	"encoding/json": {
		"bytes",           // encoding/json
		"encoding",        // encoding/json
		"encoding/base64", // encoding/json
		"errors",          // encoding/json
		"fmt",             // encoding/json
		"io",              // encoding/json
		"math",            // encoding/json
		"reflect",         // encoding/json
		"runtime",         // encoding/json
		"sort",            // encoding/json
		"strconv",         // encoding/json
		"strings",         // encoding/json
		"sync",            // encoding/json
		"sync/atomic",     // encoding/json
		"unicode",         // encoding/json
		"unicode/utf16",   // encoding/json
		"unicode/utf8",    // encoding/json
	},

	"encoding/xml": {
		"bufio",        // encoding/xml
		"bytes",        // encoding/xml
		"encoding",     // encoding/xml
		"errors",       // encoding/xml
		"fmt",          // encoding/xml
		"io",           // encoding/xml
		"reflect",      // encoding/xml
		"strconv",      // encoding/xml
		"strings",      // encoding/xml
		"sync",         // encoding/xml
		"unicode",      // encoding/xml
		"unicode/utf8", // encoding/xml
	},

	"errors": {
		"runtime", // errors
	},

	"flag": {
		"errors",  // flag
		"fmt",     // flag
		"io",      // flag
		"os",      // flag
		"reflect", // flag
		"sort",    // flag
		"strconv", // flag
		"strings", // flag
		"time",    // flag
	},

	"fmt": {
		"errors",       // fmt
		"io",           // fmt
		"math",         // fmt
		"os",           // fmt
		"reflect",      // fmt
		"strconv",      // fmt
		"sync",         // fmt
		"unicode/utf8", // fmt
	},

	"go/ast": {
		"bytes",        // go/ast
		"fmt",          // go/ast
		"go/scanner",   // go/ast
		"go/token",     // go/ast
		"io",           // go/ast
		"os",           // go/ast
		"reflect",      // go/ast
		"sort",         // go/ast
		"strconv",      // go/ast
		"strings",      // go/ast
		"unicode",      // go/ast
		"unicode/utf8", // go/ast
	},

	"go/build": {
		"bufio",         // go/build
		"bytes",         // go/build
		"errors",        // go/build
		"fmt",           // go/build
		"go/ast",        // go/build
		"go/doc",        // go/build
		"go/parser",     // go/build
		"go/token",      // go/build
		"io",            // go/build
		"io/ioutil",     // go/build
		"log",           // go/build
		"os",            // go/build
		"path",          // go/build
		"path/filepath", // go/build
		"runtime",       // go/build
		"sort",          // go/build
		"strconv",       // go/build
		"strings",       // go/build
		"unicode",       // go/build
		"unicode/utf8",  // go/build
	},

	"go/doc": {
		"go/ast",        // go/doc
		"go/token",      // go/doc
		"io",            // go/doc
		"path",          // go/doc
		"regexp",        // go/doc
		"sort",          // go/doc
		"strconv",       // go/doc
		"strings",       // go/doc
		"text/template", // go/doc
		"unicode",       // go/doc
		"unicode/utf8",  // go/doc
	},

	"go/parser": {
		"bytes",         // go/parser
		"errors",        // go/parser
		"fmt",           // go/parser
		"go/ast",        // go/parser
		"go/scanner",    // go/parser
		"go/token",      // go/parser
		"io",            // go/parser
		"io/ioutil",     // go/parser
		"os",            // go/parser
		"path/filepath", // go/parser
		"strconv",       // go/parser
		"strings",       // go/parser
		"unicode",       // go/parser
	},

	"go/scanner": {
		"bytes",         // go/scanner
		"fmt",           // go/scanner
		"go/token",      // go/scanner
		"io",            // go/scanner
		"path/filepath", // go/scanner
		"sort",          // go/scanner
		"strconv",       // go/scanner
		"unicode",       // go/scanner
		"unicode/utf8",  // go/scanner
	},

	"go/token": {
		"fmt",     // go/token
		"sort",    // go/token
		"strconv", // go/token
		"sync",    // go/token
	},

	"hash": {
		"io", // hash
	},

	"hash/adler32": {
		"hash", // hash/adler32
	},

	"internal/cpu": {
		"runtime", // internal/cpu
	},

	"internal/poll": {
		"errors",        // internal/poll
		"internal/race", // internal/poll
		"io",            // internal/poll
		"runtime",       // internal/poll
		"sync",          // internal/poll
		"sync/atomic",   // internal/poll
		"syscall",       // internal/poll
		"time",          // internal/poll
		"unicode/utf16", // internal/poll
		"unicode/utf8",  // internal/poll
	},

	"internal/race": {
		"runtime", // internal/race
	},

	"internal/singleflight": {
		"sync", // internal/singleflight
	},

	"internal/syscall/windows": {
		"internal/syscall/windows/sysdll", // internal/syscall/windows
		"syscall",                         // internal/syscall/windows
	},

	"internal/syscall/windows/registry": {
		"errors", // internal/syscall/windows/registry
		"internal/syscall/windows/sysdll", // internal/syscall/windows/registry
		"io",            // internal/syscall/windows/registry
		"syscall",       // internal/syscall/windows/registry
		"unicode/utf16", // internal/syscall/windows/registry
	},

	"internal/syscall/windows/sysdll": {
		"runtime", // internal/syscall/windows/sysdll
	},

	"io": {
		"errors", // io
		"sync",   // io
	},

	"io/ioutil": {
		"bytes",         // io/ioutil
		"io",            // io/ioutil
		"os",            // io/ioutil
		"path/filepath", // io/ioutil
		"sort",          // io/ioutil
		"strconv",       // io/ioutil
		"sync",          // io/ioutil
		"time",          // io/ioutil
	},

	"log": {
		"fmt",     // log
		"io",      // log
		"os",      // log
		"runtime", // log
		"sync",    // log
		"time",    // log
	},

	"math": {
		"internal/cpu", // math
	},

	"math/bits": {
		"runtime", // math/bits
	},

	"net/url": {
		"bytes",   // net/url
		"errors",  // net/url
		"fmt",     // net/url
		"sort",    // net/url
		"strconv", // net/url
		"strings", // net/url
	},

	"os": {
		"errors",                   // os
		"internal/poll",            // os
		"internal/syscall/windows", // os
		"io",            // os
		"runtime",       // os
		"sync",          // os
		"sync/atomic",   // os
		"syscall",       // os
		"time",          // os
		"unicode/utf16", // os
	},

	"os/exec": {
		"bytes",         // os/exec
		"context",       // os/exec
		"errors",        // os/exec
		"io",            // os/exec
		"os",            // os/exec
		"path/filepath", // os/exec
		"runtime",       // os/exec
		"strconv",       // os/exec
		"strings",       // os/exec
		"sync",          // os/exec
		"syscall",       // os/exec
	},

	"os/signal": {
		"os",      // os/signal
		"sync",    // os/signal
		"syscall", // os/signal
	},

	"path": {
		"errors",       // path
		"strings",      // path
		"unicode/utf8", // path
	},

	"path/filepath": {
		"errors",                   // path/filepath
		"internal/syscall/windows", // path/filepath
		"os",           // path/filepath
		"runtime",      // path/filepath
		"sort",         // path/filepath
		"strings",      // path/filepath
		"syscall",      // path/filepath
		"unicode/utf8", // path/filepath
	},

	"reflect": {
		"math",         // reflect
		"runtime",      // reflect
		"strconv",      // reflect
		"sync",         // reflect
		"unicode",      // reflect
		"unicode/utf8", // reflect
	},

	"regexp": {
		"bytes",         // regexp
		"io",            // regexp
		"regexp/syntax", // regexp
		"sort",          // regexp
		"strconv",       // regexp
		"strings",       // regexp
		"sync",          // regexp
		"unicode",       // regexp
		"unicode/utf8",  // regexp
	},

	"regexp/syntax": {
		"bytes",        // regexp/syntax
		"sort",         // regexp/syntax
		"strconv",      // regexp/syntax
		"strings",      // regexp/syntax
		"unicode",      // regexp/syntax
		"unicode/utf8", // regexp/syntax
	},

	"runtime": {
		"runtime/internal/atomic", // runtime
		"runtime/internal/sys",    // runtime
	},

	"runtime/internal/atomic": {
		"runtime/internal/sys", // runtime/internal/atomic
	},

	"runtime/internal/sys": {},

	"sort": {
		"reflect", // sort
	},

	"strconv": {
		"errors",       // strconv
		"math",         // strconv
		"unicode/utf8", // strconv
	},

	"strings": {
		"errors",       // strings
		"internal/cpu", // strings
		"io",           // strings
		"unicode",      // strings
		"unicode/utf8", // strings
	},

	"sync": {
		"internal/race", // sync
		"runtime",       // sync
		"sync/atomic",   // sync
	},

	"sync/atomic": {
		"runtime", // sync/atomic
	},

	"syscall": {
		"errors",                          // syscall
		"internal/race",                   // syscall
		"internal/syscall/windows/sysdll", // syscall
		"runtime",                         // syscall
		"sync",                            // syscall
		"sync/atomic",                     // syscall
		"unicode/utf16",                   // syscall
	},

	"text/template": {
		"bytes",               // text/template
		"errors",              // text/template
		"fmt",                 // text/template
		"io",                  // text/template
		"io/ioutil",           // text/template
		"net/url",             // text/template
		"path/filepath",       // text/template
		"reflect",             // text/template
		"runtime",             // text/template
		"sort",                // text/template
		"strings",             // text/template
		"sync",                // text/template
		"text/template/parse", // text/template
		"unicode",             // text/template
		"unicode/utf8",        // text/template
	},

	"text/template/parse": {
		"bytes",        // text/template/parse
		"fmt",          // text/template/parse
		"runtime",      // text/template/parse
		"strconv",      // text/template/parse
		"strings",      // text/template/parse
		"unicode",      // text/template/parse
		"unicode/utf8", // text/template/parse
	},

	"time": {
		"errors", // time
		"internal/syscall/windows/registry", // time
		"runtime",                           // time
		"sync",                              // time
		"syscall",                           // time
	},

	"unicode": {
		"runtime", // unicode
	},

	"unicode/utf16": {
		"runtime", // unicode/utf16
	},

	"unicode/utf8": {
		"runtime", // unicode/utf8
	},
}
