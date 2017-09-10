// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package x509

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	// testNameConstraintsAgainstOpenSSL can be set to true to run tests
	// against the system OpenSSL. This is disabled by default because Go
	// cannot depend on having OpenSSL installed at testing time.
	testNameConstraintsAgainstOpenSSL = false

	// debugOpenSSLFailure can be set to true, when
	// testNameConstraintsAgainstOpenSSL is also true, to cause
	// intermediate files to be preserved for debugging.
	debugOpenSSLFailure = false
)

type nameConstraintsTest struct {
	roots         []constraintsSpec
	intermediates [][]constraintsSpec
	leaf          []string
	expectedError string
	noOpenSSL     bool
}

type constraintsSpec struct {
	ok  []string
	bad []string
}

var nameConstraintsTests = []nameConstraintsTest{
	// #0: dummy test for the certificate generation process itself.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		leaf: []string{"dns:example.com"},
	},

	// #1: dummy test for the certificate generation process itself: single
	// level of intermediate.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:example.com"},
	},

	// #2: dummy test for the certificate generation process itself: two
	// levels of intermediates.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:example.com"},
	},

	// #3: matching DNS constraint in root
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:example.com"},
	},

	// #4: matching DNS constraint in intermediate.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					ok: []string{"dns:example.com"},
				},
			},
		},
		leaf: []string{"dns:example.com"},
	},

	// #5: .example.com only matches subdomains.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"dns:example.com"},
		expectedError: "\"example.com\" is not permitted",
	},

	// #6: .example.com matches subdomains.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					ok: []string{"dns:.example.com"},
				},
			},
		},
		leaf: []string{"dns:foo.example.com"},
	},

	// #7: .example.com matches multiple levels of subdomains
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:foo.bar.example.com"},
	},

	// #8: specifying a permitted list of names does not exclude other name
	// types
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"ip:10.1.1.1"},
	},

	// #9: specifying a permitted list of names does not exclude other name
	// types
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"ip:10.0.0.0/8"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:example.com"},
	},

	// #10: intermediates can try to permit other names, which isn't
	// forbidden if the leaf doesn't mention them. I.e. name constraints
	// apply to names, not constraints themselves.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					ok: []string{"dns:example.com", "dns:foo.com"},
				},
			},
		},
		leaf: []string{"dns:example.com"},
	},

	// #11: intermediates cannot add permitted names that the root doesn't
	// grant them.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					ok: []string{"dns:example.com", "dns:foo.com"},
				},
			},
		},
		leaf:          []string{"dns:foo.com"},
		expectedError: "\"foo.com\" is not permitted",
	},

	// #12: intermediates can further limit their scope if they wish.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					ok: []string{"dns:.bar.example.com"},
				},
			},
		},
		leaf: []string{"dns:foo.bar.example.com"},
	},

	// #13: intermediates can further limit their scope and that limitation
	// is effective
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					ok: []string{"dns:.bar.example.com"},
				},
			},
		},
		leaf:          []string{"dns:foo.notbar.example.com"},
		expectedError: "\"foo.notbar.example.com\" is not permitted",
	},

	// #14: roots can exclude subtrees and that doesn't affect other names.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"dns:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:foo.com"},
	},

	// #15: roots exclusions are effective.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"dns:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"dns:foo.example.com"},
		expectedError: "\"foo.example.com\" is excluded",
	},

	// #16: intermediates can also exclude names and that doesn't affect
	// other names.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					bad: []string{"dns:.example.com"},
				},
			},
		},
		leaf: []string{"dns:foo.com"},
	},

	// #17: intermediate exclusions are effective.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					bad: []string{"dns:.example.com"},
				},
			},
		},
		leaf:          []string{"dns:foo.example.com"},
		expectedError: "\"foo.example.com\" is excluded",
	},

	// #18: having an exclusion doesn't prohibit other types of names.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"dns:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:foo.com", "ip:10.1.1.1"},
	},

	// #19: IP-based exclusions are permitted and don't affect unrelated IP
	// addresses.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"ip:10.0.0.0/8"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"ip:192.168.1.1"},
	},

	// #20: IP-based exclusions are effective
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"ip:10.0.0.0/8"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"ip:10.0.0.1"},
		expectedError: "\"10.0.0.1\" is excluded",
	},

	// #21: intermediates can further constrain IP ranges.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"ip:0.0.0.0/1"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					bad: []string{"ip:11.0.0.0/8"},
				},
			},
		},
		leaf:          []string{"ip:11.0.0.1"},
		expectedError: "\"11.0.0.1\" is excluded",
	},

	// #22: when multiple intermediates are present, chain building can
	// avoid intermediates with incompatible constraints.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					ok: []string{"dns:.foo.com"},
				},
				constraintsSpec{
					ok: []string{"dns:.example.com"},
				},
			},
		},
		leaf:      []string{"dns:foo.example.com"},
		noOpenSSL: true, // OpenSSL's chain building is not informed by constraints.
	},

	// #23: (same as the previous test, but in the other order in ensure
	// that we don't pass it by luck.)
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{
					ok: []string{"dns:.example.com"},
				},
				constraintsSpec{
					ok: []string{"dns:.foo.com"},
				},
			},
		},
		leaf:      []string{"dns:foo.example.com"},
		noOpenSSL: true, // OpenSSL's chain building is not informed by constraints.
	},

	// #24: when multiple roots are valid, chain building can avoid roots
	// with incompatible constraints.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
			constraintsSpec{
				ok: []string{"dns:foo.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:      []string{"dns:example.com"},
		noOpenSSL: true, // OpenSSL's chain building is not informed by constraints.
	},

	// #25: (same as the previous test, but in the other order in ensure
	// that we don't pass it by luck.)
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:foo.com"},
			},
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:      []string{"dns:example.com"},
		noOpenSSL: true, // OpenSSL's chain building is not informed by constraints.
	},

	// #26: chain building can find a valid path even with multiple levels
	// of alternative intermediates and alternative roots.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:foo.com"},
			},
			constraintsSpec{
				ok: []string{"dns:example.com"},
			},
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
				constraintsSpec{
					ok: []string{"dns:foo.com"},
				},
			},
			[]constraintsSpec{
				constraintsSpec{},
				constraintsSpec{
					ok: []string{"dns:foo.com"},
				},
			},
		},
		leaf:      []string{"dns:bar.com"},
		noOpenSSL: true, // OpenSSL's chain building is not informed by constraints.
	},

	// #27: chain building doesn't get stuck when there is no valid path.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:foo.com"},
			},
			constraintsSpec{
				ok: []string{"dns:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
				constraintsSpec{
					ok: []string{"dns:foo.com"},
				},
			},
			[]constraintsSpec{
				constraintsSpec{
					ok: []string{"dns:bar.com"},
				},
				constraintsSpec{
					ok: []string{"dns:foo.com"},
				},
			},
		},
		leaf:          []string{"dns:bar.com"},
		expectedError: "\"bar.com\" is not permitted",
	},

	// #28: unknown name types don't cause a problem without constraints.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"unknown:"},
	},

	// #29: unknown name types are allowed even in constrained chains.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:foo.com", "dns:.foo.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"unknown:"},
	},

	// #30: without SANs, a certificate is rejected in a constrained chain.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:foo.com", "dns:.foo.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{},
		expectedError: "leaf doesn't have a SAN extension",
		noOpenSSL:     true, // OpenSSL doesn't require SANs in this case.
	},

	// #31: IPv6 addresses work in constraints: roots can permit them as
	// expected.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"ip:2000:abcd::/32"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"ip:2000:abcd:1234::"},
	},

	// #32: IPv6 addresses work in constraints: root restrictions are
	// effective.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"ip:2000:abcd::/32"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"ip:2000:1234:abcd::"},
		expectedError: "\"2000:1234:abcd::\" is not permitted",
	},

	// #33: An IPv6 permitted subtree doesn't affect DNS names.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"ip:2000:abcd::/32"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"ip:2000:abcd::", "dns:foo.com"},
	},

	// #34: IPv6 exclusions don't affect unrelated addresses.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"ip:2000:abcd::/32"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"ip:2000:1234::"},
	},

	// #35: IPv6 exclusions are effective.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"ip:2000:abcd::/32"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"ip:2000:abcd::"},
		expectedError: "\"2000:abcd::\" is excluded",
	},

	// #36: IPv6 constraints do not permit IPv4 addresses.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"ip:2000:abcd::/32"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"ip:10.0.0.1"},
		expectedError: "\"10.0.0.1\" is not permitted",
	},

	// #37: IPv4 constraints do not permit IPv6 addresses.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"ip:10.0.0.0/8"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"ip:2000:abcd::"},
		expectedError: "\"2000:abcd::\" is not permitted",
	},

	// #38: an exclusion of an unknown type doesn't affect other names.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"unknown:"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:example.com"},
	},

	// #39: a permitted subtree of an unknown type doesn't affect other
	// name types.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"unknown:"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:example.com"},
	},

	// #40: exact email constraints work
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"email:foo@example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"email:foo@example.com"},
	},

	// #41: exact email constraints are effective
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"email:foo@example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"email:bar@example.com"},
		expectedError: "\"bar@example.com\" is not permitted",
	},

	// #42: email canonicalisation works.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"email:foo@example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:      []string{"email:\"\\f\\o\\o\"@example.com"},
		noOpenSSL: true, // OpenSSL doesn't canonicalise email addresses before matching
	},

	// #43: limiting email addresses to a host works.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"email:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"email:foo@example.com"},
	},

	// #44: a leading dot matches hosts one level deep
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"email:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"email:foo@sub.example.com"},
	},

	// #45: a leading dot does not match the host itself
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"email:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"email:foo@example.com"},
		expectedError: "\"foo@example.com\" is not permitted",
	},

	// #46: a leading dot also matches two (or more) levels deep.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"email:.example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"email:foo@sub.sub.example.com"},
	},

	// #47: the local part of an email is case-sensitive
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"email:foo@example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"email:Foo@example.com"},
		expectedError: "\"Foo@example.com\" is not permitted",
	},

	// #48: the domain part of an email is not case-sensitive
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"email:foo@EXAMPLE.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"email:foo@example.com"},
	},

	// #49: the domain part of a DNS constraint is also not case-sensitive.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"dns:EXAMPLE.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"dns:example.com"},
	},

	// #50: URI constraints only cover the host part of the URI
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"uri:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{
			"uri:http://example.com/bar",
			"uri:http://example.com:8080/",
			"uri:https://example.com/wibble#bar",
		},
	},

	// #51: URIs with IPs are rejected
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"uri:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"uri:http://1.2.3.4/"},
		expectedError: "URI with IP",
	},

	// #52: URIs with IPs and ports are rejected
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"uri:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"uri:http://1.2.3.4:43/"},
		expectedError: "URI with IP",
	},

	// #53: URIs with IPv6 addresses are also rejected
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"uri:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"uri:http://[2006:abcd::1]/"},
		expectedError: "URI with IP",
	},

	// #54: URIs with IPv6 addresses with ports are also rejected
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"uri:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"uri:http://[2006:abcd::1]:16/"},
		expectedError: "URI with IP",
	},

	// #55: URI constraints are effective
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"uri:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"uri:http://bar.com/"},
		expectedError: "\"http://bar.com/\" is not permitted",
	},

	// #56: URI constraints are effective
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"uri:foo.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"uri:http://foo.com/"},
		expectedError: "\"http://foo.com/\" is excluded",
	},

	// #57: URI constraints can allow subdomains
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"uri:.foo.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"uri:http://www.foo.com/"},
	},

	// #58: excluding an IPv4-mapped-IPv6 address doesn't affect the IPv4
	// version of that address.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				bad: []string{"ip:::ffff:1.2.3.4/128"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"ip:1.2.3.4"},
	},

	// #59: a URI constraint isn't matched by a URN.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok: []string{"uri:example.com"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf:          []string{"uri:urn:example"},
		expectedError: "URI with empty host",
	},

	// #60: excluding all IPv6 addresses doesn't exclude all IPv4 addresses
	// too, even though IPv4 is mapped into the IPv6 range.
	nameConstraintsTest{
		roots: []constraintsSpec{
			constraintsSpec{
				ok:  []string{"ip:1.2.3.0/24"},
				bad: []string{"ip:::0/0"},
			},
		},
		intermediates: [][]constraintsSpec{
			[]constraintsSpec{
				constraintsSpec{},
			},
		},
		leaf: []string{"ip:1.2.3.4"},
	},

	// TODO(agl): handle empty name constraints. Currently this doesn't
	// work because empty values are treated as missing.
}

func makeConstraintsCACert(constraints constraintsSpec, name string, key *ecdsa.PrivateKey, parent *Certificate, parentKey *ecdsa.PrivateKey) (*Certificate, error) {
	var serialBytes [16]byte
	rand.Read(serialBytes[:])

	template := &Certificate{
		SerialNumber: new(big.Int).SetBytes(serialBytes[:]),
		Subject: pkix.Name{
			CommonName: name,
		},
		NotBefore:             time.Unix(1000, 0),
		NotAfter:              time.Unix(2000, 0),
		KeyUsage:              KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA: true,
	}

	if err := addConstraintsToTemplate(constraints, template); err != nil {
		return nil, err
	}

	if parent == nil {
		parent = template
	}
	derBytes, err := CreateCertificate(rand.Reader, template, parent, &key.PublicKey, parentKey)
	if err != nil {
		return nil, err
	}

	caCert, err := ParseCertificate(derBytes)
	if err != nil {
		return nil, err
	}

	return caCert, nil
}

func makeConstraintsLeafCert(sans []string, key *ecdsa.PrivateKey, parent *Certificate, parentKey *ecdsa.PrivateKey) (*Certificate, error) {
	var serialBytes [16]byte
	rand.Read(serialBytes[:])

	template := &Certificate{
		SerialNumber: new(big.Int).SetBytes(serialBytes[:]),
		Subject: pkix.Name{
			// Don't set a CommonName because OpenSSL (at least) will try to
			// match it against name constraints.
			OrganizationalUnit: []string{"Leaf"},
		},
		NotBefore:             time.Unix(1000, 0),
		NotAfter:              time.Unix(2000, 0),
		KeyUsage:              KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA: false,
	}

	for _, name := range sans {
		switch {
		case strings.HasPrefix(name, "dns:"):
			template.DNSNames = append(template.DNSNames, name[4:])

		case strings.HasPrefix(name, "ip:"):
			ip := net.ParseIP(name[3:])
			if ip == nil {
				return nil, fmt.Errorf("cannot parse IP %q", name[3:])
			}
			template.IPAddresses = append(template.IPAddresses, ip)

		case strings.HasPrefix(name, "email:"):
			template.EmailAddresses = append(template.EmailAddresses, name[6:])

		case strings.HasPrefix(name, "uri:"):
			uri, err := url.Parse(name[4:])
			if err != nil {
				return nil, fmt.Errorf("cannot parse URI %q: %s", name[4:], err)
			}
			template.URIs = append(template.URIs, uri)

		case strings.HasPrefix(name, "unknown:"):
			// This is a special case for testing unknown
			// name types. A custom SAN extension is
			// injected into the certificate.
			if len(sans) != 1 {
				panic("when using unknown name types, it must be the sole name")
			}

			template.ExtraExtensions = append(template.ExtraExtensions, pkix.Extension{
				Id: []int{2, 5, 29, 17},
				Value: []byte{
					0x30, // SEQUENCE
					3,    // three bytes
					9,    // undefined GeneralName type 9
					1,
					1,
				},
			})

		default:
			return nil, fmt.Errorf("unknown name type %q", name)
		}
	}

	if parent == nil {
		parent = template
	}

	derBytes, err := CreateCertificate(rand.Reader, template, parent, &key.PublicKey, parentKey)
	if err != nil {
		return nil, err
	}

	return ParseCertificate(derBytes)
}

func customConstraintsExtension(typeNum int, constraint []byte, isExcluded bool) pkix.Extension {
	appendConstraint := func(contents []byte, tag uint8) []byte {
		contents = append(contents, tag|32 /* constructed */ |0x80 /* context-specific */)
		contents = append(contents, byte(4+len(constraint)) /* length */)
		contents = append(contents, 0x30 /* SEQUENCE */)
		contents = append(contents, byte(2+len(constraint)) /* length */)
		contents = append(contents, byte(typeNum) /* GeneralName type */)
		contents = append(contents, byte(len(constraint)))
		return append(contents, constraint...)
	}

	var contents []byte
	if !isExcluded {
		contents = appendConstraint(contents, 0 /* tag 0 for permitted */)
	} else {
		contents = appendConstraint(contents, 1 /* tag 1 for excluded */)
	}

	var value []byte
	value = append(value, 0x30 /* SEQUENCE */)
	value = append(value, byte(len(contents)))
	value = append(value, contents...)

	return pkix.Extension{
		Id:    []int{2, 5, 29, 30},
		Value: value,
	}
}

func addConstraintsToTemplate(constraints constraintsSpec, template *Certificate) error {
	parse := func(constraints []string) (dnsNames []string, ips []*net.IPNet, emailAddrs []string, uriDomains []string, err error) {
		for _, constraint := range constraints {
			switch {
			case strings.HasPrefix(constraint, "dns:"):
				dnsNames = append(dnsNames, constraint[4:])

			case strings.HasPrefix(constraint, "ip:"):
				_, ipNet, err := net.ParseCIDR(constraint[3:])
				if err != nil {
					return nil, nil, nil, nil, err
				}
				ips = append(ips, ipNet)

			case strings.HasPrefix(constraint, "email:"):
				emailAddrs = append(emailAddrs, constraint[6:])

			case strings.HasPrefix(constraint, "uri:"):
				uriDomains = append(uriDomains, constraint[4:])

			default:
				return nil, nil, nil, nil, fmt.Errorf("unknown constraint %q", constraint)
			}
		}

		return dnsNames, ips, emailAddrs, uriDomains, err
	}

	handleSpecialConstraint := func(constraint string, isExcluded bool) bool {
		switch {
		case constraint == "unknown:":
			template.ExtraExtensions = append(template.ExtraExtensions, customConstraintsExtension(9 /* undefined GeneralName type */, []byte{1}, isExcluded))

		default:
			return false
		}

		return true
	}

	if len(constraints.ok) == 1 && len(constraints.bad) == 0 {
		if handleSpecialConstraint(constraints.ok[0], false) {
			return nil
		}
	}

	if len(constraints.bad) == 1 && len(constraints.ok) == 0 {
		if handleSpecialConstraint(constraints.bad[0], true) {
			return nil
		}
	}

	var err error
	template.PermittedDNSDomains, template.PermittedIPRanges, template.PermittedEmailAddresses, template.PermittedURIDomains, err = parse(constraints.ok)
	if err != nil {
		return err
	}

	template.ExcludedDNSDomains, template.ExcludedIPRanges, template.ExcludedEmailAddresses, template.ExcludedURIDomains, err = parse(constraints.bad)
	if err != nil {
		return err
	}

	return nil
}

func TestNameConstraintCases(t *testing.T) {
	privateKeys := sync.Pool{
		New: func() interface{} {
			priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				panic(err)
			}
			return priv
		},
	}

	for i, test := range nameConstraintsTests {
		rootPool := NewCertPool()
		rootKey := privateKeys.Get().(*ecdsa.PrivateKey)
		rootName := "Root " + strconv.Itoa(i)

		// keys keeps track of all the private keys used in a given
		// test and puts them back in the privateKeys pool at the end.
		keys := []*ecdsa.PrivateKey{rootKey}

		// At each level (root, intermediate(s), leaf), parent points to
		// an example parent certificate and parentKey the key for the
		// parent level. Since all certificates at a given level have
		// the same name and public key, any parent certificate is
		// sufficient to get the correct issuer name and authority
		// key ID.
		var parent *Certificate
		parentKey := rootKey

		for _, root := range test.roots {
			rootCert, err := makeConstraintsCACert(root, rootName, rootKey, nil, rootKey)
			if err != nil {
				t.Fatalf("#%d: failed to create root: %s", i, err)
			}

			parent = rootCert
			rootPool.AddCert(rootCert)
		}

		intermediatePool := NewCertPool()

		for level, intermediates := range test.intermediates {
			levelKey := privateKeys.Get().(*ecdsa.PrivateKey)
			keys = append(keys, levelKey)
			levelName := "Intermediate level " + strconv.Itoa(level)
			var last *Certificate

			for _, intermediate := range intermediates {
				caCert, err := makeConstraintsCACert(intermediate, levelName, levelKey, parent, parentKey)
				if err != nil {
					t.Fatalf("#%d: failed to create %q: %s", i, levelName, err)
				}

				last = caCert
				intermediatePool.AddCert(caCert)
			}

			parent = last
			parentKey = levelKey
		}

		leafKey := privateKeys.Get().(*ecdsa.PrivateKey)
		keys = append(keys, leafKey)

		leafCert, err := makeConstraintsLeafCert(test.leaf, leafKey, parent, parentKey)
		if err != nil {
			t.Fatalf("#%d: cannot create leaf: %s", i, err)
		}

		if !test.noOpenSSL && testNameConstraintsAgainstOpenSSL {
			output, err := testChainAgainstOpenSSL(leafCert, intermediatePool, rootPool)
			if err == nil && len(test.expectedError) > 0 {
				t.Errorf("#%d: unexpectedly succeeded against OpenSSL", i)
				if debugOpenSSLFailure {
					return
				}
			}

			if err != nil {
				if _, ok := err.(*exec.ExitError); !ok {
					t.Errorf("#%d: OpenSSL failed to run: %s", i, err)
				} else if len(test.expectedError) == 0 {
					t.Errorf("#%d: OpenSSL unexpectedly failed: %q", i, output)
					if debugOpenSSLFailure {
						return
					}
				}
			}
		}

		verifyOpts := VerifyOptions{
			Roots:         rootPool,
			Intermediates: intermediatePool,
			CurrentTime:   time.Unix(1500, 0),
		}
		_, err = leafCert.Verify(verifyOpts)

		logInfo := true
		if len(test.expectedError) == 0 {
			if err != nil {
				t.Errorf("#%d: unexpected failure: %s", i, err)
			} else {
				logInfo = false
			}
		} else {
			if err == nil {
				t.Errorf("#%d: unexpected success", i)
			} else if !strings.Contains(err.Error(), test.expectedError) {
				t.Errorf("#%d: expected error containing %q, but got: %s", i, test.expectedError, err)
			} else {
				logInfo = false
			}
		}

		if logInfo {
			certAsPEM := func(cert *Certificate) string {
				var buf bytes.Buffer
				pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
				return string(buf.Bytes())
			}
			t.Errorf("#%d: root:\n%s", i, certAsPEM(rootPool.certs[0]))
			t.Errorf("#%d: leaf:\n%s", i, certAsPEM(leafCert))
		}

		for _, key := range keys {
			privateKeys.Put(key)
		}
		keys = keys[:0]
	}
}

func writePEMsToTempFile(certs []*Certificate) *os.File {
	file, err := ioutil.TempFile("", "name_constraints_test")
	if err != nil {
		panic("cannot create tempfile")
	}

	pemBlock := &pem.Block{Type: "CERTIFICATE"}
	for _, cert := range certs {
		pemBlock.Bytes = cert.Raw
		pem.Encode(file, pemBlock)
	}

	return file
}

func testChainAgainstOpenSSL(leaf *Certificate, intermediates, roots *CertPool) (string, error) {
	args := []string{"verify", "-no_check_time"}

	rootsFile := writePEMsToTempFile(roots.certs)
	if debugOpenSSLFailure {
		println("roots file:", rootsFile.Name())
	} else {
		defer os.Remove(rootsFile.Name())
	}
	args = append(args, "-CAfile", rootsFile.Name())

	if len(intermediates.certs) > 0 {
		intermediatesFile := writePEMsToTempFile(intermediates.certs)
		if debugOpenSSLFailure {
			println("intermediates file:", intermediatesFile.Name())
		} else {
			defer os.Remove(intermediatesFile.Name())
		}
		args = append(args, "-untrusted", intermediatesFile.Name())
	}

	leafFile := writePEMsToTempFile([]*Certificate{leaf})
	if debugOpenSSLFailure {
		println("leaf file:", leafFile.Name())
	} else {
		defer os.Remove(leafFile.Name())
	}
	args = append(args, leafFile.Name())

	var output bytes.Buffer
	cmd := exec.Command("openssl", args...)
	cmd.Stdout = &output
	cmd.Stderr = &output

	err := cmd.Run()
	return string(output.Bytes()), err
}

var rfc2821Tests = []struct {
	in                string
	localPart, domain string
}{
	{"foo@example.com", "foo", "example.com"},
	{"@example.com", "", ""},
	{"\"@example.com", "", ""},
	{"\"\"@example.com", "", "example.com"},
	{"\"a\"@example.com", "a", "example.com"},
	{"\"\\a\"@example.com", "a", "example.com"},
	{"a\"@example.com", "", ""},
	{"foo..bar@example.com", "", ""},
	{".foo.bar@example.com", "", ""},
	{"foo.bar.@example.com", "", ""},
	{"|{}?'@example.com", "|{}?'", "example.com"},

	// Examples from RFC 3696
	{"Abc\\@def@example.com", "Abc@def", "example.com"},
	{"Fred\\ Bloggs@example.com", "Fred Bloggs", "example.com"},
	{"Joe.\\\\Blow@example.com", "Joe.\\Blow", "example.com"},
	{"\"Abc@def\"@example.com", "Abc@def", "example.com"},
	{"\"Fred Bloggs\"@example.com", "Fred Bloggs", "example.com"},
	{"customer/department=shipping@example.com", "customer/department=shipping", "example.com"},
	{"$A12345@example.com", "$A12345", "example.com"},
	{"!def!xyz%abc@example.com", "!def!xyz%abc", "example.com"},
	{"_somename@example.com", "_somename", "example.com"},
}

func TestRFC2821Parsing(t *testing.T) {
	for i, test := range rfc2821Tests {
		mailbox, ok := parseRFC2821Mailbox(test.in)
		expectedFailure := len(test.localPart) == 0 && len(test.domain) == 0

		if ok && expectedFailure {
			t.Errorf("#%d: %q unexpectedly parsed as (%q, %q)", i, test.in, mailbox.local, mailbox.domain)
			continue
		}

		if !ok && !expectedFailure {
			t.Errorf("#%d: unexpected failure for %q", i, test.in)
			continue
		}

		if !ok {
			continue
		}

		if mailbox.local != test.localPart || mailbox.domain != test.domain {
			t.Errorf("#%d: %q parsed as (%q, %q), but wanted (%q, %q)", i, test.in, mailbox.local, mailbox.domain, test.localPart, test.domain)
		}
	}
}

func TestBadNamesInConstraints(t *testing.T) {
	// Bad names in constraints should not parse.
	badNames := []string{
		"dns:foo.com.",
		"email:abc@foo.com.",
		"email:foo.com.",
		"uri:example.com.",
		"uri:1.2.3.4",
		"uri:ffff::1",
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	for _, badName := range badNames {
		_, err := makeConstraintsCACert(constraintsSpec{
			ok: []string{badName},
		}, "TestAbsoluteNamesInConstraints", priv, nil, priv)

		if err == nil {
			t.Errorf("bad name %q unexpectedly accepted in name constraint", badName)
			continue
		}

		if err != nil {
			if str := err.Error(); !strings.Contains(str, "failed to parse ") && !strings.Contains(str, "constraint") {
				t.Errorf("bad name %q triggered unrecognised error: %s", badName, str)
			}
		}
	}
}

func TestBadNamesInSANs(t *testing.T) {
	// Bad names in SANs should not parse.
	badNames := []string{
		"dns:foo.com.",
		"email:abc@foo.com.",
		"email:foo.com.",
		"uri:https://example.com./dsf",
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	for _, badName := range badNames {
		_, err := makeConstraintsLeafCert([]string{badName}, priv, nil, priv)

		if err == nil {
			t.Errorf("bad name %q unexpectedly accepted in SAN", badName)
			continue
		}

		if err != nil {
			if str := err.Error(); !strings.Contains(str, "cannot parse ") {
				t.Errorf("bad name %q triggered unrecognised error: %s", badName, str)
			}
		}
	}
}
