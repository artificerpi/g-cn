package main

import (
	"log"
	"net"
	"strings"
)

var (
	blacklist  []string
	whitelist  []string
	strictMode bool = false
)

type HostGroup struct {
	Domains     []string
	IPs         []string
	CachedHosts []string
}

func (g *HostGroup) Contains(host string) bool {
	return false
}

func init() {
	whitelist = []string{
		"google.com.hk",
	}

	blacklist = []string{
		"baidu.com",
	}

	var whiteIPs []string
	for _, domain := range whitelist {
		// ipList = append(ipList, addIP(white)...)
		ips, _ := net.LookupIP(domain)
		for _, ip := range ips {
			whiteIPs = append(whiteIPs, ip.String())
		}
	}
	whitelist = append(whitelist, whiteIPs...)
}

// nslookup  net.LookupAddr  context custom dns 8.8.8.8
// TODO reverse query to adjust whitelist and blacklist
// make a thread for this
// whitelist cache size 100, black 1000

func AccessDenied(addr string) (rejected bool) {
	log.Println(whitelist)
	log.Println(blacklist)
	// no filter for if the strategy is not strict
	if strictMode {
		rejected = true
	} else {
		rejected = false
	}

	if inWhiteList(addr) {
		log.Println(addr, "is in whitelist.")
		return false
	} else if inBlackList(addr) {
		log.Println(addr, "is in blacklist.")
		rejected = true
	}

	return
}

func inWhiteList(addr string) bool {
	for _, host := range whitelist {
		if strings.Contains(addr, host) {
			return true
		}
	}

	return false
}

func inBlackList(addr string) bool {
	for _, host := range blacklist {
		if strings.Contains(addr, host) {
			return true
		}
	}

	return false
}
