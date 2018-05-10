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