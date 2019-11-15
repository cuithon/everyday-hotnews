// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"context"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

const _WSAHOST_NOT_FOUND = syscall.Errno(11001)

func winError(call string, err error) error {
	switch err {
	case _WSAHOST_NOT_FOUND:
		return errNoSuchHost
	}
	return os.NewSyscallError(call, err)
}

func getprotobyname(name string) (proto int, err error) {
	p, err := syscall.GetProtoByName(name)
	if err != nil {
		return 0, winError("getprotobyname", err)
	}
	return int(p.Proto), nil
}

// lookupProtocol looks up IP protocol name and returns correspondent protocol number.
func lookupProtocol(ctx context.Context, name string) (int, error) {
	// GetProtoByName return value is stored in thread local storage.
	// Start new os thread before the call to prevent races.
	type result struct {
		proto int
		err   error
	}
	ch := make(chan result) // unbuffered
	go func() {
		acquireThread()
		defer releaseThread()
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		proto, err := getprotobyname(name)
		select {
		case ch <- result{proto: proto, err: err}:
		case <-ctx.Done():
		}
	}()
	select {
	case r := <-ch:
		if r.err != nil {
			if proto, err := lookupProtocolMap(name); err == nil {
				return proto, nil
			}

			dnsError := &DNSError{Err: r.err.Error(), Name: name}
			if r.err == errNoSuchHost {
				dnsError.IsNotFound = true
			}
			r.err = dnsError
		}
		return r.proto, r.err
	case <-ctx.Done():
		return 0, mapErr(ctx.Err())
	}
}

func (r *Resolver) lookupHost(ctx context.Context, name string) ([]string, error) {
	ips, err := r.lookupIP(ctx, "ip", name)
	if err != nil {
		return nil, err
	}
	addrs := make([]string, 0, len(ips))
	for _, ip := range ips {
		addrs = append(addrs, ip.String())
	}
	return addrs, nil
}

func (r *Resolver) lookupIP(ctx context.Context, network, name string) ([]IPAddr, error) {
	// TODO(bradfitz,brainman): use ctx more. See TODO below.

	var family int32 = syscall.AF_UNSPEC
	switch ipVersion(network) {
	case '4':
		family = syscall.AF_INET
	case '6':
		family = syscall.AF_INET6
	}

	getaddr := func() ([]IPAddr, error) {
		acquireThread()
		defer releaseThread()
		hints := syscall.AddrinfoW{
			Family:   family,
			Socktype: syscall.SOCK_STREAM,
			Protocol: syscall.IPPROTO_IP,
		}
		var result *syscall.AddrinfoW
		name16p, err := syscall.UTF16PtrFromString(name)
		if err != nil {
			return nil, &DNSError{Name: name, Err: err.Error()}
		}
		e := syscall.GetAddrInfoW(name16p, nil, &hints, &result)
		if e != nil {
			err := winError("getaddrinfow", e)
			dnsError := &DNSError{Err: err.Error(), Name: name}
			if err == errNoSuchHost {
				dnsError.IsNotFound = true
			}
			return nil, dnsError
		}
		defer syscall.FreeAddrInfoW(result)
		addrs := make([]IPAddr, 0, 5)
		for ; result != nil; result = result.Next {
			addr := unsafe.Pointer(result.Addr)
			switch result.Family {
			case syscall.AF_INET:
				a := (*syscall.RawSockaddrInet4)(addr).Addr
				addrs = append(addrs, IPAddr{IP: IPv4(a[0], a[1], a[2], a[3])})
			case syscall.AF_INET6:
				a := (*syscall.RawSockaddrInet6)(addr).Addr
				zone := zoneCache.name(int((*syscall.RawSockaddrInet6)(addr).Scope_id))
				addrs = append(addrs, IPAddr{IP: IP{a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9], a[10], a[11], a[12], a[13], a[14], a[15]}, Zone: zone})
			default:
				return nil, &DNSError{Err: syscall.EWINDOWS.Error(), Name: name}
			}
		}
		return addrs, nil
	}

	type ret struct {
		addrs []IPAddr
		err   error
	}

	var ch chan ret
	if ctx.Err() == nil {
		ch = make(chan ret, 1)
		go func() {
			addr, err := getaddr()
			ch <- ret{addrs: addr, err: err}
		}()
	}

	select {
	case r := <-ch:
		return r.addrs, r.err
	case <-ctx.Done():
		// TODO(bradfitz,brainman): cancel the ongoing
		// GetAddrInfoW? It would require conditionally using
		// GetAddrInfoEx with lpOverlapped, which requires
		// Windows 8 or newer. I guess we'll need oldLookupIP,
		// newLookupIP, and newerLookUP.
		//
		// For now we just let it finish and write to the
		// buffered channel.
		return nil, &DNSError{
			Name:      name,
			Err:       ctx.Err().Error(),
			IsTimeout: ctx.Err() == context.DeadlineExceeded,
		}
	}
}

func (r *Resolver) lookupPort(ctx context.Context, network, service string) (int, error) {
	if r.preferGo() {
		return lookupPortMap(network, service)
	}

	// TODO(bradfitz): finish ctx plumbing. Nothing currently depends on this.
	acquireThread()
	defer releaseThread()
	var stype int32
	switch network {
	case "tcp4", "tcp6":
		stype = syscall.SOCK_STREAM
	case "udp4", "udp6":
		stype = syscall.SOCK_DGRAM
	}
	hints := syscall.AddrinfoW{
		Family:   syscall.AF_UNSPEC,
		Socktype: stype,
		Protocol: syscall.IPPROTO_IP,
	}
	var result *syscall.AddrinfoW
	e := syscall.GetAddrInfoW(nil, syscall.StringToUTF16Ptr(service), &hints, &result)
	if e != nil {
		if port, err := lookupPortMap(network, service); err == nil {
			return port, nil
		}
		err := winError("getaddrinfow", e)
		dnsError := &DNSError{Err: err.Error(), Name: network + "/" + service}
		if err == errNoSuchHost {
			dnsError.IsNotFound = true
		}
		return 0, dnsError
	}
	defer syscall.FreeAddrInfoW(result)
	if result == nil {
		return 0, &DNSError{Err: syscall.EINVAL.Error(), Name: network + "/" + service}
	}
	addr := unsafe.Pointer(result.Addr)
	switch result.Family {
	case syscall.AF_INET:
		a := (*syscall.RawSockaddrInet4)(addr)
		return int(syscall.Ntohs(a.Port)), nil
	case syscall.AF_INET6:
		a := (*syscall.RawSockaddrInet6)(addr)
		return int(syscall.Ntohs(a.Port)), nil
	}
	return 0, &DNSError{Err: syscall.EINVAL.Error(), Name: network + "/" + service}
}

func (*Resolver) lookupCNAME(ctx context.Context, name string) (string, error) {
	// TODO(bradfitz): finish ctx plumbing. Nothing currently depends on this.
	acquireThread()
	defer releaseThread()
	var r *syscall.DNSRecord
	e := syscall.DnsQuery(name, syscall.DNS_TYPE_CNAME, 0, nil, &r, nil)
	// windows returns DNS_INFO_NO_RECORDS if there are no CNAME-s
	if errno, ok := e.(syscall.Errno); ok && errno == syscall.DNS_INFO_NO_RECORDS {
		// if there are no aliases, the canonical name is the input name
		return absDomainName([]byte(name)), nil
	}
	if e != nil {
		return "", &DNSError{Err: winError("dnsquery", e).Error(), Name: name}
	}
	defer syscall.DnsRecordListFree(r, 1)

	resolved := resolveCNAME(syscall.StringToUTF16Ptr(name), r)
	cname := syscall.UTF16ToString((*[256]uint16)(unsafe.Pointer(resolved))[:])
	return absDomainName([]byte(cname)), nil
}

func (*Resolver) lookupSRV(ctx context.Context, service, proto, name string) (string, []*SRV, error) {
	// TODO(bradfitz): finish ctx plumbing. Nothing currently depends on this.
	acquireThread()
	defer releaseThread()
	var target string
	if service == "" && proto == "" {
		target = name
	} else {
		target = "_" + service + "._" + proto + "." + name
	}
	var r *syscall.DNSRecord
	e := syscall.DnsQuery(target, syscall.DNS_TYPE_SRV, 0, nil, &r, nil)
	if e != nil {
		return "", nil, &DNSError{Err: winError("dnsquery", e).Error(), Name: target}
	}
	defer syscall.DnsRecordListFree(r, 1)

	srvs := make([]*SRV, 0, 10)
	for _, p := range validRecs(r, syscall.DNS_TYPE_SRV, target) {
		v := (*syscall.DNSSRVData)(unsafe.Pointer(&p.Data[0]))
		srvs = append(srvs, &SRV{absDomainName([]byte(syscall.UTF16ToString((*[256]uint16)(unsafe.Pointer(v.Target))[:]))), v.Port, v.Priority, v.Weight})
	}
	byPriorityWeight(srvs).sort()
	return absDomainName([]byte(target)), srvs, nil
}

func (*Resolver) lookupMX(ctx context.Context, name string) ([]*MX, error) {
	// TODO(bradfitz): finish ctx plumbing. Nothing currently depends on this.
	acquireThread()
	defer releaseThread()
	var r *syscall.DNSRecord
	e := syscall.DnsQuery(name, syscall.DNS_TYPE_MX, 0, nil, &r, nil)
	if e != nil {
		return nil, &DNSError{Err: winError("dnsquery", e).Error(), Name: name}
	}
	defer syscall.DnsRecordListFree(r, 1)

	mxs := make([]*MX, 0, 10)
	for _, p := range validRecs(r, syscall.DNS_TYPE_MX, name) {
		v := (*syscall.DNSMXData)(unsafe.Pointer(&p.Data[0]))
		mxs = append(mxs, &MX{absDomainName([]byte(syscall.UTF16ToString((*[256]uint16)(unsafe.Pointer(v.NameExchange))[:]))), v.Preference})
	}
	byPref(mxs).sort()
	return mxs, nil
}

func (*Resolver) lookupNS(ctx context.Context, name string) ([]*NS, error) {
	// TODO(bradfitz): finish ctx plumbing. Nothing currently depends on this.
	acquireThread()
	defer releaseThread()
	var r *syscall.DNSRecord
	e := syscall.DnsQuery(name, syscall.DNS_TYPE_NS, 0, nil, &r, nil)
	if e != nil {
		return nil, &DNSError{Err: winError("dnsquery", e).Error(), Name: name}
	}
	defer syscall.DnsRecordListFree(r, 1)

	nss := make([]*NS, 0, 10)
	for _, p := range validRecs(r, syscall.DNS_TYPE_NS, name) {
		v := (*syscall.DNSPTRData)(unsafe.Pointer(&p.Data[0]))
		nss = append(nss, &NS{absDomainName([]byte(syscall.UTF16ToString((*[256]uint16)(unsafe.Pointer(v.Host))[:])))})
	}
	return nss, nil
}

func (*Resolver) lookupTXT(ctx context.Context, name string) ([]string, error) {
	// TODO(bradfitz): finish ctx plumbing. Nothing currently depends on this.
	acquireThread()
	defer releaseThread()
	var r *syscall.DNSRecord
	e := syscall.DnsQuery(name, syscall.DNS_TYPE_TEXT, 0, nil, &r, nil)
	if e != nil {
		return nil, &DNSError{Err: winError("dnsquery", e).Error(), Name: name}
	}
	defer syscall.DnsRecordListFree(r, 1)

	txts := make([]string, 0, 10)
	for _, p := range validRecs(r, syscall.DNS_TYPE_TEXT, name) {
		d := (*syscall.DNSTXTData)(unsafe.Pointer(&p.Data[0]))
		s := ""
		for _, v := range (*[1 << 10]*uint16)(unsafe.Pointer(&(d.StringArray[0])))[:d.StringCount] {
			s += syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(v))[:])
		}
		txts = append(txts, s)
	}
	return txts, nil
}

func (*Resolver) lookupAddr(ctx context.Context, addr string) ([]string, error) {
	// TODO(bradfitz): finish ctx plumbing. Nothing currently depends on this.
	acquireThread()
	defer releaseThread()
	arpa, err := reverseaddr(addr)
	if err != nil {
		return nil, err
	}
	var r *syscall.DNSRecord
	e := syscall.DnsQuery(arpa, syscall.DNS_TYPE_PTR, 0, nil, &r, nil)
	if e != nil {
		return nil, &DNSError{Err: winError("dnsquery", e).Error(), Name: addr}
	}
	defer syscall.DnsRecordListFree(r, 1)

	ptrs := make([]string, 0, 10)
	for _, p := range validRecs(r, syscall.DNS_TYPE_PTR, arpa) {
		v := (*syscall.DNSPTRData)(unsafe.Pointer(&p.Data[0]))
		ptrs = append(ptrs, absDomainName([]byte(syscall.UTF16ToString((*[256]uint16)(unsafe.Pointer(v.Host))[:]))))
	}
	return ptrs, nil
}

const dnsSectionMask = 0x0003

// returns only results applicable to name and resolves CNAME entries
func validRecs(r *syscall.DNSRecord, dnstype uint16, name string) []*syscall.DNSRecord {
	cname := syscall.StringToUTF16Ptr(name)
	if dnstype != syscall.DNS_TYPE_CNAME {
		cname = resolveCNAME(cname, r)
	}
	rec := make([]*syscall.DNSRecord, 0, 10)
	for p := r; p != nil; p = p.Next {
		// in case of a local machine, DNS records are returned with DNSREC_QUESTION flag instead of DNS_ANSWER
		if p.Dw&dnsSectionMask != syscall.DnsSectionAnswer && p.Dw&dnsSectionMask != syscall.DnsSectionQuestion {
			continue
		}
		if p.Type != dnstype {
			continue
		}
		if !syscall.DnsNameCompare(cname, p.Name) {
			continue
		}
		rec = append(rec, p)
	}
	return rec
}

// returns the last CNAME in chain
func resolveCNAME(name *uint16, r *syscall.DNSRecord) *uint16 {
	// limit cname resolving to 10 in case of an infinite CNAME loop
Cname:
	for cnameloop := 0; cnameloop < 10; cnameloop++ {
		for p := r; p != nil; p = p.Next {
			if p.Dw&dnsSectionMask != syscall.DnsSectionAnswer {
				continue
			}
			if p.Type != syscall.DNS_TYPE_CNAME {
				continue
			}
			if !syscall.DnsNameCompare(name, p.Name) {
				continue
			}
			name = (*syscall.DNSPTRData)(unsafe.Pointer(&r.Data[0])).Host
			continue Cname
		}
		break
	}
	return name
}

// concurrentThreadsLimit returns the number of threads we permit to
// run concurrently doing DNS lookups.
func concurrentThreadsLimit() int {
	return 500
}
