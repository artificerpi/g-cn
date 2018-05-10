# SS-Server with Web Content Filtering

## Apply filter

``` go
// add this in startup func `run` to start a goroutine to filter network
go filterTraffic()

// add this before dial up in handleConnection func
if AccessDenied(host) {
	log.Println("Blocking host", host)
	closed = true
	return
}
```

``` go
// Add filter configuration
var filterFile string
flag.StringVar(&filterFile, "filter", "filter.json", "specify filter config file")
```

## Systemd Unit File

``` systemd
[Unit]
Description=Go Shadowsocks server daemon
Wants=network-online.target
After=network.target network-online.target

[Service]
ExecStart=/usr/local/bin/ss-server -u -c /etc/shadowsocks/config.json -filter /etc/shadowsocks/filter.json
ExecReload=/bin/kill -HUP $MAINPID
PIDFile=/var/run/shadowsocks.pid
KillMode=process
Restart=on-failure
TimeoutStopSec=10

[Install]
WantedBy=multi-user.target
Alias=ss-server.service
``` 