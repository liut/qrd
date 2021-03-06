package main

import (
	"flag"
	"image/png"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"strconv"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	envcfg "github.com/wealthworks/envflagset"
)

const (
	maxSize = 720
	minSize = 60
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
		log.Print("empty content", r.RequestURI)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size := dimension
	if s := r.FormValue("s"); s != "" {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			size = validSize(int(i))
		}
	}

	w.Header().Set("Content-Type", "image/png")

	err := genQRcode(w, str, size)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func genQRcode(w io.Writer, text string, size int) error {

	qrcode, err := qr.Encode(text, qr.Q, qr.Auto) // L,M,Q,H
	if err != nil {
		log.Printf("encode ERR %s", err)
		return err
	}
	qrcode, err = barcode.Scale(qrcode, size, size)
	if err != nil {
		log.Printf("scale ERR %s", err)
		return err
	}

	// log.Printf("generated barcode: %s", qrcode.Content())
	return png.Encode(w, qrcode)
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
	dimension = validSize(dimension)

	if addr[0] == '/' {
		l, err = net.Listen("unix", addr)
	} else {
		l, err = net.Listen("tcp", addr)
	}

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Start fcgi service at addr %s", addr)
	srv := new(httpServer)
	fcgi.Serve(l, srv)
}

func validSize(dimension int) int {
	if dimension < minSize {
		return minSize
	} else if dimension > maxSize {
		return maxSize
	}
	return dimension
}
