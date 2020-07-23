package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"bitbucket.org/qubole/wireguard/internal/router"
	"bitbucket.org/qubole/wireguard/internal/server"
	"bitbucket.org/qubole/wireguard/internal/workgroup"
	"bitbucket.org/qubole/wireguard/pkg/cache"
	"bitbucket.org/qubole/wireguard/pkg/ip"
	"bitbucket.org/qubole/wireguard/pkg/wgclient"
	"bitbucket.org/qubole/wireguard/pkg/wgserver"
)

// Config struct.
type Config struct {
	Port          int    `json:"port,omitempty"`
	SSHPrivateKey string `json:"ssh_private_key,omitempty"`
	SSHPublicKey  string `json:"ssh_public_key,omitempty"`
}

func main() {
	cfg := &Config{
		Port:          4000,
		SSHPublicKey:  "test",
		SSHPrivateKey: "test",
	}
	fs := flag.NewFlagSet("server", flag.PanicOnError)
	fs.IntVar(&cfg.Port, "port", cfg.Port, "server port")
	fs.StringVar(&cfg.SSHPublicKey, "pubkey", cfg.SSHPublicKey, "ssh public key")
	fs.StringVar(&cfg.SSHPrivateKey, "privkey", cfg.SSHPrivateKey, "ssh private key")

	c := cache.Map{}
	ipsvc := ip.NewSvc(c)

	wgs := wgserver.NewSvc(c, ipsvc, cfg.SSHPublicKey, cfg.SSHPrivateKey)

	h := &handlers{cfg: cfg, c: wgclient.NewSvc(c, ipsvc, wgs)}

	router := router.CreateRouter("gorilla")
	setRoutes(router, h)

	// setup workgroup
	g := workgroup.Group{}
	// shutdown http
	g.Add(func(stop <-chan struct{}) error {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		<-sigChan

		return fmt.Errorf("Interrupted")
	})

	s := server.New(server.Logger("info", "app", "wireguard", "type", "server"), server.Port(cfg.Port), server.NotFoundHandler(router))
	for _, fn := range s.Runnables() {
		g.Add(fn)
	}

	g.Run()
}

func setRoutes(r router.Router, h *handlers) {
	r.Handle("post", "/wgclient", h.getClientProfile())

	r.Handle("get", "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
		})
	}))
}
