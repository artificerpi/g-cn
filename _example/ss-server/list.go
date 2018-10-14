package main

import (
	"fmt"
	"net"
	"strings"
)

// CountSet stores host and count of times to access this element
type CountSet map[string]int

const MAX_SIZE = 1024

func (s *CountSet) Add(e string) {
	if !s.Contains(e) {
		map[string]int(*s)[e] = 1
	} else {
		map[string]int(*s)[e]++
	}

	if s.Size() > MAX_SIZE {
		go s.Shrink()
	}
}

func (s *CountSet) Remove(e string) {
	if s.Contains(e) {
		delete(*s, e)
	}
}

func (s *CountSet) Shrink() {
	for k, v := range *s {
		if v < 10 {
			s.Remove(k)
		}

		map[string]int(*s)[k]--
	}
}

func (s *CountSet) Contains(e string) bool {
	if map[string]int(*s)[e] > 0 {
		return true
	}

	return false
}

func (s *CountSet) Size() int {
	return len(*s)
}

// HostGroup store the domains, ips and cached hosts
type HostGroup struct {
	Domains     []string // domains and ip networks
	IPs         CountSet // ips
	CachedHosts CountSet // cached records
}

// InitDomainIPs inits the ips with existing domains
func (g *HostGroup) InitDomainIPs() {
	for _, domain := range whitelist.Domains {
		ips, _ := net.LookupIP(domain)
		for _, ip := range ips {
			whitelist.IPs.Add(ip.String())
		}
	}
}

// Cache caches the visited host
func (g *HostGroup) Cache(host string) {
	g.CachedHosts.Add(host)
}

// Contains check whether HostGroup contains the host
func (g *HostGroup) Contains(host string) bool {
	if g.CachedHosts.Contains(host) || g.IPs.Contains(host) {
		return true
	}

	for _, domain := range g.Domains {
		if strings.Contains(host, domain) {
			return true
		}
	}

	return false
}

func (g *HostGroup) Match(ip string) bool {
	addrs, err := net.LookupAddr(ip)
	if err != nil {
		for _, addr := range addrs {
			if g.Contains(addr) {
				return true
			}
		}
	}

	return false
}

func (g *HostGroup) String() {
	fmt.Sprintln("domains", g.Domains, " ips: ", g.IPs, " cache: ", g.CachedHosts)
}
