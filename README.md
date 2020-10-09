# go-soxy
Network proxy to test slow connection pipes written in go

# Usage:
```
go get github.com:xelasov/go-soxy
```

Then

```
go-soxy 
Usage of go-soxy:
  -d duration
        Packet delay [Duration] (default 100ms)
  -l string
        Port to listen on (default "localhost:8888")
  -r string
        Host:Port to proxy to (required)
  -s int
        Packet Size in bytes (default 512)


NOTE: this program will run until it's killed externally
```
