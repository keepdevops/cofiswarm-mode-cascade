package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/keepdevops/cofiswarm-mode-cascade/internal/bus"
	"github.com/keepdevops/cofiswarm-mode-sdk/pkg/mode"
	"github.com/keepdevops/cofiswarm-observer-sdk/pkg/buspresence"
	"github.com/keepdevops/cofiswarm-observer-sdk/pkg/servicecomponent"
)

const modeName = "cascade"

func main() {
	listen := flag.String("listen", "", "listen address override")
	cfg := flag.String("config", "", "mode yaml path")
	flag.Parse()
	if *cfg == "" {
		if v := os.Getenv("COFISWARM_MODE_CONFIG"); v != "" {
			*cfg = v
		} else {
			*cfg = "/etc/cofiswarm/mode-" + modeName + "/mode-" + modeName + ".yaml"
		}
	}
	srv, err := mode.NewServer("mode-"+modeName, *cfg)
	if err != nil {
		log.Fatal(err)
	}
	addr := srv.Addr()
	if *listen != "" {
		addr = *listen
	}

	// Optional: announce presence on the observer bus alongside the HTTP API (default-off).
	// COFISWARM_NATS_URL=nats://host:4222 enables it.
	var comp *servicecomponent.Component
	if url := os.Getenv("COFISWARM_NATS_URL"); url != "" {
		nc, cErr := servicecomponent.Connect(url, "cofiswarm-mode-"+modeName)
		if cErr != nil {
			log.Printf("bus connect %s: %v (running without presence)", url, cErr)
		} else {
			defer nc.Close()
			comp = servicecomponent.New(nc, "mode-"+modeName, "mode-"+modeName, bus.Routes(modeName))
			if sErr := comp.Start(); sErr != nil {
				log.Printf("bus start: %v (running without presence)", sErr)
				comp = nil
			} else {
				log.Printf("mode-%s announcing presence via %s", modeName, url)
			}
		}
	}

	// Carrier presence (broker-free, default-off via COFISWARM_BRIDGE_URL): appear in the
	// observer live roster over the zmq-bridge without needing a NATS broker. HTTP /healthz
	// + /v1/info remain the request/reply surface.
	ID := "mode-" + modeName
	stopPresence := buspresence.StartPresence(os.Getenv("COFISWARM_BRIDGE_URL"), ID, map[string]any{"name": ID})

	httpSrv := &http.Server{Addr: addr, Handler: srv.Handler()}
	go func() {
		log.Printf("mode-%s listening on %s", modeName, addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("mode-%s: server error: %v", modeName, err)
		}
	}()

	// On SIGINT/SIGTERM: say goodbye (flip offline now, not after the TTL) then drain HTTP.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Printf("mode-%s: shutting down", modeName)
	if comp != nil {
		comp.Shutdown() // goodbye -> offline
	}
	stopPresence() // carrier goodbye -> offline
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(shutCtx); err != nil {
		log.Printf("mode-%s: graceful shutdown: %v", modeName, err)
	}
}
