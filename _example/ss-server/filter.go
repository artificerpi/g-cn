package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var (
	whitelist  HostGroup
	blacklist  HostGroup
	greyIPs    []string
	strictMode bool = false

	// config FilterConfig
)

// nslookup  net.LookupAddr  context custom dns 8.8.8.8
// TODO reverse query to adjust whitelist and blacklist
// make a thread for this
// whitelist cache size 100, black 1000

func inWhiteList(addr string) bool {
	if whitelist.Contains(addr) {
		return true
	}

	return false
}

func inBlackList(addr string) bool {
	if blacklist.Contains(addr) {
		return true
	}

	return false
}

func AccessDenied(addr string) (rejected bool) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Println(err)
	}

	// dont filter dns query
	if port == "53" {
		return false
	}

	// FOR debug
	whitelist.Log()
	blacklist.Log()

	// default to false
	rejected = false

	if inWhiteList(host) {
		// add to white cache
		log.Println(addr, "is in whitelist.")
		return false
	} else if inBlackList(host) {
		// add to black cache
		log.Println(addr, "is in blacklist.")
		return true
	}

	if strictMode {
		// no filter for if the strategy is not strict
		// add to black cache
		rejected = true
		// TODO
		// addr check
		greyIPs = append(greyIPs, host)
		// update block cache check whether it's in white list
	}

	return
}

func filterTraffic() {
	whitelist = HostGroup{
		Domains: []string{
			"whitesite.com",
		},
		IPs:         CountSet{},
		CachedHosts: CountSet{},
	}

	blacklist = HostGroup{
		Domains: []string{
			"blacksite.com",
		},
		IPs:         CountSet{},
		CachedHosts: CountSet{},
	}

	whitelist.InitDomainIPs()
	blacklist.InitDomainIPs()

	log.Println("network traffic filter is started.")

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// every 5 minutes
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			go clearGreyIPs()
		case <-signalCh:
			ticker.Stop()
			break
		}
	}
}

func clearGreyIPs() {
	log.Println("grep ips", greyIPs)

	for _, ip := range greyIPs {
		if whitelist.Match(ip) {
			whitelist.CachedHosts.Add(ip)
		}
		if blacklist.Match(ip) {
			blacklist.CachedHosts.Add(ip)
		}
	}

	greyIPs = []string{}
}
