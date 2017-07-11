package main

import (
	"flag"
	"image/png"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	envcfg "github.com/wealthworks/envflagset"
)

var (
	fs        *flag.FlagSet
	addr      string
	dimension int
	version   = "dev"
)

type httpServer struct{}

func (s httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	str := r.FormValue("c")
	if str == "" {
		log.Print("empty content")
		return
	}

	size := validSize(dimension)
	if s := r.FormValue("s"); s != "" {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			size = validSize(int(i))
		}
	}

	w.Header().Set("Content-Type", "image/png")

	qrcode, err := qr.Encode(str, qr.Q, qr.Auto) // L,M,Q,H
	if err != nil {
		log.Println(err)
	} else {
		qrcode, err = barcode.Scale(qrcode, size, size)
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("generated barcode: %s", qrcode.Content())
			png.Encode(w, qrcode)
		}
	}
}

func init() {
	fs = envcfg.New("qrd", version)
	fs.StringVar(&addr, "listen", "127.0.0.1:9001", "listen address")
	fs.IntVar(&dimension, "dimension", 160, "barcode dimension")
}

func main() {
	var (
		l   net.Listener
		err error
	)
	envcfg.Parse()

	if addr[0] == '/' {
		l, err = net.Listen("unix", addr)
	} else {
		l, err = net.Listen("tcp", addr)
	}

	if err != nil {
		log.Println(err)
	}

	log.Printf("Start fcgi service at addr %s", addr)
	srv := new(httpServer)
	fcgi.Serve(l, srv)
}

func validSize(dimension int) int {
	if dimension < 60 {
		return 60
	} else if dimension > 720 {
		return 720
	}
	return dimension
}
