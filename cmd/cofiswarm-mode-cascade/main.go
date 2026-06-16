package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/keepdevops/cofiswarm-mode-sdk/pkg/mode"
)

func main() {
	listen := flag.String("listen", "", "listen address override")
	cfg := flag.String("config", "", "mode yaml path")
	flag.Parse()
	if *cfg == "" {
		if v := os.Getenv("COFISWARM_MODE_CONFIG"); v != "" {
			*cfg = v
		} else {
			*cfg = "/etc/cofiswarm/mode-cascade/mode-cascade.yaml"
		}
	}
	srv, err := mode.NewServer("mode-cascade", *cfg)
	if err != nil {
		log.Fatal(err)
	}
	addr := srv.Addr()
	if *listen != "" {
		addr = *listen
	}
	log.Printf("mode-cascade listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Handler()))
}
