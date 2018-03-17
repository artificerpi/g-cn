package main

import (
	"log"
	"net"
	"strings"
)

// CountSet stores host and count of times to access this element
type CountSet map[string]int

const MAX_SIZE = 1024 * 1024

func (s *CountSet) Add(e string) {
	if !s.Contains(e) {
		map[string]int(*s)[e] = 1
	}
}

func (s *CountSet) Remove(e string) {
	if s.Contains(e) {
		delete(*s, e)
	}
}

func (s *CountSet) RemoveMostAccessed() {

}

func (s *CountSet) RemoveLeastAccessed() {

}

func (s *CountSet) clear() {
	if s.Size() > MAX_SIZE {
		// remove at most 1/4 records if it's full
		for i := 0; i <= MAX_SIZE/8; i++ {
			s.RemoveMostAccessed()
			s.RemoveLeastAccessed()
		}
	}
}

func (s *CountSet) Contains(e string) bool {
	if map[string]int(*s)[e] > 0 {
		return false
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

func (g *HostGroup) Match(host string) bool {
	if g.Contains(host) {
		return true
	}

	ip := host
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

func (g *HostGroup) Log() {
	log.Println("domains", g.Domains, " ips: ", g.IPs, " cache: ", g.Cache)
}
