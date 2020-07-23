package wgpeer

// WGPeer structs
type WGPeer struct {
	AllowedIPS []string `json:"allowed_ips,omitempty"`
	KeepAlive  int      `json:"keepalive,omitempty"`
	PublicKey  string   `json:"public_key,omitempty"` // public key of peer`
	EndPoint   string   `json:"endpoint,omitempty"`
}
