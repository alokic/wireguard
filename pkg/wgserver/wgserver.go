package wgserver

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

// WGServer info.
type WGServer struct {
	ID        string `json:"id,omitempty"`
	PrivateIP string `json:"private_ip,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"` // public IP or cname accessible by client
	PublicKey string `json:"public_key,omitempty"`
}

// Svc struct.
type Svc struct {
	store         Store
	ip            IPSvc
	sshPublicKey  string `json:"ssh_public_key,omitempty"`
	sshPrivateKey string `json:"ssh_private_key,omitempty"`
}

// NewSvc is svc constructor.
func NewSvc(store Store, ip IPSvc, sshPublicKey, sshPrivateKey string) *Svc {
	return &Svc{store: store, ip: ip, sshPublicKey: sshPublicKey, sshPrivateKey: sshPrivateKey}
}

// CreateInput struct
type CreateInput struct {
}

// CreateOutput needs to be returned to wgserver
type CreateOutput struct {
}

// Create wgserver.
func (s *Svc) Create(ctx context.Context, in *CreateInput) (*CreateOutput, error) {
	return &CreateOutput{}, nil
}

// CronStorePeers scrape new peers (clients or servers) with wg tools and put on redis.
func (s *Svc) CronStorePeers(ctx context.Context) error {
	return nil
}

// CronSyncPeersFromStore syncs new peers (clients or servers) from store and add locally usin wg add tool.
func (s *Svc) CronSyncPeersFromStore(ctx context.Context) error {
	return nil
}

// ServerPeers returns list of server peers
func (s *Svc) ServerPeers(ctx context.Context) []wgpeer.WGPeer {
	return []wgpeer.WGPeer{{PublicKey: s.sshPublicKey, EndPoint: "1.1.1.1"}}
}

// SSHAuthorizedKeys returns list of valid sshpublickeys
func (s *Svc) SSHAuthorizedKeys(ctx context.Context) []string {
	return []string{s.sshPublicKey}
}

func (s *Svc) key(id string) string {
	return fmt.Sprintf("wgserver:%s", id)
}
