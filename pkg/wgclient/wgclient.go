package wgclient

import (
	"context"
	"fmt"

	"bitbucket.org/qubole/wireguard/pkg/wgpeer"
)

// IPSvc to fetch IP.
type IPSvc interface {
	Get(context.Context) (string, error)
}

// Store interface.
type Store interface {
	Get(context.Context, string) (interface{}, error)
	Set(context.Context, string, interface{}, ...int) error
}

// WGServer interface.
type WGServer interface {
	SSHAuthorizedKeys(context.Context) []string
	ServerPeers(context.Context) []wgpeer.WGPeer
}

// WGClient info.
type WGClient struct {
	ID         string   `json:"id,omitempty"`
	PrivateIP  string   `json:"private_ip,omitempty"` // private key of client
	PublicKey  string   `json:"public_key,omitempty"` // public key of client
	DNSServers []string `json:"dns_servers,omitempty"`
}

// Svc struct.
type Svc struct {
	store    Store
	ip       IPSvc
	wgServer WGServer
}

// NewSvc is svc constructor.
func NewSvc(store Store, ip IPSvc, wgServer WGServer) *Svc {
	return &Svc{store: store, ip: ip, wgServer: wgServer}
}

// CreateInput struct
type CreateInput struct {
	ID        string `json:"id,omitempty"`
	PublicKey string `json:"public_key,omitempty"`
}

// CreateOutput needs to be returned to wgclient
type CreateOutput struct {
	Client            *WGClient       `json:"client,omitempty"`
	SSHAuthorizedKeys []string        `json:"ssh_authorized_keys,omitempty"`
	Peers             []wgpeer.WGPeer `json:"peers,omitempty"`
}

// Create wgclient.
func (s *Svc) Create(ctx context.Context, in *CreateInput) (*CreateOutput, error) {
	v, err := s.store.Get(ctx, s.key(in.ID))
	if err != nil {
		return nil, err
	}

	var client *WGClient
	if v != nil {
		c, ok := v.(*WGClient)
		if !ok {
			return nil, fmt.Errorf("store:invalid_client")
		}

		client = c
	} else {
		pkey, err := s.store.Get(ctx, s.publicKey(in.PublicKey))
		if err != nil {
			return nil, fmt.Errorf("publickey:get:%v", err)
		}
		if pkey != nil {
			return nil, fmt.Errorf("publickey:duplicate")
		}

		i, err := s.ip.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("ip:get")
		}

		client = &WGClient{ID: in.ID, PublicKey: in.PublicKey, PrivateIP: i}
		err = s.store.Set(ctx, s.key(in.ID), client)
		if err != nil {
			return nil, fmt.Errorf("store:set:wgclient:%v", err)
		}

		err = s.store.Set(ctx, s.publicKey(in.PublicKey), client)
		if err != nil {
			return nil, fmt.Errorf("store:set:publickey:%v", err)
		}
	}

	return &CreateOutput{
		Client:            client,
		SSHAuthorizedKeys: s.wgServer.SSHAuthorizedKeys(ctx),
		Peers:             s.wgServer.ServerPeers(ctx),
	}, nil
}

func (s *Svc) key(id string) string {
	return fmt.Sprintf("wgclient:%s", id)
}

func (s *Svc) publicKey(pkey string) string {
	return fmt.Sprintf("pubkey:wgclient:%s", pkey)
}
