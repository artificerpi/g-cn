package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
)

var (
	whitelist HostGroup
	blacklist HostGroup
	// greyIPs    []string

	strictMode = false
)

// AccessDenied check wheter connection to the address is allowed
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
		// add to whitelist
		log.Println(addr, "is in whitelist.")
		whitelist.CachedHosts.Add(host)
		return false
	} else if inBlackList(host) {
		// add to blacklist
		blacklist.CachedHosts.Add(host)
		log.Println(addr, "is in blacklist.")
		return true
	}

	if strictMode {
		// no filter for if the strategy is not strict
		// add to black cache
		rejected = true

		// update block cache check whether it's in white list
		if whitelist.Match(host) {
			whitelist.CachedHosts.Add(host)
			return true
		}
	}

	return
}

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

func filterTraffic() {
	filterConfig, err := loadConfig(filterFile)
	if err != nil {
		panic(err)
	}

	whitelist = HostGroup{
		Domains:     filterConfig.Whitelist,
		IPs:         CountSet{},
		CachedHosts: CountSet{},
	}

	blacklist = HostGroup{
		Domains:     filterConfig.Blacklist,
		IPs:         CountSet{},
		CachedHosts: CountSet{},
	}

	strictMode = filterConfig.StrictMode

	whitelist.InitDomainIPs()
	blacklist.InitDomainIPs()

	log.Println("network traffic filter is started.")
}

func loadConfig(file string) (FilterConfig, error) {
	var c FilterConfig
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}

	json.Unmarshal(data, &c)
	return c, nil
}
