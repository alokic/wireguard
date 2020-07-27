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
	"bitbucket.org/qubole/wireguard/pkg/api"
	"bitbucket.org/qubole/wireguard/pkg/auth"
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
	JWTKey        string `json:"jwt_key,omitempty"`
}

func main() {
	cfg := &Config{
		Port:          4000,
		SSHPublicKey:  "test",
		SSHPrivateKey: "test",
		JWTKey:        "test",
	}
	fs := flag.NewFlagSet("server", flag.PanicOnError)
	fs.IntVar(&cfg.Port, "port", cfg.Port, "server port")
	fs.StringVar(&cfg.SSHPublicKey, "pubkey", cfg.SSHPublicKey, "ssh public key")
	fs.StringVar(&cfg.SSHPrivateKey, "privkey", cfg.SSHPrivateKey, "ssh private key")
	fs.StringVar(&cfg.JWTKey, "jwtkey", cfg.JWTKey, "jwt key")

	// set cache
	c := cache.NewMap()

	// set ipsvc
	ipsvc := ip.NewSvc(c)

	// set wireguard server service
	wgs := wgserver.NewSvc(c, ipsvc, cfg.SSHPublicKey, cfg.SSHPrivateKey)

	// set wireguard client service
	wgc := wgclient.NewSvc(c, ipsvc, wgs)

	// set jwt
	jwt := auth.NewJWT(cfg.JWTKey)

	// set REST api handler.
	rapi := &api.REST{WGC: wgc, WGS: wgs}

	// set  routes
	router := router.CreateRouter("gorilla")
	setRoutes(router, rapi, jwt)

	// start server

	//// setup workgroup
	g := workgroup.Group{}

	//// set interrupt handler
	g.Add(func(stop <-chan struct{}) error {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)

		<-sigChan

		return fmt.Errorf("Interrupted")
	})

	//// create server object
	s := server.New(server.Logger("info", "app", "wireguard", "type", "server"), server.Port(cfg.Port), server.NotFoundHandler(router))
	for _, fn := range s.Runnables() {
		g.Add(fn)
	}

	//// run workgroup
	g.Run()
}

func setRoutes(r router.Router, rapi *api.REST, jwt *auth.JWT) {
	r.Handle("post", "/wgclient", jwt.HTTPMiddleware(rapi.ClientGererateConfig()))

	r.Handle("get", "/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "OK",
		})
	}))
}
