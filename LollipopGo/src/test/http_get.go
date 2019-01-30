package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var Gmap map[int]int
var TimeOutDSQ chan map[int]int

var (
	listen = flag.String("listen", ":8080", "listen address")
	dir    = flag.String("dir", ".", "directory to serve")
)

func init() {
	Gmap = make(map[int]int)
	TimeOutDSQ = make(chan map[int]int, 1000)
	return
}

func main() {

	Gmap[1000] = 1111
	Gmap[3000] = 3333

	TimeOutDSQ <- Gmap

	data := <-TimeOutDSQ
	fmt.Println(len(data))
	for k, v := range data {
		fmt.Println(k)
		fmt.Println(v)
	}

	//--------------------------------------------------------------------------
	flag.Parse()
	log.Printf("listening on %q...", *listen)
	log.Fatal(http.ListenAndServe(*listen, http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, ".wasm") {
			resp.Header().Set("content-type", "application/wasm")
		}
		http.FileServer(http.Dir(*dir)).ServeHTTP(resp, req)
	})))
}
