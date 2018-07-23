package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	opts := getOpts()

	// listen on local address
	ln, err := net.Listen("tcp", opts.locAddr)
	fatalErr(err)
	defer ln.Close()

	latch := sync.WaitGroup{}
	for {
		locConn, err := ln.Accept()
		fatalErr(err)

		// establish remote connection
		remConn, err := net.Dial("tcp", opts.remAddr)
		fatalErr(err)

		latch.Add(2)

		go copyWithDelay(opts.pktSize, opts.delay, locConn, remConn, &latch)
		go copyWithDelay(opts.pktSize, opts.delay, remConn, locConn, &latch)

	}

	latch.Wait()
}

func copyWithDelay(size int64, delay time.Duration, from io.ReadCloser, to io.WriteCloser, mon *sync.WaitGroup) {
	defer mon.Done()

	for {
		_, err := io.CopyN(to, from, size)
		if err != nil {
			from.Close()
			to.Close()
			break
		}
		time.Sleep(delay)
	}
}

type Opts struct {
	locAddr string
	remAddr string
	pktSize int64
	delay   time.Duration
}

func getOpts() *Opts {
	rv := Opts{}
	origUsage := flag.Usage
	flag.Usage = func() {
		origUsage()
		fmt.Fprintf(flag.CommandLine.Output(), "\n\nNOTE: this program will run until it's killed externally\n\n")
	}

	flag.StringVar(&rv.locAddr, "l", "localhost:8888", "Port to listen on")
	flag.StringVar(&rv.remAddr, "r", "", "Host:Port to proxy to (required)")
	flag.Int64Var(&rv.pktSize, "s", 512, "Packet Size in bytes")
	flag.DurationVar(&rv.delay, "d", 100*time.Millisecond, "Packet delay [Duration]")
	flag.Parse()

	if rv.remAddr == "" {
		flag.Usage()
		os.Exit(2)
	}
	return &rv
}

func fatalErr(err error) {
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(111)
	}
}
