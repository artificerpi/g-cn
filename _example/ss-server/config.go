package main

// FilterConfig store the configuration for filtering traffic
type FilterConfig struct {
	Blacklist  []string `json:"blacklist"`
	Whitelist  []string `json:"whitelist"`
	StrictMode bool     `json:"strict"`
}
