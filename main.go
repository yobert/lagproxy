package main

import (
	"net"
	"log"
	"io"
	"math/rand"
	"time"
	"flag"
)

var (
	from = flag.String("from", "localhost:8082", "listen on address")
	to = flag.String("to", "localhost:8081", "proxy to address")
	min = flag.Int("min", 10, "minimum ms to sleep per block")
	max = flag.Int("max", 100, "maximum ms to sleep per block")
	block = flag.Int("block", 64, "block size in bytes")
)

func main() {
	flag.Parse()

	listen, err := net.Listen("tcp", *from)
	if err != nil {
		log.Fatal(err)
	}

	for {
		sock, err := listen.Accept()
		if err != nil {
			log.Println(err)
		} else {
			go proxy(sock)
		}
	}
}

func proxy(a net.Conn) {
	b, err := net.Dial("tcp", *to)
	if err != nil {
		log.Println(err)
		a.Close()
		return
	}

	go send(a, b)
	send(b, a)
}

func send(w net.Conn, r net.Conn) {
	defer w.Close()
	slowcopy(w, r)
}

// randomly slow version of io.copy
func slowcopy(dst io.Writer, src io.Reader) (written int64, err error) {
	buf := make([]byte, *block)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(*max - *min) + *min)) // make it slow
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}
