package main

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/ianschenck/envflag"
	"image/png"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
)

var (
	addr string
	size int
)

type httpServer struct{}

func (s httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	str := r.FormValue("c")
	if str == "" {
		log.Print("empty content")
		return
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
	envflag.StringVar(&addr, "QRD_LISTEN", "127.0.0.1:9001", "listen address")
	envflag.IntVar(&size, "QRD_SIZE", 160, "barcode dimension")
}

func main() {
	var (
		l   net.Listener
		err error
	)
	envflag.Parse()

	if size < 100 {
		size = 100
	} else if size > 720 {
		size = 720
	}

	if addr[0] == '/' {
		l, err = net.Listen("unix", addr)
	} else {
		l, err = net.Listen("tcp", addr)
	}

	if err != nil {
		log.Println(err)
	}

	log.Printf("Start x-bar service at addr %s", addr)
	srv := new(httpServer)
	fcgi.Serve(l, srv)
}
